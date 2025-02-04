package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/masintxi/blog_aggregator/internal/database"
)

func handleAgg(s *state, cmd command) error {
	err := checkArgs(cmd, 1)
	if err != nil {
		return err
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("failed to parse the duration <%v>: %w", cmd.args[0], err)
	}

	log.Printf("Collecting feeds every %s...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer s.cleanup()

	numScrappers := 4
	scrapeBatch(ctx, s, numScrappers)

	for {
		log.Printf("Aggregator waiting for next event...")
		select {
		case <-ticker.C:
			log.Printf("Ticker triggered, starting new batch...")
			scrapeBatch(ctx, s, numScrappers)
		case newFeed, ok := <-s.newFeeds:
			log.Printf("Received something on newFeeds channel, ok=%v", ok)
			if !ok {
				return nil
			}
			log.Printf("Processing new feed: %s", newFeed.Name)
			go func(feed database.Feed) {
				if err := scrapeFeed(ctx, s.db, feed); err != nil {
					log.Printf("Failed to scrape new feed: %v", err)
				}
			}(newFeed)
		case <-ctx.Done():
			log.Printf("Context cancelled, shutting down...")
			return nil
		}
	}
}

func scrapeBatch(ctx context.Context, s *state, numWorkers int) {
	log.Printf("Starting new batch scrape with %d workers", numWorkers)
	feeds, err := s.db.ListFeeds(ctx)
	if err != nil {
		log.Printf("Could not fetch feeds: %v", err)
		return
	}

	var wg sync.WaitGroup
	feedChan := make(chan database.Feed)
	remainingFeeds := make(chan database.Feed, len(feeds))

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for feedInfo := range feedChan {
				fmt.Printf("Worker %d: going for the <%s> feed\n", workerID, feedInfo.Url)

				err := scrapeFeed(ctx, s.db, feedInfo)
				if err != nil {
					remainingFeeds <- feedInfo
					continue
				}

				//fmt.Printf("Worker %d: fetched the posts of the <%s> feed\n", workerID, feedInfo.Name)
			}
		}(i)
	}

	for _, feed := range feeds {
		feedChan <- feed
	}
	close(feedChan)

	wg.Wait()

	close(remainingFeeds)

	// Process any remaining feeds
	if len(remainingFeeds) > 0 {
		log.Printf("%d feeds need to be reprocessed:", len(remainingFeeds))
		var remainingFeedsList []database.Feed
		for feed := range remainingFeeds {
			remainingFeedsList = append(remainingFeedsList, feed)
		}

		for i, remFeed := range remainingFeedsList {
			log.Printf("[%d/%d] Failed feed: %s (URL: %s)",
				i+1, len(remainingFeedsList), remFeed.Name, remFeed.Url)
		}
	}
}

func scrapeFeed(ctx context.Context, db *database.Queries, feedInfo database.Feed) error {
	err := db.MarkFeedFetched(ctx, feedInfo.ID)
	if err != nil {
		log.Printf("Could not mark feed <%s>: %v", feedInfo.Name, err)
		return err
	}
	//fmt.Printf("feed <%s> marked as fetched\n", feedInfo.Url)

	feed, err := fetchFeed(ctx, feedInfo.Url)
	if err != nil {
		log.Printf("Could not fetch the feed <%s>: %v", feedInfo.Name, err)
		return err
	}

	postCount := 0
	for _, item := range feed.Channel.Item {
		parsedPubDate := parsePubDate(item.PubDate)

		post, err := db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       sql.NullString{String: item.Title, Valid: true},
			Description: sql.NullString{String: item.Description, Valid: true},
			Url:         item.Link,
			PublishedAt: parsedPubDate,
			FeedID:      feedInfo.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("something went wrong... %v", err)
			continue
		}
		log.Printf("post <%s> created successfully", post.Url)
		postCount++
	}

	if postCount == 0 {
		log.Printf("No new posts detected in feed <%s>", feedInfo.Name)
		return nil
	}
	log.Printf("Added %d new post from <%s>", postCount, feedInfo.Name)
	return nil
}

func parsePubDate(pubDate string) sql.NullTime {
	formats := []string{time.RFC1123Z, time.RFC3339, time.RFC1123}
	for _, format := range formats {
		if t, err := time.Parse(format, pubDate); err == nil {
			return sql.NullTime{Time: t, Valid: true}
		}
	}
	return sql.NullTime{}
}

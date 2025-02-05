package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/masintxi/blog_aggregator/internal/database"
)

func getFeedsList(ctx context.Context, s *state, cmd command) error {
	feeds, err := s.db.ListFeeds(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve the list of feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found")
		return nil
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserById(ctx, feed.UserID)
		if err != nil {
			return fmt.Errorf("could not retrive the user: %w", err)
		}
		printFeedInfo(feed, user.Name)
		fmt.Println("------------------------------------------------------")
	}
	return nil
}

func handleNewFeed(ctx context.Context, s *state, cmd command, user database.User) error {
	err := checkArgs(cmd, 2)
	if err != nil {
		return err
	}

	fName := cmd.args[0]
	fUrl := cmd.args[1]

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      fName,
		Url:       fUrl,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("could not add feed <%v> into the database: %v", fName, err)
	}
	fmt.Printf("Feed <%s> created successfully:\n", fName)
	printFeedInfo(feed, user.Name)

	fFollow, err := generateFeedFollow(ctx, s, fUrl)
	if err != nil {
		return fmt.Errorf("failed following the feed: %w", err)
	}
	fmt.Printf("User <%s> is now following the feed <%s>\n", fFollow.UserName, fFollow.FeedName)

	originalOutput := log.Writer()
	log.SetOutput(io.Discard)
	err = scrapeFeed(context.Background(), s.db, feed)
	log.SetOutput(originalOutput)

	if err != nil {
		log.Printf("Failed to scrape new feed <%s>: %v", feed.Name, err)
	} else {
		log.Printf("Succesfully scrapped the new feed <%s>", feed.Name)
	}

	return nil
}

func printFeedInfo(feed database.Feed, userName string) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", userName)
	//fmt.Printf("* LastFetchedAt: %v\n", feed.LastFetchedAt.Time)
}

func deleteFeed(ctx context.Context, s *state, cmd command) error {
	err := checkArgs(cmd, 1)
	if err != nil {
		return err
	}

	fUrl := cmd.args[0]

	err = s.db.DeleteFeedByURL(ctx, fUrl)
	if err != nil {
		fmt.Printf("could not delete the <%s> feed: %v\n", fUrl, err)
	}
	fmt.Printf("feed <%s> successfully deleted\n", fUrl)
	return nil
}

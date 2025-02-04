package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/masintxi/blog_aggregator/internal/database"
)

func getFeedsList(s *state, cmd command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to retrieve the list of feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found")
		return nil
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("could not retrive the user: %w", err)
		}
		printFeedInfo(feed, user.Name)
		fmt.Println("------------------------------------------------------")
	}
	return nil
}

func handleNewFeed(s *state, cmd command, user database.User) error {
	err := checkArgs(cmd, 2)
	if err != nil {
		return err
	}

	fName := cmd.args[0]
	fUrl := cmd.args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
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

	fFollow, err := generateFeedFollow(s, fUrl)
	if err != nil {
		return fmt.Errorf("failed following the feed: %w", err)
	}
	fmt.Printf("User <%s> is now following the feed <%s>\n", fFollow.UserName, fFollow.FeedName)

	return nil
}

func printFeedInfo(feed database.Feed, userName string) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", userName)
	fmt.Printf("* LastFetchedAt: %v\n", feed.LastFetchedAt.Time)
}

func deleteFeed(s *state, cmd command) error {
	err := checkArgs(cmd, 1)
	if err != nil {
		return err
	}

	fUrl := cmd.args[0]

	err = s.db.DeleteFeedByURL(context.Background(), fUrl)
	if err != nil {
		fmt.Printf("could not delete the <%s> feed: %v\n", fUrl, err)
	}
	fmt.Printf("feed <%s> successfully deleted\n", fUrl)
	return nil
}

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/masintxi/blog_aggregator/internal/database"
)

func handleAgg(s *state, cmd command) error {
	// fUrl, err := checkArgs(cmd)
	// if err != nil {
	// 	return err
	// }

	fUrl := "https://www.wagslane.dev/index.xml"

	var info *RSSFeed
	info, err := fetchFeed(context.Background(), fUrl)
	if err != nil {
		return err
	}

	fmt.Println(info)
	return nil
}

func getFeedsList(s *state, cmd command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to retrieve the list of feeds")
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("could not retrive the user: %w", err)
		}
		printFeedInfo(feed, user.Name)
		fmt.Println("------------------------------------------------------")
		// fmt.Print("*")
		// fmt.Printf(" Name: %s |", feed.Name)
		// fmt.Printf(" URL: %s |", feed.Url)
		// userName, _ := s.db.GetUserById(context.Background(), feed.UserID)
		// fmt.Printf(" Created by: %s", userName)
		// if userName == s.cfg.CurrentUserName {
		// 	fmt.Print(" (current user)")
		// }
		// fmt.Println("")
	}
	return nil
}

func handleNewFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("not enough arguments received for the <%v> command, 2 required", cmd.name)
	}

	fName := cmd.args[0]
	fUrl := cmd.args[1]

	feedArgs := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      fName,
		Url:       fUrl,
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), feedArgs)
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
}

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/masintxi/blog_aggregator/internal/database"
)

func generateFeedFollow(s *state, fUrl string) (database.CreateFeedFollowRow, error) {
	fFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      s.cfg.CurrentUserName,
		Url:       fUrl,
	})

	if err != nil {
		return database.CreateFeedFollowRow{}, err
	}

	feed, err := s.db.GetFeedByURL(context.Background(), fUrl)
	if err == nil {
		select {
		case s.newFeeds <- feed:
			log.Printf("Successfully sent feed %s to channel", feed.Name)
		default:
			log.Printf("Warning: Channel full, feed <%s> will be processed in next batch", feed.Name)
		}
	} else {
		log.Printf("Warning: Feed <%s> not found and will be processed in next batch: %v", feed.Name, err)
	}

	return fFollow, nil
}

func handleNewFollow(s *state, cmd command, user database.User) error {
	err := checkArgs(cmd, 1)
	if err != nil {
		return err
	}

	fUrl := cmd.args[0]

	fFollow, err := generateFeedFollow(s, fUrl)
	if err != nil {
		return fmt.Errorf("failed following the feed: %w", err)
	}

	fmt.Printf("User <%s> is now following the feed <%s>\n", fFollow.UserName, fFollow.FeedName)
	return nil
}

func handleFollowsForUser(s *state, cmd command, user database.User) error {
	userFollows, err := s.db.GetFeedFollowsForUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("could not get the follow list: %w", err)
	}

	if len(userFollows) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}

	fmt.Printf("Feed follows for user <%s>:\n", s.cfg.CurrentUserName)
	for _, userFollow := range userFollows {
		fmt.Printf("* %s\n", userFollow.FeedName)
	}
	return nil
}

func handleUnfollow(s *state, cmd command, user database.User) error {
	err := checkArgs(cmd, 1)
	if err != nil {
		return err
	}

	fUrl := cmd.args[0]

	err = s.db.DeleteFollowForUser(context.Background(), database.DeleteFollowForUserParams{
		Url:  fUrl,
		Name: user.Name,
	})
	if err != nil {
		return fmt.Errorf("failed unfollowing the feed: %w", err)
	}

	fmt.Printf("user <%s> unfollowed the feed <%s>\n", user.Name, fUrl)

	return nil
}

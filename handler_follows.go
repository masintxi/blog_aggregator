package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/masintxi/blog_aggregator/internal/database"
)

func generateFeedFollow(s *state, fUrl string) (database.CreateFeedFollowRow, error) {
	return s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      s.cfg.CurrentUserName,
		Url:       fUrl,
	},
	)
}

func handleNewFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("not enough arguments received for the <%v> command, 1 required", cmd.name)
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
	if len(cmd.args) < 1 {
		return fmt.Errorf("not enough arguments received for the <%v> command, 1 required", cmd.name)
	}

	fUrl := cmd.args[0]

	err := s.db.DeleteFollowForUser(context.Background(), database.DeleteFollowForUserParams{
		Url:  fUrl,
		Name: user.Name,
	})
	if err != nil {
		return fmt.Errorf("failed unfollowing the feed: %w", err)
	}

	fmt.Printf("user <%s> unfollowed the feed <%s>\n", user.Name, fUrl)

	return nil
}

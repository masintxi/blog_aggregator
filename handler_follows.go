package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/masintxi/blog_aggregator/internal/database"
)

func generateFeedFollow(ctx context.Context, s *state, fUrl string) (database.CreateFeedFollowRow, error) {
	return s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      s.cfg.CurrentUserName,
		Url:       fUrl,
	})
}

func handleNewFollow(ctx context.Context, s *state, cmd command, user database.User) error {
	err := checkArgs(cmd, 1)
	if err != nil {
		return err
	}

	fUrl := cmd.args[0]

	fFollow, err := generateFeedFollow(ctx, s, fUrl)
	if err != nil {
		return fmt.Errorf("failed following the feed: %w", err)
	}

	fmt.Printf("User <%s> is now following the feed <%s>\n", fFollow.UserName, fFollow.FeedName)
	return nil
}

func handleFollowsForUser(ctx context.Context, s *state, cmd command, user database.User) error {
	userFollows, err := s.db.GetFeedFollowsForUser(ctx, s.cfg.CurrentUserName)
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

func handleUnfollow(ctx context.Context, s *state, cmd command, user database.User) error {
	err := checkArgs(cmd, 1)
	if err != nil {
		return err
	}

	fUrl := cmd.args[0]

	err = s.db.DeleteFollowForUser(ctx, database.DeleteFollowForUserParams{
		Url:  fUrl,
		Name: user.Name,
	})
	if err != nil {
		return fmt.Errorf("failed unfollowing the feed: %w", err)
	}

	fmt.Printf("user <%s> unfollowed the feed <%s>\n", user.Name, fUrl)

	return nil
}

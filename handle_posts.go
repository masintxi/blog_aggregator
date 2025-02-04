package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/masintxi/blog_aggregator/internal/database"
)

func handleBrowsePosts(s *state, cmd command, user database.User) error {
	pLimit := 2
	if len(cmd.args) > 0 {
		if intVal, err := strconv.Atoi(cmd.args[0]); err == nil {
			pLimit = intVal
		}
	}

	userPosts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(pLimit),
	})

	if err != nil {
		fmt.Printf("could not get posts for <%s>: %v", user.Name, err)
	}

	if len(userPosts) == 0 {
		fmt.Printf("no posts found for user <%s>\n", user.Name)
		return nil
	}

	fmt.Printf("Last %v posts for user <%v>:\n", pLimit, user.Name)
	for i, post := range userPosts {
		fmt.Printf("%d* (%s) from %s:\n", i+1, post.PublishedAt.Time.Format("02/Mar/2006 15:04:05"), post.FeedName)
		fmt.Printf(" -Title: %s\n", post.Title.String)
		fmt.Printf(" -Description:\n%s\n", post.Description.String)
		fmt.Printf(" -Link: %s\n", post.Url)
		fmt.Println()
	}

	return nil
}

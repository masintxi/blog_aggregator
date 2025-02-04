package main

import (
	"context"
	"fmt"

	"github.com/masintxi/blog_aggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}

func checkArgs(cmd command, numArgs int) error {
	if len(cmd.args) < numArgs {
		return fmt.Errorf("the <%v> command requires %d arguments", cmd.name, numArgs)
	}
	return nil
}

func (s *state) cleanup() {
	close(s.done)
	close(s.newFeeds)
}

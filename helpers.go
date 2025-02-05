package main

import (
	"context"
	"fmt"

	"github.com/masintxi/blog_aggregator/internal/database"
)

func middlewareLoggedIn(handler func(ctx context.Context, s *state, cmd command, user database.User) error) func(context.Context, *state, command) error {
	return func(ctx context.Context, s *state, cmd command) error {
		user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(ctx, s, cmd, user)
	}
}

func checkArgs(cmd command, numArgs int) error {
	if len(cmd.args) < numArgs {
		return fmt.Errorf("the <%v> command requires %d arguments", cmd.name, numArgs)
	}
	return nil
}

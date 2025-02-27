package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/masintxi/blog_aggregator/internal/database"
)

func handlerLogin(ctx context.Context, s *state, cmd command) error {
	err := checkArgs(cmd, 1)
	if err != nil {
		return err
	}
	userName := cmd.args[0]

	_, err = s.db.GetUser(ctx, userName)
	if err != nil {
		return fmt.Errorf("the user <%v> does not exists in the database", userName)
	}

	if err := s.cfg.SetUser(userName); err != nil {
		return err
	}

	fmt.Printf("Welcome %s\n", s.cfg.CurrentUserName)
	return nil
}

func registerUser(ctx context.Context, s *state, cmd command) error {
	err := checkArgs(cmd, 1)
	if err != nil {
		return err
	}
	userName := cmd.args[0]

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      userName,
	}

	user, err := s.db.CreateUser(ctx, userParams)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return fmt.Errorf("the user <%v> already exists in the database", userName)
		}
		return fmt.Errorf("error creating user <%v> in database: %v", userName, err)
	}

	if err := s.cfg.SetUser(userName); err != nil {
		return err
	}
	fmt.Printf("User created. ID: %v - Name: %v\n", user.ID, user.Name)
	return nil
}

func getUsersList(ctx context.Context, s *state, cmd command) error {
	users, err := s.db.ListUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrive the list of users")
	}

	if len(users) == 0 {
		fmt.Println("No users found.")
		return nil
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func resetUsersTable(ctx context.Context, s *state, cmd command) error {
	err := s.db.ResetDB(ctx)
	if err != nil {
		return fmt.Errorf("error reseting the database: %v", err)
	}
	fmt.Println("database successfully reseted")
	return nil
}

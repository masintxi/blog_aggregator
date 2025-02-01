package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/masintxi/blog_aggregator/internal/database"
)

func checkArgs(cmd command) (string, error) {
	if len(cmd.args) == 0 {
		return "", fmt.Errorf("no arguments received for the <%v> command", cmd.name)
	}
	return cmd.args[0], nil
}

func handlerLogin(s *state, cmd command) error {
	userName, err := checkArgs(cmd)
	if err != nil {
		return err
	}

	_, err = s.db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("the user <%v> does not exists in the database", userName)
	}

	if err := s.cfg.SetUser(userName); err != nil {
		return err
	}

	fmt.Printf("Welcome %s\n", s.cfg.CurrentUserName)
	return nil
}

func registerUser(s *state, cmd command) error {
	userName, err := checkArgs(cmd)
	if err != nil {
		return err
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      userName,
	}

	user, err := s.db.CreateUser(context.Background(), userParams)
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

func getUsersList(s *state, cmd command) error {
	users, err := s.db.ListUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to retrive the list of users")
	}

	for _, user := range users {
		if user == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}
	return nil
}

func resetUsersTable(s *state, cmd command) error {
	err := s.db.ResetDB(context.Background())
	if err != nil {
		return fmt.Errorf("error reseting the database: %v", err)
	}
	fmt.Println("database successfully reseted")
	return nil
}

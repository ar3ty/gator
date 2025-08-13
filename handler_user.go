package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ar3ty/gator/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(st *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.name)
	}
	userName := cmd.args[0]

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      userName,
	}
	user, err := st.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	fmt.Printf("User was created -\n")
	printUser(user)

	err = st.cfg.SetUser(params.Name)
	if err != nil {
		return fmt.Errorf("current user setting failed: %w", err)
	}
	fmt.Println("User has been successfully set.")
	return nil
}

func handlerLogin(st *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}
	userName := cmd.args[0]

	_, err := st.db.GetUser(context.Background(), userName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user doesn't exist: %w", err)
		}
		return fmt.Errorf("cannot find user: %w", err)
	}

	err = st.cfg.SetUser(userName)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}
	fmt.Println("User has been successfully set.")
	return nil
}

func handlerListUsers(st *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("usage: %s", cmd.name)
	}

	users, err := st.db.GetUsers(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no user exists: %w", err)
		}
		return fmt.Errorf("cannot get users: %w", err)
	}

	for _, user := range users {
		fmt.Printf("* %s", user.Name)
		if user.Name == st.cfg.CurrentUserName {
			fmt.Print(" (current)")
		}
		fmt.Println()
	}

	return nil
}

func printUser(user database.User) {
	fmt.Printf("ID:   	%v\n", user.ID)
	fmt.Printf("Name: 	%v\n", user.Name)
}

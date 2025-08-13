package main

import (
	"context"
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
	fmt.Printf("User was created -\nID: %v,\nname: %s.\n", user.ID, user.Name)

	err = st.cfg.SetUser(params.Name)
	if err != nil {
		return fmt.Errorf("current user setting failed: %w", err)
	}
	fmt.Println("User has been successfully set.")
	return nil
}

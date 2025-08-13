package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func handlerGetUsers(st *state, cmd command) error {
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

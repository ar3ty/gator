package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

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

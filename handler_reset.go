package main

import (
	"context"
	"fmt"
)

func handlerReset(st *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("usage: %s", cmd.name)
	}

	err := st.db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("cannot reset db:\n %w", err)
	}

	fmt.Println("db was successfully reset")
	return nil
}

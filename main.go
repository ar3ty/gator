package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ar3ty/gator/internal/config"
	"github.com/ar3ty/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func middlewareLoggedIn(handler func(st *state, cmd command, user database.User) error) func(*state, command) error {
	return func(st *state, cmd command) error {
		currentUserName := st.cfg.CurrentUserName
		user, err := st.db.GetUser(context.Background(), currentUserName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("no login user exists: %w", err)
			}
			return fmt.Errorf("cannot get login user: %w", err)
		}
		return handler(st, cmd, user)
	}
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error opening db connection: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	programState := state{
		db:  dbQueries,
		cfg: &cfg,
	}
	cmds := commands{
		handlersMap: map[string]func(*state, command) error{},
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAggregate)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerListFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))

	inline := os.Args
	if len(inline) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	command := command{
		name: inline[1],
		args: inline[2:],
	}

	err = cmds.run(&programState, command)
	if err != nil {
		log.Fatalf("%v", err)
	}
}

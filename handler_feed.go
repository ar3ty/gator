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

func handlerAddFeed(st *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.name)
	}
	currentUser := st.cfg.CurrentUserName
	user, err := st.db.GetUser(context.Background(), currentUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no user exists: %w", err)
		}
		return fmt.Errorf("cannot get user: %w", err)
	}

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}

	feed, err := st.db.CreateFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}

	fmt.Println("Feed has been successfully created.")
	printFeed(feed)
	return nil
}

func handlerListFeeds(st *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("usage: %s", cmd.name)
	}

	feeds, err := st.db.ListFeeds(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no feed exists: %w", err)
		}
		return fmt.Errorf("cannot get feed: %w", err)
	}

	for i, feed := range feeds {
		fmt.Printf("- - - Feed %d - - -\n", i)
		printFeed(feed)
		name, err := st.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get username by id: %w", err)
		}
		fmt.Printf("Created by %s\n", name)
		fmt.Printf("- - - - - - - - - -\n")
	}

	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("ID:   		%v\n", feed.ID)
	fmt.Printf("Name: 		%s\n", feed.Name)
	fmt.Printf("URL:   		%s\n", feed.Url)
	fmt.Printf("UserID: 	%v\n", feed.UserID)
}

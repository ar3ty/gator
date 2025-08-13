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

func handlerAddFeed(st *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.name)
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
	printFeed(feed, user)

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = st.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("cannot create feed follow: %w", err)
	}

	println("New feed follow is created:")
	printFeedFollow(user.Name, feed.Name)
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

	if len(feeds) == 0 {
		fmt.Println("No feed found.")
		return nil
	}

	fmt.Printf("Found %d feeds.\n", len(feeds))
	for i, feed := range feeds {
		fmt.Printf("- - - Feed %d - - -\n", i)
		user, err := st.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get user by id: %w", err)
		}
		printFeed(feed, user)
		fmt.Printf("- - - - - - - - - -\n")
	}

	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("ID:   		%v\n", feed.ID)
	fmt.Printf("Name: 		%s\n", feed.Name)
	fmt.Printf("URL:   		%s\n", feed.Url)
	fmt.Printf("User: 		%s\n", user.Name)
}

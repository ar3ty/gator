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

func handlerFollow(st *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.name)
	}

	feedToFollow, err := st.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("feed doesn't exist: %w", err)
		}
		return fmt.Errorf("cannot find feed by name: %w", err)
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feedToFollow.ID,
	}

	_, err = st.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("cannot create feed follow: %w", err)
	}

	fmt.Println("New feed follow is created:")
	printFeedFollow(st.cfg.CurrentUserName, feedToFollow.Name)
	return nil
}

func handlerFollowing(st *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: %s", cmd.name)
	}

	feedsFollowed, err := st.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("follows don't exist: %w", err)
		}
		return fmt.Errorf("cannot find follows: %w", err)
	}

	if len(feedsFollowed) == 0 {
		fmt.Println("No followed feeds found for current user.")
		return nil
	}

	fmt.Printf("Feeds %s is following:\n", user.Name)
	for i, feedf := range feedsFollowed {
		fmt.Printf(" %d - %s\n", i, feedf.FeedName)
	}

	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("User: %s", username)
	fmt.Printf("Feed: %s\n", feedname)
}

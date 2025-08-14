package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ar3ty/gator/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	feed := RSSFeed{}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create request:\n %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot get response:\n %w", err)
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response:\n %w", err)
	}

	err = xml.Unmarshal(responseBody, &feed)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshaling:\n %w", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		feed.Channel.Item[i] = item
	}

	return &feed, nil
}

func scrapeFeeds(st *state) error {
	nextFeed, err := st.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch next feed in query:\n %w", err)
	}

	feed, err := st.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return fmt.Errorf("failed to mark feed fetched:\n %w", err)
	}

	newFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	log.Printf("Feed %s is fetched:\n", newFeed.Channel.Title)

	for _, item := range newFeed.Channel.Item {
		timeParsed, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return fmt.Errorf("error timeparsing:\n %w", err)
		}
		publ_at := sql.NullTime{
			Time:  timeParsed,
			Valid: true,
		}
		descr := sql.NullString{
			String: item.Description,
			Valid:  true,
		}
		params := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: descr,
			PublishedAt: publ_at,
			FeedID:      feed.ID,
		}
		_, err = st.db.CreatePost(context.Background(), params)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("error creating post: %v\n", err)
		}
	}

	log.Printf("Total posts: %d\n", len(newFeed.Channel.Item))
	return nil
}

func handlerAggregate(st *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <time_between_requests (1m|1s|1h etc)>", cmd.name)
	}

	time_between_reqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("invalid duration:\n %w", err)
	}
	fmt.Printf("Collecting feeds every %v...\n", time_between_reqs)

	ticker := time.NewTicker(time_between_reqs)

	for ; ; <-ticker.C {
		err = scrapeFeeds(st)
		if err != nil {
			fmt.Printf("error during scraping feeds: %v\n", err)
		}
	}
}

func handlerBrowse(st *state, cmd command, user database.User) error {
	if len(cmd.args) > 1 {
		return fmt.Errorf("usage: %s <limit(optional)>", cmd.name)
	}

	var limit int
	var err error
	if len(cmd.args) == 1 {
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("limit is provided, not recognised:\n %w", err)
		}
	} else {
		limit = 2
	}

	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}

	posts, err := st.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("cannot get posts for current user:\n %w", err)
	}

	if len(posts) == 0 {
		fmt.Println("No posts found for current user.")
		return nil
	}

	fmt.Printf("Posts found for current user - %d:", len(posts))
	for i, post := range posts {
		fmt.Printf("Post %d:\n", i)
		printPost(post)
	}

	return nil
}

func printPost(post database.Post) {
	fmt.Printf("\tPublished:	%v\n", post.PublishedAt.Time.Format("Mon Jan 2"))
	fmt.Printf("\tURL:		%v\n", post.Url)
	fmt.Printf("\tTitle: 		%v\n", post.Title)
	fmt.Printf("\tDescription:	%v\n", post.Description.String)
}

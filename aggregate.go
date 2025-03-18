package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/GLobyNew/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("command 'agg' expects one argument: <time between requests> (e.g '1m')")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests.String())
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

}

func scrapeFeeds(s *state) error {
	feedToFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	err = s.db.MarkFeedFetched(context.Background(), feedToFetch.ID)
	if err != nil {
		return err
	}

	fetchedFeed, err := fetchFeed(context.Background(), feedToFetch.Url)
	if err != nil {
		return err
	}

	for _, item := range fetchedFeed.Channel.Item {
		pubTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return err
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			PublishedAt: pubTime,
			Title: item.Title,
			Url: item.Link,
			Description: item.Description,
			FeedID: feedToFetch.ID,
		})
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" { // 23505 is the PostgreSQL error code for unique violations
				continue
			} else {
				return err
			} 
		}
	}

	return nil
}

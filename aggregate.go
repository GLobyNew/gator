package main

import (
	"context"
	"errors"
	"fmt"
	"time"
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
		fmt.Printf("* %s\n", item.Title)
	}

	return nil
}

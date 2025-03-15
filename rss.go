package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
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
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}

	req.Header.Set("User-Agent", "gator")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}

	rBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	var rFeed RSSFeed
	err = xml.Unmarshal(rBody, &rFeed)
	if err != nil {
		return &RSSFeed{}, err
	}
	unescapeFeed(&rFeed)

	return &rFeed, nil
}

func unescapeFeed(feed *RSSFeed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for _, entry := range feed.Channel.Item {
		entry.Title = html.UnescapeString(entry.Title)
		entry.Description = html.UnescapeString(entry.Description)
	}
}

func printFeed(feed *RSSFeed) {
	fmt.Printf("Channel:\n")
	fmt.Printf("Title: %v\n", feed.Channel.Title)
	fmt.Printf("Link: %v\n", feed.Channel.Link)
	fmt.Printf("Description: %v\n", feed.Channel.Description)
	fmt.Printf("Items:\n")
	for _, entry := range feed.Channel.Item {
		fmt.Printf("- Title: %v\n", entry.Title)
		fmt.Printf("  Link: %v\n", entry.Link)
		fmt.Printf("  Description: %v\n", entry.Description)
		fmt.Printf("  Publication Date: %v\n", entry.PubDate)
	}
}

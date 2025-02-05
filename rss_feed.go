package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
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
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("wrong status code: %v", res.Status)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, err
	}

	feed.Channel.Title = extractAndCombineHTML(feed.Channel.Title)
	feed.Channel.Description = extractAndCombineHTML(feed.Channel.Description)
	for i, item := range feed.Channel.Item {
		item.Title = extractAndCombineHTML(item.Title)
		item.Description = extractAndCombineHTML(item.Description)
		feed.Channel.Item[i] = item
	}

	return &feed, nil
}

// Extracts and cleans content while removing only the first <div> block after each <!-- BREAK X -->
func extractAndCombineHTML(input string) string {
	// Regular expression to find break markers like <!-- BREAK 1 -->, <!-- BREAK 2 -->, etc.
	breakRegex := regexp.MustCompile(`<!-- BREAK \d+ -->`)

	// Split the input string into sections based on the break markers
	splits := breakRegex.Split(input, -1)

	// Slice to collect cleaned sections
	var sections []string

	// Process each section after the markers
	for i, section := range splits {
		if i == 0 {
			// First section (before any marker), keep as is
			sections = append(sections, cleanHTML(section, false))
			continue
		}

		// Remove only the first <div> block after each marker
		cleanedSection := cleanHTML(section, true)
		sections = append(sections, cleanedSection)
	}

	// Join all sections with a newline
	return strings.Join(sections, "\n")
}

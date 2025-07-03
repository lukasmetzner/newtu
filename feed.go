package main

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

// FetchArticles fetches articles from all configured RSS feeds.
// Individual feed errors are collected but do not stop processing.
// Returns all successfully parsed articles and any errors encountered.
func FetchArticles(feeds []RSSFeed) ([]Article, error) {
	parser := gofeed.NewParser()
	var articles []Article
	var errs []error

	for _, feed := range feeds {
		parsed, err := parser.ParseURL(feed.URL)
		if err != nil {
			errs = append(errs, fmt.Errorf("fetching %s (%s): %w", feed.Source, feed.URL, err))
			continue
		}

		for _, item := range parsed.Items {
			if item.Published == "" {
				continue
			}

			dt, err := parseTime(item.Published)
			if err != nil {
				errs = append(errs, fmt.Errorf("parsing date for %q from %s: %w", item.Title, feed.Source, err))
				continue
			}

			articles = append(articles, Article{
				Source:   feed.Source,
				Title:    item.Title,
				Datetime: dt,
				Link:     item.Link,
			})
		}
	}

	if len(errs) > 0 && len(articles) == 0 {
		return nil, fmt.Errorf("all feeds failed: %v", errs)
	}

	return articles, nil
}

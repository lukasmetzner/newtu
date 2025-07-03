package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/dustin/go-humanize"
)

type Article struct {
	Source   string
	Title    string
	Datetime time.Time
	Link     string
}

var supportedTimeFormats = []string{
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
}

func parseTime(raw string) (time.Time, error) {
	for _, layout := range supportedTimeFormats {
		if t, err := time.Parse(layout, raw); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported date format: %q", raw)
}

func (a Article) Row(position int) table.Row {
	return table.Row{
		strconv.Itoa(position),
		a.Source,
		humanize.Time(a.Datetime),
		a.Title,
	}
}

// ArticlesToRows converts a slice of articles to table rows with 1-based positions.
func ArticlesToRows(articles []Article) []table.Row {
	rows := make([]table.Row, len(articles))
	for i, a := range articles {
		rows[i] = a.Row(i + 1)
	}
	return rows
}

// SortByDateDesc sorts articles in-place, newest first.
func SortByDateDesc(articles []Article) {
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Datetime.After(articles[j].Datetime)
	})
}

// MergeArticles deduplicates two slices of articles by link.
// When both slices contain the same link, the article from a takes precedence.
func MergeArticles(a, b []Article) []Article {
	seen := make(map[string]struct{}, len(a)+len(b))
	merged := make([]Article, 0, len(a)+len(b))

	for _, article := range a {
		if _, ok := seen[article.Link]; !ok {
			seen[article.Link] = struct{}{}
			merged = append(merged, article)
		}
	}
	for _, article := range b {
		if _, ok := seen[article.Link]; !ok {
			seen[article.Link] = struct{}{}
			merged = append(merged, article)
		}
	}

	return merged
}

// FilterByTitle returns articles whose title contains the query (case-insensitive).
func FilterByTitle(articles []Article, query string) []Article {
	if query == "" {
		return articles
	}

	q := strings.ToLower(query)
	filtered := make([]Article, 0, len(articles))
	for _, a := range articles {
		if strings.Contains(strings.ToLower(a.Title), q) {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

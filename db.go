package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

func OpenDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if err := createSchema(db); err != nil {
		closeErr := db.Close()
		return nil, errors.Join(fmt.Errorf("creating schema: %w", err), closeErr)
	}

	return db, nil
}

func createSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS articles (
			id       INTEGER NOT NULL PRIMARY KEY,
			source   TEXT,
			title    TEXT,
			datetime INTEGER,
			link     TEXT UNIQUE
		);
	`)
	return err
}

func SaveArticles(db *sql.DB, articles []Article) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // rollback after commit is a no-op

	stmt, err := tx.Prepare(`
		INSERT OR IGNORE INTO articles (source, title, datetime, link)
		VALUES (?, ?, ?, ?);
	`)
	if err != nil {
		return fmt.Errorf("preparing statement: %w", err)
	}
	defer stmt.Close() //nolint:errcheck

	for _, a := range articles {
		if _, err := stmt.Exec(a.Source, a.Title, a.Datetime.UnixMilli(), a.Link); err != nil {
			return fmt.Errorf("inserting article %q: %w", a.Title, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func LoadArticles(db *sql.DB) ([]Article, error) {
	rows, err := db.Query(`
		SELECT source, title, datetime, link
		FROM articles
		ORDER BY datetime DESC;
	`)
	if err != nil {
		return nil, fmt.Errorf("querying articles: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	var articles []Article
	for rows.Next() {
		var (
			source, title, link string
			millis              int64
		)
		if err := rows.Scan(&source, &title, &millis, &link); err != nil {
			return nil, fmt.Errorf("scanning article row: %w", err)
		}

		articles = append(articles, Article{
			Source:   source,
			Title:    title,
			Datetime: time.UnixMilli(millis),
			Link:     link,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating article rows: %w", err)
	}

	return articles, nil
}

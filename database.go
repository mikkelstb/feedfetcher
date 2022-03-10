package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mikkelstb/feed_fetcher/feed"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) String() string {
	return "Hello!"
}

func (r *SQLiteRepository) InsertNewsItem(a feed.NewsItem) error {
	_, err := r.db.Exec("INSERT INTO newsitem (docdate, id, source, headline, story, url) values(?,?,?,?,?,?)", a.Docdate.UTC().Format("200601021504"), a.Id(), a.FeedId, a.Headline, a.Story, a.Url)
	if err != nil {
		return err
	}
	return nil
}

func (r *SQLiteRepository) All() ([][]string, error) {
	rows, err := r.db.Query("SELECT * FROM source")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all_lines [][]string
	for rows.Next() {
		var value []string = make([]string, 2)
		if err := rows.Scan(&value[0], &value[1]); err != nil {
			return nil, err
		}
		all_lines = append(all_lines, value)
	}
	return all_lines, nil
}

func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}

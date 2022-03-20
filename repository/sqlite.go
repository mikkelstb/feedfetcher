package repository

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mikkelstb/feedfetcher/feed"
)

type SQLite struct {
	db *sql.DB
}

func NewSQLite(filename string) (*SQLite, error) {
	sq := new(SQLite)

	db_config, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	sq.db = db_config
	return sq, nil
}

func (r *SQLite) String() string {
	return "sqlite db"
}

func (r *SQLite) WriteSingle(a feed.NewsItem) (string, error) {
	_, err := r.db.Exec("INSERT INTO newsitem (docdate, id, source, headline, story, url) values(?,?,?,?,?,?)", a.Docdate, a.Id(), a.FeedId, a.Headline, a.Story, a.Url)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (r *SQLite) Close() error {
	return r.db.Close()
}

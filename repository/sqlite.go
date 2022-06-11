package repository

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mikkelstb/feedfetcher/feed"
)

type SQLite struct {
	db       *sql.DB
	max_days int
}

func NewSQLite(filename string, max_days int) (*SQLite, error) {
	sq := new(SQLite)

	sq.max_days = max_days

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
	_, err := r.db.Exec("INSERT INTO newsitem (docdate, id, source, headline, story, url) values(?,?,?,?,?,?)", a.Docdate, a.GetId(), a.FeedId, a.Headline, a.Story, a.Url)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (r *SQLite) Close() error {
	return r.db.Close()
}

func (r *SQLite) EraseOldArticles() (int, error) {
	erasedate := time.Now().AddDate(0, 0, r.max_days*-1)
	res, err := r.db.Exec("DELETE FROM newsitem where docdate<?", erasedate.Format("2006-01-02"))
	if err != nil {
		return 0, err
	}
	rows, _ := res.RowsAffected()
	return int(rows), nil
}

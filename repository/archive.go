package repository

import "github.com/mikkelstb/feedfetcher/feed"

type Archive interface {
	WriteSingle(feed.NewsItem) (string, error)
	String() string
	Close() error
}

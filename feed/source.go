package feed

import (
	"fmt"

	"time"

	"github.com/mikkelstb/feedfetcher/config"
)

type Source struct {
	config config.SourceConfig
	feed   Feed
}

func (s Source) Name() string {
	return s.config.Name
}

func NewSource(cfg config.SourceConfig) *Source {
	s := Source{config: cfg}
	s.feed = &RSSFeed{}
	s.feed.Init(s.config.Feed)
	return &s
}

func (s *Source) Process() error {

	err := s.feed.Read()
	if err != nil {
		return err
	}

	err = s.feed.Parse()
	if err != nil {
		return err
	}

	return nil
}

func (f *Source) GetNewsitems() []NewsItem {

	var articles []NewsItem
	for f.feed.HasNext() {

		article, err := f.feed.GetNext()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		article.FeedId = f.config.Id
		article.Source = f.config.Screen_name
		article.Feed = f.config.Name
		article.Mediatype = f.config.Mediatype
		article.Country = f.config.Country
		article.Language = f.config.Language
		article.Docdate, err = time.Parse(f.config.Feed["docdate_layout"], article.Docdatestring)

		if err != nil {
			fmt.Println(err)
			continue
		}

		articles = append(articles, *article)
	}
	return articles
}

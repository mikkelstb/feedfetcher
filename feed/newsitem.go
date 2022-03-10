package feed

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

type NewsItem struct {
	Feed          string
	Source        string
	FeedId        int
	Mediatype     string
	Headline      string
	Story         string
	Local_id      string
	Url           string
	Docdatestring string
	Docdate       time.Time
	FetchTime     time.Time
	Language      string
	Country       string
}

/*
	ID() returns a "unique" four letter id based on headline and story
*/
func (ni NewsItem) Id() string {
	id := md5.New()
	io.WriteString(id, strings.Join([]string{ni.Headline, ni.Story}, ""))
	return fmt.Sprintf(hex.EncodeToString(id.Sum(nil))[0:4])
}

func (ni NewsItem) ToJson() ([]byte, error) {

	type jsn struct {
		Feed         string `json:"feed"`
		Source       string `json:"source"`
		FeedId       int    `json:"feedId"`
		Mediatype    string `json:"mediatype"`
		Local_id     string `json:"localId"`
		Language     string `json:"language"`
		Country      string `json:"country"`
		Id           string `json:"id"`
		Headline     string `json:"headline"`
		Story        string `json:"story"`
		Url          string `json:"url"`
		Docdate      string `json:"docdate"`
		LocalDocdate string `json:"localDocdate"`
		FetchTime    string `json:"fetchTime"`
	}

	jsn_struct := jsn{

		Feed:         ni.Feed,
		Source:       ni.Source,
		FeedId:       ni.FeedId,
		Mediatype:    ni.Mediatype,
		Language:     ni.Language,
		Country:      ni.Country,
		Headline:     ni.Headline,
		Story:        ni.Story,
		Local_id:     ni.Local_id,
		Docdate:      ni.Docdate.UTC().Format("2006-01-02T15:04:05-07:00"),
		LocalDocdate: ni.Docdate.Format("2006-01-02T15:04:05-07:00"),
		FetchTime:    ni.FetchTime.UTC().Format("2006-01-02T15:04:05-07:00"),
		Id:           ni.Id(),
		Url:          ni.Url,
	}

	buffer := bytes.Buffer{}
	e := json.NewEncoder(&buffer)
	e.SetEscapeHTML(false)
	e.SetIndent("", " ")
	err := e.Encode(jsn_struct)

	return buffer.Bytes(), err
	//return json.MarshalIndent(jsn_struct, "", " ")
}

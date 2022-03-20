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
	Feed      string `json:"feed"`
	Source    string `json:"source"`
	FeedId    int    `json:"feedId"`
	Mediatype string `json:"mediatype"`
	Headline  string `json:"headline"`
	Story     string `json:"story"`
	Local_id  string `json:"localId"`
	Url       string `json:"url"`
	Language  string `json:"language"`
	Country   string `json:"country"`
	Docdate   string `json:"docdate"`
	FetchTime string `json:"fetchTime"`
}

/*
	ID() returns a "unique" four letter id based on headline and story
*/
func (ni NewsItem) Id() string {
	id := md5.New()
	io.WriteString(id, strings.Join([]string{ni.Headline, ni.Story}, ""))
	return fmt.Sprintf(hex.EncodeToString(id.Sum(nil))[0:4])
}

func (ni NewsItem) GetDocdate() time.Time {
	dd, _ := time.Parse(time.RFC3339, ni.Docdate)
	return dd
}

func (ni NewsItem) ToJson() ([]byte, error) {

	buffer := bytes.Buffer{}
	e := json.NewEncoder(&buffer)
	e.SetEscapeHTML(false)
	e.SetIndent("", " ")
	err := e.Encode(ni)

	return buffer.Bytes(), err
}

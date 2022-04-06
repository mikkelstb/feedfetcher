package feed

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type NewsItem struct {
	Feed      string `json:"feed"`
	Source    string `json:"source"`
	FeedId    int    `json:"feedId"`
	Mediatype string `json:"mediatype"`
	Headline  string `json:"headline"`
	Story     string `json:"story"`
	Url       string `json:"url"`
	Language  string `json:"language"`
	Country   string `json:"country"`
	Docdate   string `json:"docdate"`
	FetchTime string `json:"fetchTime"`
	Id        string `json:"id"`
}

/*
	ID() returns a "unique" four letter id based on headline and story
*/

func (ni NewsItem) GetId() string {
	id := md5.New()
	io.WriteString(id, ni.Headline)
	return fmt.Sprintf("%02d%v%v", ni.FeedId, ni.GetDocdate().Format("0601021504"), hex.EncodeToString(id.Sum(nil))[0:4])
}

func (ni NewsItem) GetDocdate() time.Time {
	dd, _ := time.Parse(time.RFC3339, ni.Docdate)
	return dd
}

func (ni NewsItem) ToJson() ([]byte, error) {

	buffer := bytes.Buffer{}
	jsn_encoder := json.NewEncoder(&buffer)

	jsn_encoder.SetEscapeHTML(false)
	jsn_encoder.SetIndent("", " ")

	ni.Id = ni.GetId()

	err := jsn_encoder.Encode(ni)

	return buffer.Bytes(), err
}

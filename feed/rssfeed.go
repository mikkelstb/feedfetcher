package feed

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

type RSSFeed struct {
	url            string
	docdate_layout string
	fetch_time     time.Time
	rss_data       []byte

	XMLName     xml.Name `xml:"rss"`
	Title       string   `xml:"channel>title"`
	Description string   `xml:"channel>description"`

	Items []struct {
		Headline string `xml:"title"`
		Story    string `xml:"description"`
		LocalId  string `xml:"guid"`
		Url      string `xml:"link"`
		Docdate  string `xml:"pubDate"`
	} `xml:"channel>item"`

	html_cleaner   bluemonday.Policy
	newlinepattern regexp.Regexp
	spaces         regexp.Regexp
	startbracket   regexp.Regexp
	endword        regexp.Regexp
}

/*
	feed.Init sets up regexpes and configuration
*/

func (feed *RSSFeed) Init(config map[string]string) {

	feed.html_cleaner = *bluemonday.StrictPolicy()
	feed.newlinepattern = *regexp.MustCompile(`\n+`)
	feed.spaces = *regexp.MustCompile(`\s{2,}`)
	feed.startbracket = *regexp.MustCompile(`\[.+?\]`)
	feed.endword = *regexp.MustCompile(`\w+…$`)

	feed.url = config["url"]
	feed.docdate_layout = config["docdate_layout"]
	feed.fetch_time = time.Now()
}

/*
	Connects to http server defined by url
	Function returns an error if either connection or reading of data fails
	On success the response will be stored into feed.rss_data
*/

func (feed *RSSFeed) Read() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, feed.url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("error response from server: %s", response.Status)
	}

	feed.rss_data, err = io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return nil
}

/*
	Parse reads the rss_data, and uses xml.Unmarshal to set itself up
*/

func (feed *RSSFeed) Parse() error {
	err := xml.Unmarshal(feed.rss_data, feed)
	return err
}

func (feed *RSSFeed) HasNext() bool {
	return len(feed.Items) != 0
}

/*
	GetNext returns the first newsitem in the Rssfeed
	It runs sanitize on headline and story, inserts fetchtime from feed
	If storytext is smaller than 16 bytes returns nil with errormessage
	Rerurns nil, error if feed.Items is empty
*/

func (feed *RSSFeed) GetNext() (*NewsItem, error) {

	if len(feed.Items) == 0 {
		panic(fmt.Errorf("warning: no more items in feed"))
	}

	nextitem := feed.Items[len(feed.Items)-1]

	// Slice off current item
	feed.Items = feed.Items[0 : len(feed.Items)-1]

	n := new(NewsItem)

	n.Headline = feed.sanitize(nextitem.Headline)
	n.Story = feed.sanitize(nextitem.Story)
	n.Url = nextitem.Url
	n.localId = nextitem.LocalId

	dd, err := time.Parse(feed.docdate_layout, nextitem.Docdate)
	if err != nil {
		return nil, err
	}
	n.Docdate = dd.UTC().Format(time.RFC3339)
	n.FetchTime = feed.fetch_time.UTC().Format(time.RFC3339)

	if len(n.Story) < 16 {
		return nil, fmt.Errorf("story text too small")
	}

	return n, nil
}

/*
	private function to clean up most unnessesary symbols and html tags
*/

func (feed *RSSFeed) sanitize(field string) string {
	clean_field := feed.newlinepattern.ReplaceAllString(field, " ")
	clean_field = feed.html_cleaner.Sanitize(clean_field)
	clean_field = html.UnescapeString(clean_field)
	clean_field = feed.startbracket.ReplaceAllString(clean_field, "")
	clean_field = feed.endword.ReplaceAllString(clean_field, "")
	clean_field = feed.spaces.ReplaceAllString(clean_field, " ")

	clean_field = strings.Trim(clean_field, " :,-.")

	return clean_field
}

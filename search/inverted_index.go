package search

import (
	"bufio"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/mikkelstb/feedfetcher/feed"
)

/*
	InvertedIndex is a collection of unique terms, and their assosiated document-ids
	The Index can be fed continuously with the function AddDocument(NewsItem)
*/

type InvertedIndex struct {
	terms      map[string][]string
	textreader strings.Reader
	docscanner bufio.Scanner

	puctuations regexp.Regexp
}

func NewInvertedIndex() *InvertedIndex {
	ii := InvertedIndex{}
	ii.terms = make(map[string][]string)
	ii.puctuations = *regexp.MustCompile(`[,.:;]`)
	return &ii
}

func (ii *InvertedIndex) AddDocument(doc *feed.NewsItem) {

	text := doc.Headline + " " + doc.Story
	text = strings.ToLower(text)
	text = ii.puctuations.ReplaceAllString(text, " ")

	ii.textreader = *strings.NewReader(text)
	ii.docscanner = *bufio.NewScanner(&ii.textreader)

	ii.docscanner.Split(bufio.ScanWords)

	var nextword string
	var document_terms map[string]int = map[string]int{}

	for ii.docscanner.Scan() {
		nextword = ii.docscanner.Text()
		term := strings.Trim(nextword, "‘’”“+!?() \",.;:'-_")
		term = strings.TrimSuffix(term, "'s")
		term = strings.TrimSuffix(term, "’s")

		if utf8.RuneCountInString(term) < 4 {
			continue
		} else {
			document_terms[term]++
		}
	}

	for term := range document_terms {
		ii.terms[term] = append(ii.terms[term], doc.Id)
	}
}

func (ii InvertedIndex) PrintAllTerms() {

	type pair struct {
		term string
		docs int
	}

	terms := make([]pair, len(ii.terms))
	i := 0
	for k, v := range ii.terms {
		terms[i] = pair{k, len(v)}
		i++
	}

	sort.Slice(terms, func(i, j int) bool { return terms[i].docs < terms[j].docs })

	fmt.Println(terms)

	// for term, docs := range ii.terms {
	// 	fmt.Printf("%s: %v\n", term, docs)
	// }
	// fmt.Printf("Terms: %d\n", len(ii.terms))
}

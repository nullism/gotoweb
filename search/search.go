package search

import (
	"encoding/json"
	"regexp"
	"strings"
)

type Document struct {
	Href  string `json:"href"`
	Title string `json:"title"`
}

type Search struct {
	CurrentId  int              `json:"i"`
	KeywordMap map[string][]int `json:"idx"`
	TagMap     map[string][]int `json:"tags"`
	DocMap     map[int]Document `json:"docs"`
}

var htmlTagRe = regexp.MustCompile(`(?i)<[^>]*>|&[a-z0-9]+;`)

func New() *Search {
	return &Search{
		CurrentId:  0,
		DocMap:     make(map[int]Document),
		TagMap:     make(map[string][]int),
		KeywordMap: make(map[string][]int),
	}
}

func (s *Search) Index(href, title, body string, tags []string) error {
	s.DocMap[s.CurrentId] = Document{
		Href:  href,
		Title: title,
	}

	for _, tag := range tags {
		s.TagMap[tag] = append(s.TagMap[tag], s.CurrentId)
	}

	body = htmlTagRe.ReplaceAllString(body, " ")
	words := regexp.MustCompile(`\w+`).FindAllString(body+" "+title, -1)
	for _, w := range words {
		if len(w) < 3 {
			continue
		}
		s.KeywordMap[strings.ToLower(w)] = append(s.KeywordMap[w], s.CurrentId)
	}

	s.CurrentId++
	return nil
}

func (s *Search) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

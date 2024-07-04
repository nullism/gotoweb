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
	CurrentId  int                    `json:"i"`
	KeywordMap map[string][]int       `json:"idx"`
	KwMap      map[string]map[int]int `json:"kw"`
	TagMap     map[string][]int       `json:"tags"`
	DocMap     map[int]Document       `json:"docs"`
}

var htmlTagRe = regexp.MustCompile(`(?i)<[^>]*>|&[a-z0-9]+;`)
var wordRe = regexp.MustCompile(`\w+`)

func New() *Search {
	return &Search{
		CurrentId:  1,
		DocMap:     make(map[int]Document),
		TagMap:     make(map[string][]int),
		KeywordMap: make(map[string][]int),
		KwMap:      make(map[string]map[int]int),
	}
}

func (s *Search) getIdByHref(href string) int {
	for id, doc := range s.DocMap {
		if doc.Href == href {
			return id
		}
	}
	return -1
}

func (s *Search) Index(href, title, body string, tags []string) error {
	if id := s.getIdByHref(href); id != -1 {
		println("SKIPPING ", href)
		return nil // already indexed
	}

	id := s.CurrentId

	s.DocMap[id] = Document{
		Href:  href,
		Title: title,
	}

	for _, tag := range tags {
		s.TagMap[tag] = append(s.TagMap[tag], id)
	}

	body = htmlTagRe.ReplaceAllString(body, " ")
	titleWords := wordRe.FindAllString(title, -1)
	for _, w := range titleWords {
		if len(w) < 3 {
			continue
		}
		if _, ok := s.KwMap[strings.ToLower(w)]; !ok {
			s.KwMap[strings.ToLower(w)] = make(map[int]int)
		}
		s.KwMap[strings.ToLower(w)][id] += 2

	}

	words := wordRe.FindAllString(body, -1)
	for _, w := range words {
		if len(w) < 3 {
			continue
		}
		if _, ok := s.KwMap[strings.ToLower(w)]; !ok {
			s.KwMap[strings.ToLower(w)] = make(map[int]int)
		}
		s.KwMap[strings.ToLower(w)][id] += 1

	}

	s.CurrentId += 1
	return nil
}

func (s *Search) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

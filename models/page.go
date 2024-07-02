package models

type Page struct {
	Posts  []*Post
	Number int
	Total  int
}

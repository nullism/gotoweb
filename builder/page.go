package builder

type Page struct {
	Posts        []*Post
	Number       int
	Total        int
	NextHref     string
	PreviousHref string
}

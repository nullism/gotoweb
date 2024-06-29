package models

type Post struct {
	Title      string
	Body       string
	SourcePath string
	DestPath   string
}

func PostFromSource(soursePath string) (*Post, error) {
	return &Post{Title: "unimplemented", Body: "<b>STUFF</b>"}, nil
}

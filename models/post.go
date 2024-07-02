package models

import (
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Post struct {
	Title      string
	Body       string
	SourcePath string
	DestPath   string
}

func PostFromSource(sourcePath string) (*Post, error) {
	bs, err := os.ReadFile(sourcePath)
	if err != nil {
		return nil, err
	}
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(bs)

	htmlFlags := html.CommonFlags
	renderer := html.NewRenderer(html.RendererOptions{Flags: htmlFlags})
	htmlBytes := markdown.Render(doc, renderer)

	return &Post{Title: "unimplemented", Body: string(htmlBytes)}, nil
}

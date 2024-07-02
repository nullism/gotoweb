package models

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Post struct {
	Title string
	Body  string
	Href  string
}

func PostFromSource(sourcePath string) (*Post, error) {
	bs, err := os.ReadFile(sourcePath)
	if err != nil {
		return nil, err
	}
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.Tables
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(bs)

	htmlFlags := html.CommonFlags
	renderer := html.NewRenderer(html.RendererOptions{Flags: htmlFlags})
	htmlBytes := markdown.Render(doc, renderer)

	title := filepath.Base(strings.TrimSuffix(sourcePath, ".md"))

	return &Post{Title: title, Body: string(htmlBytes)}, nil
}

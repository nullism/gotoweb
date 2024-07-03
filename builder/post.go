package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"gopkg.in/yaml.v3"
)

type Post struct {
	Title    string
	Body     string
	Href     string
	Tags     []string
	markdown string
}

type PostConfig struct {
	Title string
	Tags  []string
}

// postRe is a real mess. It grabs the yaml block, then a random newline (non-greedy), then the separator, then the rest of the file.
var postRe = regexp.MustCompile(`(?m)((.|\r?\n)*?)(^<![-]+\s+[-]+>$)((.|\r?\n)*)`)

func (b *Builder) postFromSource(sourcePath string) (*Post, error) {

	bs, err := os.ReadFile(sourcePath)
	if err != nil {
		return nil, err
	}
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.Tables
	p := parser.NewWithExtensions(extensions)

	pc := PostConfig{}
	pc.Title = filepath.Base(strings.TrimSuffix(sourcePath, ".md"))

	matches := postRe.FindStringSubmatch(string(bs))
	if len(matches) > 2 {
		bs = []byte(matches[4])
		err := yaml.Unmarshal([]byte(matches[1]), &pc)
		if err != nil {
			return nil, fmt.Errorf("could not parse yaml config from %v: %w", sourcePath, err)
		}
	}

	doc := p.Parse(bs)

	htmlFlags := html.CommonFlags
	renderer := html.NewRenderer(html.RendererOptions{Flags: htmlFlags})

	htmlBytes := markdown.Render(doc, renderer)

	post := &Post{
		Title:    pc.Title,
		Tags:     pc.Tags,
		Body:     string(htmlBytes),
		markdown: string(bs),
	}

	href := strings.Replace(strings.TrimPrefix(sourcePath, b.site.SourceDir), ".md", ".html", 1)
	post.Href = href
	return post, err
}

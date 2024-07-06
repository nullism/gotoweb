package builder

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"gopkg.in/yaml.v3"
)

type Post struct {
	Title       string
	Body        string
	Blurb       string
	Href        string
	Tags        []string
	SkipIndex   bool `yaml:"skip_index"`
	SkipPublish bool `yaml:"skip_publish"`
	markdown    string
}

// postRe requires the start (---) to be on the first line.
var postRe = regexp.MustCompile(`(?m)([ -~\n]*?)^---$((.|\r?\n)*?)^---$((.|\r?\n)*)`)

// regex to strip html tags
var tagRe = regexp.MustCompile(`<[^>]*>`)

func parsePostConfig(post *Post, body []byte) (*Post, []byte, error) {
	text := body
	matches := postRe.FindStringSubmatch(string(body))
	if matches != nil {
		if matches[1] == "" {

			if len(matches) > 2 {
				text = []byte(matches[4])
				err := yaml.Unmarshal([]byte(matches[2]), post)
				if err != nil {
					return nil, []byte(""), fmt.Errorf("could not parse yaml config: %w", err)
				}
			}
		}
	}
	return post, text, nil
}

func (b *Builder) postFromBytes(bs []byte, sourcePath string) (*Post, error) {
	var err error
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.Tables
	p := parser.NewWithExtensions(extensions)

	post := &Post{}
	post.Title = b.files.Base(strings.TrimSuffix(sourcePath, ".md"))
	post, bs, err = parsePostConfig(post, bs)
	if err != nil {
		return nil, err
	}

	doc := p.Parse(bs)

	htmlFlags := html.CommonFlags
	renderer := html.NewRenderer(html.RendererOptions{Flags: htmlFlags})

	htmlBytes := markdown.Render(doc, renderer)

	post.Body = string(htmlBytes)
	post.markdown = string(bs)

	if post.Blurb == "" {
		blurbBytes := tagRe.ReplaceAll(htmlBytes, []byte(" "))
		post.Blurb = string(blurbBytes[:min(200, len(blurbBytes))])
	}

	href := b.site.Prefix + strings.Replace(strings.TrimPrefix(sourcePath, b.site.SourceDir), ".md", ".html", 1)
	post.Href = href
	return post, err
}

func (b *Builder) postFromSource(sourcePath string) (*Post, error) {

	bs, err := os.ReadFile(sourcePath)
	if err != nil {
		return nil, err
	}
	p, err := b.postFromBytes(bs, sourcePath)
	if err != nil {
		return nil, fmt.Errorf("error parsing source %v: %v", sourcePath, err)
	}
	return p, nil
}

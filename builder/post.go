package builder

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/nullism/gotoweb/config"
	"gopkg.in/yaml.v3"
)

type Post struct {
	Title       string
	Body        string
	Blurb       string
	Href        string
	Tags        []string
	SkipIndex   bool      `yaml:"skip_index"`
	SkipPublish bool      `yaml:"skip_publish"`
	CreatedAt   time.Time `yaml:"created_at"`
	UpdatedAt   time.Time `yaml:"updated_at"`
	markdown    string
}

// postRe requires the start (---) to be on the first line.
var postRe = regexp.MustCompile(`(?m)([ -~\n]*?)^---$((.|\r?\n)*?)^---$((.|\r?\n)*)`)

// regex to strip html tags
var tagRe = regexp.MustCompile(`<[^>]*>`)

// UpdateFromConfig updates post properties fron a config.PostConfig object.
func (p *Post) UpdateFromConfig(pc *config.PostConfig) {
	if pc.Title != "" {
		p.Title = pc.Title
	}
	if pc.Blurb != "" {
		p.Blurb = pc.Blurb
	}
	if pc.Tags != nil {
		p.Tags = pc.Tags
	}
	if pc.SkipIndex {
		p.SkipIndex = true
	}
	if pc.SkipPublish {
		p.SkipPublish = true
	}
	if !pc.CreatedAt.IsZero() {
		p.CreatedAt = pc.CreatedAt
	}
	if !pc.UpdatedAt.IsZero() {
		p.UpdatedAt = pc.UpdatedAt
	}
}

func postConfigFromBytes(body []byte) (*config.PostConfig, string, error) {
	pconf := &config.PostConfig{}
	matches := postRe.FindStringSubmatch(string(body))
	if matches != nil {
		if matches[1] == "" {

			if len(matches) > 2 {
				err := yaml.Unmarshal([]byte(matches[2]), pconf)
				if err != nil {
					return nil, "", fmt.Errorf("could not parse yaml config: %w", err)
				}
				return pconf, matches[4], nil
			}
		}
	}
	return &config.PostConfig{}, string(body), nil
}

func (b *Builder) postConfigFromPath(path string) (*config.PostConfig, string, error) {
	bs, err := b.files.ReadFile(path)
	if err != nil {
		return nil, "", fmt.Errorf("could not read source %v: %v", path, err)
	}
	return postConfigFromBytes(bs)
}

func (b *Builder) postFromBytes(bs []byte, sourcePath string) (*Post, error) {
	var err error
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.Tables
	p := parser.NewWithExtensions(extensions)

	post := &Post{}
	post.Title = b.files.Base(strings.TrimSuffix(sourcePath, ".md"))

	if fi, err := b.files.Stat(sourcePath); err == nil {
		post.CreatedAt = fi.ModTime()
		post.UpdatedAt = fi.ModTime()
	}

	pconf, body, err := postConfigFromBytes(bs)
	if err != nil {
		return nil, err
	}
	bs = []byte(body)
	post.UpdateFromConfig(pconf)

	href, err := b.prefix(strings.Replace(strings.TrimPrefix(sourcePath, b.site.SourceDir), ".md", ".html", 1))
	if err != nil {
		return nil, err
	}
	post.Href = href

	newCtx := b.context // copy the context so we don't modify the original with something like `plink`
	if b.context != nil {
		newCtx.Post = post
	}
	str, err := b.Render(sourcePath, b.getSourceFuncMap(), bs, newCtx)
	if err != nil {
		return nil, err
	}
	bs = []byte(str)

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

	return post, err
}

func (b *Builder) postFromPath(sourcePath string) (*Post, error) {

	bs, err := b.files.ReadFile(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("could not read source %v: %v", sourcePath, err)
	}

	p, err := b.postFromBytes(bs, sourcePath)
	if err != nil {
		return nil, fmt.Errorf("error parsing source %v: %v", sourcePath, err)
	}
	return p, nil
}

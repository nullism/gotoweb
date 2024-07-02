package builder

import (
	"strings"

	"github.com/nullism/gotoweb/models"
)

func (b *Builder) postFromSource(sourcePath string) (*models.Post, error) {
	bs, err := models.PostFromSource(sourcePath)
	if err != nil {
		return nil, err
	}
	href := strings.Replace(strings.TrimPrefix(sourcePath, b.site.SourceDir), ".md", ".html", 1)
	bs.Href = href
	return bs, err
}

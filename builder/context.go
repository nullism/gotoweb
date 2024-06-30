package builder

import "github.com/nullism/gotoweb/models"

type RenderContext struct {
	Site *models.SiteConfig
	Post *models.Post
}

func (r *RenderContext) Reset() {
	r.Post = nil
}

package builder

import "github.com/nullism/gotoweb/models"

type RenderContext struct {
	Site  *models.SiteConfig
	Post  *models.Post
	Posts []*models.Post
	Page  *models.Page
	Args  map[string]any
}

func (r *RenderContext) Reset() {
	r.Post = nil
	r.Posts = nil
	r.Page = nil
}

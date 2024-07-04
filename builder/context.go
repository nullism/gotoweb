package builder

import "github.com/nullism/gotoweb/config"

type RenderContext struct {
	Site  *config.SiteConfig
	Post  *Post
	Posts []*Post
	Page  *Page
	Args  map[string]any
}

func (r *RenderContext) Reset() {
	r.Post = nil
	r.Posts = nil
	r.Page = nil
}

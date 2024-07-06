package builder

import "github.com/nullism/gotoweb/config"

type RenderContext struct {
	Site   *config.SiteConfig
	Post   *Post
	Posts  []*Post
	Page   *Page
	Args   map[string]any
	TagMap map[string]int
}

// AddTags add tags and their counts to the context.
func (r *RenderContext) AddTags(tags ...string) {
	if r.TagMap == nil {
		r.TagMap = make(map[string]int)
	}
	for _, tag := range tags {
		if _, ok := r.TagMap[tag]; !ok {
			r.TagMap[tag] = 1
		} else {
			r.TagMap[tag]++
		}
	}
}

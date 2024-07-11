package builder

import (
	"testing"

	"github.com/nullism/gotoweb/config"
	"github.com/stretchr/testify/assert"
)

func TestRenderContext_AddTags(t *testing.T) {
	b := &Builder{site: &config.SiteConfig{SourceDir: "/foo"}}
	r := &RenderContext{Site: b.site}
	tags := []string{"a", "b", "c", "a"}
	r.AddTags(tags...)
	assert.Equal(t, 3, len(r.TagMap))
	assert.Equal(t, 2, r.TagMap["a"])
	assert.Equal(t, 1, r.TagMap["b"])
}

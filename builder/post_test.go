package builder

import (
	"reflect"
	"testing"

	"github.com/nullism/gotoweb/config"
	"github.com/nullism/gotoweb/fsys"
	"github.com/stretchr/testify/assert"
)

func Test_postFromBytes(t *testing.T) {
	path := "/foo/test.md"
	bytes := []byte("---\ntitle: test\n---\nHello World")
	b := &Builder{site: &config.SiteConfig{SourceDir: "/foo"}, files: &fsys.OsFileSystem{}}
	post, err := b.postFromBytes(bytes, path)
	assert.NoError(t, err)
	assert.NotEmpty(t, post)
	assert.Equal(t, "test", post.Title)
}

func Test_getPostConfig(t *testing.T) {
	tests := []struct {
		name     string
		body     []byte
		want     *config.PostConfig
		wantText string
		wantErr  bool
	}{
		{
			"simple test with title",
			[]byte("---\ntitle: \"Hello World\"\n---\nASDFASDF"),
			&config.PostConfig{Title: "Hello World"},
			"\nASDFASDF",
			false,
		},
		{
			"test with tags",
			[]byte("---\ntags: [a, b, c]\n---\nASDFASDF"),
			&config.PostConfig{Tags: []string{"a", "b", "c"}},
			"\nASDFASDF",
			false,
		},
		{
			"test skip publish and skip index",
			[]byte("---\nskip_publish: true\nskip_index: true\n---\nASDFASDF"),
			&config.PostConfig{SkipPublish: true, SkipIndex: true},
			"\nASDFASDF",
			false,
		},
		{
			"test with unparsable header",
			[]byte("\n---\ntitle: [123]\n---\nASDFASDF"),
			&config.PostConfig{},
			"\n---\ntitle: [123]\n---\nASDFASDF",
			false,
		},
		{
			"empty test no match",
			[]byte(`ASD FASDF`),
			&config.PostConfig{},
			`ASD FASDF`,
			false,
		},
		{
			"test with invalid title",
			[]byte("---\ntitle: [123]\n---\nASD FASDF"),
			nil,
			"",
			true,
		},
		{
			"test double header",
			[]byte("---\ntitle: foo\n---\nASD FASDF\n\n#foo\n---\ntitle: bar\n---\n"),
			&config.PostConfig{Title: "foo"},
			"\nASD FASDF\n\n#foo\n---\ntitle: bar\n---\n",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, text, err := postConfigFromBytes(tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPostConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPostConfig() = %v, want %v", got, tt.want)
			}
			if text != tt.wantText {
				t.Errorf("getPostConfig() = `%v`, want `%v`", text, tt.wantText)
			}
		})
	}
}

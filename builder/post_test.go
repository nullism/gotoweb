package builder

import (
	"reflect"
	"testing"

	"github.com/nullism/gotoweb/config"
	"github.com/stretchr/testify/assert"
)

func Test_postFromBytes(t *testing.T) {
	path := "/foo/test.md"
	bytes := []byte("---\ntitle: test\n---\nHello World")
	b := &Builder{site: &config.SiteConfig{SourceDir: "/foo"}}
	post, err := b.postFromBytes(bytes, path)
	assert.NoError(t, err)
	assert.NotEmpty(t, post)
	assert.Equal(t, "test", post.Title)
}

func Test_parsePostConfig(t *testing.T) {

	tests := []struct {
		name       string
		post       *Post
		body       string
		want       *Post
		wantString string
		wantErr    bool
	}{
		{
			"simple test with title",
			&Post{Title: "replaceme"},
			"---\ntitle: \"Hello World\"\n---\nASDFASDF",
			&Post{Title: "Hello World"},
			"\nASDFASDF",
			false,
		},
		{
			"test with unparsable header",
			&Post{Title: "asdf"},
			"\n---\ntitle: [123]\n---\nASDFASDF",
			&Post{Title: "asdf"},
			"\n---\ntitle: [123]\n---\nASDFASDF",
			false,
		},
		{
			"empty test no match",
			&Post{Title: "foo"},
			`ASD FASDF`,
			&Post{Title: "foo"},
			`ASD FASDF`,
			false,
		},
		{
			"test with invalid title",
			&Post{Title: "foo"},
			"---\ntitle: [123]\n---\nASD FASDF",
			nil,
			``,
			true,
		},
		{
			"test double header",
			&Post{Title: "foo"},
			"---\ntitle: foo\n---\nASD FASDF\n\n#foo\n---\ntitle: bar\n---\n",
			&Post{Title: "foo"},
			"\nASD FASDF\n\n#foo\n---\ntitle: bar\n---\n",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parsePostConfig(tt.post, []byte(tt.body))
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePostConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePostConfig() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(string(got1), string(tt.wantString)) {
				t.Errorf("parsePostConfig() got1 = %v, want %v", string(got1), tt.wantString)
			}
		})
	}
}

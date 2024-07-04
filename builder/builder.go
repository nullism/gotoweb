package builder

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/nullism/gotoweb/config"
	"github.com/nullism/gotoweb/logging"
	"github.com/nullism/gotoweb/search"
	"github.com/nullism/gotoweb/theme"
	cp "github.com/otiai10/copy"
)

var log = logging.GetLogger()

type Builder struct {
	site    *config.SiteConfig
	context *RenderContext
	search  *search.Search
}

func New(conf *config.SiteConfig) (*Builder, error) {
	if _, err := os.Stat(conf.SourceDir); err != nil {
		return nil, err
	}
	log.Info("Creating builder", "source", conf.SourceDir)
	s := search.New(conf.Search.MinKeywordLength, conf.Search.StopWords)

	return &Builder{
		site: conf,
		context: &RenderContext{
			Site: conf,
		},
		search: s,
	}, nil
}

func (b *Builder) BuildOne(tplPath, outPath string) error {

	if b.context.Post != nil && b.context.Post.SkipPublish {
		log.Warn("post is marked as skip publish", "title", b.context.Post.Title)
		return nil
	}

	log.Debug("building "+filepath.Base(outPath), "from", tplPath, "to", outPath)
	out, err := b.Render(tplPath, b.context)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(outPath), 0755)
	if err != nil {
		return err
	}
	err = os.WriteFile(outPath, []byte(out), 0755)
	if err != nil {
		return err
	}

	return err
}

func (b *Builder) BuildExtraPages() error {
	var err error
	for _, tpl := range theme.ExtraPageNames {
		tplPath := filepath.Join(b.site.ThemeDir(), tpl+config.TemplateExt)
		outPath := filepath.Join(b.site.PublicDir, tpl+".html")
		sourcePath := filepath.Join(b.site.SourceDir, tpl+".md")

		if _, err := os.Stat(sourcePath); err == nil {
			p, err := b.postFromSource(sourcePath)
			if err != nil {
				log.Error("Could not render template", "error", err)
				return err
			}
			b.context.Post = p
		} else {
			b.context.Post = nil
		}

		err = b.BuildOne(tplPath, outPath)

		if err != nil {
			log.Error("Could not render template", "error", err)
			return err
		}
	}
	return nil
}

func (b *Builder) BuildPostPreview(outPath string) error {
	tplPath := filepath.Join(b.site.ThemeDir(), "post-preview"+config.TemplateExt)
	outPath = outPath + ".preview"
	err := b.BuildOne(tplPath, outPath)
	return err

}

func (b *Builder) BuildPosts() error {
	tplPath := filepath.Join(b.site.ThemeDir(), "post.html")
	_, err := os.Stat(b.site.PublicDir)
	if err != nil {
		err2 := os.Mkdir(b.site.PublicDir, 0755)
		if err2 != nil {
			return err2
		}
	}
	err = filepath.WalkDir(b.site.SourceDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// name with subdirectories included, leading slash removed
		subName := strings.TrimPrefix(strings.TrimPrefix(path, b.site.SourceDir), "/")

		if filepath.Ext(d.Name()) != ".md" {
			// TODO: copy files over?
			return nil
		}
		// if theme.IsExtraPage(strings.TrimSuffix(subName, ".md")) {
		// 	log.Warn("skipping extra page", "page", subName)
		// 	// skip built-in pages (they are built with the theme)
		// 	return nil
		// }
		plain := strings.TrimSuffix(subName, filepath.Ext(d.Name()))

		outPath := filepath.Join(b.site.PublicDir, plain+".html")
		post, err := b.postFromSource(path)
		if err != nil {
			return err
		}
		if post.SkipPublish {
			return nil
		}
		if !post.SkipIndex {
			err = b.search.Index(post.Href, post.Title, post.Body, post.Tags)
			if err != nil {
				return err
			}
		}
		b.context.Post = post
		b.context.Posts = append(b.context.Posts, post)
		err = b.BuildOne(tplPath, outPath)
		if err != nil {
			return err
		}
		err = b.BuildPostPreview(outPath)
		return err
	})
	return err
}

func (b *Builder) BuildPostLists() error {

	postCount := len(b.context.Posts)
	postsPerPage := b.site.PostsPerPage

	pageCount := int(math.Ceil(float64(postCount) / float64(postsPerPage)))

	log.Info("Building post pages", "posts per page", postsPerPage, "total pages", pageCount, "total posts", postCount)
	postI := 0
	for pnum := range pageCount {
		b.context.Page = &Page{
			Number: pnum + 1,
			Total:  pageCount,
			Posts:  []*Post{},
		}
		for i := 0; i < postsPerPage; i++ {
			if postI >= postCount {
				break
			}
			b.context.Page.Posts = append(b.context.Page.Posts, b.context.Posts[postI])
			postI++
		}

		tplPath := filepath.Join(b.site.ThemeDir(), "posts.html")
		outPath := filepath.Join(b.site.PublicDir, fmt.Sprintf("posts-%d.html", pnum+1))
		err := b.BuildOne(tplPath, outPath)
		if err != nil {
			return err
		}
	}
	return nil
}

// BuildAll builds all the pages and posts for the site.
func (b *Builder) BuildAll() error {
	log.Info("building project", "site", b.site.Name)

	err := b.RemovePublic()
	if err != nil {
		log.Error("could not clean public directory", "error", err)
		return err
	}

	err = b.CheckDirectories()
	if err != nil {
		log.Error("directory check failed", "error", err)
		return err
	}

	// build posts
	err = b.BuildPosts()
	if err != nil {
		return err
	}

	// build pages
	err = b.BuildExtraPages()
	if err != nil {
		return err
	}

	err = b.BuildPostLists()
	if err != nil {
		return err
	}

	idxPath := filepath.Join(b.site.PublicDir, "index.json")
	log.Info("writing search index", "outfile", idxPath)
	idx, err := b.search.ToJson()
	if err != nil {
		return err
	}

	err = os.WriteFile(idxPath, idx, 0755)
	if err != nil {
		return err
	}
	err = cp.Copy(filepath.Join(b.site.ThemeDir(), "dist"), filepath.Join(b.site.PublicDir, "dist"), cp.Options{AddPermission: 0755})

	// println(string(idx))
	return err
}

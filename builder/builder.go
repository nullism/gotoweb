package builder

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/nullism/gotoweb/config"
	"github.com/nullism/gotoweb/fsys"
	"github.com/nullism/gotoweb/logging"
	"github.com/nullism/gotoweb/search"
	"github.com/nullism/gotoweb/theme"
)

var log = logging.GetLogger()

type Builder struct {
	site    *config.SiteConfig
	context *RenderContext
	search  *search.Search
	files   fsys.FileSystem
}

func New(conf *config.SiteConfig, files fsys.FileSystem) (*Builder, error) {
	if _, err := files.Stat(conf.SourceDir); err != nil {
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
		files:  files,
	}, nil
}

func (b *Builder) BuildOne(tplPath, outPath string) error {

	if b.context.Post != nil && b.context.Post.SkipPublish {
		log.Warn("post is marked as skip publish", "title", b.context.Post.Title)
		return nil
	}

	log.Debug("building "+b.files.Base(outPath), "from", tplPath, "to", outPath)
	out, err := b.Render(tplPath, b.context)
	if err != nil {
		return err
	}
	err = b.files.MkdirAll(b.files.Dir(outPath), 0755)
	if err != nil {
		return err
	}
	err = b.files.WriteFile(outPath, []byte(out), 0755)
	if err != nil {
		return err
	}

	return err
}

func (b *Builder) BuildExtraPages() error {
	var err error
	for _, tpl := range theme.ExtraPageNames {
		tplPath := b.files.Join(b.site.ThemeDir(), tpl+config.TemplateExt)
		outPath := b.files.Join(b.site.PublicDir, tpl+".html")
		sourcePath := b.files.Join(b.site.SourceDir, tpl+".md")

		if _, err := b.files.Stat(sourcePath); err == nil {
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
	tplPath := b.files.Join(b.site.ThemeDir(), "post-preview"+config.TemplateExt)
	outPath = outPath + ".preview"
	err := b.BuildOne(tplPath, outPath)
	return err

}

func (b *Builder) BuildPosts() error {
	tplPath := b.files.Join(b.site.ThemeDir(), "post.html")
	_, err := b.files.Stat(b.site.PublicDir)
	if err != nil {
		err2 := b.files.Mkdir(b.site.PublicDir, 0755)
		if err2 != nil {
			return err2
		}
	}
	err = b.files.WalkDir(b.site.SourceDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// name with subdirectories included, leading slash removed
		subName := strings.TrimPrefix(strings.TrimPrefix(path, b.site.SourceDir), "/")

		if b.files.Ext(d.Name()) != ".md" {
			// TODO: copy files over?
			return nil
		}
		// if theme.IsExtraPage(strings.TrimSuffix(subName, ".md")) {
		// 	log.Warn("skipping extra page", "page", subName)
		// 	// skip built-in pages (they are built with the theme)
		// 	return nil
		// }
		plain := strings.TrimSuffix(subName, b.files.Ext(d.Name()))

		outPath := b.files.Join(b.site.PublicDir, plain+".html")
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

		tplPath := b.files.Join(b.site.ThemeDir(), "posts.html")
		outPath := b.files.Join(b.site.PublicDir, fmt.Sprintf("posts-%d.html", pnum+1))
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

	idxPath := b.files.Join(b.site.PublicDir, "index.json")
	log.Info("writing search index", "outfile", idxPath)
	idx, err := b.search.ToJson()
	if err != nil {
		return err
	}

	err = b.files.WriteFile(idxPath, idx, 0755)
	if err != nil {
		return err
	}
	err = b.files.Copy(b.files.Join(b.site.ThemeDir(), "dist"), b.files.Join(b.site.PublicDir, "dist"), 0755)

	// println(string(idx))
	return err
}

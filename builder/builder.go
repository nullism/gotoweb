package builder

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/nullism/gotoweb/logging"
	"github.com/nullism/gotoweb/models"
	"github.com/nullism/gotoweb/theme"
)

var log = logging.GetLogger()

type Builder struct {
	site    *models.SiteConfig
	context *RenderContext
}

func New(conf *models.SiteConfig) (*Builder, error) {
	if _, err := os.Stat(conf.SourceDir); err != nil {
		return nil, err
	}
	log.Info("Creating builder", "source", conf.SourceDir)
	return &Builder{
		site: conf,
		context: &RenderContext{
			Site: conf,
		},
	}, nil
}

// Render renders a template with the given content.
func (b *Builder) Render(tplPath string, content *RenderContext) (string, error) {
	bs, err := os.ReadFile(tplPath)
	if err != nil {
		return "", err
	}

	tpl, err := template.New(tplPath).Funcs(b.getFuncMap()).Option("missingkey=error").Parse(string(bs))

	if err != nil {
		return "", err
	}

	sb := strings.Builder{}
	err = tpl.Execute(&sb, content)
	return sb.String(), err
}

func (b *Builder) BuildOne(tplPath, outPath string) error {
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
	return err
}

func (b *Builder) BuildExtraPages() error {
	var err error
	for _, tpl := range theme.ExtraPageNames {
		tplPath := filepath.Join(b.site.ThemeDir, tpl+".html")
		outPath := filepath.Join(b.site.PublicDir, tpl+".html")
		sourcePath := filepath.Join(b.site.SourceDir, tpl+".md")

		if _, err := os.Stat(sourcePath); err == nil {
			p, err := models.PostFromSource(sourcePath)
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

func (b *Builder) BuildPosts() error {
	tplPath := filepath.Join(b.site.ThemeDir, "post.html")
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
		if theme.IsExtraPage(strings.TrimSuffix(subName, ".md")) {
			log.Warn("skipping extra page", "page", subName)
			// skip built-in pages (they are built with the theme)
			return nil
		}
		plain := strings.TrimSuffix(subName, filepath.Ext(d.Name()))

		outPath := filepath.Join(b.site.PublicDir, plain+".html")
		post, err := models.PostFromSource(path)
		if err != nil {
			return err
		}
		b.context.Post = post
		b.context.Posts = append(b.context.Posts, post)
		return b.BuildOne(tplPath, outPath)
	})
	return err
}

func (b *Builder) BuildPostLists() error {

	postCount := len(b.context.Posts)
	postsPerPage := 5

	pageCount := int(math.Ceil(float64(postCount) / float64(postsPerPage)))

	log.Info("Building post pages", "total pages", pageCount, "total posts", postCount, "posts per page", postsPerPage)

	for pnum := range pageCount {
		b.context.Page = &models.Page{
			Number: pnum + 1,
			Total:  pageCount,
			Posts:  []*models.Post{},
		}
		for i := pnum; i < (pnum+1)*postsPerPage; i++ {
			if i >= postCount {
				break
			}
			b.context.Page.Posts = append(b.context.Page.Posts, b.context.Posts[i])
		}

		tplPath := filepath.Join(b.site.ThemeDir, "posts.html")
		outPath := filepath.Join(b.site.PublicDir, fmt.Sprintf("posts-%d.html", pnum+1))
		err := b.BuildOne(tplPath, outPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) checkDirectories() error {
	// Source directory
	sd, err := os.Stat(b.site.SourceDir)
	if err != nil || !sd.IsDir() {
		log.Error("source directory inaccessible", "directory", b.site.SourceDir)
		return err
	}

	// Public directory
	pd, err := os.Stat(b.site.PublicDir)
	if err != nil {
		err2 := os.Mkdir(b.site.PublicDir, 0755)
		if err2 != nil {
			log.Error("could not stat or create public directory", "directory", b.site.PublicDir, "mkdir error", err2, "error", err)
			return err2
		}
	} else if !pd.IsDir() {
		log.Error("public exists and is not a directory", "directory", b.site.PublicDir)
		return fmt.Errorf("public exists and is not a directory")
	}
	return nil
}

// BuildAll builds all the pages and posts for the site.
func (b *Builder) BuildAll() error {
	log.Info("Building project", "site", b.site.Name)

	err := b.checkDirectories()
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
	return nil
}

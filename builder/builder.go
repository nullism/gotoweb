package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/nullism/gotoweb/logging"
	"github.com/nullism/gotoweb/models"
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

func (b *Builder) Render(tplPath string, content any) (string, error) {
	bs, err := os.ReadFile(tplPath)
	if err != nil {
		return "", err
	}

	tpl, err := template.New("foo").Funcs(b.getFuncMap()).Parse(string(bs))

	if err != nil {
		return "", err
	}

	sb := strings.Builder{}
	err = tpl.Execute(&sb, content)
	return sb.String(), err
}

func (b *Builder) renderOne(tplPath, mdPath string, content *RenderContext) (string, error) {
	fi, err := os.Stat(mdPath)
	if err == nil && !fi.IsDir() {
		p, err := models.PostFromSource(mdPath)
		if err != nil {
			return "", err
		}
		log.Debug("rendering "+mdPath, "template", tplPath, "post", p)
		content.Post = p
	} else {
		log.Warn("no source file for page", "source", mdPath, "template", tplPath)
	}

	return b.Render(tplPath, content)
}

func (b *Builder) buildOne(tplPath, mdPath, outPath string, content *RenderContext) error {
	out, err := b.renderOne(tplPath, mdPath, content)
	if err != nil {
		return err
	}
	err = os.WriteFile(outPath, []byte(out), 0644)
	if err != nil {
		return err
	}
	if content != nil {
		content.Reset()
	}
	log.Debug("rendered template", "from", tplPath, "to", outPath)
	return nil
}

func (b *Builder) buildPosts() error {
	pubPostsDir := filepath.Join(b.site.PublicDir, models.PostsDir)
	tplPath := filepath.Join(b.site.ThemeDir, "post.html")
	_, err := os.Stat(pubPostsDir)
	if err != nil {
		err2 := os.Mkdir(pubPostsDir, 0755)
		if err2 != nil {
			return err2
		}
	}
	err = filepath.WalkDir(filepath.Join(b.site.SourceDir, models.PostsDir), func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		log.Debug("building post", "file", path, "name", d.Name())
		plain := strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))
		outPath := filepath.Join(pubPostsDir, plain+".html")

		return b.buildOne(tplPath, path, outPath, b.context)
	})
	return err
}

func (b *Builder) checkDirectories() error {
	sd, err := os.Stat(b.site.SourceDir)
	if err != nil || !sd.IsDir() {
		log.Error("source directory inaccessible", "directory", b.site.SourceDir)
		return err
	}
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

func (b *Builder) Build() error {
	log.Info("Building project", "site", b.site.Name)

	err := b.checkDirectories()
	if err != nil {
		log.Error("directory check failed", "error", err)
		return err
	}

	for _, tpl := range []string{"index", "about"} {
		err = b.buildOne(
			filepath.Join(b.site.ThemeDir, tpl+".html"),
			filepath.Join(b.site.SourceDir, tpl+".md"),
			filepath.Join(b.site.PublicDir, tpl+".html"),
			b.context)

		if err != nil {
			log.Error("Could not render template", "error", err)
			return err
		}
	}

	err = b.buildPosts()
	if err != nil {
		return err
	}
	b.context.Post = nil

	// TODO: Send entire context to all templates (all posts, site config, etc.)
	// but also, iterate posts and render each one.
	// Does this make sense? Post pagination?

	return nil
}

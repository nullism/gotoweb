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

	tpl, err := template.New("foo").Funcs(b.getFuncMap()).Option("missingkey=error").Parse(string(bs))

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
	err = os.WriteFile(outPath, []byte(out), 0644)
	return err
}

func (b *Builder) BuildPosts() error {
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
		if filepath.Ext(d.Name()) != ".md" {
			// TODO: copy files over?
			return nil
		}
		log.Debug("building post", "file", path, "name", d.Name())
		plain := strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))
		outPath := filepath.Join(pubPostsDir, plain+".html")
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
	for _, tpl := range []string{"index", "about"} {
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

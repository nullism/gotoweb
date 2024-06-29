package builder

import (
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

func (b *Builder) Render(htmlPath string, content any) (string, error) {
	tpl, err := template.ParseFiles(htmlPath)
	if err != nil {
		return "", err
	}
	tpl = tpl.Funcs(b.getFuncMap())

	sb := strings.Builder{}
	err = tpl.Execute(&sb, content)
	return sb.String(), err
}

func (b *Builder) buildOne(tplFname, mdFname string, content *RenderContext) (string, error) {
	full := filepath.Join(b.site.SourceDir, mdFname)
	_, err := os.Stat(full)
	if err == nil {
		p, err := models.PostFromSource(full)
		if err != nil {
			return "", err
		}
		log.Debug("Rendering "+full, "from", mdFname, "post", p)
		content.Post = p
	}

	return b.Render(*b.site.ThemeDir+"/"+tplFname, content)
}

func (b *Builder) Build() {
	log.Info("Building project", "site", b.site.Name)
	// out, err := b.Render(*b.site.ThemeDir+"/index.html", b.context)
	out, err := b.buildOne("index.html", "index.md", b.context)
	if err != nil {
		log.Error("Could not render template", "error", err)
		return
	}
	log.Debug("Rendered template", "output", out)

	// TODO: Send entire context to all templates (all posts, site config, etc.)
	// but also, iterate posts and render each one.
	// Does this make sense? Post pagination?
}

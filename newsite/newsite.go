package newsite

import (
	"fmt"
	"strings"
	"text/template"

	_ "embed"

	"github.com/nullism/gotoweb/config"
	"github.com/nullism/gotoweb/fsys"
	"github.com/nullism/gotoweb/logging"
)

type NewSite struct {
	Name  string
	conf  config.SiteConfig
	files fsys.FileSystem
}

type configContext struct {
	Title string
}

var log = logging.GetLogger()

//go:embed sample.config.yaml
var sampleYaml []byte

func New(name string, files fsys.FileSystem) (*NewSite, error) {

	path, err := files.Abs(name)
	if err != nil {
		return nil, err
	}

	if files.Exists(path) {
		return nil, fmt.Errorf("directory %s already exists", name)
	}

	log.Info("creating new site", "name", name, "path", path)
	err = files.MkdirAll(files.Join(path, config.SourceDir), 0755)
	if err != nil {
		return nil, err
	}

	tpl, err := template.New("config.yaml").Parse(string(sampleYaml))
	if err != nil {
		return nil, fmt.Errorf("could not parse sample yaml: %w", err)
	}

	sb := strings.Builder{}
	err = tpl.Execute(&sb, configContext{
		Title: name,
	})
	if err != nil {
		return nil, fmt.Errorf("could not execute template: %w", err)
	}

	err = files.WriteFile(files.Join(path, "config.yaml"), []byte(sb.String()), 0755)
	if err != nil {
		return nil, fmt.Errorf("could not create config.yaml: %w", err)
	}

	ns := &NewSite{
		Name:  name,
		files: files,
	}

	return ns, nil
}

func (s *NewSite) Create() error {
	return nil
}

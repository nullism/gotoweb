package newsite

import (
	"github.com/nullism/gotoweb/config"
	"github.com/nullism/gotoweb/fsys"
)

type NewSite struct {
	Name  string
	conf  config.SiteConfig
	files fsys.FileSystem
}

func New(name string, files fsys.FileSystem) (*NewSite, error) {
	return &NewSite{
		Name:  name,
		files: files,
	}, nil
}

func (s *NewSite) Create() error {
	return nil
}

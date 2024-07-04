package newsite

import "github.com/nullism/gotoweb/config"

type NewSite struct {
	Name string
	conf config.SiteConfig
}

func New(name string) (*NewSite, error) {
	return &NewSite{Name: name}, nil
}

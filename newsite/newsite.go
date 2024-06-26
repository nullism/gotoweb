package newsite

import "github.com/nullism/gotoweb/models"

type NewSite struct {
	Name string
	conf models.SiteConfig
}

func New(name string) (*NewSite, error) {
	return &NewSite{Name: name}, nil
}

package builder

import (
	"github.com/nullism/gotoweb/logging"
	"github.com/nullism/gotoweb/models"
)

var log = logging.GetLogger()

type Builder struct {
	conf *models.SiteConfig
}

func New(conf *models.SiteConfig) (*Builder, error) {
	return &Builder{
		conf: conf,
	}, nil
}

func (b *Builder) Build() {
	log.Info("Building project", "site", b.conf.Name)
}

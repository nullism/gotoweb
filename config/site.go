package config

import (
	"fmt"
	"os"

	"github.com/creasty/defaults"
	"github.com/nullism/gotoweb/fsys"
	"github.com/nullism/gotoweb/logging"
	"gopkg.in/yaml.v3"
)

var log = logging.GetLogger()

type SiteConfig struct {
	Name       string
	Theme      ThemeConfig
	Copyright  string
	ConfigPath string // path to config.yaml
	Index      IndexConfig
	PublicDir  string
	RootDir    string
	SourceDir  string
	// ThemeDir   string `yaml:"theme_directory,ignore"`
	Prefix       string `yaml:"uri_prefix"`
	Search       SearchConfig
	Version      string
	PostsPerPage int `default:"2"`
	files        fsys.FileSystem
}

func (s SiteConfig) ThemeDir() string {
	return s.files.Join(s.RootDir, s.Theme.Path)
}

func SiteFromConfig(files fsys.FileSystem) (*SiteConfig, error) {
	confPath, err := files.FindInParent("config.yaml", 3)
	if err != nil {
		return nil, err
	}
	bs, err := os.ReadFile(confPath)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	var s SiteConfig
	err = yaml.Unmarshal(bs, &s)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal config: %w", err)
	}
	s.files = files

	err = defaults.Set(&s)
	if err != nil {
		return nil, fmt.Errorf("could not set defaults: %w", err)
	}

	if s.Version == "" {
		return nil, fmt.Errorf("config version not set")
	}

	if s.Index.MinKeywordLength == 0 {
		s.Index.MinKeywordLength = 3
	}

	s.ConfigPath = confPath
	s.RootDir = s.files.Dir(confPath)
	s.SourceDir = s.files.Join(s.RootDir, SourceDir)

	if s.PublicDir == "" {
		sd := s.files.Join(s.RootDir, PublicDir)
		s.PublicDir = sd
	} else if !s.files.IsAbs(s.PublicDir) {
		// make it absolute and preserve relative paths, e.g /root/../../public
		pd := s.files.Join(s.RootDir, s.PublicDir)
		s.PublicDir = pd
	}

	log.Info("Loaded config", "site name", s.Name, "public dir", s.PublicDir, "source dir", s.SourceDir, "theme dir", s.ThemeDir)
	return &s, nil
}

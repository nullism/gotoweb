package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/creasty/defaults"
	"github.com/nullism/gotoweb/fsys"
	"github.com/nullism/gotoweb/logging"
	"gopkg.in/yaml.v3"
)

var log = logging.GetLogger()

type SiteConfig struct {
	TimeFormat string `default:"2006-01-02 15:04:05"`
	Title      string
	Theme      ThemeConfig
	Copyright  string
	ConfigPath string // path to config.yaml
	Homepage   string
	Index      IndexConfig
	PublicDir  string
	RootDir    string
	SourceDir  string
	// ThemeDir   string `yaml:"theme_directory,ignore"`
	Prefix       string `yaml:"uri_prefix"`
	Search       SearchConfig
	Version      string
	PostsPerPage int        `default:"2"`
	Menu         MenuConfig `yaml:"menu"`
	files        fsys.FileSystem
}

func (s SiteConfig) ThemeDir() string {
	return s.files.Join(s.RootDir, s.Theme.Path)
}

// UpdateMenuPrefixes recursively adds the site prefix to the menu items.
func (s *SiteConfig) UpdateMenuPrefixes(item *MenuItem) {
	if !s.Menu.AutoPrefix {
		return
	}
	if item.Href != "" && !strings.HasPrefix(item.Href, "http://") && !strings.HasPrefix(item.Href, "https://") {
		href, _ := url.JoinPath(s.Prefix, item.Href)
		item.Href = href
	}
	if item.Children != nil {
		for _, child := range item.Children {
			s.UpdateMenuPrefixes(child)
		}
	}
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

	// Add site prefix if necessary
	if s.Menu.AutoPrefix && s.Menu.Items != nil {
		for _, itm := range s.Menu.Items {
			s.UpdateMenuPrefixes(itm)
		}
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

	log.Info("Loaded config", "site name", s.Title, "public dir", s.PublicDir, "source dir", s.SourceDir, "theme dir", s.ThemeDir())
	return &s, nil
}

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/creasty/defaults"
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
}

func (s SiteConfig) ThemeDir() string {
	return filepath.Join(s.RootDir, s.Theme.Path)
}

func SiteFromConfig() (*SiteConfig, error) {
	confPath, err := getConfigPath(0)
	if err != nil {
		return nil, err
	}
	bs, err := os.ReadFile(confPath)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	var sc SiteConfig
	err = yaml.Unmarshal(bs, &sc)

	if err != nil {
		return nil, fmt.Errorf("could not unmarshal config: %w", err)
	}
	err = defaults.Set(&sc)
	if err != nil {
		return nil, fmt.Errorf("could not set defaults: %w", err)
	}

	if sc.Version == "" {
		return nil, fmt.Errorf("config version not set")
	}

	if sc.Index.MinKeywordLength == 0 {
		sc.Index.MinKeywordLength = 3
	}

	sc.ConfigPath = confPath
	sc.RootDir = filepath.Dir(confPath)
	sc.SourceDir = filepath.Join(sc.RootDir, SourceDir)
	// if sc.ThemeDir == "" {
	// 	tp := filepath.Join(sc.RootDir, ThemesDir, sc.Theme.Path)
	// 	sc.ThemeDir = tp
	// } else if !filepath.IsAbs(sc.ThemeDir) {
	// 	// make it absolute and preserve relative paths, e.g /root/../../themes/foo
	// 	td := filepath.Join(sc.RootDir, sc.ThemeDir)
	// 	sc.ThemeDir = td
	// }

	if sc.PublicDir == "" {
		sd := filepath.Join(sc.RootDir, PublicDir)
		sc.PublicDir = sd
	} else if !filepath.IsAbs(sc.PublicDir) {
		// make it absolute and preserve relative paths, e.g /root/../../public
		pd := filepath.Join(sc.RootDir, sc.PublicDir)
		sc.PublicDir = pd
	}

	log.Info("Loaded config", "site name", sc.Name, "public dir", sc.PublicDir, "source dir", sc.SourceDir, "theme dir", sc.ThemeDir)
	return &sc, nil
}

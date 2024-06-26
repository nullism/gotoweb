package models

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type SiteConfig struct {
	Name       string
	Theme      ThemeConfig
	ConfigPath string // path to config.yaml
	RootDir    string
	SourceDir  string
	ThemeDir   *string `yaml:"theme_root_directory"`
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

	println("config: ", string(bs))
	sc.ConfigPath = confPath
	sc.RootDir = filepath.Dir(confPath)
	sc.SourceDir = filepath.Join(sc.RootDir, "source")
	if sc.ThemeDir == nil {
		tp := filepath.Join(sc.RootDir, "themes", sc.Theme.Name)
		sc.ThemeDir = &tp
	}
	return &sc, nil
}

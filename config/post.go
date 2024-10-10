package config

import "time"

// PostConfig contains configuration information for a given post.
type PostConfig struct {
	Title       string
	Blurb       string
	Tags        []string
	SkipIndex   bool `yaml:"skip_index"`
	SkipPublish bool `yaml:"skip_publish"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

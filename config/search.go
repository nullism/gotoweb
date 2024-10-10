package config

type SearchConfig struct {
	// MinKeywordLength is the minimum length of a keyword to be indexed.
	MinKeywordLength int      `yaml:"min_keyword_length" default:"3"`
	StopWords        []string `yaml:"stop_words"`
}

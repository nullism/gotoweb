package config

type ThemeConfig struct {
	Path string
	Args map[string]any
}

type ThemeAuthor struct {
	Name    string
	Email   string
	Website string
}

type ThemeInfo struct {
	Name        string
	Description string
	Author      ThemeAuthor
}

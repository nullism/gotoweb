package config

type MenuConfig struct {
	Items      []*MenuItem
	AutoPrefix bool `default:"true"`
}

type MenuItem struct {
	Title    string
	Href     string
	Children []*MenuItem
}

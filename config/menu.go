package config

type MenuConfig struct {
	Items      []*MenuItem
	AutoPrefix bool `default:"true"`
}

type MenuItem struct {
	Target   string `default:"_self"`
	Title    string
	Href     string
	Children []*MenuItem
}

package theme

var ExtraPageNames = []string{"about", "index", "search", "404"}

func IsExtraPage(name string) bool {
	for _, n := range ExtraPageNames {
		if n == name {
			return true
		}
	}
	return false
}

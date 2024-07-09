package theme

var ExtraPageNames = []string{"search"}

func IsExtraPage(name string) bool {
	for _, n := range ExtraPageNames {
		if n == name {
			return true
		}
	}
	return false
}

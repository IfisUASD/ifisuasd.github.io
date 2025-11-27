package components

import "strings"

func PrefixPath(path, lang string) string {
	if strings.HasPrefix(path, "/assets") {
		return path
	}
	if lang == "en" {
		return "/en" + path
	}
	return path
}

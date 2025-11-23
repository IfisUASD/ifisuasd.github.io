package components

func prefixPath(path, lang string) string {
	if lang == "en" {
		return "/en" + path
	}
	return path
}

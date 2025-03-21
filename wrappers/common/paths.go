package common

// combines a base path with a new path, returning a valid url
func ComposePath(base string, path string) string {
	// ensure base path ends in /
	if base[len(base)-1] != '/' {
		base += "/"
	}

	// if path is / return early
	if path == "/" {
		return base
	}

	// ensure new path ends in /
	if path[len(path)-1] != '/' && path[len(path)-1] != '}' {
		path += "/"
	}

	// ensure base path does not start with /
	if path[0] == '/' {
		path = path[1:]
	}

	return base + path
}

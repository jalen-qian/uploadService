package mPath

import "strings"

func Base(path string) string {
	if path == "" {
		return "."
	}
	// Strip trailing slashes.
	for len(path) > 0 && path[len(path)-1] == '\\' {
		path = path[0 : len(path)-1]
	}
	// Find the last element
	if i := strings.LastIndex(path, "\\"); i >= 0 {
		path = path[i+1:]
	}
	// If empty now, it had only slashes.
	if path == "" {
		return "\\"
	}
	return path
}

func Ext(path string) string {
	for i := len(path) - 1; i >= 0 && path[i] != '\\'; i-- {
		if path[i] == '.' {
			return path[i:]
		}
	}
	return ""
}

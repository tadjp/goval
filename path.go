package goval

import "strings"

type Path []string

// parse parsing path.
func parse(path string) Path {
	return strings.Split(path, ".")
}

package utils

import (
	"slices"
	"strings"
)

// Path is a split HTTP request path (segments after the leading slash).
type Path []string

// NewPathWithString creates a Path from an HTTP request path string.
func NewPathWithString(path string) Path {
	output := []string{}
	if path == "" || path == "/" {
		return output
	}
	if path[0] != '/' {
		return output
	}
	output = strings.Split(path, "/")[1:]
	for l := len(output) - 1; l >= 0; l-- {
		if output[l] != "" {
			return output[:l+1]
		}
	}
	return []string{}
}

// RemovePrefixes strips API prefixes such as "/api/<apiPath>/v.../"
// and returns the remaining resource path. apiPath is the configured
// path segment (e.g. "lightwell-next" or "lightwell").
func (p Path) RemovePrefixes(apiPath string) Path {
	output := []string(p)
	if len(output) < 4 {
		return []string{}
	}
	if output[0] != "api" {
		return []string{}
	}
	if output[1] != apiPath {
		return []string{}
	}
	if output[2] == "" || output[2][0] != 'v' {
		return []string{}
	}
	return output[3:]
}

// StartWithResources reports whether the path begins with any of the given
// resource segment lists. Combine with RemovePrefixes to match API resources.
func (p Path) StartWithResources(resources ...[]string) bool {
	for _, r := range resources {
		if len(p) < len(r) {
			continue
		}
		if slices.Equal([]string(p[:len(r)]), r) {
			return true
		}
	}
	return false
}

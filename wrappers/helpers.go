package wrappers

import (
	"strings"
)

func CurlyToColon(path string) string {
	path = strings.Replace(path, "{", ":", -1)
	path = strings.Replace(path, "}", "", -1)
	return path
}

package common

import (
	"path"
	"strings"
)

func Url(elems ...string) string {
	for i, elem := range elems {
		elems[i] = strings.Trim(elem, "/")
	}
	elems = append([]string{""}, elems...)
	return strings.Join(elems, "/")
}

func Path(elems ...string) string {
	return path.Join(elems...)
}

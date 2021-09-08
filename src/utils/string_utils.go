package utils

import (
	"regexp"
	"strings"
)

var (
	matchFirstCap      = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap        = regexp.MustCompile("([a-z0-9])([A-Z])")
	matchFirstCapCamel = regexp.MustCompile("(.)([a-z]+)")
	matchAllCapCamel   = regexp.MustCompile("([a-z])([A-Z])")
)

func ToSnakeCase(text string) string {
	snake := matchFirstCap.ReplaceAllString(text, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func ToCamelCase(text string) string {
	return strings.Title(
		strings.ReplaceAll(
			strings.ToLower(text),
			"_", " "))
}

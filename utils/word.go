package utils

import (
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

func Underscore(words string) string {
	words = strings.ToLower(words)

	expr := regexp.MustCompile("[[:space:][:blank:]]")
	strByte := expr.ReplaceAll([]byte(words), []byte("_"))

	expr = regexp.MustCompile("`[^a-z0-9]`i")
	strByte = expr.ReplaceAll(strByte, []byte("_"))

	expr = regexp.MustCompile("[!/']")
	strByte = expr.ReplaceAll(strByte, []byte("_"))

	words = strings.TrimPrefix(string(strByte), "_")
	words = strings.TrimSuffix(words, "_")

	return words
}

func Camelcase(word string) string {
	return strcase.ToCamel(word)
}

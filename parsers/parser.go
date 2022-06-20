package parsers

type Parser interface {
	Parse(dir string) []string
}

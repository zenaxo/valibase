package gen

import (
	_ "embed"
	"strings"
)

//go:embed templates/imports.ts.txt
var importsRaw string

//go:embed templates/helpers.ts.txt
var helpersRaw string

//go:embed templates/tail.ts.txt
var tailRaw string

func imports() string     { return normalize(importsRaw) }
func typeHelpers() string { return normalize(helpersRaw) }
func tail() string        { return normalize(tailRaw) }

func normalize(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	if s == "" || strings.HasSuffix(s, "\n") {
		return s
	}
	return s + "\n"
}

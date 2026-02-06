package utils

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// ToPascalCase converts s to PascalCase.
//
// It splits on non-alphanumeric characters and also detects common camel-case
// boundaries (e.g. "fooBar" -> "FooBar").
//
// Examples:
//
//	"hello world"     -> "HelloWorld"
//	"foo_bar-baz"     -> "FooBarBaz"
//	"fooBar"          -> "FooBar"
//	"JSONData"        -> "JsonData"
//	"api/v2/user_id"  -> "ApiV2UserId"
//	"  --a*b^c  "     -> "ABC"
func ToPascalCase(s string) string {
	if s == "" {
		return s
	}

	var words []string
	var curr []rune

	runes := []rune(s)
	for i := range runes {
		r := runes[i]
		prev := rune(0)
		next := rune(0)

		if i > 0 {
			prev = runes[i-1]
		}
		if i+1 < len(runes) {
			next = runes[i+1]
		}

		isAlnum := unicode.IsLetter(r) || unicode.IsDigit(r)
		if !isAlnum {
			if len(curr) > 0 {
				words = append(words, string(curr))
				curr = curr[:0]
			}
			continue
		}

		if len(curr) > 0 && camelBoundary(prev, r, next) {
			words = append(words, string(curr))
			curr = curr[:0]
		}

		curr = append(curr, r)
	}

	if len(curr) > 0 {
		words = append(words, string(curr))
	}

	var b strings.Builder
	for _, w := range words {
		b.WriteString(titleWord(w))
	}
	return b.String()
}

// ToLowerCamelCase converts s to lowerCamelCase.
//
// It reuses ToPascalCase and lowercases the first rune.
//
// Examples:
//
//	"hello world"     -> "helloWorld"
//	"foo_bar-baz"     -> "fooBarBaz"
//	"fooBar"          -> "fooBar"
//	"JSONData"        -> "jsonData"
//	"api/v2/user_id"  -> "apiV2UserId"
//	"  --a*b^c  "     -> "aBC"
func ToLowerCamelCase(s string) string {
	pascal := ToPascalCase(s)
	if pascal == "" {
		return pascal
	}

	runes := []rune(pascal)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

var irregularSingulars = map[string]string{
	"people":   "person",
	"men":      "man",
	"women":    "woman",
	"children": "child",
	"teeth":    "tooth",
	"feet":     "foot",
	"mice":     "mouse",
	"geese":    "goose",
}

// ToSingular attempts to convert a simple English plural to its singular form.
//
// This is a best-effort helper intended for identifier-like words such as
// "users", "todos", and "stories". It is not a full inflection engine.
//
// Examples:
//
//	"users"    -> "user"
//	"todos"    -> "todo"
//	"stories"  -> "story"
//	"classes"  -> "class"
//	"boxes"    -> "box"
//	"people"   -> "person"
func ToSingular(s string) string {
	if s == "" {
		return s
	}

	lower := strings.ToLower(s)

	if sg, ok := irregularSingulars[lower]; ok {
		return sg
	}

	if strings.HasSuffix(lower, "ies") && len(lower) > 3 {
		return s[:len(s)-3] + "y"
	}

	if strings.HasSuffix(lower, "ses") ||
		strings.HasSuffix(lower, "shes") ||
		strings.HasSuffix(lower, "ches") ||
		strings.HasSuffix(lower, "xes") ||
		strings.HasSuffix(lower, "zes") {
		return s[:len(s)-2]
	}

	if strings.HasSuffix(lower, "s") && !strings.HasSuffix(lower, "ss") {
		return s[:len(s)-1]
	}

	return s
}

func camelBoundary(prev, curr, next rune) bool {
	isL := unicode.IsLetter
	isU := func(r rune) bool { return unicode.IsUpper(r) }
	isLo := func(r rune) bool { return unicode.IsLower(r) }
	isD := unicode.IsDigit

	// digit <-> letter boundary: "v2user" or "user2"
	if (isD(prev) && isL(curr)) || (isL(prev) && isD(curr)) {
		return true
	}

	// lower -> upper boundary: "fooBar"
	if isLo(prev) && isU(curr) {
		return true
	}

	// acronym boundary: "JSONData" -> "JSON" + "Data"
	if isU(prev) && isU(curr) && isLo(next) && isL(next) {
		return true
	}

	return false
}

func titleWord(w string) string {
	if w == "" {
		return w
	}
	runes := []rune(w)

	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	return string(runes)
}

// ToQuotedStringArray returns the elements of s as a comma-separated list of
// quoted strings (using strconv.Quote).
//
// Example output:
//
//	"facebook.com", "instagram.com"
func ToQuotedStringArray(s []string) string {
	quoted := make([]string, len(s))
	for i, d := range s {
		quoted[i] = strconv.Quote(d)
	}
	return strings.Join(quoted, ", ")
}

// SizeExpression returns a Go-like constant expression for n (bytes) and a
// human-readable label.
//
// It prefers powers of 1024 that divide n exactly. For example, 1048576
// becomes ("1024 * 1024", "1 MB").
func SizeExpression(n int64) (expr string, label string) {
	if n <= 0 {
		return "0", "0 bytes"
	}

	units := []string{"bytes", "KB", "MB", "GB", "TB", "PB"}

	pow := 0
	size := int64(1)

	for pow+1 < len(units) {
		nextSize := size * 1024
		if nextSize <= 0 || n%nextSize != 0 {
			break
		}
		size = nextSize
		pow++
	}

	value := n / size

	// Build expr like "1024 * 1024 * 10"
	var parts []string
	for range pow {
		parts = append(parts, "1024")
	}
	if pow == 0 {
		parts = append(parts, fmt.Sprintf("%d", value))
	} else if value != 1 {
		parts = append(parts, fmt.Sprintf("%d", value))
	}
	expr = strings.Join(parts, " * ")

	unit := units[pow]
	if pow == 0 {
		if value == 1 {
			label = "1 byte"
		} else {
			label = fmt.Sprintf("%d bytes", value)
		}
	} else {
		label = fmt.Sprintf("%d %s", value, unit)
	}

	return
}

var tsReservedWords = map[string]struct{}{
	// JS keywords
	"break": {}, "case": {}, "catch": {}, "class": {}, "const": {}, "continue": {},
	"debugger": {}, "default": {}, "delete": {}, "do": {}, "else": {}, "export": {},
	"extends": {}, "finally": {}, "for": {}, "function": {}, "if": {}, "import": {},
	"in": {}, "instanceof": {}, "new": {}, "return": {}, "super": {}, "switch": {},
	"this": {}, "throw": {}, "try": {}, "typeof": {}, "var": {}, "void": {},
	"while": {}, "with": {}, "yield": {},

	// JS literals / restricted
	"null": {}, "true": {}, "false": {},

	// Future reserved / strict mode
	"enum": {}, "implements": {}, "interface": {}, "let": {}, "package": {},
	"private": {}, "protected": {}, "public": {}, "static": {}, "await": {},

	// TypeScript-specific / contextual that commonly cause pain as keys
	"type": {}, "readonly": {}, "abstract": {}, "as": {}, "asserts": {},
	"any": {}, "unknown": {}, "never": {}, "boolean": {}, "number": {}, "string": {},
	"symbol": {}, "bigint": {}, "object": {}, "keyof": {}, "infer": {}, "is": {},
	"namespace": {}, "declare": {}, "module": {}, "global": {}, "override": {},
	"require": {}, "from": {}, "of": {}, "satisfies": {},
}

func isASCIIIdentStart(r rune) bool {
	return (r >= 'A' && r <= 'Z') ||
		(r >= 'a' && r <= 'z') ||
		r == '_' || r == '$'
}

func isASCIIIdentPart(r rune) bool {
	return isASCIIIdentStart(r) || (r >= '0' && r <= '9')
}

func isValidTSIdentifier(name string) bool {
	if name == "" {
		return false
	}
	for i, r := range name {
		if i == 0 {
			if !isASCIIIdentStart(r) {
				return false
			}
			continue
		}
		if !isASCIIIdentPart(r) {
			return false
		}
	}
	return true
}

func quoteTSPropertyKey(name string) string {
	// Use single-quoted TS string literal and escape what's needed.
	// (We don't try to be a full JS string escaper; this is sufficient for keys.)
	escaped := strings.ReplaceAll(name, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `'`, `\'`)
	return "'" + escaped + "'"
}

// SanitizeFieldName returns name as a valid TypeScript property key.
//
// If name is not a valid identifier (e.g. contains '-', starts with a digit,
// includes spaces, etc.) or is a reserved word (e.g. "default", "class"),
// it is returned as a quoted property key.
func SanitizeFieldName(name string) string {
	if name == "" {
		return quoteTSPropertyKey(name)
	}

	if isValidTSIdentifier(name) {
		if _, reserved := tsReservedWords[name]; !reserved {
			return name
		}
	}

	return quoteTSPropertyKey(name)
}

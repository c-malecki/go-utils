package pstring

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

var personNameCleanRegexp = regexp.MustCompile(`[^A-Za-z0-9 .,()"'-\[\]]+`)

func RemoveNonAlphaDigitPunc(s string) string {
	return personNameCleanRegexp.ReplaceAllString(s, "")
}

func RemoveDiacritics(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	t := norm.NFD.String(s)
	b := make([]rune, 0, len(t))
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		b = append(b, r)
	}
	return string(b)
}

// Replaces diacritics with plain English letters and non-alphanumeric characters that are not .,()"'-[]\/
func CleanPersonName(s string) string {
	s = RemoveDiacritics(s)
	s = RemoveNonAlphaDigitPunc(s)
	return s
}

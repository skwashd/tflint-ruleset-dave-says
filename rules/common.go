package rules

import (
	"regexp"
	"strings"
)

// Common regex patterns
var (
	SnakeRegex = regexp.MustCompile(`^[a-z0-9_]+$`)
	KebabRegex = regexp.MustCompile(`^[a-z0-9-]+$`)
)

// SplitWordsOnUnderscore splits a string on underscores only (for labels)
func SplitWordsOnUnderscore(s string) []string {
	return strings.Split(s, "_")
}

// SplitWordsOnDash splits a string on dashes only (for name attributes)
func SplitWordsOnDash(s string) []string {
	return strings.Split(s, "-")
}

// SplitWords splits a string on underscores and dashes (for resource types)
func SplitWords(s string) []string {
	return strings.FieldsFunc(s, func(c rune) bool {
		return c == '_' || c == '-'
	})
}

// ContainsAnyWord checks if any word from needle slice appears in haystack slice
func ContainsAnyWord(haystack, needle []string) (bool, string) {
	for _, n := range needle {
		for _, h := range haystack {
			if strings.EqualFold(n, h) {
				return true, n
			}
		}
	}
	return false, ""
}

package logserver

import "strings"

type filter func(string) bool

func HasPrefix(prefix string) filter {
	return func(s string) bool {
		return strings.HasPrefix(s, prefix)
	}
}

func HasSuffix(suffix string) filter {
	return func(s string) bool {
		return strings.HasSuffix(s, suffix)
	}
}

func NoDuplicate() filter {
	var prevLine string
	return func(s string) bool {
		if s == prevLine {
			return false
		}
		prevLine = s
		return true
	}
}

package main

import (
	"strings"
)

func upper(s string) string {
	return strings.ToUpper(s)
}

func getSuffix(s, prefix string) (string, bool) {
	lenPre := len(prefix)
	sub := safeSubstring(s, lenPre)

	if strings.HasPrefix(upper(s), upper(sub)) {
		return s[lenPre:], true
	} else {
		return "", false
	}
}

func contains(a []string, s string) bool {
	for _, e := range a {
		if e == s {
			return true
		}
	}

	return false
}

func safeSubstring(s string, n int) string {
	if len(s) < n {
		return s
	} else {
		return s[:n]
	}
}

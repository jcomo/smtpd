package main

import (
	"strings"
)

func getSuffix(s, prefix string) (string, bool) {
	lenPre := len(prefix)
	sub := safeSubstring(s, lenPre)

	capS := strings.ToUpper(s)
	capSub := strings.ToUpper(sub)

	if !strings.HasPrefix(capS, capSub) {
		return "", false
	} else {
		return s[lenPre:], true
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

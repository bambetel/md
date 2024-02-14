package main

import (
	"fmt"
	"strings"
)

func blockType(l string) string {
	if hn := isHeading(l); hn > 0 {
		return fmt.Sprintf("h%d", hn)
	} else if isFence(l) != "" {
		return "fence"
	} else if isHR(l) {
		return "hr"
	}
	return "p"
}

func isFence(l string) string {
	if len(l) < 3 {
		return ""
	}
	if l[0:3] == "```" || l[0:3] == "~~~" {
		lang := strings.TrimSpace(l[3:])
		return lang
	}
	return ""
}

func isHeading(l string) int {
	hn := 0
	for _, c := range l {
		if c != '#' {
			break
		}
		hn++
	}
	if hn < 1 {
		return 0
	}
	if len(l) > hn {
		if l[hn] == ' ' {
			// TODO a word next
			return hn
		}
	}
	return 0
}

func isHR(l string) bool {
	var marker rune
	count := 0
	for _, c := range l {
		if c != ' ' {
			marker = c
			break
		}
	}
	for _, c := range l {
		if c != marker && c != ' ' {
			return false
		} else if c == marker {
			count++
		}
	}
	return (count >= 3)
}

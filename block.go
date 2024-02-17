package md

import (
	"fmt"
	"regexp"
	"strings"
)

// Functions below take a line with normalized spacing
// The indentations up to 3 spaces are already trimmed
// as non-relevant (always?)
func BlockType(l string) string {
	if hn := isHeading(l); hn > 0 {
		return fmt.Sprintf("h%d", hn)
	} else if isFence(l) != "" {
		return "fence"
	} else if isHR(l) {
		return "hr"
	} else if li := isLi(l); li != 0 {
		return fmt.Sprintf("li%c", li)
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

func isLi(l string) byte {
	i := 0
	if len(l) < 3 { // minimal list item is `- a`
		return 0
	}
	if l[1] == ' ' && isListPunctor(l[0]) {
		if len(l) >= 7 { // `- [?] a`
			// TODO same checklist punctor required (?)
			if l[2] == '[' && l[4] == ']' { // TODO variable spacing
				return 'x'
			}
		}
		return l[0]
	}
	// check OL `a.`
	if isOLLower(l) {
		return 'a'
	}
	if isOLUpper(l) {
		return 'A'
	}
	// check for OL first
	// minimal OL item: `1. a`
	for i < len(l)-1 && isDigit(l[i]) {
		i++
	}
	if i != 0 {
		if l[i] == '.' && l[i+1] == ' ' { // OL punctor after numbers
			// TODO require any list item content!
			return '1'
		}
	}
	return 0
}

var (
	reOLLower, reOLUpper *regexp.Regexp
)

func init() {
	reOLLower = regexp.MustCompile("^[a-z]\\.\\s+\\S+")
	reOLUpper = regexp.MustCompile("^[A-Z]\\.\\s+\\S+")
}

func isOLLower(l string) bool {
	return reOLLower.MatchString(l)
}

func isOLUpper(l string) bool {
	return reOLUpper.MatchString(l)
}

func isListPunctor(c byte) bool {
	return c == '-' || c == '*' || c == '+'
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func isSpace(c byte) bool { // TODO or rune?
	return c == ' ' || c == '\t'
}

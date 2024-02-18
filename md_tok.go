package md

import (
	"fmt"
	"regexp"
	"strings"
)

// Md tokenizer, parser functions
// Assumes (at least leading, meaningful in Md) spacing normalized to space
// characters using tabstop value.

type mdLine struct {
	Line   int
	Join   bool
	Prefix string
	Marker string
	Text   string
}

func BlankLine(nr int) mdLine {
	return mdLine{nr, false, "", "", ""}
}

// For now, just playground
func MdTok(lines []string, pre string) []mdLine {
	fmt.Printf("MdTok() %d %q\n", len(lines), pre)

	out := make([]mdLine, 0, 16)

	for i := 0; i < len(lines); i++ {
		l := lines[i]

		// blank line
		if isBlankLine(l) {
			out = append(out, BlankLine(i))
			continue
		}

		// // isolate a blockquote
		// if l[0] == '>' {
		// 	j := i // to be set as an element after list end
		// 	for ; j < len(lines); j++ {
		// 		if len(lines[j]) < 1 {
		// 			break
		// 		}
		// 		if lines[j][0] != '>' {
		// 			break
		// 		}
		// 	}
		//
		// 	block := make([]string, 0, 10)
		// 	for k := i; k < j; k++ {
		// 		block = append(block, lines[k][min(len(lines[k]), 2):])
		// 	}
		//
		// 	MdTok(block, pre+">>>>")
		// 	i = j - 1
		// 	continue
		// }

		join := false
		p := getLinePrefix(l)
		l = l[len(p):]
		mark, l := stripLineMark(l)

		if mark == "```" {
			lang := strings.TrimSpace(l)
			fmt.Println("Found block in language:", lang)
			start := i
			i++
			code := ""
			for i < len(lines) && strings.HasPrefix(getLinePrefix(lines[i]), p) && strings.Index(lines[i], "```") == -1 {
				// TODO error if unclosed fence (within line prefix)
				code += (lines[i][len(p):] + "\n")
				i++
			}
			fmt.Printf("Found fenced code block: %q\n", code)
			item := mdLine{start, false, p, mark, code}
			out = append(out, item)
			continue
		} else {
			if i > 0 && mark == "" {
				if out[len(out)-1].Text != "" && strings.HasPrefix(p, out[len(out)-1].Prefix) {
					join = true
				}
			}
			item := mdLine{i, join, p, mark, l}
			out = append(out, item)
		}
	}

	return out
}

// only a block mark or also blockquote?
// handling: HR, H1..6, LI 1., -, a., A., - [x],
// (?): ~~~/```
// TODO (?) return possible tag (?)
func stripLineMark(line string) (string, string) {
	if len(line) < 2 {
		return "", line
	}
	mark := ""
	reH := regexp.MustCompile("^(#+)\\s+")
	reLi := regexp.MustCompile("^(\\d+\\.|[a-zA-Z]\\.|[-+*]|[-+*]\\s+\\[[ x]\\])\\s+")

	if strings.HasPrefix(line, "```") || strings.HasPrefix(line, "~~~") {
		mark = "```" // normalize line[0:3]
	} else if isHR(line) {
		return line, ""
	} else if m := reH.FindString(line); len(m) > 0 {
		mark = m
	} else if m := reLi.FindString(line); len(m) > 0 {
		mark = m
	}

	if len(mark) == len(line) {
		// block type indicators (except for hr) require non-empty content
		// otherwise it is just a line of text
		return "", line
	}

	return line[:len(mark)], line[len(mark):]
}

func getLineMark(line string) string {
	if isHR(line) {
		return "---"
	}
	if isHeading(line) > 0 {
		return "H"
	}
	if listType := isLi(line); listType != 0 {
		return string(listType)
	}
	return ""
}

// Get Md line prefix that consists of whitespace and blockquote markers.
// Meaning can be relative to the previous line.
func getLinePrefix(l string) string {
	i := 0
	if len(l) < 1 {
		return ""
	}
	for ; i < len(l); i++ {
		if !isMdLinePrefixChar(l[i]) {
			break
		}
	}
	return l[0:i]
}

func isMdLinePrefixChar(c byte) bool {
	return c == ' ' || c == '>'
}

// replace characters in a string based on a map
func tr(s string, replace map[rune]rune) string {
	sb := strings.Builder{}
	for i, c := range s {
		if v, found := replace[c]; found {
			sb.WriteRune(v)
		} else {
			sb.WriteByte(s[i])
		}
	}
	return sb.String()
}

func isBlankLine(l string) bool {
	for _, c := range l {
		if c != ' ' && c != '\t' {
			return false
		}
	}
	return true
}

package md

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Md tokenizer, parser functions
// Assumes (at least leading, meaningful in Md) spacing normalized to space
// characters using tabstop value.

type mdLine struct {
	Nr     int
	Join   bool
	Prefix string
	Marker string
	Text   string
	Tag    string
}

// TODO: needed (?), rules, which tags, where (in pre blank lines are preserved literally)
func (ml *mdLine) IsBlank() bool {
	return ml.Text == "" && ml.Tag == ""
}

// TODO: needed?
func (ml *mdLine) LimitPrefix(l int) string {
	if len(ml.Prefix) <= l {
		return ""
	}
	return ml.Prefix[l:]
}

// TODO: function description

func MdTok(r io.Reader, parentPrefix string) []mdLine {
	scanner := bufio.NewScanner(r)
	// TODO: li prefix vs bq prefix

	// i - input line index
	lines := make([]string, 0, 32)
	for i := 0; scanner.Scan(); i++ {
		lines = append(lines, normalizeWS(scanner.Text(), tabstop))
	}
	// out := mdTokR(lines, "", 0)

	out := make([]mdLine, len(lines))

	// lastBlock (?)
	// lastTag := ""
	lastToken := ""
	lastPrefix := ""

	for i := 0; i < len(lines); i++ {
		l := lines[i]
		prefix := ""
		prefixLen := 0
		cutSpace := 0

		// de novo prefix
		if len(l) == 0 {
			fmt.Printf("%3d: NULL LINE @@@@@@@@@@\n", i+1)
			out[i] = mdLine{Nr: i}
			lastToken = ""
			continue
		}
		if len(l) > 0 { // condition order, nesting?
			k := 0
			for ; k < len(l) && k < len(lastPrefix); k++ {
				if l[k] != lastPrefix[k] {
					break
				}
			}
			prefixLen = k
			for ; prefixLen < len(l); prefixLen++ {
				if strings.HasPrefix(l[prefixLen:], "    ") {
					break
				}
				if l[prefixLen] == ' ' {
					continue
				}
				if l[prefixLen] == '>' {
					continue
				} else {
					break
				}
			}
			if prefixLen == len(l) {
				// blank line
				// update last prefix when differs in terms of '>' (?)
				// TODO: when update lastToken?
				fmt.Printf("%3d: %s @@@@@@@@@@\n", i+1, l[:prefixLen])
				out[i] = mdLine{Nr: i}
				lastToken = ""
				// TODO: close block here
				continue
			}
			// cutting trailing spaces - TODO: if after a li, 4 are acceptable!
			checkStart := k
			// if strings.HasPrefix(prefix, lastPrefix) {
			// 	checkStart = len(lastPrefix)
			// }
			if prefixLen > checkStart {
				for prefixLen >= checkStart+1 {
					if l[prefixLen-1] == ' ' {
						prefixLen--
					} else {
						break
					}
				}
				if prefixLen > 0 {
					cutSpace = 1
				}
			}
			prefix = l[:prefixLen]
		}
		if prefixLen > len(lastPrefix) {
			fmt.Printf("DETECTED INSIDE:\n   OUT: %s@\n   IN:  %s@\n\n", lastPrefix, prefix)
		}
		if strings.HasPrefix(lastToken, "li") {
			if strings.HasPrefix(l[prefixLen:], "    ") {
				fmt.Printf("DETECTED INSIDE LIST:\n   OUT: %s\n   IN:  %s\n\n", lines[i-1], lines[i])
				prefix += "    "
				prefixLen = len(prefix)
			}
		}

		mark, token := getLineMark(l[prefixLen:])
		join := false
		if mark == "" && lastToken != "" {
			join = true
		}

		cutLen := min(len(l), prefixLen+cutSpace+len(mark))
		text := l[cutLen:]
		fmt.Printf("%3d: %s@%s   last=%s\n", i+1, prefix, text, lastToken)
		if join {
			text += " %" + lastToken
		}

		out[i] = mdLine{Nr: i, Tag: token, Marker: mark, Prefix: prefix, Text: text, Join: join}

		if join {
			continue
		}

		lastToken = token
		lastPrefix = prefix
	}

	return out
}

func unescapeLine(l string) string {
	if len(l) >= 2 { // escaped char, actually to be meaningful, requires len of 3, 4?
		// TODO: handle only block marks escaping!
		if l[0] == '\\' {
			return l[1:]
		}
	}
	return l
}

// Tell if an element can be a multiline block, default true, but
// HR, ### Headings 1..6 and extension DL>DD are always single line
func isBreakable(tag string) bool {
	if tag == "dd" || tag == "hr" {
		return false
	}
	if len(tag) == 2 {
		if tag[0] == 'h' && '1' <= tag[1] && tag[1] <= '6' {
			return false
		}
	}

	return true
}

func getLineMark(line string) (mark string, tag string) {
	if len(line) < 2 || isBlankLine(line) { // impossible
		return "", ""
	}
	// TODO: allow optional beginning 1-3 spaces?

	switch {
	case strings.HasPrefix(line, "```"), strings.HasPrefix(line, "~~~"):
		return line[0:3], "pre"
	case reSettextUnderH1.MatchString(line): // TODO: TEST settext h1 if only '=' and after a regular p candidate
		return line, "h1"
	case strings.HasPrefix(line, ": "): // extension dl > (dt + dd+)+
		return ": ", "dd"
	case isHR(line):
		return line, "hr"
	case line[0] == '#':
		if m := reH.FindString(line); len(m) > 0 {
			return m, fmt.Sprintf("h%d", len(strings.TrimSpace(m)))
		}
	case line[0] == '[':
		if m := reRef.FindString(line); len(m) > 0 {
			return m, "li:ref"
		}
	default:
		if m := reLiNum.FindString(line); len(m) > 0 {
			return m, "li:1"
		} else if m := reLiLower.FindString(line); len(m) > 0 {
			return m, "li:a"
		} else if m := reLiUpper.FindString(line); len(m) > 0 {
			return m, "li:A"
		} else if m := reLiRomanLower.FindString(line); len(m) > 0 {
			// TODO: roman vs alpha ambiguous, interchangable
			return m, "li:i"
		} else if m := reLiRomanUpper.FindString(line); len(m) > 0 {
			return m, "li:I"
		} else if m := reLiCheck.FindString(line); len(m) > 0 {
			return m, "li:x"
		} else if m := reLiUL.FindString(line); len(m) > 0 {
			return m, fmt.Sprintf("li:ul%c", m[0])
		}
	}

	return "", ""
}

func isBlankLine(l string) bool {
	for _, c := range l {
		if c != ' ' && c != '\t' {
			return false
		}
	}
	return true
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
	if marker != '-' && marker != '_' && marker != '*' {
		return false
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

package md

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
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
	lines := make([]string, 0, 32)
	for i := 0; scanner.Scan(); i++ {
		lines = append(lines, normalizeWS(scanner.Text(), tabstop))
	}

	out := make([]mdLine, len(lines))

	lastPrefix := mdPrefix{}
	lastToken := "---"
	blockStart := 0

	for i := 0; i < len(lines); i++ {
		l := lines[i]
		mark := ""
		token := ""
		pushLi := false
		join := false

		prefix, prefixLen := lastPrefix.Common(l)
		// Grow Bq prefix if new present
		// TODO: bq marker indentation tolerance?
		for strings.HasPrefix(l[prefixLen:], ">") {
			cut := 1
			if strings.HasPrefix(l[prefixLen:], "> ") {
				cut = 2
			}
			prefix.PushBq()
			prefixLen += cut
		}
		samePrefix := lastPrefix.Equals(prefix) || lastPrefix.EqualsN(prefix, lastPrefix.Len()-1) && lastPrefix.PeekKind() == mdPrefixLi
		if samePrefix {
			fmt.Printf("%3d: prefix merged: [%v]==[%v]\n", i, lastPrefix, prefix)
		}

		text, cutTolerance := lineTolerance(l[prefixLen:])
		if isBlankLine(text) {
			token = "---"
			goto blank
		}
		if mark, token = getLineMark(text); strings.HasPrefix(token, "li") {
			pushLi = true
		}
		if token == "dd" && lastToken != "p" {
			mark = ""
			token = ""
		}

		// fmt.Printf("%d: [%s]\n\t%v\n\t%v --> %v\n\n", i, l, lastPrefix, prefix, samePrefix)

		if samePrefix && mark == "" && lastToken != "---" && isBreakable(lastToken) {
			join = true
			token = lastToken
			// check if exit pre with no trailing blank line when indent decreased
			if token == "pre" && !strings.HasPrefix(text, "    ") {
				token = "p"
				join = false
			}
		} else if token == "" {
			// open pre/p
			if strings.HasPrefix(text, "    ") {
				token = "pre"
			} else {
				token = "p"
			}
		}

		if token == "h1set" || token == "h2set" || token == "dd" {
			setTag := ""
			switch token {
			case "h1set":
				setTag = "h1"
			case "h2set":
				setTag = "h2"
			case "dd":
				setTag = "dt"
			}
			if lastToken == "p" {
				for j := blockStart; j < i; j++ {
					out[j].Tag = setTag
				}
			} else {
				if token == "h2set" {
					token = "hr"
				} else {
					mark = ""
					token = "p"
				}
			}
		}

		if !join {
			lastPrefix = prefix
			if pushLi {
				lastPrefix.PushLi()
			}
		}
	blank:
		text = unescapeLine(l[prefixLen+cutTolerance+len(mark):]) // + fmt.Sprintf(" cut=%d", cutTolerance)
		lastToken = token
		if !join {
			blockStart = i
		}
		out[i] = mdLine{Nr: i, Tag: token, Marker: mark, Prefix: prefix.String(), Text: text + fmt.Sprintf("[%d]", blockStart), Join: join}
	}

	return out
}

var reTolerate3Sp = regexp.MustCompile("^ {0,3}\\S")

func lineTolerance(s string) (string, int) {
	cut := 0
	if reTolerate3Sp.MatchString(s) {
		for ; cut < len(s) && cut <= 3; cut++ {
			if s[cut] != ' ' {
				break
			}
		}
		return s[cut:], cut
	}
	return s, 0
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
	if tag == "h1set" || tag == "h2set" {
		return false
	}
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
		return line[0:3], "pre:fence"
	case reSettextUnderH1.MatchString(line): // TODO: TEST settext h1 if only '=' and after a regular p candidate
		return line, "h1set"
	case reSettextUnderH2.MatchString(line): // TODO: TEST settext h1 if only '=' and after a regular p candidate
		return line, "h2set"
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

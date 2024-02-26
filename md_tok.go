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
	Line   int
	Join   bool
	Prefix string
	Marker string
	Text   string
	Tag    string
}

func (ml *mdLine) IsBlank() bool {
	return ml.Text == ""
}

func (ml *mdLine) LimitPrefix(l int) string {
	if len(ml.Prefix) <= l {
		return ""
	}
	return ml.Prefix[l:]
}

func BlankLine(nr int) mdLine {
	return mdLine{nr, false, "", "", "", ""}
}

// For now, just playground
func MdTok(r io.Reader, pre string) []mdLine {
	out := make([]mdLine, 0, 16)
	scanner := bufio.NewScanner(r)
	nextLine := func() string {
		return normalizeWS(scanner.Text(), tabstop)
	}

	for i := 0; scanner.Scan(); i++ {
		l := nextLine()

		// TODO Needed (?), mark previous element end (?)
		if isBlankLine(l) {
			out = append(out, BlankLine(i))
			continue
		}

		join := false
		p := getLinePrefix(l)
		l = l[len(p):]
		mark, l, tag := stripLineMark(l)

		// TODO could use also for reading literal text when indented code detected.
		// Can it be done here? Would save spoiling quoted (apparent) Markdown.
		if mark == "```" || mark == "~~~" {
			lang := strings.TrimSpace(l)
			fmt.Println("Found block in language:", lang) // TODO Just to use lang before handling added
			start := i
			i++
			code := ""
			for scanner.Scan() {
				l := nextLine()
				if !strings.HasPrefix(getLinePrefix(l), p) {
					break
				}
				if strings.Index(strings.TrimSpace(l), mark) != -1 {
					break
				}
				// TODO error if unclosed fence (within line prefix)
				code += (l[len(p):] + "\n")
				i++
			}

			item := mdLine{start, false, p, mark, code, "pre>code"}
			out = append(out, item)
			continue
		} else {
			if i > 0 && mark == "" && l != "" {
				// TODO cleaner: join when same prefix, same kind not separated with blank
				// BUT: extension-dl - dd can be only single Md line
				if out[len(out)-1].Tag != "dd" && out[len(out)-1].Text != "" && strings.HasPrefix(p, out[len(out)-1].Prefix) && equalQuote(p, out[len(out)-1].Prefix) {
					join = true
				}
			}
			item := mdLine{i, join, p, mark, l, tag}
			out = append(out, item)
		}
	}

	return out
}

// only a block mark or also blockquote?
// handling: HR, H1..6, LI 1., -, a., A., - [x],
// (?): ~~~/```
// TODO (?) return possible tag (?)
func stripLineMark(line string) (mark, text, tag string) {
	if len(line) < 2 {
		return "", line, ""
	}
	reH := regexp.MustCompile("^(#+)\\s+")
	reLi := regexp.MustCompile("^(\\d+\\.|^[a-zA-Z]\\.|^[-+*]|^[-+*]\\s+\\[[ x]\\])\\s+|^[ivx]+\\.|^[IVX]+\\.")
	reRef := regexp.MustCompile("^\\[\\w+\\]:\\s+")

	switch {
	case strings.HasPrefix(line, "```"), strings.HasPrefix(line, "~~~"):
		mark, tag = line[:3], "pre>code" // normalize line[0:3]
	case strings.HasPrefix(line, ": "): // extension dl > (dt + dd+)+
		mark, tag = ": ", "dd"
	case isHR(line):
		return line, "", "hr"
	case line[0] == '#':
		if m := reH.FindString(line); len(m) > 0 {
			mark, tag = m, fmt.Sprintf("h%d", len(strings.TrimSpace(m)))
		}
	case line[0] == '[':
		if m := reRef.FindString(line); len(m) > 0 {
			mark, tag = m, "ref"
		}
	default:
		if m := reLi.FindString(line); len(m) > 0 {
			mark, tag = m, "li" // TODO list type
		}
	}

	// block type indicators (except for hr) require non-empty content
	// otherwise it is just a line of text (block type depending on context)
	if len(mark) == len(line) {
		return "", line, ""
	}

	return line[:len(mark)], line[len(mark):], tag
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

func isBlankLine(l string) bool {
	for _, c := range l {
		if c != ' ' && c != '\t' {
			return false
		}
	}
	return true
}

// Check if prefixes are the same on the blockquote level
func equalQuote(pre1, pre2 string) bool {
	pre1 = strings.TrimRight(pre1, " ")
	pre2 = strings.TrimRight(pre2, " ")

	if len(pre1) != len(pre2) {
		return false
	}

	for i := range pre1 {
		if pre1[i] != pre2[i] {
			return false
		}
	}
	return true
}

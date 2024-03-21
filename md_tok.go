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

	out := make([]mdLine, len(lines))

	lastToken := "" // actually: last block kind
	lastPrefix := ""
	lastChild := ""
	lastBlock := 0

	for i := 0; i < len(lines); i++ {
		l := lines[i]
		prefixLen := 0

		// de novo prefix
		if len(l) == 0 {
			fmt.Printf("%3d: NULL LINE @@@@@@@@@@\n", i+1)
			out[i] = mdLine{Nr: i}
			lastToken = ""
			continue
		}
		if len(l) > 0 { // condition order, nesting?
			keep := 0 // common prefix length with lastPrefix; keep it
			for ; keep < len(l) && keep < len(lastChild); keep++ {
				if l[keep] != lastChild[keep] {
					break
				}
			}
			prefixLen = keep
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
				out[i] = mdLine{Nr: i, Prefix: l}
				lastToken = "" // keep container, reset block
				continue
			}
			// cutting trailing spaces -
			// TODO: (here?) if after a li, 4 are acceptable!
			checkStart := keep
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
			}
		}
		if prefixLen > len(lastPrefix) {
			fmt.Printf("DETECTED INSIDE:\n   OUT: %s@\n   IN:  %s@\n\n", lastPrefix, l[:prefixLen])
		}
		if prefixLen > 0 {
			if l[prefixLen-1] == '>' && prefixLen < len(l) {
				if l[prefixLen] == ' ' {
					prefixLen++
				}
			}
		}

		// TODO: should pre lines be marked as joined here at all?
		join := false

		mark, token := getLineMark(l[prefixLen:])
		if token == "pre:fence" {
			fmt.Printf("FENCE: %s\n", lines[i])
			out[i] = mdLine{Nr: i, Prefix: lines[i][:prefixLen], Tag: token, Marker: lines[i][prefixLen:], Join: false}
			prefix := lines[i][:prefixLen]
			for i++; i < len(lines); i++ {
				if len(lines[i]) < prefixLen {
					i--
					break // not closed!
				}
				if !strings.HasPrefix(lines[i], prefix) {
					i--
					break // not closed!
				}
				if strings.HasPrefix(lines[i][prefixLen:], mark) {
					out[i] = mdLine{Nr: i, Prefix: lines[i][:prefixLen], Tag: token, Text: "", Join: true}
					break // closing
				}
				out[i] = mdLine{Nr: i, Prefix: lines[i][:prefixLen], Tag: token, Text: lines[i][prefixLen:], Join: true}
			}
			lastToken = ""
			continue
		}

		// special fixes - look back
		setTag := ""
		if token == "h1set" {
			if lastToken == "p" {
				setTag = "h1"
			} else {
				mark = ""
				token = ""
			}
		} else if lastToken == "p" && token == "hr" &&
			strings.HasPrefix(strings.TrimSpace(l[prefixLen:]), "---") { // TODO: no spaces
			setTag = "h2"
		} else if token == "dd" {
			if lastToken == "p" {
				setTag = "dt"
			} else if lastToken != "dt" {
				if lastToken == "" {
					token = "" // no join check; would start a new block anyway?
				} else {
					token = ""
				}
				mark = ""
			}
		}
		if setTag != "" { // loop for edge cases; reasonable use would be a single line in case of h1, h2, dt
			for j := i - 1; j >= 0 && j >= lastBlock; j-- {
				out[j].Tag = setTag
				if j != lastBlock {
					out[j].Join = true
				}
			}
			join = true
			lastToken = ""
			goto joinSet
		}

		if mark != "" ||
			prefixLen < len(lastPrefix) ||
			len(strings.TrimRight(l[:prefixLen], " ")) != len(strings.TrimRight(lastPrefix, " ")) ||
			lastToken == "" {
			// new block
			join = false
		} else if mark == "" && lastToken != "" {
			join = true
		}
		if !join && token == "" {
			token = "p"
		}
		if strings.HasPrefix(l[prefixLen:], "    ") {
			if !join && token == "p" {
				token = "pre"
			} else if lastToken == "pre" {
				join = true
				token = "pre"
			}
		} else if lastToken == "pre" {
			token = "p"
			join = false
		}
	joinSet:

		// if (lastToken == "" || lastToken == "pre") && strings.HasPrefix(l[prefixLen:], "    ") {
		// 	token = "pre"
		// }
		if join && token == "" {
			token = lastToken
		}
		prefix := l[:prefixLen]

		cutLen := min(len(l), prefixLen+len(mark))
		text := unescapeLine(l[cutLen:])
		fmt.Printf("%3d: %s@%s   last=%s\n", i+1, prefix, text, lastToken)

		// // TODO: block marking specification
		// if token == "" && strings.HasPrefix(text, "    ") && (lastToken == "" || lastToken == "pre") {
		// 	token = "pre"
		// } else if join && token == "" {
		// 	join = false
		// 	token = "p"

		// save lastToken here
		if !isBreakable(token) {
			lastToken = ""
		} else if !join {
			lastToken = token
		}
		// }
		// if token == "" {
		// 	if strings.HasPrefix(text, "    ") {
		// 		token = "pre"
		// 	} else {
		// 		token = "p"
		// 	}
		// }
		out[i] = mdLine{Nr: i, Tag: token, Marker: mark, Prefix: prefix, Text: text, Join: join}

		if !join {
			// new block
			lastPrefix = prefix
			if strings.HasPrefix(token, "li") {
				lastChild = prefix + "    "
			} else {
				lastChild = prefix
			}
			lastBlock = i
		}
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
		return line[0:3], "pre:fence"
	case reSettextUnderH1.MatchString(line): // TODO: TEST settext h1 if only '=' and after a regular p candidate
		return line, "h1set"
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

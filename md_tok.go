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

func (ml *mdLine) IsBlank() bool {
	return ml.Text == "" && ml.Tag == ""
}

func (ml *mdLine) LimitPrefix(l int) string {
	if len(ml.Prefix) <= l {
		return ""
	}
	return ml.Prefix[l:]
}

func NewBlankMdLine(nr int) mdLine {
	return mdLine{nr, false, "", "", "", ""}
}

// Return the document lines with annotations, what they are in terms of block
// elements. It might be useful for syntax highlighting.
//
// No container elements are hinted, but the meaning of each line should be
// accurate apart from indentation levels that can change p/code or
// inconsistent list indentation and formatting.
func MdTok(r io.Reader, parentPrefix string) []mdLine {
	scanner := bufio.NewScanner(r)
	// TODO: li prefix vs bq prefix

	// i - input line index
	lines := make([]string, 0, 32)
	for i := 0; scanner.Scan(); i++ {
		lines = append(lines, normalizeWS(scanner.Text(), tabstop))
	}
	out := mdTokR(lines, "", 0)

	return out
}

func mdTokR(inlines []string, pre string, shift int) []mdLine {
	out := make([]mdLine, 0, 16)
	lines := make([]string, len(inlines))
	for i := range inlines {
		if len(inlines[i]) < shift {
			lines[i] = ""
		} else {
			lines[i] = inlines[i][shift:]
		}
	}
	fmt.Printf("mdTokR shift=%d\n", shift)
	fmt.Printf("received lines: %q\n", lines)

	for i := 0; i < len(lines); i++ {
		join := false
		container := []mdLine{}
		baseLine := i
		blockEnd := baseLine

		if isBlankLine(lines[i]) {
			fmt.Printf("%s   --%d--\n", pre, i+1)
			line := NewBlankMdLine(i)
			line.Prefix = pre
			out = append(out, line)
			continue
		}

		if strings.HasPrefix(lines[i], "    ") {
			// isolate indented pre
			item := mdLine{Nr: i, Text: lines[i], Tag: "pre", Prefix: pre}
			out = append(out, item)
			continue
		}

		// isolate bq container
		if lines[i][0] == '>' {
			// Assumes every single hard-wrapped line of a bq starts with `>`
			// with an equal spacing.
			// If GFM lazy principle was used, a breakable elemenet would
			// continue unless line was empty or on its own would start a new block.
			var s int
			for s = i + 1; s < len(lines); s++ {
				if len(lines[s]) == 0 {
					break
				}
				if lines[s][0] != '>' {
					break
				}
			}
			fmt.Printf("Found bq: [%d-%d]\n", i+1, s+1)
			block := mdTokR(lines[i:s], pre+">>>>", 1) // TODO: 2 assumes obligatory space
			out = append(out, block...)
			i = s - 1
			continue
		}

		// block processing
		mark, tagHeur := getLineMark(lines[i])
		tag := tagHeur
		if tag == "" {
			tag = "p"
		}

		// lookforward
		j := i + 1
		for ; j < len(lines) && isBreakable(tag); j++ {
			var nm string
			if isBlankLine(lines[j]) {
				break
			}
			if lines[j][0] == '>' {
				break
			}
			if strings.HasPrefix(lines[j], "     ") {

			} else if strings.HasPrefix(lines[j], "    ") {
				nm, _ = getLineMark(strings.TrimPrefix(lines[j], "    "))
				if nm != "" && strings.HasPrefix(tagHeur, "li") {
					break
				}
			} else {
				nm, _ = getLineMark(lines[j])
				if nm != "" { // TODO: nested inside a li
					break
				}
			}
		}
		if j > i+1 {
			fmt.Printf("found multiline block: %q\n", lines[i:j])
			i = j - 1
			blockEnd = j - 1
		}

		if strings.HasPrefix(tagHeur, "li") {
			tag = "li"
			l := i + 1
			if i < len(lines)-1 {
				firstBlank := false // phat items consume following blank lines
				blankOnly := true
				for l < len(lines) {
					if isBlankLine(lines[l]) {
						if l == i+1 {
							firstBlank = true
						}
						l++
					} else if strings.HasPrefix(lines[l], "    ") {
						blankOnly = false
						l++
					} else {
						break
					}
				}
				if !firstBlank || blankOnly {
					// return trailing blank lines if no phat item
					for isBlankLine(lines[l-1]) {
						l--
					}
				}
			}

			if l > i+1 {
				fmt.Printf("Found li content (%d:%d): %qEOT\n", i+1, l+1, lines[i:l])
				container = mdTokR(lines[i+1:l], pre+"....", 4)
				i = l - 1
			}
		}

		line := mdLine{Nr: baseLine, Tag: tag, Text: lines[baseLine][len(mark):], Prefix: pre, Join: join, Marker: mark}
		out = append(out, line)
		for ln := baseLine + 1; ln <= blockEnd; ln++ {
			line := mdLine{Nr: ln, Tag: tag, Text: lines[ln], Prefix: pre, Join: true}
			out = append(out, line)
		}

		if len(container) > 0 {
			out = append(out, container...)
		}

		fmt.Printf("%s %3d: %s\n", pre, i+1, lines[i])
	}
	return out
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
	reH := regexp.MustCompile("^(#{1,6})\\s+")
	// checklist before ul!
	// reLi := regexp.MustCompile("^(\\d+\\.|^[a-zA-Z]\\.|[-+*]\\s+\\[[ x]\\]\\s+|[-+*]\\s+|[ivx]+\\.|[IVX]+\\.)")
	reLiNum := regexp.MustCompile("^\\d+\\.\\s+")
	reLiLower := regexp.MustCompile("^[a-z]\\.\\s+")
	reLiUpper := regexp.MustCompile("^[A-Z]\\.\\s+")
	reLiRomanLower := regexp.MustCompile("^[ivx]+\\.")
	reLiRomanUpper := regexp.MustCompile("^[IVX]+\\.")
	reLiUL := regexp.MustCompile("^[-+*]\\s+")
	reLiCheck := regexp.MustCompile("^[-+*]\\s+\\[[ x]\\]\\s+")
	reRef := regexp.MustCompile("^\\[\\w+\\]:\\s+")
	reSettextUnderH1 := regexp.MustCompile("^={3,}\\s*$") // TODO handling trailing spaces?

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
			mark, tag = m, fmt.Sprintf("h%d", len(strings.TrimSpace(m)))
		}
	case line[0] == '[':
		if m := reRef.FindString(line); len(m) > 0 {
			mark, tag = m, "li:ref"
		}
	default:
		if m := reLiNum.FindString(line); len(m) > 0 {
			mark, tag = m, "li:1"
		} else if m := reLiLower.FindString(line); len(m) > 0 {
			mark, tag = m, "li:a"
		} else if m := reLiUpper.FindString(line); len(m) > 0 {
			mark, tag = m, "li:A"
		} else if m := reLiRomanLower.FindString(line); len(m) > 0 {
			// TODO: roman vs alpha ambiguous
			mark, tag = m, "li:i"
		} else if m := reLiRomanUpper.FindString(line); len(m) > 0 {
			mark, tag = m, "li:I"
		} else if m := reLiCheck.FindString(line); len(m) > 0 {
			mark, tag = m, "li:x"
		} else if m := reLiUL.FindString(line); len(m) > 0 {
			mark = m
			tag = fmt.Sprintf("li:ul%c", mark[0])
		}
	}

	// block type indicators (except for hr) require non-empty content
	// otherwise it is just a line of text (block type depending on context)
	if len(mark) == len(line) {
		return "", ""
	}

	return line[:len(mark)], tag
}

// only a block mark or also blockquote?
// handling: HR, H1..6, LI 1., -, a., A., - [x],
func stripLineMark(line string) (mark, text, tag string) {
	if len(line) < 2 {
		return "", line, ""
	}
	reH := regexp.MustCompile("^(#{1,6})\\s+")
	// checklist before ul!
	// reLi := regexp.MustCompile("^(\\d+\\.|^[a-zA-Z]\\.|[-+*]\\s+\\[[ x]\\]\\s+|[-+*]\\s+|[ivx]+\\.|[IVX]+\\.)")
	reLiNum := regexp.MustCompile("^\\d+\\.\\s+")
	reLiLower := regexp.MustCompile("^[a-z]\\.\\s+")
	reLiUpper := regexp.MustCompile("^[A-Z]\\.\\s+")
	reLiRomanLower := regexp.MustCompile("^[ivx]+\\.")
	reLiRomanUpper := regexp.MustCompile("^[IVX]+\\.")
	reLiUL := regexp.MustCompile("^[-+*]\\s+")
	reLiCheck := regexp.MustCompile("^[-+*]\\s+\\[[ x]\\]\\s+")
	reRef := regexp.MustCompile("^\\[\\w+\\]:\\s+")
	reSettextUnderH1 := regexp.MustCompile("^={3,}\\s*$") // TODO handling trailing spaces?

	switch {
	case strings.HasPrefix(line, "```"), strings.HasPrefix(line, "~~~"):
		return line[:3], line[3:], "pre" // normalize line[0:3], TODO pre > code
	case reSettextUnderH1.MatchString(line): // TODO: TEST settext h1 if only '=' and after a regular p candidate
		return line, "", "h1"
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
		if m := reLiNum.FindString(line); len(m) > 0 {
			mark, tag = m, "li:1"
		} else if m := reLiLower.FindString(line); len(m) > 0 {
			mark, tag = m, "li:a"
		} else if m := reLiUpper.FindString(line); len(m) > 0 {
			mark, tag = m, "li:A"
		} else if m := reLiRomanLower.FindString(line); len(m) > 0 {
			// TODO: roman vs alpha ambiguous
			mark, tag = m, "li:i"
		} else if m := reLiRomanUpper.FindString(line); len(m) > 0 {
			mark, tag = m, "li:I"
		} else if m := reLiCheck.FindString(line); len(m) > 0 {
			mark, tag = m, "li:x"
		} else if m := reLiUL.FindString(line); len(m) > 0 {
			mark = m
			tag = fmt.Sprintf("li:ul%c", mark[0])
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
// Finishes on `>` character
func getLinePrefixQ(l string) string {
	i := 0
	if len(l) < 1 {
		return ""
	}
	for ; i < len(l); i++ {
		if !isMdLinePrefixChar(l[i]) {
			break
		}
	}
	return strings.TrimRight(l[0:i], " ")
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

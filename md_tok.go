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
	out := make([]mdLine, 0, 16)
	scanner := bufio.NewScanner(r)
	// TODO: li prefix vs bq prefix

	// i - input line index
	lines := make([]string, 0, 32)
	for i := 0; scanner.Scan(); i++ {
		for scanner.Scan() {
			lines = append(lines, normalizeWS(scanner.Text(), tabstop))
		}
	}
	mdTokR(lines, "", 0)

	return out
}

func mdTokR(inlines []string, pre string, shift int) {
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
		if isBlankLine(lines[i]) {
			fmt.Printf("%s   --%d--\n", pre, i+1)
			continue
		}
		if lines[i][0] == '>' {
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
			mdTokR(lines[i:s], pre+">>>>", 1) // TODO: 2 assumes obligatory space
			i = s - 1
			continue
		}

		fmt.Printf("%s %3d: %s\n", pre, i+1, lines[i])
	}
}

func MdTok2(r io.Reader, pre string) []mdLine {
	out := make([]mdLine, 0, 16)
	scanner := bufio.NewScanner(r)
	nextLine := func() string {
		return normalizeWS(scanner.Text(), tabstop)
	}
	blockStart := 0

	for i := 0; scanner.Scan(); i++ {
		l := nextLine()

		// TODO Needed (?), mark previous element end (?)
		if isBlankLine(l) {
			out = append(out, NewBlankMdLine(i))
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
			i++
			for scanner.Scan() {
				l := nextLine()
				if !strings.HasPrefix(getLinePrefix(l), p) {
					break
				}
				if strings.Index(strings.TrimSpace(l), mark) != -1 {
					break
				}
				// TODO: error if unclosed fence (within line prefix)
				// TODO: repeat mark?
				out = append(out, mdLine{i, false, p, mark, l[len(p):], "pre"})
				i++
			}
		} else {
			if i > 0 && mark == "" && l != "" {
				// TODO cleaner: join when same prefix, same kind not separated with blank
				// BUT: extension-dl - dd can be only single Md line
				prev := &out[len(out)-1]
				if (isBreakable(prev.Tag) && !prev.IsBlank()) && strings.HasPrefix(p, prev.Prefix) && equalQuote(p, prev.Prefix) {
					join = true
				}
			}
			if i > 0 {
				prev := &out[len(out)-1]
				if tag == "dd" && (prev.IsBlank() || (prev.Tag != "" && prev.Tag != "dd" && prev.Tag != "dt")) {
					// apparent dd fix; also would be invalid after anything apart from p candidate
					// if prev.Join {
					// 	join = true
					// }
					tag = "p"    // determine it is a p for any dl checks
					l = mark + l // todo wiser
					mark = ""    // regular p
				} else if tag == "dd" && !prev.IsBlank() && prev.Tag != "dt" && prev.Tag != "dd" {
					// definition list fix
					// edge case: previous would normally be a hard-wrapped paragraph, but dt is expected to be one line
					if blockStart < len(out) {
						out[blockStart].Tag = "dt"
					} else {
						fmt.Println("WTF?!", i, blockStart, len(out))
					}
					// prev.Tag = "dt"
				} else if strings.HasPrefix(mark, "===") && tag == "h1" {
					// settext h1 fix
					prev.Tag = tag
					join = true // quick fix to skip empty <h1> or <h1> for the underline
				} else if strings.HasPrefix(mark, "---") && strings.Index(strings.TrimSpace(mark), " ") == -1 {
					// settext h2 or hr
					if prev.Tag == "" && !prev.IsBlank() { // a regular p candidate
						// force hr to part of h2
						prev.Tag = "h2"
						tag = "h2"
						join = true // quick fix to skip empty <h1> or <h1> for the underline
					}
				}
			}
			if !join {
				blockStart = i
			}
			item := mdLine{blockStart, join, p, mark, l, tag}
			out = append(out, item)
		}
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

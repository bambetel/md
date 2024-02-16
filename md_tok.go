package md

import (
	"fmt"
	"strings"

	"github.com/bambetel/colo"
)

func MdTok(lines []string, pre string) {
	fmt.Printf("MdTok() %d %q\n", len(lines), pre)
	g := colo.Green.Fmt
	r := colo.Red.Fmt
	c := colo.Cyan.Fmt

	for i := 0; i < len(lines); i++ {
		l := lines[i]
		fmt.Print(c(fmt.Sprintf("[%s:%d]", pre, i)))
		// blank line
		if len(l) < 1 || isBlankLine(l) {
			fmt.Println(i, "↔")
			continue
		}

		// consume a blockquote
		if l[0] == '>' {
			j := i                                           // to be set as an element after list end
			for ; j < len(lines) && len(lines[j]) > 0; j++ { // TODO safeguard len(lines) - 1 ???
				if lines[j][0] != '>' {
					break
				}
			}

			block := make([]string, 0, 10)
			for k := i; k < j; k++ {
				block = append(block, lines[k][2:]) // TODO handling `>`, len=1
			}
			MdTok(block, pre+">>>>")
			i = j - 1
			continue
		}

		t := BlockType(l)
		p := getLinePrefix(l)
		fmt.Print(pre)
		fmt.Print(fmt.Sprintf("%3d", i), r(fmt.Sprintf("%5s ", t)))
		fmt.Print(g(tr(p, map[rune]rune{' ': '_'})))
		fmt.Println(l[len(p):])
	}
	return
}

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

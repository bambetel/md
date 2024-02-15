package md

import (
	"fmt"
	"strings"

	"github.com/bambetel/colo"
)

func MdTok(lines []string) {
	for i, l := range lines {
		t := BlockType(l)
		if len(l) < 1 {
			fmt.Println("↔")
			continue
		}
		if l[0] == '>' {
			// TODO blockquote container handling
			fmt.Print(colo.BriYellow.Fmt("BQ"))
		}
		p := getLinePrefix(l)
		fmt.Print(fmt.Sprintf("%3d", i), colo.Red.Fmt(fmt.Sprintf("%5s ", t)))
		c := colo.Green
		fmt.Print(c.Fmt(tr(p, map[rune]rune{' ': '_'})))
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

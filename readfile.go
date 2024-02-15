package md

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// Read a markdown file and normalize spacing to spaces, tabs converted using tabstop

const (
	tabstop = 4
)

func ReadMdFile(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	lines, err := ReadMd(f)
	if err != nil {
		return nil, err
	}

	return lines, nil
}

func ReadMd(r io.Reader) ([]string, error) {
	s := bufio.NewScanner(r)
	f := make([]string, 0, 100)

	for s.Scan() {
		line := normalizeWS(s.Text(), tabstop)
		f = append(f, line)
	}

	return f, nil
}

// TODO any other spacing normalization needed?
func normalizeWS(s string, tabstop int) string {
	sb := strings.Builder{}
	i := 0
	for _, c := range s {
		if c == '\t' {
			dst := (i + tabstop) - (i+tabstop)%tabstop
			for ; i < dst; i++ {
				sb.WriteByte(' ')
			}
		} else {
			sb.WriteRune(c)
			i++
		}
	}
	return sb.String()
}

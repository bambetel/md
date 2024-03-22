package md

import (
	"fmt"
	"strings"
)

func getMdPrefix(line string) mdPrefix {
	mp := mdPrefix{}

	return mp
}

type mdPrefixType uint8

const (
	mdPrefixNone mdPrefixType = iota
	mdPrefixBq
	mdPrefixLi
)

type mdPrefixPart struct {
	Kind mdPrefixType
	Str  string
}

type mdPrefix struct {
	parts []mdPrefixPart
}

func (mp *mdPrefix) String() string {
	join := strings.Builder{}
	for i := range mp.parts {
		join.WriteString(mp.parts[i].Str)
	}
	return join.String()
}

func (mp *mdPrefix) Push(s string) {
	t := mdPrefixNone
	if strings.IndexByte(s, '>') != -1 {
		t = mdPrefixBq
	} else {
		t = mdPrefixLi
	}
	if t != mdPrefixNone {
		mp.parts = append(mp.parts, mdPrefixPart{t, s})
	}
}

func (mp *mdPrefix) PushLi() {
	mp.parts = append(mp.parts, mdPrefixPart{mdPrefixLi, "    "})
}
func (mp *mdPrefix) PushBq() {
	mp.parts = append(mp.parts, mdPrefixPart{mdPrefixBq, "> "}) // TODO: `> ` OR `>`?
}

func (mp *mdPrefix) Common(l string) (prefix mdPrefix, prefixLen int) {
match:
	for _, part := range mp.parts {
		switch part.Kind {
		case mdPrefixLi:
			if strings.HasPrefix(l[prefixLen:], "    ") {
				prefixLen += 4
				prefix.PushLi()
			} else {
				break match
			}
		case mdPrefixBq:
			if strings.HasPrefix(l[prefixLen:], ">") {
				cut := 1
				if strings.HasPrefix(l[prefixLen:], "> ") {
					cut = 2
				}
				prefixLen += cut
				prefix.PushBq()
			} else {
				break match
			}
		default:
			break match
		}
	}
	return
}

func (mp *mdPrefix) Match(s string) (n int, total bool) {
	for _, p := range mp.parts {
		switch p.Kind {
		case mdPrefixLi:
			if !strings.HasPrefix(s, "    ") {
				break
			}
			s = s[4:]
		case mdPrefixBq:
			cut := 1
			if s[0] != '>' {
				break
			}
			if len(s) >= 2 {
				if s[1] == ' ' {
					cut = 2
				}
			}
			s = s[cut:]
		default:
			break
		}
		n++
		if len(s) < 1 {
			return n, true
		}
	}
	return n, false
}
func (mp *mdPrefix) Same(s string) bool {
	if n, total := mp.Match(s); n == len(mp.parts) && total == true {
		return true
	}
	return false
}
func (mp *mdPrefix) HasPrefix(s string) bool {
	if _, total := mp.Match(s); total {
		return true
	}
	return false
}

func (mp *mdPrefix) Len() int {
	return len(mp.parts)
}

func (mp *mdPrefix) Pop() (res string) {
	if len(mp.parts) > 0 {
		res = mp.parts[len(mp.parts)-1].Str
		mp.parts = mp.parts[:len(mp.parts)-1]
	}
	return
}

func (mp *mdPrefix) Peek() (str string, kind mdPrefixType) {
	if len(mp.parts) > 0 {
		item := &mp.parts[len(mp.parts)-1]
		str = item.Str
		kind = item.Kind
	}
	return
}
func (mp *mdPrefix) PeekKind() (kind mdPrefixType) {
	if len(mp.parts) == 0 {
		return mdPrefixNone
	}
	return mp.parts[len(mp.parts)-1].Kind
}

func (mp *mdPrefix) New(s string, prev mdPrefix) {
	spaces := 0
	i := 0
	for ; i < len(s); i++ {
		if s[i] == '>' {
			spaces = 4 // or 5?
			mp.parts = append(mp.parts, mdPrefixPart{Kind: mdPrefixBq, Str: ""})
			continue
		}
		if s[i] == ' ' {
			if spaces == 0 {
				break
			}
			spaces--
		} else {
			break
		}
	}
	fmt.Printf("New: [%s]\n", s[:i])
}

func isMdPrefixChar(c byte) bool {
	switch c {
	case ' ', '>':
		return true
	default:
		return false
	}
}

func getMaxMdPrefix(s string) string {
	for i := range s {
		if !isMdPrefixChar(s[i]) {
			return s[:i]
		}
	}
	return s
}

func getNewMdPrefix(s string) string {
	spaces := 3
	i := 0
	for ; i < len(s); i++ {
		if s[i] == '>' {
			spaces = 4 // or 5?
			continue
		}
		if s[i] == ' ' {
			if spaces == 0 {
				break
			}
			spaces--
		} else {
			break
		}
	}
	return s[:i]
}

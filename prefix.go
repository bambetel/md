package md

import (
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
	Str  string // TODO: necessary? Maybe handling only normalized prefixes
}

type mdPrefix struct {
	parts []mdPrefixPart
}

// Return normalized string representing the Markdown prefix
func (mp *mdPrefix) String() string {
	join := strings.Builder{}
	for i := range mp.parts {
		join.WriteString(mp.parts[i].Str)
	}
	return join.String()
}

// Push normalized prefix for a container element
func (mp *mdPrefix) PushLi() {
	mp.parts = append(mp.parts, mdPrefixPart{mdPrefixLi, "    "})
}
func (mp *mdPrefix) PushBq() {
	mp.parts = append(mp.parts, mdPrefixPart{mdPrefixBq, "> "}) // TODO: `> ` OR `>`?
}

// Check if both prefixes are logically equivalent
func (mp *mdPrefix) Equals(other mdPrefix) bool {
	if len(mp.parts) != len(other.parts) {
		return false
	}
	for i := range mp.parts {
		if mp.parts[i].Kind != other.parts[i].Kind {
			return false
		}
	}
	return true
}

// Check if the first n elements of both prefixes are logically equivalent
// if any of the prefixes is too short, false is returned.
func (mp *mdPrefix) EqualsN(other mdPrefix, n int) bool {
	if len(mp.parts) < n || len(other.parts) < n {
		return false
	}
	for i := 0; i < n; i++ {
		if mp.parts[i].Kind != other.parts[i].Kind {
			return false
		}
	}
	return true
}

// Get the longest possible common prefix for mp and the string
// not in terms of apparent text columns but logical containers,
// for example " >  > " means the same as ">>".
// Return the number of characters consumed (final marker is either '>' or '> ').
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
			retry := 0
			for ; retry < 4 && prefixLen < len(l); retry++ {
				if l[prefixLen+retry] == '>' {
					break
				}
				if l[prefixLen+retry] != ' ' {
					break
				}
			}
			if retry <= 3 {
				prefixLen += retry
			}
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

// Get blockquote markers (exclude literal '>') from the string.
// Return number of characters consumed (final marker is either '>' or '> ').
func (mp *mdPrefix) AppendBqs(l string) (prefixLen int) {
	for {
		retry := 0
		ok := false
		for ; retry < 4 && prefixLen+retry < len(l); retry++ {
			if strings.HasPrefix(l[prefixLen+retry:], ">") {
				ok = true
				break
			}
		}
		if ok {
			prefixLen += retry
		} else {
			break
		}
		cut := 1
		if strings.HasPrefix(l[prefixLen:], "> ") {
			cut = 2
		}
		mp.PushBq()
		prefixLen += cut
	}
	return prefixLen
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

func (mp *mdPrefix) PeekKind() (kind mdPrefixType) {
	if len(mp.parts) == 0 {
		return mdPrefixNone
	}
	return mp.parts[len(mp.parts)-1].Kind
}

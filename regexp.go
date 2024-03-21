package md

import "regexp"

// Regular expression for particular block type beginnings
var (
	reH              = regexp.MustCompile("^(#{1,6})\\s+")
	reLiNum          = regexp.MustCompile("^\\d+\\.\\s+")
	reLiLower        = regexp.MustCompile("^[a-z]\\.\\s+")
	reLiUpper        = regexp.MustCompile("^[A-Z]\\.\\s+")
	reLiRomanLower   = regexp.MustCompile("^[ivx]+\\.")
	reLiRomanUpper   = regexp.MustCompile("^[IVX]+\\.")
	reLiUL           = regexp.MustCompile("^[-+*]\\s+")
	reLiCheck        = regexp.MustCompile("^[-+*]\\s+\\[[ x]\\]\\s+")
	reRef            = regexp.MustCompile("^\\[\\w+\\]:\\s+")
	reSettextUnderH1 = regexp.MustCompile("^={3,}\\s*$") // TODO: handling trailing spaces?
	reSettextUnderH2 = regexp.MustCompile("^-{3,}\\s*$") // TODO: and beginning up to 3?
)

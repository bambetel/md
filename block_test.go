package md

import (
	"testing"
)

func TestBlockTypeP(t *testing.T) {
	exp := "p"
	got := blockType("###Not a heading!")
	if got != exp {
		t.Errorf("invalid type, expected %q, got %q\n", exp, got)
	}
}

func TestBlockTypeH3(t *testing.T) {
	exp := "h3"
	got := blockType("### It is a Heading 3 3 3!")
	if got != exp {
		t.Errorf("invalid type, expected %q, got %q\n", exp, got)
	}
}

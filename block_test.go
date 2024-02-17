package md

import (
	"testing"
)

func TestBlockTypeP(t *testing.T) {
	exp := "p"
	got := BlockType("###Not a heading!")
	if got != exp {
		t.Errorf("invalid type, expected %q, got %q\n", exp, got)
	}
}

func TestBlockTypeH3(t *testing.T) {
	exp := "h3"
	got := BlockType("### It is a Heading 3 3 3!")
	if got != exp {
		t.Errorf("invalid type, expected %q, got %q\n", exp, got)
	}
}

func TestRegexpOLLower(t *testing.T) {
	if !reOLLower.MatchString("a. ol lower") {
		t.Errorf("Should match reOLLower")
	}
}

func TestIsOLLower(t *testing.T) {
	if got := BlockType("b. Bi"); got != "lia" {
		t.Errorf("invalid type, expected lia, got %q\n", got)
	}
}

func TestBlockTypeOLLower(t *testing.T) {
	exp := "lia"
	got := BlockType("a. This is an OL lowercase item")
	if got != exp {
		t.Errorf("invalid type, expected %q, got %q\n", exp, got)
	}
}

func TestBlockTypeOLUpper(t *testing.T) {
	exp := "liA"
	got := BlockType("A. This is an OL uppercase item")
	if got != exp {
		t.Errorf("invalid type, expected %q, got %q\n", exp, got)
	}
}

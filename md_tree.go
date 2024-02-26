package md

import (
	"fmt"
	"strings"
)

func MdTree(lines []mdLine, depth int, tag string) *MdNode {
	root := MdNode{
		Tag:  tag,
		Type: Element,
	}
	var prev *MdNode

	for i := 0; i < len(lines); i++ {
		l := lines[i]
		fmt.Printf("MdTree line %3d [%s] %v \n", i, l.Prefix[depth:], l)
		if l.Text == "" {
			continue
		}

		// blockquote handling
		j := i
		for j < len(lines) && strings.HasPrefix(lines[j].Prefix[depth:], ">") {
			j++
		}
		if j-i > 0 {
			// fmt.Printf("Found BQ to recurse, lines: %d-%d\n", i, j)
			res := MdTree(lines[i:j], depth+2, "blockquote")
			// fmt.Printf("MdTree returned (%d): %v", len(res.Children), res)
			root.Children = append(root.Children, *res)
			i = j - 1
			continue
		}

		// list handling
		// - if new li kind (type, prefix) - close if previous, open new list
		// - build a li, consume line and container if exists

		if l.Join && prev != nil {
			fmt.Printf("MdTree JOIN LINES %d: %s %s\n", i, prev.Text, l.Text)
			prev.JoinString(l.Text)
			continue
		}
		tag := "p"
		if l.Tag != "" {
			tag = l.Tag
		}

		n := MdNode{
			Type: Element,
			Tag:  tag,
			Text: l.Text,
		}

		root.Children = append(root.Children, n)
		prev = &root.Children[len(root.Children)-1]
	}
	return &root
}

func ProcessBQ(lines []mdLine, depth int) {
	fmt.Printf("ProcessBQ(%d)\n", len(lines))
	for i, l := range lines {
		fmt.Printf(">>>> %3d [%s] %s\n", i, l.Prefix[depth:], l.Text)
	}
}

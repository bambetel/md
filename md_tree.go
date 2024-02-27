package md

import (
	"fmt"
	"strings"
)

func MdTree(lines []mdLine, depth int, rootTag string) *MdNode {
	root := MdNode{
		Tag:  rootTag,
		Type: Element,
	}
	var prev *MdNode // for line joining TODO even necessary? no joining: dt, dd
	var currList *MdNode

	addChildNode := func(n MdNode) {
		root.Children = append(root.Children, n)
		if n.Tag == "ol" || n.Tag == "ul" || n.Tag == "dl" {
			currList = &root.Children[len(root.Children)-1]
		} else if isBreakable(n.Tag) && len(n.Text) > 0 {
			prev = &root.Children[len(root.Children)-1]
		} else if n.Tag != "li" && n.Tag != "dl" && n.Tag != "dt" { // any else?
			currList = nil
		}
	}

	// TODO: also differentiate punctor style
	requireList := func(tag string, addNode MdNode) {
		if currList != nil {
			if currList.Tag != tag {
				currList = nil
			}
		}
		if currList == nil {
			addChildNode(NewMdNodeElement(tag))
		}
		prev = nil
		currList.Children = append(currList.Children, addNode)
		if isBreakable(tag) {
			prev = &currList.Children[len(currList.Children)-1]
		}
	}

	for i := 0; i < len(lines); i++ {
		l := lines[i]
		fmt.Printf("MdTree line %3d [%s] %v \n", i, l.LimitPrefix(depth), l)
		if l.IsBlank() {
			if currList != nil {
				if currList.Tag != "dl" {
					currList = nil
				}
			}
			continue
		}

		// blockquote handling
		j := i
		for j < len(lines) && strings.HasPrefix(lines[j].LimitPrefix(depth), ">") {
			j++
		}
		if j-i > 0 {
			addChildNode(*MdTree(lines[i:j], depth+2, "blockquote"))
			i = j - 1
			continue
		}

		// TODO where should line joining be?
		// - should it be lookforward?
		if l.Join && prev != nil {
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

		if lines[i].Tag == "dt" || lines[i].Tag == "dd" {
			requireList("dl", n)
		} else if lines[i].Tag == "li" {
			// list handling
			// - if new li kind (type, prefix) - close if previous, open new list
			// - build a li, consume line and container if exists

			// *** Test before solving line joining

			// (???) TODO abstract func consumeLi() -> (simple|compound, container) either simple or compound (???)
			fmt.Println("A LI detected!")
			// next: 1. blank - lookforward for li.container, +tab after blank
			//       2. child list item (indent +1..4 or 7 spaces)
			//       3. sibling li
			//       4. something else - end of list

			// action: 1. try to find and isolate li.container, add as child,
			//         prepare for the next sibling li
			//    else 2. process either: a) child list, b) sibling li, c) end list

			// test simple nesting - next line is a child
			if i < len(lines)-1 && lines[i+1].Tag == "li" && prefixInside4s(lines[i].Prefix, lines[i+1].Prefix) {
				fmt.Println("Simple list nesting, skipping compound item check")
				// simplified: TODO
				// now just checking prefix inside to take into simple item with sublist
				// TODO: what about nested compound items?
				j := i + 1
				fmt.Printf("prefixInside4s(%q, %q) %v\n", lines[i].Prefix, lines[j].Prefix, prefixInside4s(lines[i].Prefix, lines[j].Prefix))
				for j < len(lines) && (lines[j].IsBlank() || prefixInside4s(lines[i].Prefix, lines[j].Prefix)) {
					fmt.Printf("prefixInside4s(%q, %q) %v\n", lines[i].Prefix, lines[j].Prefix, prefixInside4s(lines[i].Prefix, lines[j].Prefix))
					j++
				}
				if j-i > 0 {
					fmt.Println("Found simple nesting")
					fmt.Printf("---- lines: %d-%d %v\n", i, j, lines[i+1:j])
					res := MdTree(lines[i+1:j], depth+4, "ol") // TODO depth detected, list type
					n.Children = res.Children
				}
				i = j - 1
			} else if i < len(lines)-2 { // li.container even possible
				// simplified (?) just one blank line TODO - what with two???
				// TODO prefix offset handling for nesting!!!
				// NOTE: if li.block line hard-wrapping allowed, line merging should be done before the lookforward below
				if lines[i+1].IsBlank() && prefixInside4s(lines[i].Prefix, lines[i+2].Prefix) {
					fmt.Println("--- a compound li")
					// compound li handling - consumes any adjacent blank lines (!)
					// while either blank or prefix inside, put to the li.container
					j := i + 1
					// TODO: check prefixInside at least 4 spaces?
					for j < len(lines) && (lines[j].IsBlank() || prefixInside4s(lines[i].Prefix, lines[j].Prefix)) {
						j++
					}
					if j-i > 0 {
						fmt.Printf("Found LI to recurse, lines: %d-%d %v\n", i, j, lines[i+1:j])
						res := MdTree(lines[i+1:j], depth+4, "div")
						fmt.Printf("MdTree returned (%d): %v\n", len(res.Children), res)
						n.Children = res.Children
						i = j - 1
					}
				}
			}
			requireList("ol", n) // TODO list type, punctor style etc.
		} else { // regular block; not a list/dl item
			addChildNode(n)
			// currList = nil
			// root.Children = append(root.Children, n)
			prev = &root.Children[len(root.Children)-1]
		}
	}
	return &root
}

func prefixInside(in, pre string) bool {
	if len(in) >= len(pre) {
		return false
	}
	if pre[:len(in)] != in {
		return false
	}
	return true
}

func prefixInside4s(in, pre string) bool {
	if len(in) >= len(pre) {
		return false
	}
	if pre[:len(in)] != in {
		return false
	}
	diff := pre[len(in):]
	if strings.HasPrefix(diff, "    ") { // a TAB
		return true
	}
	return false
}

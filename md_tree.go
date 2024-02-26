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
	var prev *MdNode // for line joining TODO even necessary?
	var currList *MdNode

	for i := 0; i < len(lines); i++ {
		l := lines[i]
		fmt.Printf("MdTree line %3d [%s] %v \n", i, l.LimitPrefix(depth), l)
		if l.IsBlank() {
			currList = nil
			continue
		}

		// blockquote handling
		j := i
		for j < len(lines) && strings.HasPrefix(lines[j].LimitPrefix(depth), ">") {
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

		// TODO where should line joining be?
		// li handling needs lookforward to tell li type, so block lines could be joined here
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

		if lines[i].Tag != "li" {
			currList = nil
			root.Children = append(root.Children, n)
			prev = &root.Children[len(root.Children)-1]
		} else {
			// list handling
			// - if new li kind (type, prefix) - close if previous, open new list
			// - build a li, consume line and container if exists

			// *** Test before solving line joining

			// (???) TODO abstract func consumeLi() -> (simple|compound, container) either simple or compound (???)
			fmt.Println("A LI detected!")
			if currList == nil {
				fmt.Println("Creating a list TODO type")
				listNode := MdNode{
					Type: Element,
					Tag:  "ol",
				}
				root.Children = append(root.Children, listNode)
				currList = &root.Children[len(root.Children)-1]
			}
			// next: 1. blank - lookforward for li.container, +tab after blank
			//       2. child list item (indent +1..4 or 7 spaces)
			//       3. sibling li
			//       4. something else - end of list

			// action: 1. try to find and isolate li.container, add as child,
			//         prepare for the next sibling li
			//    else 2. process either: a) child list, b) sibling li, c) end list

			// test simple nesting - next line is a child
			if i < len(lines)-1 && lines[i+1].Tag == "li" && prefixInside(lines[i].Prefix, lines[i+1].Prefix) {
				fmt.Println("Simple list nesting, skipping compound item check")
				// simplified: TODO
				// now just checking prefix inside to take into simple item with sublist
				// TODO: what about nested compound items?
				j := i + 1
				for j < len(lines) && (lines[j].IsBlank() || prefixInside(lines[i].Prefix, lines[j].Prefix)) {
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
				if lines[i+1].IsBlank() && prefixInside(lines[i].Prefix, lines[i+2].Prefix) {
					fmt.Println("--- a compound li")
					// compound li handling - consumes any adjacent blank lines (!)
					// while either blank or prefix inside, put to the li.container
					j := i + 1
					for j < len(lines) && (lines[j].IsBlank() || prefixInside(lines[i].Prefix, lines[j].Prefix)) {
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
			currList.Children = append(currList.Children, n)
			prev = &currList.Children[len(currList.Children)-1]
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

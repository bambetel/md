package md

import (
	"strings"
)

func MdTree(lines []mdLine, depth int, rootTag string) *MdNode {
	root := MdNode{
		Tag:  rootTag,
		Type: Element,
	}
	var currList *MdNode

	addChildNode := func(n MdNode) {
		root.Children = append(root.Children, n)
		if n.Tag == "ol" || n.Tag == "ul" || n.Tag == "dl" {
			currList = &root.Children[len(root.Children)-1]
		} else if !strings.HasPrefix(n.Tag, "li") && n.Tag != "dl" && n.Tag != "dt" { // any else?
			currList = nil
		}
	}

	// TODO: also differentiate punctor style
	currListType := ""
	requireList := func(mark string, addNode MdNode) {
		if strings.HasPrefix(addNode.Tag, "li") {
			addNode.Tag = "li" // NOTE: not a reference
		}
		if currList != nil {
			if currListType != mark {
				currList = nil
			}
		}
		tag := "ol"
		if currList == nil {
			switch {
			case mark == "dl":
				tag = "dl"
			case strings.HasPrefix(mark, "li:ul"):
				tag = "ul"
			// TODO: checklist
			default:
				tag = "ol"
			}
			addChildNode(NewMdNodeElement(tag))
			currListType = mark
		}
		currList.Children = append(currList.Children, addNode)
	}
	// fmt.Printf("MD TREE BEGIN %s\n", rootTag)

	for i := 0; i < len(lines); i++ {
		l := lines[i]
		// fmt.Printf("MdTree line %3d [%s] %v \n", i, l.LimitPrefix(depth), l)
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
		for j < len(lines) && strings.HasPrefix(lines[j].LimitPrefix(depth), "> ") { // assume normalized to '> '
			// TODO valid bq mark check `> word` or `>`
			j += 1
		}
		if j-i > 0 {
			addChildNode(*MdTree(lines[i:j], depth+2, "blockquote"))
			i = j - 1
			continue
		}

		tag := "p"       // default tag; todo check meaningful indentation
		if l.Tag != "" { // TODO: check actual tag vs token here
			tag = l.Tag
		}
		var joinText string
		if tag == "pre" {
			for j := i + 1; j < len(lines); j++ {
				if lines[j].Tag != "pre" {
					break
				}
				joinText += "\n" + lines[j].Text
				i++
			}

		} else { // regular block - normalize ws
			// JOIN next lines if needed, advance i to tell next block
			// TODO: check if needs isBreakable checks
			for j = i + 1; j < len(lines); j++ {
				if !lines[j].Join {
					break
				}
				// fmt.Printf("JOIN TEXT %d %s\n", i, lines[j].Text)
				// TODO block text join wrapped lines here
				// TODO: trailing ws handling: strip/add
				joinText += lines[j].Text
				i++
			}
		}

		n := MdNode{
			Type: Element,
			Tag:  tag,
			Text: l.Text + joinText,
		}

		if l.Tag == "dt" || l.Tag == "dd" {
			requireList("dl", n)
		} else if strings.HasPrefix(l.Tag, "li") {
			// (???) TODO abstract func consumeLi() -> (simple|compound, container) either simple or compound (???)
			// next: 1. blank - lookforward for li.container, +tab after blank
			//       2. child list item (indent +1..4 or 7 spaces)
			//       3. sibling li
			//       4. something else - end of list

			if i < len(lines)-1 && strings.HasPrefix(lines[i+1].Tag, "li") && prefixInside4s(l.Prefix, lines[i+1].Prefix) {
				// Simple item with a nested list
				// next line is not blank.

				// simplified: TODO
				// now just checking prefix inside to take into simple item with sublist
				// TODO: what about nested compound items?
				j := i + 1
				// fmt.Printf("prefixInside4s(%q, %q) %v\n", l.Prefix, lines[j].Prefix, prefixInside4s(lines[i].Prefix, lines[j].Prefix))
				for j < len(lines) && (lines[j].IsBlank() || prefixInside4s(l.Prefix, lines[j].Prefix)) {
					// fmt.Printf("prefixInside4s(%q, %q) %v\n", l.Prefix, lines[j].Prefix, prefixInside4s(l.Prefix, lines[j].Prefix))
					j++
				}
				if j-i > 0 {
					// fmt.Println("Found simple nesting")
					// fmt.Printf("---- lines: %d-%d %v\n", i, j, lines[i+1:j])
					res := MdTree(lines[i+1:j], depth+4, "ol") // TODO depth detected, list type
					n.Children = res.Children
				}
				i = j - 1
			} else if i < len(lines)-2 { // li.container even possible
				// simplified (?) just one blank line TODO - what with two???
				// TODO prefix offset handling for nesting!!!
				// NOTE: if li.block line hard-wrapping allowed, line merging should be done before the lookforward below
				if lines[i+1].IsBlank() && prefixInside4s(l.Prefix, lines[i+2].Prefix) {
					// fmt.Println("--- a compound li")
					// compound li handling - consumes any adjacent blank lines (!)
					// while either blank or prefix inside, put to the li.container
					j := i + 1
					for j < len(lines) && (lines[j].IsBlank() || prefixInside4s(l.Prefix, lines[j].Prefix)) {
						j++
					}
					if j-i > 0 {
						// fmt.Printf("Found LI to recurse, lines: %d-%d %v\n", i, j, lines[i+1:j])
						res := MdTree(lines[i+1:j], depth+4, "div")
						// fmt.Printf("MdTree returned (%d): %v\n", len(res.Children), res)
						n.Children = res.Children
						i = j - 1
					}
				}
			}
			requireList(n.Tag, n) // TODO list type, punctor style etc.
		} else { // regular block; not a list/dl item
			addChildNode(n)
		}
	}
	return &root
}

// func prefixInside(in, pre string) bool {
// 	if len(in) >= len(pre) {
// 		return false
// 	}
// 	if pre[:len(in)] != in {
// 		return false
// 	}
// 	return true
// }

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

// Check how much spaces (till any other character) the second prefix is more indented than the first
// func prefixInsideN(in, pre string) int {
// 	if len(in) >= len(pre) {
// 		return 0
// 	}
// 	if pre[:len(in)] != in {
// 		return 0
// 	}
// 	diff := pre[len(in):]
// 	n := 0
// 	for ; n < len(diff); n++ {
// 		if diff[n] != '.' {
// 			break
// 		}
// 	}
// 	return n
// }

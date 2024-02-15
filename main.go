package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("vim-go")
	n := MdNode{Tag: "a", Children: []MdNode{
		{Tag: "strong", Text: "A text node in an anchor"},
	}}
	WriteHTML(n, os.Stdout)

	fmt.Println("\n===")
	l := "- <b>HTML tag<b> Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat."
	fmt.Println(blockType(l), " <- ", l)
	l = "### Heading 3 example ###"
	fmt.Println(blockType(l), " <- ", l)
	l = "###Not a heading!"
	fmt.Println(blockType(l), " <- ", l)
	l = "```javascript"
	fmt.Println(blockType(l), " <- ", l)
	l = "```  GO "
	fmt.Println(blockType(l), " <- ", l)
	l = "* * * "
	fmt.Println(blockType(l), " <- ", l)
	l = "______"
	fmt.Println(blockType(l), " <- ", l)
	l = "  ---- - -  -"
	fmt.Println(blockType(l), " <- ", l)
	l = "  ---- - -  "
	fmt.Println(blockType(l), " <- ", l)
	l = "--"
	fmt.Println(blockType(l), " <- ", l)
	l = " * * "
	fmt.Println(blockType(l), " <- ", l)
	l = "1. Ordered list"
	fmt.Println(blockType(l), " <- ", l)
	l = "1."
	fmt.Println(blockType(l), " <- ", l)
	l = "1. "
	fmt.Println(blockType(l), " <- ", l)
	l = "1.Not a list item"
	fmt.Println(blockType(l), " <- ", l)
	l = "-Not a list!"
	fmt.Println(blockType(l), " <- ", l)
	l = "- A list!"
	fmt.Println(blockType(l), " <- ", l)
	l = "- [x] Checked TODO item"
	fmt.Println(blockType(l), " <- ", l)
	l = "- [ ] Not checked"
	fmt.Println(blockType(l), " <- ", l)

	fmt.Println("\n*****")

	f, err := ReadMdFile("testdata/tabs2.md")
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}
	// fmt.Printf("MD file:\n%q\n", f)
	for i, l := range f {
		fmt.Printf("%03d: %s\n", i, l)
	}
}

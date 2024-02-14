package main

import "fmt"

func main() {
	fmt.Println("vim-go")
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
}

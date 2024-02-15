package main

import (
	"fmt"
	"io"
)

func WriteHTML(n MdNode, w io.Writer) error {
	// TODO if using MdNode.Ready, check if ready/parsed, else error
	// TODO get attrs from []Attr or child Attr nodes
	attr := map[string]string{"href": "https://www.wp.pl/", "title": "Wirtualna Polska"}
	w.Write([]byte(fmt.Sprintf("<%s", n.Tag)))
	for k, v := range attr {
		w.Write([]byte(fmt.Sprintf(" %s=\"%s\"", k, v)))
	}
	w.Write([]byte(">"))

	// innerHTML
	// TODO how to determine if use MdNode.Text
	// - or a function MdNode.HTML()
	if n.Children == nil {
		w.Write([]byte(n.Text))
	} else {
		for _, c := range n.Children {
			WriteHTML(c, w)
		}
	}

	w.Write([]byte(fmt.Sprintf("</%s>", n.Tag)))
	return nil
}

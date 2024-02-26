package md

import (
	"fmt"
	"io"
)

func WriteHTML(n MdNode, w io.Writer) error {
	// TODO if using MdNode.Ready, check if ready/parsed, else error
	attr := make(map[string]string, 0)
	for _, c := range n.Children {
		if c.Type == Attr {
			attr[c.Tag] = c.Text
		}
	}
	w.Write([]byte(fmt.Sprintf("<%s", n.Tag)))
	for k, v := range attr {
		w.Write([]byte(fmt.Sprintf(" %s=\"%s\"", k, v)))
	}
	w.Write([]byte(">"))

	// innerHTML
	// TODO how to determine (by MdNodeType?) if use MdNode.Text
	// - or a function MdNode.HTML()

	if n.Children == nil || n.Tag == "li" {
		w.Write([]byte(n.Text))
	}
	if n.Children != nil {
		for _, c := range n.Children {
			if c.Type == Element {
				WriteHTML(c, w)
			}
		}
	}

	w.Write([]byte(fmt.Sprintf("</%s>\n", n.Tag)))
	return nil
}

package md

type MdNode struct {
	Type     MdNodeType
	Children []MdNode
	Tag      string // or attribute name for attr node
	Text     string // only for text nodes
	Ready    bool   // already parsed
}

type MdNodeType int

const (
	Unknown MdNodeType = 0
	Element            = 1
	Attr               = 2 // how to use?
	Text               = 3
	Comment            = 8 // how to in md?
)

func NewMdNodeElement(tag string) MdNode {
	return MdNode{Type: Element, Tag: tag}
}

// TODO necessary?
func (n *MdNode) JoinString(str string) {
	// TODO separate with ws, a newline when trailing double space
	n.Text += " " + str
}

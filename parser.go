package main

type Node struct {
	Data      string
	Parent    *Node
	Children  []*Node
	IsVar     bool
	IsLiteral bool
}

func (n *Node) Clean() {
	for _, child := range n.Children {
		if child.Data == "" && len(child.Children) == 0 {
			idx := index(n.Children, child)
			n.Children = delete(n.Children, idx, idx+1)
			child.Parent = nil
			continue
		}
		child.Clean()
	}
}

func (n Node) String() string {
	if len(n.Children) == 0 {
		return "\"" + n.Data + "\""
	}
	var out = "(" + n.Data + " "
	for _, child := range n.Children {
		out += child.String() + " "
	}
	out += ")"
	return out
}

func (n *Node) AddChild(node *Node) {
	n.Children = append(n.Children, node)
	node.Parent = n
}

type Context struct {
	RootNode         *Node
	CurrentText      string
	CurrentLiteral   string
	InLiteral        bool
	StartLiteralChar byte
	InEscape         bool
	AddTo            *Node
	FirstArg         bool
	IsVar            bool
	IsLiteral        bool
	Buffer           []byte // Buffer to hold unread characters
	BufferPos        int    // Current position in the buffer
}

func CreateContext() Context {
	var out = Context{}
	out.RootNode = &Node{}
	out.AddTo = out.RootNode
	return out
}

func (ctx *Context) Clean() {
	ctx.RootNode.Clean()
	if ctx.RootNode == nil {
		ctx.RootNode = &Node{}
	}
}

func (ctx Context) String() string {
	return ctx.RootNode.String()
}

func ParseEscape(b byte) string {
	switch b {
	case "\\"[0]:
		return "\\"
	case "n"[0]:
		return "\n"
	case "t"[0]:
		return "\t"
	case "a"[0]:
		return "\a"
	case "\""[0]:
		return "\""
	case "'"[0]:
		return "'"
	case "b"[0]:
		return "\b"
	case "f"[0]:
		return "\f"
	case "r"[0]:
		return "\r"
	case "v"[0]:
		return "\v"
	default:
		return ""
	}
}

func (ctx *Context) AddChar(b byte) {
	switch {
	case ctx.InEscape:
		// Handle escape sequence
		ctx.CurrentLiteral += ParseEscape(b)
		ctx.InEscape = false

	case ctx.InLiteral:
		// If we encounter an escape character within a literal
		switch b {
		case '\\':
			ctx.InEscape = true
		case ctx.StartLiteralChar:
			// End of the literal
			ctx.InLiteral = false
			ctx.CurrentText += ctx.CurrentLiteral
			ctx.CurrentLiteral = ""
			ctx.IsVar = false
			ctx.IsLiteral = true
		default:
			// Add byte to the current literal
			ctx.CurrentLiteral += string(b)
		}

	case b == '(':
		// Start a new Node
		newNode := &Node{}
		ctx.AddTo.AddChild(newNode)
		ctx.AddTo = newNode
		ctx.FirstArg = true

	case b == ')':
		// Close current Node
		if ctx.FirstArg {
			// Only store the first argument in the Data field
			ctx.AddTo.Data = ctx.CurrentText
		} else {
			// Create a child node for any remaining text
			child := &Node{Data: ctx.CurrentText, IsVar: ctx.IsVar, IsLiteral: ctx.IsLiteral}
			ctx.AddTo.AddChild(child)
		}
		ctx.CurrentText = ""
		if ctx.AddTo.Parent != nil {
			ctx.AddTo = ctx.AddTo.Parent
		}

	case b == '"' || b == '\'':
		// Start a literal
		ctx.InLiteral = true
		ctx.StartLiteralChar = b

	case b == '\\':
		// Start an escape sequence
		ctx.InEscape = true

	case b == ' ' || b == '\t' || b == '\n':
		// Handle spaces as argument separators
		if ctx.CurrentText != "" {
			if ctx.FirstArg {
				ctx.AddTo.Data = ctx.CurrentText
				ctx.FirstArg = false
			} else {
				// Create a new child node for each subsequent argument
				child := &Node{Data: ctx.CurrentText, IsVar: ctx.IsVar, IsLiteral: ctx.IsLiteral}
				ctx.AddTo.AddChild(child)
			}
			ctx.CurrentText = ""
		}

	default:
		// Add byte to the current text
		ctx.CurrentText += string(b)
		var isLiteral = false
		for _, digit := range []byte(digits) {
			if b == digit {
				isLiteral = true
				break
			}
		}
		ctx.IsVar = !isLiteral
		ctx.IsLiteral = isLiteral
	}
}

func (ctx *Context) ParseBuffer() {
	for _, b := range ctx.Buffer {
		if b == 0 {
			break
		}
		ctx.AddChar(b)
	}
}

// Unread removes a character from the buffer
func (ctx *Context) Unread() {
	if ctx.BufferPos <= 0 {
		ctx.BufferPos = 0
		return
	}
	ctx.BufferPos--
}

func (ctx *Context) ReadChar(b byte) {
	if ctx.BufferPos >= len(ctx.Buffer)-1 {
		// If the buffer is full, expand it as needed
		ctx.Buffer = append(ctx.Buffer, b)
	} else {
		ctx.Buffer[ctx.BufferPos] = b
	}
	ctx.BufferPos++
}

func (ctx *Context) IsInvalid() bool {
	// Check if there is an unclosed literal
	if ctx.InLiteral {
		Error("invalid context: unclosed literal")
		return true
	}

	// Check if there is an unclosed escape sequence
	if ctx.InEscape {
		Error("invalid context: unended escape")
		return true
	}

	// Check if the current node is the root but it still has children, meaning unclosed parentheses
	if ctx.AddTo != ctx.RootNode {
		Error("invalid context: unclosed paren")
		return true
	}

	// If none of the above conditions are met, the context is considered valid
	return false
}

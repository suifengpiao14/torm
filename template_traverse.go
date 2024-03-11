package torm

import "text/template/parse"

func Traverse(cur parse.Node, visitor func(parse.Node)) {
	switch node := cur.(type) {
	case *parse.ActionNode:
		if node.Pipe != nil {
			Traverse(node.Pipe, visitor)
		}
	case *parse.BoolNode:
	case *parse.BranchNode:
		if node.Pipe != nil {
			Traverse(node.Pipe, visitor)
		}
		if node.List != nil {
			Traverse(node.List, visitor)
		}
		if node.ElseList != nil {
			Traverse(node.ElseList, visitor)
		}
	case *parse.BreakNode:
	case *parse.ChainNode:
	case *parse.CommandNode:
		if node.Args != nil {
			for _, arg := range node.Args {
				Traverse(arg, visitor)
			}
		}
	case *parse.CommentNode:
	case *parse.ContinueNode:
	case *parse.DotNode:
	case *parse.FieldNode:
	case *parse.IdentifierNode:
	case *parse.IfNode:
		Traverse(&node.BranchNode, visitor)
	case *parse.ListNode:
		if node.Nodes != nil {
			for _, child := range node.Nodes {
				Traverse(child, visitor)
			}
		}
	case *parse.NilNode:
	case *parse.NumberNode:
	case *parse.PipeNode:
		if node.Cmds != nil {
			for _, cmd := range node.Cmds {
				Traverse(cmd, visitor)
			}
		}
		if node.Decl != nil {
			for _, decl := range node.Decl {
				Traverse(decl, visitor)
			}
		}
	case *parse.RangeNode:
		Traverse(&node.BranchNode, visitor)
	case *parse.StringNode:
	case *parse.TemplateNode:
		if node.Pipe != nil {
			Traverse(node.Pipe, visitor)
		}
	case *parse.TextNode:
	case *parse.VariableNode:
	case *parse.WithNode:
		Traverse(&node.BranchNode, visitor)
	}
	visitor(cur)
}

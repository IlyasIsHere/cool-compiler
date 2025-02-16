package parser

import (
	"cool-compiler/ast"
	"fmt"
	"strings"
)

func SerializeExpression(exp ast.Expression) string {
	switch node := exp.(type) {
	case *ast.IntegerLiteral:
		return fmt.Sprintf("%d", node.Value)
	case *ast.StringLiteral:
		return fmt.Sprintf("%q", node.Value)
	case *ast.BooleanLiteral:
		return fmt.Sprintf("%t", node.Value)
	case *ast.ObjectIdentifier:
		return node.Value
	case *ast.UnaryExpression:
		return fmt.Sprintf("(%s %s)", node.Operator, SerializeExpression(node.Right))
	case *ast.BinaryExpression:
		return fmt.Sprintf("(%s %s %s)", SerializeExpression(node.Left), node.Operator, SerializeExpression(node.Right))
	case *ast.IfExpression:
		return fmt.Sprintf("if %s then %s else %s fi", SerializeExpression(node.Condition), SerializeExpression(node.Consequence), SerializeExpression(node.Alternative))
	case *ast.WhileExpression:
		return fmt.Sprintf("while %s loop %s pool", SerializeExpression(node.Condition), SerializeExpression(node.Body))
	case *ast.BlockExpression:
		var sb strings.Builder
		sb.WriteString("{ ")
		for i, expr := range node.Expressions {
			sb.WriteString(SerializeExpression(expr))
			if i < len(node.Expressions)-1 {
				sb.WriteString("; ")
			} else {
				sb.WriteString(";")
			}
		}
		sb.WriteString(" }")
		return sb.String()
	case *ast.LetExpression:
		var sb strings.Builder
		sb.WriteString("let ")
		for i, binding := range node.Bindings {
			sb.WriteString(binding.Identifier.Value)
			sb.WriteString(" : ")
			sb.WriteString(binding.Type.Value)
			if binding.Init != nil {
				sb.WriteString(" <- ")
				sb.WriteString(SerializeExpression(binding.Init))
			}
			if i < len(node.Bindings)-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString(" in ")
		sb.WriteString(SerializeExpression(node.In))
		return sb.String()
	case *ast.NewExpression:
		return fmt.Sprintf("new %s", node.Type.Value)
	case *ast.IsVoidExpression:
		return fmt.Sprintf("isvoid %s", SerializeExpression(node.Expression))
	case *ast.CallExpression:
		var sb strings.Builder
		sb.WriteString(SerializeExpression(node.Function))
		sb.WriteString("(")
		for i, arg := range node.Arguments {
			sb.WriteString(SerializeExpression(arg))
			if i < len(node.Arguments)-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString(")")
		return sb.String()
	case *ast.DotCallExpression:
		var sb strings.Builder
		sb.WriteString(SerializeExpression(node.Object))
		if node.Type != nil {
			sb.WriteString("@")
			sb.WriteString(node.Type.Value)
		}
		sb.WriteString(".")
		sb.WriteString(node.Method.Value)
		sb.WriteString("(")
		for i, arg := range node.Arguments {
			sb.WriteString(SerializeExpression(arg))
			if i < len(node.Arguments)-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString(")")
		return sb.String()
	case *ast.CaseExpression:
		var sb strings.Builder
		sb.WriteString("case ")
		sb.WriteString(SerializeExpression(node.Expression))
		sb.WriteString(" of")
		for _, branch := range node.Branches {
			sb.WriteString(" ")
			sb.WriteString(branch.Identifier.Value)
			sb.WriteString(" : ")
			sb.WriteString(branch.Type.Value)
			sb.WriteString(" => ")
			sb.WriteString(SerializeExpression(branch.Expression))
			sb.WriteString(";")
		}
		sb.WriteString(" esac")
		return sb.String()
	case *ast.AssignmentExpression:
		return fmt.Sprintf("(%s <- %s)", node.Identifier.Value, SerializeExpression(node.Expression))
	default:
		return fmt.Sprintf("%t", node)
	}
}

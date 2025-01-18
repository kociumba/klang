package generator

import (
	"fmt"
	"strings"

	"github.com/kociumba/klang/parser"
)

func (cg *CodeGen) generateExpression(expr *parser.Expression) string {
	if expr == nil {
		return ""
	}

	// First, handle the left-side Atom
	leftStr := cg.generateAtom(expr.Atom)

	// If there's no operator or right side, just return the Atom
	if expr.Op == nil || expr.Right == nil {
		return leftStr
	}

	// Handle the operator and right side
	rightStr := cg.generateExpression(expr.Right)

	// Special handling for assignment operators
	switch *expr.Op {
	case "+=", "-=", "*=", "/=", "%=", "&=", "|=", "^=", "<<=", ">>=":
		// These can be passed directly to C
		return fmt.Sprintf("%s %s %s", leftStr, *expr.Op, rightStr)

	// Comparison operators
	case "==", "!=", "<", ">", "<=", ">=":
		return fmt.Sprintf("%s %s %s", leftStr, *expr.Op, rightStr)

	// Logical operators
	case "&&", "||":
		return fmt.Sprintf("%s %s %s", leftStr, *expr.Op, rightStr)

	// Arithmetic operators
	case "+", "-", "*", "/", "%":
		return fmt.Sprintf("%s %s %s", leftStr, *expr.Op, rightStr)

	default:
		// Unknown operator - could panic or log error
		return fmt.Sprintf("/* ERROR: Unknown operator %s */ %s ? %s", *expr.Op, leftStr, rightStr)
	}
}

func (cg *CodeGen) generateAtom(atom *parser.Atom) string {
	if atom == nil {
		return ""
	}

	// Handle each possible type of Atom
	switch {
	case atom.Number != nil:
		return fmt.Sprintf("%d", *atom.Number)

	case atom.String != nil:
		// Strings need to be handled carefully in C
		return fmt.Sprintf("\"%s\"", *atom.String)

	case atom.FuncCall != nil:
		return cg.generateFunctionCall(atom.FuncCall)

	case atom.Ident != nil:
		// Check if this identifier has been replaced
		if replacement, exists := cg.replacements[*atom.Ident]; exists {
			return replacement
		}
		return *atom.Ident

	case atom.SubExpr != nil:
		// Wrap sub-expressions in parentheses
		return fmt.Sprintf("(%s)", cg.generateExpression(atom.SubExpr))

	default:
		return "/* ERROR: Invalid Atomic expression */"
	}
}

func (cg *CodeGen) generateFunctionCall(call *parser.FunctionCall) string {
	if call == nil {
		return ""
	}

	// Generate arguments
	args := make([]string, 0, len(call.Arguments))
	for _, arg := range call.Arguments {
		args = append(args, cg.generateExpression(arg))
	}

	// Special handling for built-in functions
	switch call.Name {
	case "println":
		// Convert println to printf
		return fmt.Sprintf("printf(%s\\n)", strings.Join(args, ", "))

	default:
		// Regular function call
		return fmt.Sprintf("%s(%s)", call.Name, strings.Join(args, ", "))
	}
}

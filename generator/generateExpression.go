package generator

import "github.com/kociumba/klang/parser"

func (cg *CodeGen) generateExpression(expr *parser.Expression) string {
	// 	if expr == nil {
	// 		return ""
	// 	}

	// 	leftStr := cg.generateAtom(expr.Left)
	// 	if expr.Op == nil || expr.Right == nil {
	// 		// Handle safe access without an operator
	// 		if expr.SafeAccess {
	// 			return fmt.Sprintf("(%s.is_present ? %s.value : make_empty_Int())",
	// 				leftStr, leftStr)
	// 		}
	// 		// Handle elvis operator
	// 		if expr.ElvisOp && expr.DefaultValue != nil {
	// 			defaultStr := cg.generateExpression(expr.DefaultValue)
	// 			return fmt.Sprintf("(%s.is_present ? %s.value : %s)",
	// 				leftStr, leftStr, defaultStr)
	// 		}
	// 		return leftStr
	// 	}

	// 	rightStr := cg.generateExpression(expr.Right)

	// 	// Handle assignment operators
	// 	switch *&expr.Op.Operator {
	// 	case "+=", "-=", "*=", "/=", "%=":
	// 		if expr.SafeAccess {
	// 			return fmt.Sprintf("(%s.is_present ? (%s.value %s %s.value) : make_empty_Int())",
	// 				leftStr,
	// 				leftStr,
	// 				strings.TrimRight(*&expr.Op.Operator, "="),
	// 				rightStr)
	// 		}
	// 		return fmt.Sprintf("%s.value %s %s.value",
	// 			leftStr,
	// 			*expr.Op,
	// 			rightStr)

	// 	case "+", "-", "*", "/", "%":
	// 		if expr.SafeAccess {
	// 			return fmt.Sprintf("(%s.is_present && %s.is_present ? make_nullable_Int(%s.value %s %s.value) : make_empty_Int())",
	// 				leftStr,
	// 				rightStr,
	// 				leftStr,
	// 				*expr.Op,
	// 				rightStr)
	// 		}
	// 		return fmt.Sprintf("make_nullable_Int(%s.value %s %s.value)",
	// 			leftStr,
	// 			*expr.Op,
	// 			rightStr)

	// 	case "==", "!=", "<", ">", "<=", ">=":
	// 		// Comparison operators should always check for null
	// 		return fmt.Sprintf("(%s.is_present && %s.is_present && %s.value %s %s.value)",
	// 			leftStr,
	// 			rightStr,
	// 			leftStr,
	// 			*expr.Op,
	// 			rightStr)

	// 	case "&&", "||":
	// 		// Logical operators should always check for null
	// 		return fmt.Sprintf("(%s.is_present && %s.is_present && %s.value %s %s.value)",
	// 			leftStr,
	// 			rightStr,
	// 			leftStr,
	// 			*expr.Op,
	// 			rightStr)

	// 	default:
	// 		return fmt.Sprintf("/* ERROR: Unknown operator %s */", *expr.Op)
	// 	}

	return ""
}

// func (cg *CodeGen) generateAtom(atom *parser.Atom) string {
// 	if atom == nil {
// 		return ""
// 	}

// 	switch {
// 	case atom.Number != nil:
// 		// Numbers are always wrapped in nullable type
// 		return fmt.Sprintf("make_nullable_Int(%d)", *atom.Number)

// 	case atom.String != nil:
// 		return fmt.Sprintf("\"%s\"", *atom.String)

// 	case atom.FuncCall != nil:
// 		// Function calls return nullable types by default
// 		result := cg.generateFunctionCall(atom.FuncCall)
// 		return fmt.Sprintf("make_nullable_Int(%s)", result)

// 	case atom.Ident != nil:
// 		if replacement, exists := cg.replacements[*atom.Ident]; exists {
// 			return replacement
// 		}
// 		// Identifiers are assumed to be nullable
// 		return *atom.Ident

// 	case atom.SubExpr != nil:
// 		return fmt.Sprintf("(%s)", cg.generateExpression(atom.SubExpr))

// 	default:
// 		return "/* ERROR: Invalid Atomic expression */"
// 	}
// }

// func (cg *CodeGen) generateFunctionCall(call *parser.FunctionCall) string {
// 	if call == nil {
// 		return ""
// 	}

// 	args := make([]string, 0, len(call.Arguments))
// 	for _, arg := range call.Arguments {
// 		args = append(args, cg.generateExpression(arg))
// 	}

// 	switch call.Name {
// 	case "println":
// 		// Special handling for println - unwrap values
// 		unwrappedArgs := make([]string, len(args))
// 		for i, arg := range args {
// 			unwrappedArgs[i] = fmt.Sprintf("%s.value", arg)
// 		}
// 		return fmt.Sprintf("printf(%s\\n)", strings.Join(unwrappedArgs, ", "))
// 	default:
// 		return fmt.Sprintf("%s(%s)", call.Name, strings.Join(args, ", "))
// 	}
// }

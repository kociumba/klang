package generator

import (
	"fmt"
	"strings"

	"github.com/kociumba/klang/parser"
)

func (cg *CodeGen) generateVarDecl(decl *parser.VarDecl) string {
	if decl == nil {
		return ""
	}

	if decl.Shorthand != nil {
		value := cg.generateExpression(decl.Shorthand.Init)
		var varType string
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			varType = "string"
		} else if strings.Contains(value, ".") {
			varType = "double"
		} else if value == "true" || value == "false" {
			varType = "bool"
		} else {
			varType = "int"
		}
		return fmt.Sprintf("    %s %s = %s;\n", varType, decl.Shorthand.Name, cg.generateExpression(decl.Shorthand.Init))
	} else if decl.Explicit != nil {
		return fmt.Sprintf("    %s %s = %s;\n", *decl.Explicit.Type, decl.Explicit.Name, cg.generateExpression(decl.Explicit.Init))
	}

	return ""
	// varType := decl.Type
	// if decl.Init == nil {
	// 	return fmt.Sprintf("    NULLABLE_TYPE(%s) %s = {0, 0};\n", varType, decl.Name)
	// }

	// return fmt.Sprintf("    NULLABLE_TYPE(%s) %s = {%s, 1};\n",
	// 	varType,
	// 	decl.Name,
	// 	cg.generateExpression(decl.Init))
}

package generator

import (
	"fmt"

	"github.com/kociumba/klang/parser"
)

func (cg *CodeGen) generateVarDecl(decl *parser.VarDecl) string {
	if decl == nil {
		return ""
	}

	// Create a nullable wrapper for uninitialized variables
	varType := decl.Type
	if decl.Init == nil {
		return fmt.Sprintf("NULLABLE_TYPE(%s) %s = {0, false};\n", *varType, decl.Name)
	}

	// For initialized variables, use direct initialization
	return fmt.Sprintf("%s %s = %s;\n", *varType, decl.Name, cg.generateExpression(decl.Init))
}

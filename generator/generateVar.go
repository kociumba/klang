package generator

import (
	"fmt"

	"github.com/kociumba/klang/parser"
)

func (cg *CodeGen) generateVarDecl(decl *parser.VarDecl) string {
	if decl == nil {
		return ""
	}

	if decl.Init != nil {
		return fmt.Sprintf("%s %s = %s;", decl.Type, decl.Name, cg.generateExpression(decl.Init))
	}

	return fmt.Sprintf("%s %s;", decl.Type, decl.Name)
}

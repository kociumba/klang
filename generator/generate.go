package generator

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/kociumba/klang/parser"
)

type CodeGen struct {
	// Return type checking ?
	functions map[string]string // map[funcName]returnType
	// Return type checking ?
	currentFunction string
}

func NewCodeGen() *CodeGen {
	return &CodeGen{
		functions:       make(map[string]string),
		currentFunction: "",
	}
}

// Simple plan here:
// - generate macros for defines
// - generate function prototypes
// - generate function implementations
func (cg *CodeGen) Generate(root *parser.Root) string {
	var output strings.Builder
	tm, err := NewTemplateManager()
	if err != nil {
		log.Fatal(err)
	}

	// headers
	if err := tm.GenerateToBuffer("headers", nil, "headers"); err != nil {
		log.Warn(err)
	}

	// typedefs
	//
	// TODO: move this to nodes when i add parsing for it
	// if err := tm.GenerateToBuffer("typedefs.go.tmpl", nil, "typedefs"); err != nil {
	// 	log.Warn(err)
	// }

	// parse anything extracted from the ast
	for _, node := range root.Nodes {
		switch node.(type) {
		case parser.Typedecl:
			t := node.(parser.Typedecl)
			if err := tm.GenerateToBuffer("type", t, "typedefs"); err != nil {
				log.Warn(err)
			}
		case parser.VarDecl:
			v := node.(parser.VarDecl)
			if err := tm.GenerateToBuffer("var_decl", v, "variables"); err != nil {
				log.Warn(err)
			}
		case parser.Definedecl:
			def := node.(parser.Definedecl)
			if err := tm.GenerateToBuffer("macro", def, "macros"); err != nil {
				log.Warn(err)
			}
		case parser.FunctionDef:
			fn := node.(parser.FunctionDef)
			// pp.Println(fn.Body.Statements)
			if err := tm.GenerateToBuffer("function_prototype", fn, "function_prototypes"); err != nil {
				log.Warn(err)
			}
			if err := tm.GenerateToBuffer("function", fn, "functions"); err != nil {
				log.Warn(err)
			}
		}
	}

	if err := tm.WriteBuffersToBuilder(&output); err != nil {
		log.Fatal(err)
	}

	return output.String()
}

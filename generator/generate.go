package generator

import (
	"fmt"
	"strings"

	"github.com/kociumba/klang/parser"
)

// CodeGen represents the C code generator
type CodeGen struct {
	// Track declared functions for type checking return values
	functions map[string]string // map[funcName]returnType
	// Track replacements for keyword substitution
	replacements map[string]string
	// You might want to track current function for return type checking
	currentFunction string
}

func NewCodeGen() *CodeGen {
	return &CodeGen{
		functions:       make(map[string]string),
		replacements:    make(map[string]string),
		currentFunction: "",
	}
}

// Example method to generate C code from the AST
func (cg *CodeGen) Generate(program *parser.Program) string {
	var output strings.Builder

	// Add standard C headers
	output.WriteString("#include <stdio.h>\n")
	output.WriteString("#include <stdlib.h>\n")
	output.WriteString("#include <stdbool.h>\n\n")
	// Add your string typedef
	output.WriteString("typedef char* string;\n\n")

	// Process replacements first
	for _, decl := range program.Declarations {
		if decl.Replace != nil {
			cg.replacements[decl.Replace.From] = decl.Replace.To
		}
	}

	// Generate function prototypes first (for forward declarations)
	for _, decl := range program.Declarations {
		if decl.Define != nil || decl.Fun != nil {
			var fn *parser.FunctionDef
			if decl.Define != nil {
				fn = decl.Define
				// Handle define differently - maybe as macro?
				output.WriteString(cg.generateMacro(fn))
			} else {
				fn = decl.Fun
				// Add to function map for return type checking
				if fn.ReturnType != nil {
					cg.functions[fn.Name] = *fn.ReturnType
				}
				output.WriteString(cg.generateFunctionPrototype(fn))
			}
		}
	}

	// Generate function implementations
	for _, decl := range program.Declarations {
		if decl.Fun != nil {
			cg.currentFunction = decl.Fun.Name
			output.WriteString(cg.generateFunction(decl.Fun))
			cg.currentFunction = ""
		}
	}

	return output.String()
}

func (cg *CodeGen) generateBlock(block *parser.Block) string {
	if block == nil {
		return ""
	}

	var output strings.Builder
	for _, stmt := range block.Statements {
		output.WriteString(cg.generateStatement(&stmt))
	}
	return output.String()
}

func (cg *CodeGen) generateStatement(stmt *parser.Statement) string {
	if stmt == nil {
		return ""
	}

	if stmt.Return != nil {
		return fmt.Sprintf("    return %s;\n", cg.generateExpression(stmt.Return.Value))
	}

	if stmt.VarDecl != nil {
		return cg.generateVarDecl(stmt.VarDecl)
	}

	if stmt.ForLoop != nil {
		return cg.generateForLoop(stmt.ForLoop)
	}

	if stmt.IfStmt != nil {
		return cg.generateIfStatement(stmt.IfStmt)
	}

	if stmt.Expr != nil {
		// This will handle function calls through the Expression struct
		return fmt.Sprintf("    %s;\n", cg.generateExpression(stmt.Expr))
	}

	return ""
}

func (cg *CodeGen) generateForLoop(loop *parser.ForLoop) string {
	// Handle the range function specially
	return fmt.Sprintf("    for(int %s = %s; %s < %s; %s++) {\n%s    }\n",
		loop.Iterator,
		cg.generateExpression(loop.Start),
		loop.Iterator,
		cg.generateExpression(loop.End),
		loop.Iterator,
		cg.generateBlock(loop.Body))
}

func (cg *CodeGen) generateIfStatement(ifStmt *parser.IfStatement) string {
	return fmt.Sprintf("    if(%s) {\n%s    }\n",
		cg.generateExpression(ifStmt.Condition),
		cg.generateBlock(ifStmt.Body))
}

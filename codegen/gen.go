package codegen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kociumba/klang/parser"
)

type Generator struct {
	output strings.Builder
	indent int
}

func New() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(program *parser.Program) string {
	g.writeLine("#include <stdio.h>")
	g.writeLine("#include <stdlib.h>")
	g.writeLine("")

	for _, stmt := range program.Statements {
		g.generateStatement(stmt)
	}

	return g.output.String()
}

func (g *Generator) generateStatement(stmt parser.Statement) {
	switch s := stmt.(type) {
	case *parser.FunctionLiteral:
		g.generateFunction(s)
	case *parser.VarStatement:
		g.generateVarStatement(s)
	case *parser.ReturnStatement:
		g.generateReturnStatement(s)
	}
}

func (g *Generator) generateFunction(fun *parser.FunctionLiteral) {
	g.write("%s %s(", fun.ReturnType, fun.Name)

	for i, p := range fun.Parameters {
		if i > 0 {
			g.write(", ")
		}
		g.write("%s %s", p.Type, p.Name)
	}
	g.writeLine(") {")

	g.indent++
	for _, stmt := range fun.Body.Statements {
		g.generateStatement(stmt)
	}
	g.indent--
	g.writeLine("}")
	g.writeLine("")
}

func (g *Generator) generateVarStatement(stmt *parser.VarStatement) {
	if stmt.Value != nil {
		g.writeLine("%s %s = %s;", stmt.Type, stmt.Name, g.generateExpression(stmt.Value))
	} else {
		g.writeLine("%s %s;", stmt.Type, stmt.Name)
	}
}

func (g *Generator) generateReturnStatement(stmt *parser.ReturnStatement) {
	g.writeLine("return %s;", g.generateExpression(stmt.ReturnValue))
}

func (g *Generator) generateExpression(expr parser.Expression) string {
	switch e := expr.(type) {
	case *parser.Identifier:
		return e.Value
	case *parser.IntegerLiteral:
		return strconv.FormatInt(e.Value, 10)
	default:
		return ""
	}
}

func (g *Generator) write(format string, args ...interface{}) {
	g.output.WriteString(strings.Repeat("    ", g.indent))
	fmt.Fprintf(&g.output, format, args...)
}

func (g *Generator) writeLine(format string, args ...interface{}) {
	g.write(format+"\n", args...)
}

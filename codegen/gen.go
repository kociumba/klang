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
	g.writeLine("#include <string.h>")
	g.writeLine("")

	g.writeLine("typedef char* string;")
	g.writeLine("")

	for _, stmt := range program.Statements {
		if fn, ok := stmt.(*parser.FunctionLiteral); ok {
			g.generateFunction(fn)
		}
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
	case *parser.ForStatement:
		g.generateForStatement(s)
	case *parser.IfStatement:
		g.generateIfStatement(s)
	case *parser.ExpressionStatement:
		g.generateExpressionStatement(s)
	}
}

func (g *Generator) generateExpressionStatement(stmt *parser.ExpressionStatement) {
	g.writeLine("%s;", g.generateExpression(stmt.Expression))
}

func (g *Generator) generateExpression(expr parser.Expression) string {
	switch e := expr.(type) {
	case *parser.Identifier:
		return e.Value
	case *parser.IntegerLiteral:
		return strconv.FormatInt(e.Value, 10)
	case *parser.CallExpression:
		return g.generateCallExpression(e)
	case *parser.InfixExpression:
		return g.generateInfixExpression(e)
	default:
		return ""
	}
}

func (g *Generator) generateIfStatement(stmt *parser.IfStatement) {
	g.writeLine("if (%s) {", g.generateExpression(stmt.Condition))
	g.indent++
	for _, s := range stmt.Body.Statements {
		g.generateStatement(s)
	}
	g.indent--
	g.writeLine("}")
}

func (g *Generator) generateInfixExpression(expr *parser.InfixExpression) string {
	left := g.generateExpression(expr.Left)
	right := g.generateExpression(expr.Right)
	return fmt.Sprintf("%s %s %s", left, expr.Operator, right)
}

func (g *Generator) generateCallExpression(expr *parser.CallExpression) string {
	args := make([]string, 0, len(expr.Arguments))
	for _, arg := range expr.Arguments {
		args = append(args, g.generateExpression(arg))
	}
	fn := g.generateExpression(expr.Function)
	return fmt.Sprintf("%s(%s)", fn, strings.Join(args, ", "))
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

func (g *Generator) generateForStatement(stmt *parser.ForStatement) {
	g.writeLine("for (int %s = %s; %s < %s; %s++) {",
		stmt.Iterator,
		g.generateExpression(stmt.Start),
		stmt.Iterator,
		g.generateExpression(stmt.End),
		stmt.Iterator)

	g.indent++
	for _, s := range stmt.Body.Statements {
		g.generateStatement(s)
	}
	g.indent--
	g.writeLine("}")
}

func (g *Generator) write(format string, args ...interface{}) {
	g.output.WriteString(strings.Repeat("    ", g.indent))
	fmt.Fprintf(&g.output, format, args...)
}

func (g *Generator) writeLine(format string, args ...interface{}) {
	g.write(format+"\n", args...)
}

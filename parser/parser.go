package parser

import (
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/charmbracelet/log"
	"github.com/davecgh/go-spew/spew"
)

func Parse(input string, fileName string) *Root {
	participle.Trace(os.Stdout)

	rules := []lexer.SimpleRule{
		{Name: "Comment", Pattern: `//[^\n]*`},
		{Name: "Whitespace", Pattern: `\s+`},
		{Name: "Replace", Pattern: `replace [a-zA-Z_]\w* -> [a-zA-Z_]\w*`},
		{Name: "String", Pattern: `"[^"]*"`},
		{Name: "Int", Pattern: `(-)?[0-9]+`},
		{Name: "Float", Pattern: `(-)?[0-9]+(\.[0-9]+)?`},
		{Name: "Bool", Pattern: `true|false`},
		{Name: "Ident", Pattern: `[a-zA-Z_]\w*`},
		{Name: "Arrow", Pattern: `->`},
		{Name: "Fun", Pattern: `fun`},
		{Name: "Var", Pattern: `var`},
		{Name: "Define", Pattern: `define`},
		{Name: "AssignOp", Pattern: `\+=|-=|\*=|/=|%=|&=|\|=|\^=|<<=|>>=|=`},
		{Name: "IncDec", Pattern: `\+\+|--`},
		{Name: "ShiftOp", Pattern: `<<|>>`},
		{Name: "CompareOp", Pattern: `<=|>=|==|!=|<|>`},
		{Name: "LogicalOp", Pattern: `&&|\|\|`},
		{Name: "BitwiseOp", Pattern: `[&|^]`},
		{Name: "MathOp", Pattern: `[+\-*/%]`},
		{Name: "UnaryOp", Pattern: `[!~]`},
		{Name: "Punct", Pattern: `[(),{}.?:]`},
	}

	lex := lexer.MustSimple(rules)

	parser := participle.MustBuild[Root](
		participle.Lexer(lex),
		participle.Elide("Comment", "Whitespace", "Replace"),
		participle.UseLookahead(1024),
		participle.Union[Node](VarDecl{}, Definedecl{}, FunctionDef{}),
	)

	program, err := parser.ParseString(fileName, input)
	if err != nil {
		spew.Dump(parser)
		log.Fatal(err)
	}

	return program
}

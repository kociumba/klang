package parser

import (
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// AST structures
type Program struct {
	Declarations []Declaration `@@*`
}

type Declaration struct {
	Replace *ReplaceStmt `  @@`
	Define  *FunctionDef `| "define" @@` // Add explicit "define" match
	Fun     *FunctionDef `| "fun" @@`    // Add explicit "fun" match
}

type ReplaceStmt struct {
	From string `"replace" @Ident`
	To   string `"->" @(String|Ident)`
}

type FunctionDef struct {
	Name       string      `@Ident`
	Params     []Parameter `"(" (@@ ("," @@)*)? ")"`
	ReturnType *string     `("->" (@Ident | "void")?)?`
	Body       *Block      `@@`
}

type Parameter struct {
	Name string `@Ident`
	Type string `":" @Ident`
}

type Block struct {
	Statements []Statement `"{" @@* "}"`
}

type Statement struct {
	VarDecl *VarDecl     `  @@`
	ForLoop *ForLoop     `| @@`
	IfStmt  *IfStatement `| @@`
	Return  *ReturnStmt  `| @@`
	Expr    *Expression  `| @@`
}

type VarDecl struct {
	Name string      `"var" @Ident`
	Type string      `":" @Ident`
	Init *Expression `("=" @@)?`
}

type ForLoop struct {
	Iterator string      `("for"|"gabagool"|"ل") @Ident`
	In       string      `"in"`
	Start    *Expression `"range" "(" @@ ","`
	End      *Expression `@@ ")"`
	Body     *Block      `@@`
}

type IfStatement struct {
	Condition *Expression `"if" @@`
	Body      *Block      `@@`
}

type ReturnStmt struct {
	Value *Expression `"return" @@`
}

type Expression struct {
	Atom  *Atom       `@@`
	Op    *string     `@("+" | "-" | "*" | "/" | "%" | "==" | "!=" | "<" | ">" | "<=" | ">=" | "&&" | "||" | "+=" | "-=" | "*=" | "/=" | "%=" | "&=" | "|=" | "^=" | "<<=" | ">>="`
	Right *Expression `@@)?`
}

// New type to handle atomic expressions
type Atom struct {
	Number   *int64        `  @Int`
	String   *string       `| @String`
	FuncCall *FunctionCall `| @@`
	Ident    *string       `| @Ident`
	SubExpr  *Expression   `| "(" @@ ")"`
}

type FunctionCall struct {
	Name      string        `@Ident`
	Arguments []*Expression `"(" (@@ ("," @@)*)? ")"`
}

func Parse(input *os.File) *Program {
	lex := lexer.MustSimple([]lexer.SimpleRule{
		{Name: "Comment", Pattern: `//[^\n]*`},
		{Name: "Whitespace", Pattern: `\s+`},
		{Name: "String", Pattern: `"[^"]*"`},
		{Name: "Int", Pattern: `[0-9]+`},
		{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
		{Name: "Arrow", Pattern: `->`},
		{Name: "Operator", Pattern: `\+=|-=|\*=|/=|%=|&=|\|=|\^=|<<=|>>=|&&|\|\||==|!=|<=|>=|[+\-*/%&|^<>]=?`},
		{Name: "Punct", Pattern: `[(),{}:;]`},
	})

	parser := participle.MustBuild[Program](
		participle.Lexer(lex),
		participle.Elide("Comment", "Whitespace"),
		participle.UseLookahead(2),
	)

	r, err := parser.Parse("test/grammar.k", input)
	if err != nil {
		panic(err)
	}

	return r
}

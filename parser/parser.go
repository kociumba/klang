package parser

import (
	"bufio"
	"os"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/charmbracelet/log"
)

// AST structures
type Program struct {
	Declarations []Declaration `@@*`
}

type Declaration struct {
	Replace *ReplaceStmt `  @@`
	Define  *FunctionDef `| "define" @@`
	Fun     *FunctionDef `| "fun" @@`
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
	Block   *Block       `| @@`
	Expr    *Expression  `| @@`
}
type VarDecl struct {
	Name string      `("var" @Ident | @Ident ":=")` // Added support for :=
	Type *string     `(":" @Ident)?`                // Made Type optional for := syntax
	Init *Expression `(("=" | ":=") @@)?`           // Support both = and :=
}

type ForLoop struct {
	Keyword  string      `(@Ident | "for")` // Captures "for" or its replacement
	Iterator string      `@Ident`           // Captures just the iterator variable name
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
	Value *Expression `"return" @@?` // Made Value optional for void returns
}

type Expression struct {
	Atom  *Atom       `@@`
	Op    *string     `@("+" | "-" | "*" | "/" | "%" | "==" | "!=" | "<" | ">" | "<=" | ">=" | "&&" | "||" | "+=" | "-=" | "*=" | "/=" | "%=" | "&=" | "|=" | "^=" | "<<=" | ">>=")? `
	Right *Expression `@@?`
}

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

type Replacement struct {
	From string
	To   string
}

func CollectReplacements(input *os.File) ([]Replacement, error) {
	var replacements []Replacement
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "replace") {
			parts := strings.Fields(line)
			if len(parts) >= 4 && parts[2] == "->" {
				from := parts[1]
				to := strings.Trim(parts[3], "\"")
				replacements = append(replacements, Replacement{From: from, To: to})
			}
		}
	}

	input.Seek(0, 0)
	return replacements, scanner.Err()
}

func Parse(input *os.File) *Program {
	replacements, err := CollectReplacements(input)
	if err != nil {
		log.Fatal(err)
	}

	replacementPatterns := make([]string, len(replacements))
	for i, r := range replacements {
		replacementPatterns[i] = r.From
	}

	rules := []lexer.SimpleRule{
		{Name: "Comment", Pattern: `//[^\n]*`},
		{Name: "Whitespace", Pattern: `\s+`},
		{Name: "String", Pattern: `"[^"]*"`},
		{Name: "Int", Pattern: `[0-9]+`},
		{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
		{Name: "Arrow", Pattern: `->`},
		{Name: "Operator", Pattern: `\+=|-=|\*=|/=|%=|&=|\|=|\^=|<<=|>>=|&&|\|\||==|!=|<=|>=|[+\-*/%&|^<>]=?`},
		{Name: "Punct", Pattern: `[(),{}:;]`},
	}

	lex := lexer.MustSimple(rules)

	parser := participle.MustBuild[Program](
		participle.Lexer(lex),
		participle.Elide("Comment", "Whitespace"),
		participle.UseLookahead(2),
	)

	input.Seek(0, 0)
	program, err := parser.Parse(input.Name(), input)
	if err != nil {
		log.Fatal(err)
	}

	return program
}

package parser

type Root struct {
	Nodes []Node `@@*`
}

type Node interface{ node() }

type VarDecl struct {
	Keyword    string      `@"var"`
	Name       string      `@Ident`
	Type       string      `(":" @Ident)?`
	Assignment string      `"="`
	Init       *Expression `@@`
}

func (VarDecl) node() {}

// defines must have to have a return type
type Definedecl struct {
	Keyword    string  `@"define"`
	Name       string  `@Ident`
	Params     []Arg   `"(" (@@ ("," @@)*)? ")"`
	ReturnType *string `("->" @Ident)`
	Body       *Block  `@@`
}

func (Definedecl) node() {}

type FunctionDef struct {
	Keyword    string  `@"fun"`
	Name       string  `@Ident`
	Params     []Arg   `"(" (@@ ("," @@)*)? ")"`
	ReturnType *string `("->" (@Ident | "void")?)`
	Body       *Block  `@@`
}

func (FunctionDef) node() {}

type Arg struct {
	Name string `@Ident`
	Type string `":" @Ident`
}

type FuncCall struct {
	Name string        `@Ident`
	Args []*Expression `"(" ( @@ ( "," @@ )* )? ")"`
}

type Block struct {
	Statements []Statement `"{" @@* "}"`
}

type Statement struct {
	VarDecl    *VarDecl     `  @@`
	Assignment *Assignment  `| @@`
	Return     *ReturnStmt  `| "return" @@`
	ForLoop    *ForLoop     `| "for" @@`
	IfStmt     *IfStatement `| "if" @@`
	FuncCall   *FuncCall    `| @@`
	Expr       *Expression  `| @@`
}

type Assignment struct {
	Name  string      `@Ident`
	Op    string      `@("=" | "+=" | "-=" | "*=" | "/=" | "%=")`
	Value *Expression `@@`
}

type ReturnStmt struct {
	Value *Expression `@@`
}

type ForLoop struct {
	Iterator string      `@Ident`
	In       string      `"in"`
	Start    *Expression `"range" "(" @@ ","`
	End      *Expression `@@ ")"`
	Body     *Block      `@@`
}

type IfStatement struct {
	Condition *Expression `@@`
	Body      *Block      `@@`
	Else      *ElseBranch `( "else" @@ )?`
}

type ElseBranch struct {
	ElseIf *IfStatement `@@`
	Else   *Block       `| @@`
}

type Expression struct {
	Equality *Equality `@@`
}

type Equality struct {
	Comparison *Comparison `@@`
	Op         string      `( @( "!" "=" | "=" "=" | "<" "=" | ">" "=" )`
	Next       *Equality   `  @@ )*`
}

type Comparison struct {
	Addition *Addition   `@@`
	Op       string      `( @( ">" | ">=" | "<" | "<=" )`
	Next     *Comparison `  @@ )*`
}

type Addition struct {
	Multiplication *Multiplication `@@`
	Op             string          `( @( "-" | "+" )`
	Next           *Addition       `  @@ )*`
}

type Multiplication struct {
	Unary *Unary          `@@`
	Op    string          `( @( "/" | "*" | "%" )`
	Next  *Multiplication `  @@ )*`
}

type Unary struct {
	Op      string   `  ( @( "!" | "-" )`
	Unary   *Unary   `    @@ )`
	Primary *Primary `| @@`
}

type Primary struct {
	Number        *float64    `  @Float | @Int`
	String        *string     `| @String`
	Bool          *Boolean    `| @( "true" | "false" )`
	Nil           bool        `| @"nil"`
	FuncCall      *FuncCall   `| @@`
	Ident         *string     `| @Ident`
	SubExpression *Expression `| "(" @@ ")" `
}

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "true"
	return nil
}

package parser

type Root struct {
	Nodes []Node `@@*`
}

type Node interface{ node() }

type VarDecl struct {
	Keyword    string      `@"var"`
	Name       string      `@Ident`
	Type       *Type       `(":" @@)?`
	Assignment string      `"="`
	Init       *Expression `@@`
}

func (VarDecl) node() {}

// defines must have to have a return type
type Definedecl struct {
	Keyword    string `@"define"`
	Name       string `@Ident`
	Params     []Arg  `"(" (@@ ("," @@)*)? ")"`
	ReturnType *Type  `("->" @@)`
	Body       *Block `@@`
}

func (Definedecl) node() {}

type FunctionDef struct {
	Keyword    string `@"fun"`
	Name       string `@Ident`
	Params     []Arg  `"(" (@@ ("," @@)*)? ")"`
	ReturnType *Type  `("->" @@)?`
	Body       *Block `@@`
}

func (FunctionDef) node() {}

// Struct declaration
type StructDecl struct {
	Keyword string        `@"struct"`
	Name    string        `@Ident`
	Fields  []StructField `"{" @@* "}"` // Fields separated by semicolons
}

func (StructDecl) node() {}

// Type alias (e.g., "type MyInt = int")
type TypeDecl struct {
	Keyword string `@"type"`
	Name    string `@Ident`
	Type    Type   `"=" @@`
}

func (TypeDecl) node() {}

// Struct field
type StructField struct {
	Name string `@Ident`
	Type Type   `":" @@`
}

type Arg struct {
	Name string `@Ident`
	Type *Type  `":" @@`
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
	WhileLoop  *WhileLoop   `| "while" @@`
	IfStmt     *IfStatement `| "if" @@`
	FuncCall   *FuncCall    `| @@`
	Expr       *Expression  `| @@`
}

type Assignment struct {
	Target *LValue     `@@`
	Op     string      `@("=" | "+=" | "-=" | "*=" | "/=" | "%=")`
	Value  *Expression `@@`
}

type LValue struct {
	Ident       *string      `  @Ident`
	IndexAccess *IndexAccess `| @@`
	Dereference *Dereference `| @@`
}

type Dereference struct {
	Modifiers []Modifier `@@"*"+`
	Base      *Primary   `@@`
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

type WhileLoop struct {
	Condition *Expression `"while" @@`
	Body      *Block      `@@`
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

type Type struct {
	Base      string     `@Ident`
	Modifiers []Modifier `@@*` // Track modifiers in parsed order
}

type Modifier struct {
	Pointer *PointerMod `  "*"@@`
	Array   *ArrayDim   `| @@`
}

type PointerMod struct {
	Level string `@("*")*`
} // Exists just to capture "*" tokens

type ArrayDim struct {
	Size *Expression `"[" @@? "]"`
}

type ArrayLiteral struct {
	Elements []*Expression `"{" ( @@ ( "," @@ )* )? "}"`
}

type IndexAccess struct {
	Index *Expression `"[" @@ "]"`
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
	Base      *PrimaryBase `@@`
	Postfixes []*Postfix   `@@*`
}

type PrimaryBase struct {
	Number        *float64      `  @Float | @Int`
	String        *string       `| @String`
	Bool          *Boolean      `| @( "true" | "false" )`
	Nil           bool          `| @"nil"`
	FuncCall      *FuncCall     `| @@`
	Ident         *string       `| @Ident`
	SubExpression *Expression   `| "(" @@ ")"`
	ArrayLiteral  *ArrayLiteral `| @@`
}

type Postfix struct {
	IndexAccess *IndexAccess `@@` // Index access is a postfix operation
}

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "true"
	return nil
}

package parser

import (
	"fmt"
	"strconv"

	"github.com/kociumba/klang/lexer"
)

// AST Nodes
type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

type FunctionLiteral struct {
	Token      lexer.Token
	Name       string
	Parameters []*Parameter
	ReturnType string
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) statementNode()       {}

type Parameter struct {
	Name string
	Type string
}

type BlockStatement struct {
	Token      lexer.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

type VarStatement struct {
	Token lexer.Token
	Name  string
	Type  string
	Value Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Token.Literal }

type ReturnStatement struct {
	Token       lexer.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

type ForStatement struct {
	Token    lexer.Token
	Iterator string
	Start    Expression
	End      Expression
	Body     *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }

type Identifier struct {
	Token lexer.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type IntegerLiteral struct {
	Token lexer.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

type ExpressionStatement struct {
	Token      lexer.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

type CallExpression struct {
	Token     lexer.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

// Parser
type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	// Read two tokens to initialize curToken and peekToken
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{Statements: []Statement{}}

	for p.curToken.Type != lexer.TOKEN_EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case lexer.TOKEN_FUN:
		return p.parseFunctionStatement()
	case lexer.TOKEN_VAR:
		return p.parseVarStatement()
	case lexer.TOKEN_RETURN:
		return p.parseReturnStatement()
	case lexer.TOKEN_FOR:
		return p.parseForStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseCallExpression(function Expression) *CallExpression {
	exp := &CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []Expression {
	args := []Expression{}

	if p.peekToken.Type == lexer.TOKEN_RPAREN {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression())

	for p.peekToken.Type == lexer.TOKEN_COMMA {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression())
	}

	if !p.expectPeek(lexer.TOKEN_RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression()
	return stmt
}

func (p *Parser) parseFunctionStatement() *FunctionLiteral {
	fun := &FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(lexer.TOKEN_IDENT) {
		return nil
	}
	fun.Name = p.curToken.Literal

	if !p.expectPeek(lexer.TOKEN_LPAREN) {
		return nil
	}

	fun.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(lexer.TOKEN_ARROW) {
		return nil
	}

	if !p.expectPeek(lexer.TOKEN_IDENT) {
		return nil
	}
	fun.ReturnType = p.curToken.Literal

	if !p.expectPeek(lexer.TOKEN_LBRACE) {
		return nil
	}

	fun.Body = p.parseBlockStatement()

	return fun
}

func (p *Parser) parseFunctionParameters() []*Parameter {
	params := []*Parameter{}

	if p.peekToken.Type == lexer.TOKEN_RPAREN {
		p.nextToken()
		return params
	}

	p.nextToken()

	param := &Parameter{
		Name: p.curToken.Literal,
	}

	if !p.expectPeek(lexer.TOKEN_COLON) {
		return nil
	}

	if !p.expectPeek(lexer.TOKEN_IDENT) {
		return nil
	}
	param.Type = p.curToken.Literal

	params = append(params, param)

	for p.peekToken.Type == lexer.TOKEN_COMMA {
		p.nextToken()
		p.nextToken()

		param := &Parameter{
			Name: p.curToken.Literal,
		}

		if !p.expectPeek(lexer.TOKEN_COLON) {
			return nil
		}

		if !p.expectPeek(lexer.TOKEN_IDENT) {
			return nil
		}
		param.Type = p.curToken.Literal

		params = append(params, param)
	}

	if !p.expectPeek(lexer.TOKEN_RPAREN) {
		return nil
	}

	return params
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{Token: p.curToken}
	block.Statements = []Statement{}

	p.nextToken()

	for p.curToken.Type != lexer.TOKEN_RBRACE && p.curToken.Type != lexer.TOKEN_EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseVarStatement() *VarStatement {
	stmt := &VarStatement{Token: p.curToken}

	if !p.expectPeek(lexer.TOKEN_IDENT) {
		return nil
	}

	stmt.Name = p.curToken.Literal

	if !p.expectPeek(lexer.TOKEN_COLON) {
		return nil
	}

	if !p.expectPeek(lexer.TOKEN_IDENT) {
		return nil
	}

	stmt.Type = p.curToken.Literal

	if p.peekToken.Type == lexer.TOKEN_ASSIGN {
		p.nextToken()
		p.nextToken()
		stmt.Value = p.parseExpression()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression()

	return stmt
}

func (p *Parser) parseForStatement() *ForStatement {
	stmt := &ForStatement{Token: p.curToken}

	// Parse: for <iterator> in range(<start>, <end>)
	if !p.expectPeek(lexer.TOKEN_IDENT) {
		return nil
	}
	stmt.Iterator = p.curToken.Literal

	if !p.expectPeek(lexer.TOKEN_IN) {
		return nil
	}

	if !p.expectPeek(lexer.TOKEN_RANGE) {
		return nil
	}

	if !p.expectPeek(lexer.TOKEN_LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Start = p.parseExpression()

	if !p.expectPeek(lexer.TOKEN_COMMA) {
		return nil
	}

	p.nextToken()
	stmt.End = p.parseExpression()

	if !p.expectPeek(lexer.TOKEN_RPAREN) {
		return nil
	}

	if !p.expectPeek(lexer.TOKEN_LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseExpression() Expression {
	switch p.curToken.Type {
	case lexer.TOKEN_IDENT:
		return &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	case lexer.TOKEN_INT:
		value, _ := strconv.ParseInt(p.curToken.Literal, 10, 64)
		return &IntegerLiteral{Token: p.curToken, Value: value}
	default:
		return nil
	}
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	p.errors = append(p.errors, fmt.Sprintf("expected next token to be %v, got %v instead", t, p.peekToken.Type))
	return false
}

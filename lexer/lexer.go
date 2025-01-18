package lexer

import (
	"bufio"
	"fmt"
	"strings"
)

//go:generate stringer -type=TokenType
type TokenType int

const (
	TOKEN_INVALID TokenType = iota
	TOKEN_EOF

	// Keywords
	TOKEN_FUN
	TOKEN_VAR
	TOKEN_CONST
	TOKEN_RETURN
	TOKEN_IF
	TOKEN_ELSE
	TOKEN_IN
	TOKEN_RANGE
	TOKEN_REPLACE
	TOKEN_DEFINE
	TOKEN_BREAK
	TOKEN_CONTINUE
	TOKEN_WHILE
	TOKEN_FOR

	// Identifiers and literals
	TOKEN_IDENT  // variables, functions
	TOKEN_INT    // 123
	TOKEN_STRING // "hello"
	TOKEN_FLOAT  // 123.45

	// Operators
	TOKEN_ASSIGN  // =
	TOKEN_PLUS    // +
	TOKEN_MINUS   // -
	TOKEN_ARROW   // ->
	TOKEN_COLON   // :
	TOKEN_SEMI    // ;
	TOKEN_PLUSEQ  // +=
	TOKEN_MINUSEQ // -=
	TOKEN_MULEQ   // *=
	TOKEN_DIVEQ   // /=
	TOKEN_MODEQ   // %=
	TOKEN_ANDEQ   // &=
	TOKEN_OREQ    // |=
	TOKEN_XOREQ   // ^=
	TOKEN_SHLEQ   // <<=
	TOKEN_SHREQ   // >>=

	// Arithmetic
	TOKEN_ASTERISK // *
	TOKEN_SLASH    // /
	TOKEN_PERCENT  // %

	// Bitwise
	TOKEN_AND // &
	TOKEN_OR  // |
	TOKEN_XOR // ^
	TOKEN_SHL // <<
	TOKEN_SHR // >>
	TOKEN_NOT // ~

	// Logical
	TOKEN_LAND // &&
	TOKEN_LOR  // ||
	TOKEN_LNOT // !
	TOKEN_EQ   // ==
	TOKEN_NEQ  // !=
	TOKEN_LT   // <
	TOKEN_GT   // >
	TOKEN_LTE  // <=
	TOKEN_GTE  // >=

	// Delimiters
	TOKEN_COMMA    // ,
	TOKEN_DOT      // .
	TOKEN_LPAREN   // (
	TOKEN_RPAREN   // )
	TOKEN_LBRACE   // {
	TOKEN_RBRACE   // }
	TOKEN_LBRACKET // [
	TOKEN_RBRACKET // ]
)

type Token struct {
	Type         TokenType
	Literal      string
	Line         int
	Column       int
	ErrorMsg     string
	OriginalText string
}

type Lexer struct {
	input        string
	position     int  // current position in input
	readPosition int  // next position to read
	ch           byte // current character
	line         int
	column       int
	keywords     map[string]TokenType
	replacements map[string]string // stores keyword replacements
	unicodeMap   map[string]string // map for Unicode -> C identifier conversion
	nextUniqueId int
}

func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 1,
		keywords: map[string]TokenType{
			"fun":      TOKEN_FUN,
			"var":      TOKEN_VAR,
			"const":    TOKEN_CONST,
			"return":   TOKEN_RETURN,
			"if":       TOKEN_IF,
			"else":     TOKEN_ELSE,
			"in":       TOKEN_IN,
			"range":    TOKEN_RANGE,
			"replace":  TOKEN_REPLACE,
			"break":    TOKEN_BREAK,
			"continue": TOKEN_CONTINUE,
			"while":    TOKEN_WHILE,
			"for":      TOKEN_FOR,
			"define":   TOKEN_DEFINE,
		},
		replacements: make(map[string]string),
		unicodeMap:   make(map[string]string),
	}
	l.readChar()
	return l
}

// Preprocessor scans for replacement directives
func (l *Lexer) Preprocess() error {
	scanner := bufio.NewScanner(strings.NewReader(l.input))
	var processedLines []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "//") {
			processedLines = append(processedLines, line)
			continue
		}

		// Handle replace directives
		if strings.HasPrefix(line, "replace") {
			parts := strings.Fields(line)
			if len(parts) >= 4 && parts[2] == "->" {
				original := parts[1]
				replacement := parts[3]
				l.AddReplacement(original, replacement)
			} else {
				return fmt.Errorf("invalid replace directive: %s", line)
			}
		}

		processedLines = append(processedLines, line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

	if l.ch == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	// Skip comments
	if l.ch == '/' && l.peekChar() == '/' {
		l.skipLineComment()
		return l.NextToken()
	}

	// Store current position info
	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case '"':
		tok.Type = TOKEN_STRING
		tok.Literal = l.readString()
	case '=':
		tok = l.handleEquals()
	case '+':
		tok = l.handlePlus()
	case '-':
		tok = l.handleMinus()
	case '*':
		tok = l.handleAsterisk()
	case '/':
		tok = l.handleSlash()
	case '%':
		tok = l.handlePercent()
	case '&':
		tok = l.handleAnd()
	case '|':
		tok = l.handleOr()
	case '^':
		tok = l.handleXor()
	case '<':
		tok = l.handleLessThan()
	case '>':
		tok = l.handleGreaterThan()
	case '!':
		tok = l.handleBang()
	case '~':
		tok = Token{Type: TOKEN_NOT, Literal: string(l.ch)}
	case ':':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TOKEN_ASSIGN, Literal: ":="}
		} else {
			tok = Token{Type: TOKEN_COLON, Literal: string(l.ch)}
		}
	case ';':
		tok = Token{Type: TOKEN_SEMI, Literal: string(l.ch)}
	case '.':
		tok = Token{Type: TOKEN_DOT, Literal: string(l.ch)}
	case '(':
		tok = Token{Type: TOKEN_LPAREN, Literal: string(l.ch)}
	case ')':
		tok = Token{Type: TOKEN_RPAREN, Literal: string(l.ch)}
	case '{':
		tok = Token{Type: TOKEN_LBRACE, Literal: string(l.ch)}
	case '}':
		tok = Token{Type: TOKEN_RBRACE, Literal: string(l.ch)}
	case ',':
		tok = Token{Type: TOKEN_COMMA, Literal: string(l.ch)}
	case 0:
		tok.Literal = ""
		tok.Type = TOKEN_EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			// Store original text before any replacements
			tok.OriginalText = tok.Literal

			// Check replacements
			if replacement, ok := l.replacements[tok.Literal]; ok {
				tok.Literal = replacement
				// If this is a Unicode replacement, store the original
				if original, ok := l.unicodeMap[replacement]; ok {
					tok.OriginalText = original
				}
			}

			// Check keywords
			if tokenType, ok := l.keywords[tok.Literal]; ok {
				tok.Type = tokenType
			} else {
				tok.Type = TOKEN_IDENT
			}
			return tok
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readString() string {
	position := l.position + 1 // Skip the opening quote
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n', 't', 'r', '\\', '"', '0',
				'a', 'b', 'f', 'v', '\'':
				// Valid C escape sequences
				continue
			case 'x':
				// Hex escape sequence \xHH
				l.readChar() // First hex digit
				l.readChar() // Second hex digit
			default:
				// Invalid escape sequence for C
			}
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipLineComment() Token {
	startLine := l.line
	startColumn := l.column
	comment := "//"

	for l.ch != '\n' && l.ch != 0 {
		comment += string(l.ch)
		l.readChar()
	}

	return Token{
		Type:    TOKEN_INVALID, // or TOKEN_COMMENT if you want to preserve comments
		Literal: comment,
		Line:    startLine,
		Column:  startColumn,
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() (string, TokenType) {
	startPos := l.position
	// sawDot := false
	isFloat := false

	// Read the integer part
	for isDigit(l.ch) {
		l.readChar()
	}

	// Check for decimal point
	if l.ch == '.' {
		if isDigit(l.peekChar()) {
			// sawDot = true
			l.readChar() // consume the dot

			// Read the decimal part
			for isDigit(l.ch) {
				l.readChar()
			}
			isFloat = true
		}
	}

	// Handle scientific notation (e.g., 1e-10, 1.5e+5)
	if l.ch == 'e' || l.ch == 'E' {
		peek := l.peekChar()
		if isDigit(peek) || peek == '+' || peek == '-' {
			l.readChar() // consume 'e' or 'E'
			if l.ch == '+' || l.ch == '-' {
				l.readChar() // consume '+' or '-'
			}

			// Must have at least one digit after 'e'/'E'
			if !isDigit(l.ch) {
				// Invalid scientific notation
				return l.input[startPos:l.position], TOKEN_INVALID
			}

			// Read the exponent
			for isDigit(l.ch) {
				l.readChar()
			}
			isFloat = true
		}
	}

	if isFloat {
		return l.input[startPos:l.position], TOKEN_FLOAT
	}
	return l.input[startPos:l.position], TOKEN_INT
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// AddReplacement adds a keyword replacement
func (l *Lexer) AddReplacement(original, replacement string) {
	if isASCII(replacement) {
		l.replacements[original] = replacement
	} else {
		safeIdent := l.generateSafeIdentifier(replacement)
		l.replacements[original] = safeIdent
		l.unicodeMap[safeIdent] = replacement
	}
}

func (l *Lexer) generateSafeIdentifier(unicode string) string {
	l.nextUniqueId++
	return fmt.Sprintf("unicode_id_%d", l.nextUniqueId)
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > 127 {
			return false
		}
	}
	return true
}

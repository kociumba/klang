package lexer

// Handler functions for multi-character operators
func (l *Lexer) handleEquals() Token {
	if l.peekChar() == '=' {
		l.readChar()
		return Token{Type: TOKEN_EQ, Literal: "=="}
	}
	return Token{Type: TOKEN_ASSIGN, Literal: "="}
}

func (l *Lexer) handlePlus() Token {
	if l.peekChar() == '=' {
		l.readChar()
		return Token{Type: TOKEN_PLUSEQ, Literal: "+="}
	}
	return Token{Type: TOKEN_PLUS, Literal: "+"}
}

// needs to also handle ->
func (l *Lexer) handleMinus() Token {
	if l.peekChar() == '>' {
		l.readChar()
		return Token{Type: TOKEN_ARROW, Literal: "->"}
	} else if l.peekChar() == '=' {
		l.readChar()
		return Token{Type: TOKEN_MINUSEQ, Literal: "-="}
	}
	return Token{Type: TOKEN_MINUS, Literal: "-"}
}

func (l *Lexer) handleAsterisk() Token {
	if l.peekChar() == '=' {
		l.readChar()
		return Token{Type: TOKEN_MULEQ, Literal: "*="}
	}
	return Token{Type: TOKEN_ASTERISK, Literal: "*"}
}

func (l *Lexer) handleSlash() Token {
	if l.peekChar() == '=' {
		l.readChar()
		return Token{Type: TOKEN_DIVEQ, Literal: "/="}
	}
	return Token{Type: TOKEN_SLASH, Literal: "/"}
}

func (l *Lexer) handlePercent() Token {
	if l.peekChar() == '=' {
		l.readChar()
		return Token{Type: TOKEN_MODEQ, Literal: "%="}
	}
	return Token{Type: TOKEN_PERCENT, Literal: "%"}
}

func (l *Lexer) handleAnd() Token {
	if l.peekChar() == '=' {
		l.readChar()
		return Token{Type: TOKEN_ANDEQ, Literal: "&="}
	} else if l.peekChar() == '&' {
		l.readChar()
		return Token{Type: TOKEN_LAND, Literal: "&&"}
	}
	return Token{Type: TOKEN_AND, Literal: "&"}
}

func (l *Lexer) handleOr() Token {
	if l.peekChar() == '=' {
		l.readChar()
		return Token{Type: TOKEN_OREQ, Literal: "|="}
	} else if l.peekChar() == '|' {
		l.readChar()
		return Token{Type: TOKEN_LOR, Literal: "||"}
	}
	return Token{Type: TOKEN_OR, Literal: "|"}
}

func (l *Lexer) handleXor() Token {
	if l.peekChar() == '=' {
		l.readChar()
		return Token{Type: TOKEN_XOREQ, Literal: "^="}
	}
	return Token{Type: TOKEN_XOR, Literal: "^"}
}

func (l *Lexer) handleLessThan() Token {
	if l.peekChar() == '=' {
		l.readChar()
		return Token{Type: TOKEN_LTE, Literal: "<="}
	} else if l.peekChar() == '<' {
		l.readChar()
		if l.peekChar() == '=' {
			l.readChar()
			return Token{Type: TOKEN_SHLEQ, Literal: "<<="}
		}
		return Token{Type: TOKEN_SHL, Literal: "<<"}
	}
	return Token{Type: TOKEN_LT, Literal: "<"}
}

func (l *Lexer) handleGreaterThan() Token {
	if l.peekChar() == '=' {
		l.readChar()
		return Token{Type: TOKEN_GTE, Literal: ">="}
	} else if l.peekChar() == '>' {
		l.readChar()
		if l.peekChar() == '=' {
			l.readChar()
			return Token{Type: TOKEN_SHREQ, Literal: ">>="}
		}
		return Token{Type: TOKEN_SHR, Literal: ">>"}
	}
	return Token{Type: TOKEN_GT, Literal: ">"}
}

func (l *Lexer) handleBang() Token {
	if l.peekChar() == '=' {
		l.readChar()
		return Token{Type: TOKEN_NEQ, Literal: "!="}
	}
	return Token{Type: TOKEN_LNOT, Literal: "!"}
}

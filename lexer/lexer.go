package lexer

import "interpreter/token"

type Lexer struct {
	// The position is used when we want to check identifiers or numbers
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	// Create a new Lexer instance with the input string
	// and initialize the position and readPosition to 0
	// For example, if the input is "let x = 5;", set the position and readPosition to 0
	// and read the first character
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	// Read the next character from the input string
	// and update the position and readPosition
	// For example, if the input is "let x = 5;", read the characters one by one
	// and update the position and readPosition accordingly
	if l.readPosition >= len(l.input) {
		// EOF (end of file) reached
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	// Peek the next character without advancing the read position
	// For example, if the input is "let x = 5;", peek the next character after reading "let"
	// and return the character ' ' (space) without advancing the read position
	if l.readPosition >= len(l.input) {
		// EOF (end of file) reached
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() token.Token {
	// Read the next token from the input string
	// For example, if the input is "let x = 5;", return the tokens for "let", "x", "=", "5", ";"
	// In this case, the tokens would be: LET, IDENT, ASSIGN, INT, SEMICOLON
	// The tokens are created using the newToken function
	// and the token type is determined based on the character read
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		// Check if the next character is also '=' for the EQ token
		if l.peekChar() == '=' {
			// Keep the current character to a variable, so we can use the character even after advancing to the next character
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		// Check if the next character is also '=' for the NOT_EQ token
		if l.peekChar() == '=' {
			// Keep the current character to a variable, so we can use the character even after advancing to the next character
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			// Read the identifier
			tok.Literal = l.readIdentifier()

			// Check if the identifier is a keyword
			// and set the token type accordingly
			// For example, if the identifier is "fn", set the token type to FUNCTION
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			// Read the number and set the token type to INT
			// For example, if the number is "123", set the token type to INT
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	// Read the next character to advance the position
	// For example, if the input is "let x = 5;", after reading "let", advance to the next character ' '
	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	// Create a new token with the given type and literal value
	// For example, if the token type is ASSIGN and the character is '=', create a token with type ASSIGN and literal '='
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	// Read the identifier from the input string
	// For example, if the identifier is "let", return "let"
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	// Read the number from the input string
	// For example, if the number is "123", return "123"
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	// Check if the character is a letter (a-z, A-Z) or an underscore (_)
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	// Check if the character is a digit (0-9)
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skipWhitespace() {
	// Skip whitespace characters (space, tab, newline, carriage return)
	// For example, if the input is "let x = 5;", skip the whitespace before "x"
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

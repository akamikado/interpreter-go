package lexer

import (
	"fmt"

	t "interpreter/token"
)

type Lexer struct {
	input        []byte
	position     int
	readPosition int
	char         byte
}

// Creates new lexer
func NewLexer(input []byte) *Lexer {
	lexer := &Lexer{input: input}
	return lexer
}

// Peeks at next character in input
func (l *Lexer) PeekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// Reads identifier from input
func (l *Lexer) ReadIdentifier() []byte {
	for ; l.readPosition < len(l.input) && isAlphabet(l.input[l.readPosition]); l.readPosition += 1 {
	}
	return l.input[l.position:l.readPosition]
}

// Reads number from input
func (l *Lexer) ReadNumber() []byte {
	for ; l.readPosition < len(l.input) && isDigit(l.input[l.readPosition]); l.readPosition += 1 {
	}
	return l.input[l.position:l.readPosition]
}

// Checks if byte is an alphabet
func isAlphabet(c byte) bool {
	if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') {
		return true
	}
	return false
}

// Checks if byte is a digit
func isDigit(c byte) bool {
	if '0' <= c && c <= '9' {
		return true
	}
	return false
}

// Returns next token from input
func (l *Lexer) GetToken() t.Token {
	var tok t.Token

	if l.readPosition >= len(l.input) {
		tok.Type = t.EOF
		tok.Literal = []byte{0}
		return tok
	}

	for ; l.readPosition < len(l.input) && l.input[l.readPosition] == ' '; l.readPosition += 1 {
	}

	if l.readPosition >= len(l.input) {
		tok.Type = t.EOF
		tok.Literal = []byte{0}
		return tok
	}

	l.position = l.readPosition
	l.char = l.input[l.position]
	l.readPosition += 1

	switch l.char {
	case ';':
		tok.Type = t.SEMICOLON
		tok.Literal = []byte{';'}
	case ',':
		tok.Type = t.COMMA
		tok.Literal = []byte{','}
	case '(':
		tok.Type = t.LPAREN
		tok.Literal = []byte{'('}
	case ')':
		tok.Type = t.RPAREN
		tok.Literal = []byte{')'}
	case '{':
		tok.Type = t.LBRACE
		tok.Literal = []byte{'{'}
	case '}':
		tok.Type = t.RBRACE
		tok.Literal = []byte{'}'}
	case '+':
		tok.Type = t.PLUS
		tok.Literal = []byte{'+'}
	case '-':
		tok.Type = t.MINUS
		tok.Literal = []byte{'-'}
	case '*':
		tok.Type = t.ASTERISK
		tok.Literal = []byte{'*'}
	case '/':
		tok.Type = t.SLASH
		tok.Literal = []byte{'/'}
	case '%':
		tok.Type = t.MODULO
		tok.Literal = []byte{'%'}
	case '!':
		if l.PeekChar() == '=' {
			l.readPosition += 1
			tok.Type = t.NOT_EQ
			tok.Literal = []byte{'!', '='}
		} else {
			tok.Type = t.BANG
			tok.Literal = []byte{'!'}
		}
	case '<':
		tok.Type = t.LT
		tok.Literal = []byte{'<'}
	case '>':
		tok.Type = t.GT
		tok.Literal = []byte{'>'}
	case '=':
		if l.PeekChar() == '=' {
			l.readPosition += 1
			tok.Type = t.EQ
			tok.Literal = []byte{'=', '='}
		} else {
			tok.Type = t.ASSIGN
			tok.Literal = []byte{'='}
		}
	case '\n':
		tok = l.GetToken()
	default:
		if isAlphabet(l.char) {
			tok.Literal = l.ReadIdentifier()
			tok.Type = t.LookupIdentifier(string(tok.Literal))
		} else if isDigit(l.char) {
			tok.Literal = l.ReadNumber()
			tok.Type = t.INT
		} else {
			tok.Type = t.ILLEGAL
		}
	}

	return tok
}

// Returns slice of tokens from input
func Tokenize(input []byte) []t.Token {
	lexer := NewLexer(input)

	var tokens []t.Token

	tokens = append(tokens, lexer.GetToken())

	for tokens[len(tokens)-1].Type != t.EOF {
		tokens = append(tokens, lexer.GetToken())
	}
	for _, token := range tokens {
		fmt.Printf("Token type: %s, Token literal: %s\n", token.Type, string(token.Literal))
	}

	return tokens
}

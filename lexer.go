package main

import (
	"strconv"
)

// Lexer ...
type Lexer struct {
	// client string input, e.g. "3 + 5", "12 - 5", etc
	Text string
	// Pos is an index into Text
	Pos         int
	CurrentChar byte
}

// NewLexer ...
func NewLexer(input string) *Lexer {
	l := &Lexer{}
	l.Text = input
	l.Pos = 0
	l.CurrentChar = l.Text[l.Pos]
	return l
}

func (l *Lexer) Error() {
	panic("error lexing input")
}

// ReservedWords ...
var ReservedWords = map[string]Token{
	"PROGRAM": Token{Type: PROGRAM},
	"VAR":     Token{Type: VAR},
	"DIV":     Token{Type: INTEGERDIV},
	"INTEGER": Token{Type: INTEGER},
	"REAL":    Token{Type: REAL},
	"BEGIN":   Token{Type: BEGIN},
	"END":     Token{Type: END},
}

// ID ...
// Handle identifiers and reserved keywords
func (l *Lexer) ID() Token {
	buffer := make([]byte, 0)
	for l.CurrentChar != 0 && isAlphaNumeric(l.CurrentChar) {
		buffer = append(buffer, l.CurrentChar)
		l.Advance()
	}
	str := string(buffer)
	if tok, exists := ReservedWords[str]; exists {
		return tok
	}
	return Token{Type: IDENT, Svalue: str}
}

// Advance ...
func (l *Lexer) Advance() {
	// Advance the 'pos' pointer and set the 'current_char' variable.
	l.Pos++
	if l.Pos > len(l.Text)-1 {
		l.CurrentChar = 0 // Indicates end of input
	} else {
		l.CurrentChar = l.Text[l.Pos]
	}
}

// Peek ...
func (l *Lexer) Peek() byte {
	pos := l.Pos + 1
	if pos > len(l.Text)-1 {
		return 0
	}
	return l.Text[pos]
}

// SkipWhitespace ...
func (l *Lexer) SkipWhitespace() {
	for (l.CurrentChar == ' ' ||
		l.CurrentChar == '\n' ||
		l.CurrentChar == '\r' ||
		l.CurrentChar == '\t') && l.CurrentChar != 0 {
		l.Advance()
	}
}

// SkipComment ...
func (l *Lexer) SkipComment() {
	for l.CurrentChar != '}' {
		l.Advance()
	}
	l.Advance() // For closing }
}

// Number ...
// Return a (multidigit) integer or float consumed from the input.
func (l *Lexer) Number() Token {
	var buffer []byte
	for isDigit(l.CurrentChar) {
		buffer = append(buffer, l.CurrentChar)
		l.Advance()
	}
	if l.CurrentChar == '.' {
		buffer = append(buffer, l.CurrentChar)
		l.Advance()
		for isDigit(l.CurrentChar) {
			buffer = append(buffer, l.CurrentChar)
			l.Advance()
		}
		val, _ := strconv.ParseFloat(string(buffer), 64)
		return Token{Type: REALCONST, Value: val}
	}
	val, _ := strconv.ParseFloat(string(buffer), 64)
	return Token{Type: INTEGERCONST, Value: val}
}

// GetNextToken ...
// Lexical analyzer (also known as scanner or tokenizer)
//
// This method is responsible for breaking a sentence
// apart into tokens. One token at a time.
func (l *Lexer) GetNextToken() Token {
	for l.CurrentChar != 0 {
		if l.CurrentChar == ' ' ||
			l.CurrentChar == '\n' ||
			l.CurrentChar == '\r' ||
			l.CurrentChar == '\t' {
			l.SkipWhitespace()
			continue
		}
		if isAlpha(l.CurrentChar) {
			return l.ID()
		}
		if isDigit(l.CurrentChar) {
			return l.Number()
		}
		if l.CurrentChar == ':' && l.Peek() == '=' {
			l.Advance()
			l.Advance()
			return Token{Type: ASSIGN}
		}
		switch l.CurrentChar {
		case '{':
			l.Advance()
			l.SkipComment()
			continue
		case ';':
			l.Advance()
			return Token{Type: SEMI}
		case ':':
			l.Advance()
			return Token{Type: COLON}
		case ',':
			l.Advance()
			return Token{Type: COMMA}
		case '+':
			l.Advance()
			return Token{Type: PLUS}
		case '-':
			l.Advance()
			return Token{Type: MINUS}
		case '*':
			l.Advance()
			return Token{Type: MUL}
		case '/':
			l.Advance()
			return Token{Type: FLOATDIV}
		case '(':
			l.Advance()
			return Token{Type: LPAREN}
		case ')':
			l.Advance()
			return Token{Type: RPAREN}
		case '.':
			l.Advance()
			return Token{Type: DOT}
		default:
			l.Error()
		}
	}
	return Token{Type: EOF}
}

func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

func isAlphaNumeric(b byte) bool {
	switch {
	case '0' <= b && b <= '9':
		return true
	case 'a' <= b && b <= 'z':
		return true
	case 'A' <= b && b <= 'Z':
		return true
	}
	return false
}

func isAlpha(b byte) bool {
	switch {
	case 'a' <= b && b <= 'z':
		return true
	case 'A' <= b && b <= 'Z':
		return true
	}
	return false
}

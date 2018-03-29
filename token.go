package main

import "fmt"

// Token types
//
// EOF (end-of-file) token is used to indicate that
// there is no more input left for lexical analysis
const (
	INTEGER = iota + 1
	REAL
	INTEGERCONST
	REALCONST
	PLUS
	MINUS
	MUL
	INTEGERDIV
	FLOATDIV
	LPAREN
	RPAREN
	IDENT
	ASSIGN
	BEGIN
	END
	SEMI
	DOT
	PROGRAM
	VAR
	COLON
	COMMA
	EOF
)

var (
	// TokenStr ...
	TokenStr = []string{
		"",
		"int",
		"real",
		"int const",
		"real const",
		"+",
		"-",
		"*",
		"int /",
		"float /",
		"(",
		")",
		"identifier",
		":=",
		"begin",
		"end",
		";",
		".",
		"program",
		"var",
		":",
		",",
		"eof",
	}

	strmap = []string{
		"",
		"INT",
		"REAL",
		"INT CONST",
		"REAL CONST",
		"PLUS",
		"MINUS",
		"MUL",
		"INT DIV",
		"FLOAT DIV",
		"LPAREN",
		"RPAREN",
		"IDENTIFIER",
		"ASSIGN",
		"BEGIN",
		"END",
		"SEMI",
		"DOT",
		"PROGRAM",
		"VAR",
		"COLON",
		"COMMA",
		"EOF",
	}
)

// Token ...
type Token struct {
	Type   int
	Value  float64
	Svalue string
}

// String representation of the class instance.
// 	Examples:
// 		Token(INTEGER, 3)
// 		Token(PLUS '+')
func (t Token) String() string {
	return fmt.Sprintf("Token(%s, %f)", strmap[t.Type], t.Value)
}

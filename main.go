package main

import "fmt"

// program : PROGRAM variable SEMI block DOT
//
//     block : declarations compound_statement
//
//     declarations : VAR (variable_declaration SEMI)+
//                  | empty
//
//     variable_declaration : ID (COMMA ID)* COLON type_spec
//
//     type_spec : INTEGER
//
//     compound_statement : BEGIN statement_list END
//
//     statement_list : statement
//                    | statement SEMI statement_list
//
//     statement : compound_statement
//               | assignment_statement
//               | empty
//
//     assignment_statement : variable ASSIGN expr
//
//     empty :
//
//     expr : term ((PLUS | MINUS) term)*
//
//     term : factor ((MUL | INTEGER_DIV | FLOAT_DIV) factor)*
//
//     factor : PLUS factor
//            | MINUS factor
//            | INTEGER_CONST
//            | REAL_CONST
//            | LPAREN expr RPAREN
//            | variable
//
// 	variable: ID

var pascalsample1 = `PROGRAM Part10;
VAR
   number     : INTEGER;
   a, b, c, x : INTEGER;
   y          : REAL;

BEGIN {Part10}
   BEGIN
      number := 2;
      a := number;
      b := 10 * a + 10 * number DIV 4;
      c := a - - b
   END;
   x := 11;
   y := 20 / 7 + 3.14;
   { writeln('a = ', a); }
   { writeln('b = ', b); }
   { writeln('c = ', c); }
   { writeln('number = ', number); }
   { writeln('x = ', x); }
   { writeln('y = ', y); }
END.  {Part10}`

var pascalsample2 = `PROGRAM Part10Sample2;
VAR
   a, b : INTEGER;
   y    : REAL;

BEGIN {Part10AST}
   a := 2;
   b := 10 * a + 10 * a DIV 4;
   y := 20 / 7 + 3.14;
END.  {Part10AST}`

func main() {
	lexer := NewLexer(pascalsample2)
	parser := NewParser(lexer)
	tree := parser.Parse()

	interpreter := NewInterpreter()
	interpreter.Interpret(tree)

	fmt.Printf("GLOBALSCOPE Table\n")
	fmt.Printf("-----------------\n")

	for str, val := range interpreter.GLOBALSCOPE {
		fmt.Printf("%-10s | %f\n", str, val)
	}

	av := NewASTVisualizer()
	av.Generate(tree)
}

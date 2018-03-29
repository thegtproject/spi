package main

// Parser ...
type Parser struct {
	CurrentToken Token
	lexer        *Lexer
}

// NewParser ...
func NewParser(l *Lexer) *Parser {
	p := &Parser{}
	p.lexer = l
	return p
}

func (p *Parser) Error() {
	panic("invalid syntax")
}

// Parse ...
func (p *Parser) Parse() Node {
	p.CurrentToken = p.lexer.GetNextToken()
	node := p.Program()
	if p.CurrentToken.Type != EOF {
		p.Error()
	}
	return node
}

// Eat ...
// compare the current token type with the passed token
// type and if they match then "eat" the current token
// and assign the next token to the in.CurrentToken,
// otherwise panic
func (p *Parser) Eat(tokenType int) {
	if p.CurrentToken.Type == tokenType {
		p.CurrentToken = p.lexer.GetNextToken()
	} else {
		p.Error()
	}
}

// Expr ...
// Arithmetic expression parser / interpreter.
//
// >  14 + 2 * 3 - 6 / 2
// =  17
//
// expr   : term ((PLUS | MINUS) term)*
// term   : factor ((MUL | DIV) factor)*
// factor : INTEGER
func (p *Parser) Expr() Node {
	node := p.Term()
	for p.CurrentToken.Type == PLUS ||
		p.CurrentToken.Type == MINUS {
		token := p.CurrentToken
		switch token.Type {
		case PLUS:
			p.Eat(PLUS)
		case MINUS:
			p.Eat(MINUS)
		}
		node = NewBinOp(node, token.Type, p.Term())
	}
	return node
}

// Term ...
// term : factor ((MUL | INTEGER_DIV | FLOAT_DIV) factor)*
func (p *Parser) Term() Node {
	node := p.Factor()
	for p.CurrentToken.Type == MUL ||
		p.CurrentToken.Type == INTEGERDIV ||
		p.CurrentToken.Type == FLOATDIV {
		token := p.CurrentToken
		switch token.Type {
		case MUL:
			p.Eat(MUL)
		case INTEGERDIV:
			p.Eat(INTEGERDIV)
		case FLOATDIV:
			p.Eat(FLOATDIV)
		}
		node = NewBinOp(node, token.Type, p.Factor())
	}
	return node
}

// Factor ...
// factor : PLUS  factor
//        | MINUS factor
//        | INTEGERCONST
//        | REALCONST
//        | LPAREN expr RPAREN
//        | variable
func (p *Parser) Factor() Node {
	token := p.CurrentToken
	switch token.Type {
	case PLUS:
		p.Eat(PLUS)
		return NewUnaryOp(PLUS, p.Factor())
	case MINUS:
		p.Eat(MINUS)
		return NewUnaryOp(MINUS, p.Factor())
	case INTEGERCONST:
		p.Eat(INTEGERCONST)
		return NewNum(token)
	case REALCONST:
		p.Eat(REALCONST)
		return NewNum(token)
	case LPAREN:
		p.Eat(LPAREN)
		node := p.Expr()
		p.Eat(RPAREN)
		return node
	default:
		return p.Variable()
	}
}

// Program ...
// program : PROGRAM variable SEMI block DOT
func (p *Parser) Program() Node {
	p.Eat(PROGRAM)
	varnode := p.Variable()
	programname := varnode.(*Var).Value
	p.Eat(SEMI)
	blocknode := p.Block()
	programnode := NewProgram(programname, blocknode)
	p.Eat(DOT)
	return programnode
}

// Block ...
// block : declarations compound_statement
func (p *Parser) Block() Node {
	declnodes := p.Declarations()
	compoundstatementnode := p.CompoundStatement()
	return NewBlock(declnodes, compoundstatementnode)
}

// Declarations ...
// declarations : VAR (variabledeclaration SEMI)+
//              | empty
func (p *Parser) Declarations() []Node {
	var declnodes []Node
	if p.CurrentToken.Type == VAR {
		p.Eat(VAR)
		for p.CurrentToken.Type == IDENT {
			vardecl := p.VariableDeclaration()
			declnodes = append(declnodes, vardecl...)
			p.Eat(SEMI)
		}
	}
	return declnodes
}

// VariableDeclaration ...
// variabledeclaration : IDENT (COMMA IDENT)* COLON typespec
func (p *Parser) VariableDeclaration() []Node {
	varnodes := []*Var{NewVar(p.CurrentToken, p.CurrentToken.Svalue)}
	p.Eat(IDENT)
	for p.CurrentToken.Type == COMMA {
		p.Eat(COMMA)
		varnodes = append(varnodes, NewVar(p.CurrentToken, p.CurrentToken.Svalue))
		p.Eat(IDENT)
	}
	p.Eat(COLON)
	typenode := p.TypeSpec()
	vardeclarations := []Node{}
	for _, varnode := range varnodes {
		vardeclarations = append(vardeclarations, NewVarDecl(varnode, typenode))
	}
	return vardeclarations
}

// TypeSpec ...
// type_spec : INTEGER
//           | REAL
func (p *Parser) TypeSpec() Node {
	token := p.CurrentToken
	switch token.Type {
	case INTEGER:
		p.Eat(INTEGER)
	case REAL:
		fallthrough
	default:
		p.Eat(REAL)
	}
	return NewTypeN(token)
}

// CompoundStatement ...
// compoundstatement: BEGIN statement_list END
func (p *Parser) CompoundStatement() Node {
	p.Eat(BEGIN)
	nodes := p.StatementList()
	p.Eat(END)
	node := NewCompound(nodes...)
	return node
}

// StatementList ...
//   statementlist : statement
// | statement SEMI statementlist
func (p *Parser) StatementList() []Node {
	results := []Node{
		p.Statement(),
	}
	for p.CurrentToken.Type == SEMI {
		p.Eat(SEMI)
		results = append(results, p.Statement())
	}
	if p.CurrentToken.Type == IDENT {
		p.Error()
	}
	return results
}

// Statement ...
// statement : compoundstatement
// | assignmentstatement
// | empty
func (p *Parser) Statement() Node {
	if p.CurrentToken.Type == BEGIN {
		return p.CompoundStatement()
	} else if p.CurrentToken.Type == IDENT {
		return p.AssignmentStatement()
	}
	return p.Empty()
}

// AssignmentStatement ...
// assignmentstatement : variable ASSIGN expr
func (p *Parser) AssignmentStatement() Node {
	left := p.Variable()
	token := p.CurrentToken
	p.Eat(ASSIGN)
	right := p.Expr()
	return NewAssign(left, token.Type, right)
}

// Variable ...
// variable : IDENT
func (p *Parser) Variable() Node {
	node := NewVar(p.CurrentToken, p.CurrentToken.Svalue)
	p.Eat(IDENT)
	return node
}

// Empty ...
// An empty production
func (p *Parser) Empty() Node {
	return NewNoOp()
}

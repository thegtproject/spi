package main

// Interpreter ...
type Interpreter struct {
	GLOBALSCOPE map[string]float64
	VisitMap    map[NodeType]func(n Node) float64
	parser      *Parser
}

// NewInterpreter ...
func NewInterpreter() *Interpreter {
	in := &Interpreter{}
	in.GLOBALSCOPE = make(map[string]float64)
	in.VisitMap = make(map[NodeType]func(n Node) float64)
	in.VisitMap[BinOpNode] = in.VisitBinOp
	in.VisitMap[UnaryOpNode] = in.VisitUnaryOp
	in.VisitMap[NumNode] = in.VisitNum
	in.VisitMap[CompoundNode] = in.VisitCompound
	in.VisitMap[AssignNode] = in.VisitAssign
	in.VisitMap[VarNode] = in.VisitVar
	in.VisitMap[NoOpNode] = in.VisitNoOp
	in.VisitMap[ProgramNode] = in.VisitProgram
	in.VisitMap[BlockNode] = in.VisitBlock
	in.VisitMap[VarDeclNode] = in.VisitVarDecl
	in.VisitMap[TypeNode] = in.VisitType
	return in
}

// Interpret ...
func (in *Interpreter) Interpret(n Node) {
	in.Visit(n)
}

func (in *Interpreter) Error() {
	panic("unknown node")
}

// VisitProgram ...
func (in *Interpreter) VisitProgram(n Node) float64 {
	node := n.(*Program)
	return in.Visit(node.BlockNode)
}

// VisitVarDecl ...
func (in *Interpreter) VisitVarDecl(n Node) float64 {
	return 0
}

// VisitType ...
func (in *Interpreter) VisitType(n Node) float64 {
	return 0
}

// VisitBlock ...
func (in *Interpreter) VisitBlock(n Node) float64 {
	node := n.(*Block)
	for _, declaration := range node.Decls {
		in.Visit(declaration)
	}
	in.Visit(node.CompoundStmt)
	return 0
}

// VisitBinOp ...
func (in *Interpreter) VisitBinOp(n Node) float64 {
	node := n.(*BinOp)
	switch node.Op {
	case PLUS:
		return in.Visit(node.Left) + in.Visit(node.Right)
	case MINUS:
		return in.Visit(node.Left) - in.Visit(node.Right)
	case MUL:
		return in.Visit(node.Left) * in.Visit(node.Right)
	case INTEGERDIV:
		return in.Visit(node.Left) / in.Visit(node.Right)
	case FLOATDIV:
		return in.Visit(node.Left) / in.Visit(node.Right)
	}
	in.Error()
	return 0
}

// VisitUnaryOp ...
func (in *Interpreter) VisitUnaryOp(n Node) float64 {
	node := n.(*UnaryOp)
	switch node.Op {
	case PLUS:
		return +in.Visit(node.Expr)
	case MINUS:
		return -in.Visit(node.Expr)
	}
	in.Error()
	return 0
}

// VisitNum ...
func (in *Interpreter) VisitNum(n Node) float64 {
	return n.(*Num).Value
}

// VisitCompound ...
func (in *Interpreter) VisitCompound(n Node) float64 {
	node := n.(*Compound)
	for i := range node.Children {
		in.Visit(node.Children[i])
	}
	return 0
}

// VisitAssign ...
func (in *Interpreter) VisitAssign(n Node) float64 {
	node := n.(*Assign)
	varname := node.Left.(*Var).Value
	in.GLOBALSCOPE[varname] = in.Visit(node.Right)
	return 0
}

// VisitVar ...
func (in *Interpreter) VisitVar(n Node) float64 {
	node := n.(*Var)
	varname := node.Value
	if varvalue, exists := in.GLOBALSCOPE[varname]; exists {
		return varvalue
	}
	in.Error()
	return 0
}

// VisitNoOp ...
func (in *Interpreter) VisitNoOp(n Node) float64 { return 0 }

// Visit ...
func (in *Interpreter) Visit(n Node) float64 {
	return in.VisitMap[n.Type()](n)
}

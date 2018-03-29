package main

// Node ...
type Node interface {
	String() string
	Type() NodeType
}

// NodeType ...
type NodeType int

// Node Types
const (
	BinOpNode = iota
	UnaryOpNode
	NumNode
	CompoundNode
	AssignNode
	VarNode
	NoOpNode
	ProgramNode
	BlockNode
	VarDeclNode
	TypeNode
)

// Type ...
func (nt NodeType) Type() NodeType {
	return nt
}

// BinOp ...
type BinOp struct {
	NodeType
	Op          int
	Left, Right Node
}

// NewBinOp ...
func NewBinOp(left Node, op int, right Node) *BinOp {
	return &BinOp{
		NodeType: BinOpNode,
		Left:     left,
		Op:       op,
		Right:    right,
	}
}

func (n *BinOp) String() string {
	return "binop"
}

// UnaryOp ...
type UnaryOp struct {
	NodeType
	Op   int
	Expr Node
}

// NewUnaryOp ...
func NewUnaryOp(op int, expr Node) *UnaryOp {
	return &UnaryOp{
		NodeType: UnaryOpNode,
		Op:       op,
		Expr:     expr,
	}
}

func (n *UnaryOp) String() string {
	return "UnaryOp"
}

// Num ...
type Num struct {
	NodeType
	Tok   Token
	Value float64
}

// NewNum ...
func NewNum(token Token) *Num {
	return &Num{
		NodeType: NumNode,
		Tok:      token,
		Value:    token.Value,
	}
}

func (n *Num) String() string {
	return "Num"
}

// Compound ...
type Compound struct {
	NodeType
	Children []Node
}

// NewCompound ...
func NewCompound(children ...Node) *Compound {
	node := &Compound{
		NodeType: CompoundNode,
	}
	if len(children) > 0 {
		node.Children = append(node.Children, children...)
	}
	return node
}

// Add ...
func (n *Compound) Add(node Node) {
	n.Children = append(n.Children, node)
}

func (n *Compound) String() string {
	return "Compound"
}

// Assign ...
type Assign struct {
	NodeType
	Op          int
	Left, Right Node
}

// NewAssign ...
func NewAssign(left Node, op int, right Node) *Assign {
	return &Assign{
		NodeType: AssignNode,
		Left:     left,
		Op:       op,
		Right:    right,
	}
}

func (n *Assign) String() string {
	return "Assign"
}

// Var ...
type Var struct {
	NodeType
	Tok   Token
	Value string
}

// NewVar ...
func NewVar(tok Token, value string) *Var {
	return &Var{
		NodeType: VarNode,
		Tok:      tok,
		Value:    value,
	}
}

func (n *Var) String() string {
	return "Var"
}

// NoOp ...
type NoOp struct {
	NodeType
}

// NewNoOp ...
func NewNoOp() *NoOp {
	return &NoOp{
		NodeType: NoOpNode,
	}
}

func (n *NoOp) String() string {
	return "NoOp"
}

// Program ...
type Program struct {
	NodeType
	Name      string
	BlockNode Node
}

// NewProgram ...
func NewProgram(name string, blocknode Node) *Program {
	return &Program{
		NodeType:  ProgramNode,
		Name:      name,
		BlockNode: blocknode,
	}
}

func (n *Program) String() string {
	return "Program"
}

// Block ...
type Block struct {
	NodeType
	Decls        []Node
	CompoundStmt Node
}

// NewBlock ...
func NewBlock(decls []Node, compoundstmt Node) *Block {
	return &Block{
		NodeType:     BlockNode,
		Decls:        decls,
		CompoundStmt: compoundstmt,
	}
}

func (n *Block) String() string {
	return "Block"
}

// VarDecl ...
type VarDecl struct {
	NodeType
	VNode Node
	TNode Node
}

// NewVarDecl ...
func NewVarDecl(vnode Node, tnode Node) *VarDecl {
	return &VarDecl{
		NodeType: VarDeclNode,
		VNode:    vnode,
		TNode:    tnode,
	}
}

func (n *VarDecl) String() string {
	return "VarDecl"
}

// TypeN ...
type TypeN struct {
	NodeType
	Tok   Token
	Value float64
}

// NewTypeN ...
func NewTypeN(tok Token) *TypeN {
	return &TypeN{
		NodeType: TypeNode,
		Tok:      tok,
		Value:    tok.Value,
	}
}

func (n *TypeN) String() string {
	return "TypeN"
}

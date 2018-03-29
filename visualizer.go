package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// ASTVisualizer ...
type ASTVisualizer struct {
	ID       int
	VisitMap map[NodeType]func(n Node) int

	buffer bytes.Buffer
}

// NewASTVisualizer ...
func NewASTVisualizer() *ASTVisualizer {
	av := &ASTVisualizer{}
	av.VisitMap = make(map[NodeType]func(n Node) int)
	av.VisitMap[BinOpNode] = av.VisitBinOp
	av.VisitMap[UnaryOpNode] = av.VisitUnaryOp
	av.VisitMap[NumNode] = av.VisitNum
	av.VisitMap[CompoundNode] = av.VisitCompound
	av.VisitMap[AssignNode] = av.VisitAssign
	av.VisitMap[VarNode] = av.VisitVar
	av.VisitMap[NoOpNode] = av.VisitNoOp
	av.VisitMap[ProgramNode] = av.VisitProgram
	av.VisitMap[BlockNode] = av.VisitBlock
	av.VisitMap[VarDeclNode] = av.VisitVarDecl
	av.VisitMap[TypeNode] = av.VisitType
	return av
}

// Generate ...
func (av *ASTVisualizer) Generate(n Node) {
	av.Visit(n)
	f, err := os.Create("astgraph.dot")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(DotHeader)
	f.Write(av.buffer.Bytes())
	f.WriteString("}")
	av.RunDot()
}

// RunDot ...
func (av *ASTVisualizer) RunDot() {
	_, err := exec.LookPath("dot")
	if err != nil {
		fmt.Printf("could not locate program \"dot\"\n")
		return
	}
	cmddot := exec.Command("dot", "-Tpng", "-oast.png", "astgraph.dot")
	out, err := cmddot.CombinedOutput()
	if err != nil {
		fmt.Println("error:", err)
	}
	if len(out) > 0 {
		fmt.Println("dot output:", string(out))
	}
}

// Visit ...
func (av *ASTVisualizer) Visit(n Node) int {
	return av.VisitMap[n.Type()](n)
}

// VisitProgram ...
func (av *ASTVisualizer) VisitProgram(n Node) int {
	node := n.(*Program)
	id := av.ID
	av.ID++
	s := fmt.Sprintf("Node%d [label=\"Program\n%s\"]\n", id, node.Name)
	av.buffer.WriteString(s)
	childid := av.Visit(node.BlockNode)
	s = fmt.Sprintf("Node%d -> Node%d\n", id, childid)
	av.buffer.WriteString(s)
	return id
}

// VisitBlock ...
func (av *ASTVisualizer) VisitBlock(n Node) int {
	node := n.(*Block)
	id := av.ID
	av.ID++
	s := fmt.Sprintf("Node%d [label=\"%s\"]\n", id, node.String())
	av.buffer.WriteString(s)
	for _, child := range node.Decls {
		childid := av.Visit(child)
		s = fmt.Sprintf("Node%d -> Node%d\n", id, childid)
		av.buffer.WriteString(s)
	}
	childid := av.Visit(node.CompoundStmt)
	s = fmt.Sprintf("Node%d -> Node%d\n", id, childid)
	av.buffer.WriteString(s)
	return id
}

// VisitVarDecl ...
func (av *ASTVisualizer) VisitVarDecl(n Node) int {
	node := n.(*VarDecl)
	id := av.ID
	av.ID++
	s := fmt.Sprintf("Node%d [label=\"%s\"]\n", id, "VarDecl")
	av.buffer.WriteString(s)
	lid := av.Visit(node.VNode)
	rid := av.Visit(node.TNode)
	s = fmt.Sprintf("Node%d -> Node%d\nNode%d -> Node%d\n", id, lid, id, rid)
	av.buffer.WriteString(s)
	return id
}

// VisitType ...
func (av *ASTVisualizer) VisitType(n Node) int {
	id := av.ID
	av.ID++
	s := fmt.Sprintf("Node%d [label=\"%s\"]\n", id, TokenStr[n.(*TypeN).Tok.Type])
	av.buffer.WriteString(s)
	return id
}

// VisitBinOp ...
func (av *ASTVisualizer) VisitBinOp(n Node) int {
	node := n.(*BinOp)
	op := TokenStr[node.Op]
	id := av.ID
	av.ID++
	s := fmt.Sprintf("Node%d [label=\"%s\"]\n", id, op)
	av.buffer.WriteString(s)
	lid := av.Visit(node.Left)
	rid := av.Visit(node.Right)
	s = fmt.Sprintf("Node%d -> Node%d\nNode%d -> Node%d\n", id, lid, id, rid)
	av.buffer.WriteString(s)
	return id
}

// VisitUnaryOp ...
func (av *ASTVisualizer) VisitUnaryOp(n Node) int {
	node := n.(*UnaryOp)
	op := TokenStr[node.Op]
	id := av.ID
	av.ID++
	s := fmt.Sprintf("Node%d [label=\"%s\"]\n", id, op)
	av.buffer.WriteString(s)
	childid := av.Visit(node.Expr)
	s = fmt.Sprintf("Node%d -> Node%d\n", id, childid)
	av.buffer.WriteString(s)
	return id
}

// VisitNum ...
func (av *ASTVisualizer) VisitNum(n Node) int {
	node := n.(*Num)
	id := av.ID
	av.ID++
	s := fmt.Sprintf("Node%d [label=\"%v\"]\n", id, node.Value)
	av.buffer.WriteString(s)
	return id
}

// VisitCompound ...
func (av *ASTVisualizer) VisitCompound(n Node) int {
	node := n.(*Compound)
	id := av.ID
	av.ID++
	s := fmt.Sprintf("Node%d [label=\"%s\"]\n", id, "Compound")
	av.buffer.WriteString(s)
	for i := range node.Children {
		cid := av.Visit(node.Children[i])
		s = fmt.Sprintf("Node%d -> Node%d\n", id, cid)
		av.buffer.WriteString(s)
	}
	return id
}

// VisitAssign ...
func (av *ASTVisualizer) VisitAssign(n Node) int {
	node := n.(*Assign)
	op := ":="
	id := av.ID
	av.ID++
	s := fmt.Sprintf("Node%d [label=\"%s\"]\n", id, op)
	av.buffer.WriteString(s)
	lid := av.Visit(node.Left)
	rid := av.Visit(node.Right)
	s = fmt.Sprintf("Node%d -> Node%d\nNode%d -> Node%d\n", id, lid, id, rid)
	av.buffer.WriteString(s)
	return id
}

// VisitVar ...
func (av *ASTVisualizer) VisitVar(n Node) int {
	node := n.(*Var)
	id := av.ID
	av.ID++
	s := fmt.Sprintf("Node%d [label=\"%s\"]\n", id, node.Value)
	av.buffer.WriteString(s)
	return id
}

// VisitNoOp ...
func (av *ASTVisualizer) VisitNoOp(n Node) int {
	id := av.ID
	av.ID++
	s := fmt.Sprintf("Node%d [label=\"%s\"]\n", id, "noop")
	av.buffer.WriteString(s)
	return id
}

var (
	// DotHeader ...
	DotHeader = `digraph astgraph {
  node [shape=circle, fontsize=12, fontname="Courier", height=.1];
  ranksep=.3;
  edge [arrowsize=.5]
`
)

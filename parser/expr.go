package parser

import "glox/scanner"

type VisitorExpr interface {
	VisitTernaryExpr(t TernaryExpr) any
	VisitBinaryExpr(b BinaryExpr) any
	VisitGroupingExpr(g GroupingExpr) any
	VisitLiteralExpr(l LiteralExpr) any
	VisitUnaryExpr(u UnaryExpr) any
	VisitVariableExpr(v VariableExpr) any
}

type Expr interface {
	Accept(visitor VisitorExpr) any
}

type TernaryExpr struct {
	Condition Expr
	Left      Expr
	Right     Expr
}

func (t TernaryExpr) Accept(visitor VisitorExpr) any {
	return visitor.VisitTernaryExpr(t)
}

type BinaryExpr struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

func (b BinaryExpr) Accept(visitor VisitorExpr) any {
	return visitor.VisitBinaryExpr(b)
}

type GroupingExpr struct {
	Expr Expr
}

func (g GroupingExpr) Accept(visitor VisitorExpr) any {
	return visitor.VisitGroupingExpr(g)
}

type LiteralExpr struct {
	Value any
}

func (l LiteralExpr) Accept(visitor VisitorExpr) any {
	return visitor.VisitLiteralExpr(l)
}

type UnaryExpr struct {
	Operator scanner.Token
	Right    Expr
}

func (u UnaryExpr) Accept(visitor VisitorExpr) any {
	return visitor.VisitUnaryExpr(u)
}

type VariableExpr struct {
	Name scanner.Token
}

func (v VariableExpr) Accept(visitor VisitorExpr) any {
	return visitor.VisitVariableExpr(v)
}

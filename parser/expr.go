package parser

import "glox/scanner"

type VisitorExpr interface {
	VisitTernaryExpr(TernaryExpr) any
	VisitAssignmentExpr(AssignmentExpr) any
	VisitBinaryExpr(BinaryExpr) any
	VisitGroupingExpr(GroupingExpr) any
	VisitLiteralExpr(LiteralExpr) any
	VisitUnaryExpr(UnaryExpr) any
	VisitVariableExpr(VariableExpr) any
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

type AssignmentExpr struct {
	Name  scanner.Token
	Value Expr
}

func (a AssignmentExpr) Accept(visitor VisitorExpr) any {
	return visitor.VisitAssignmentExpr(a)
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

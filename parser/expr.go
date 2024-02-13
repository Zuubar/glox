package parser

import "glox/scanner"

type VisitorExpr interface {
	VisitTernaryExpr(TernaryExpr) (any, error)
	VisitAssignmentExpr(AssignmentExpr) (any, error)
	VisitLogicalExpr(LogicalExpr) (any, error)
	VisitBinaryExpr(BinaryExpr) (any, error)
	VisitGroupingExpr(GroupingExpr) (any, error)
	VisitLiteralExpr(LiteralExpr) (any, error)
	VisitUnaryExpr(UnaryExpr) (any, error)
	VisitVariableExpr(VariableExpr) (any, error)
}

type Expr interface {
	Accept(visitor VisitorExpr) (any, error)
}

type TernaryExpr struct {
	Condition Expr
	Left      Expr
	Right     Expr
}

func (t TernaryExpr) Accept(visitor VisitorExpr) (any, error) {
	return visitor.VisitTernaryExpr(t)
}

type AssignmentExpr struct {
	Name  scanner.Token
	Value Expr
}

func (a AssignmentExpr) Accept(visitor VisitorExpr) (any, error) {
	return visitor.VisitAssignmentExpr(a)
}

type LogicalExpr struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

func (l LogicalExpr) Accept(visitor VisitorExpr) (any, error) {
	return visitor.VisitLogicalExpr(l)
}

type BinaryExpr struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

func (b BinaryExpr) Accept(visitor VisitorExpr) (any, error) {
	return visitor.VisitBinaryExpr(b)
}

type GroupingExpr struct {
	Expr Expr
}

func (g GroupingExpr) Accept(visitor VisitorExpr) (any, error) {
	return visitor.VisitGroupingExpr(g)
}

type LiteralExpr struct {
	Value any
}

func (l LiteralExpr) Accept(visitor VisitorExpr) (any, error) {
	return visitor.VisitLiteralExpr(l)
}

type UnaryExpr struct {
	Operator scanner.Token
	Right    Expr
}

func (u UnaryExpr) Accept(visitor VisitorExpr) (any, error) {
	return visitor.VisitUnaryExpr(u)
}

type VariableExpr struct {
	Name scanner.Token
}

func (v VariableExpr) Accept(visitor VisitorExpr) (any, error) {
	return visitor.VisitVariableExpr(v)
}

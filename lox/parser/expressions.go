package parser

import (
	"glox/lox/scanner"
)

type Visitor interface {
	VisitLiteral(literal LiteralExpr) any
	VisitGrouping(grouping GroupingExpr) any
	VisitUnary(unary UnaryExpr) any
	VisitBinary(binary Binary) any
	VisitTernary(ternary TernaryExpr) any
}

type Expr interface {
	Accept(visitor Visitor) any
}

type LiteralExpr struct {
	Value any
}

func (l LiteralExpr) Accept(visitor Visitor) any {
	return visitor.VisitLiteral(l)
}

type GroupingExpr struct {
	Expr Expr
}

func (g GroupingExpr) Accept(visitor Visitor) any {
	return visitor.VisitGrouping(g)
}

type UnaryExpr struct {
	Operator scanner.Token
	Right    Expr
}

func (u UnaryExpr) Accept(visitor Visitor) any {
	return visitor.VisitUnary(u)
}

type Binary struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

func (b Binary) Accept(visitor Visitor) any {
	return visitor.VisitBinary(b)
}

type TernaryExpr struct {
	Condition Expr
	Left      Expr
	Right     Expr
}

func (t TernaryExpr) Accept(visitor Visitor) any {
	return visitor.VisitTernary(t)
}

type ErrorExpr struct {
}

func (e Error) Accept(_ Visitor) any {
	return nil
}

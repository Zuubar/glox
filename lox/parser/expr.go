package parser

import (
	"glox/lox/scanner"
)

type Visitor interface {
	VisitLiteral(literal Literal) any
	VisitGrouping(grouping Grouping) any
	VisitUnary(unary Unary) any
	VisitBinary(binary Binary) any
	VisitTernary(ternary Ternary) any
}

type Expr interface {
	Accept(visitor Visitor) any
}

type Literal struct {
	Value any
}

func (l Literal) Accept(visitor Visitor) any {
	return visitor.VisitLiteral(l)
}

type Grouping struct {
	Expr Expr
}

func (g Grouping) Accept(visitor Visitor) any {
	return visitor.VisitGrouping(g)
}

type Binary struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

func (b Binary) Accept(visitor Visitor) any {
	return visitor.VisitBinary(b)
}

type Ternary struct {
	Left   Expr
	Middle Expr
	Right  Expr
}

func (t Ternary) Accept(visitor Visitor) any {
	return visitor.VisitTernary(t)
}

type Unary struct {
	Operator scanner.Token
	Right    Expr
}

func (u Unary) Accept(visitor Visitor) any {
	return visitor.VisitUnary(u)
}

type Error struct {
}

func (e Error) Accept(_ Visitor) any {
	return nil
}

package parser

import "glox/scanner"

type VisitorExpr interface {
	VisitArrayExpr(ArrayExpr) (any, error)
	VisitTernaryExpr(TernaryExpr) (any, error)
	VisitAssignmentExpr(AssignmentExpr) (any, error)
	VisitLogicalExpr(LogicalExpr) (any, error)
	VisitSetExpr(SetExpr) (any, error)
	VisitArraySetExpr(ArraySetExpr) (any, error)
	VisitSuperExpr(SuperExpr) (any, error)
	VisitBinaryExpr(BinaryExpr) (any, error)
	VisitGroupingExpr(GroupingExpr) (any, error)
	VisitLiteralExpr(LiteralExpr) (any, error)
	VisitUnaryExpr(UnaryExpr) (any, error)
	VisitGetExpr(GetExpr) (any, error)
	VisitArrayGetExpr(ArrayGetExpr) (any, error)
	VisitCallExpr(CallExpr) (any, error)
	VisitLambdaExpr(LambdaExpr) (any, error)
	VisitThisExpr(ThisExpr) (any, error)
	VisitVariableExpr(VariableExpr) (any, error)
}

type Expr interface {
	Accept(visitor VisitorExpr) (any, error)
}

type ArrayExpr struct {
	Elements []Expr
	Bracket  scanner.Token
}

func (a ArrayExpr) Accept(visitor VisitorExpr) (any, error) {
	return visitor.VisitArrayExpr(a)
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

type SetExpr struct {
	Object Expr
	Name   scanner.Token
	Value  Expr
}

func (s SetExpr) Accept(visitor VisitorExpr) (any, error) {
	return visitor.VisitSetExpr(s)
}

type ArraySetExpr struct {
	Array   Expr
	Bracket scanner.Token
	Index   Expr
	Value   Expr
}

func (a ArraySetExpr) Accept(visitor VisitorExpr) (any, error) {
	return visitor.VisitArraySetExpr(a)
}

type SuperExpr struct {
	Keyword scanner.Token
	Method  scanner.Token
}

func (s SuperExpr) Accept(visitor VisitorExpr) (any, error) {
	return visitor.VisitSuperExpr(s)
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

type GetExpr struct {
	Object Expr
	Name   scanner.Token
}

func (g GetExpr) Accept(visitor VisitorExpr) (any, error) {
	return visitor.VisitGetExpr(g)
}

type ArrayGetExpr struct {
	Array   Expr
	Bracket scanner.Token
	Index   Expr
}

func (a ArrayGetExpr) Accept(visitor VisitorExpr) (any, error) {
	return visitor.VisitArrayGetExpr(a)
}

type CallExpr struct {
	Callee      Expr
	Parenthesis scanner.Token
	Arguments   []Expr
}

func (c CallExpr) Accept(visitor VisitorExpr) (any, error) {
	return visitor.VisitCallExpr(c)
}

type LambdaExpr struct {
	Parenthesis scanner.Token
	Parameters  []scanner.Token
	Body        []Stmt
}

func (l LambdaExpr) Accept(visitor VisitorExpr) (any, error) {
	return visitor.VisitLambdaExpr(l)
}

type ThisExpr struct {
	Keyword scanner.Token
}

func (t ThisExpr) Accept(visitor VisitorExpr) (any, error) {
	return visitor.VisitThisExpr(t)
}

type VariableExpr struct {
	Name scanner.Token
}

func (v VariableExpr) Accept(visitor VisitorExpr) (any, error) {
	return visitor.VisitVariableExpr(v)
}

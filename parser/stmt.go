package parser

import "glox/scanner"

type VisitorStmt interface {
	VisitExpressionStmt(e ExpressionStmt) any
	VisitPrintStmt(p PrintStmt) any
	VisitVarStmt(v VarStmt) any
}

type Stmt interface {
	Accept(visitor VisitorStmt) any
}

type ExpressionStmt struct {
	Expression Expr
}

func (e ExpressionStmt) Accept(visitor VisitorStmt) any {
	return visitor.VisitExpressionStmt(e)
}

type PrintStmt struct {
	Expression Expr
}

func (p PrintStmt) Accept(visitor VisitorStmt) any {
	return visitor.VisitPrintStmt(p)
}

type VarStmt struct {
	Name        scanner.Token
	Initializer Expr
}

func (v VarStmt) Accept(visitor VisitorStmt) any {
	return visitor.VisitVarStmt(v)
}

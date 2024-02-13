package parser

import "glox/scanner"

type VisitorStmt interface {
	VisitExpressionStmt(ExpressionStmt) (any, error)
	VisitPrintStmt(PrintStmt) (any, error)
	VisitVarStmt(VarStmt) (any, error)
	VisitBlockStmt(BlockStmt) (any, error)
	VisitIfStmt(IfStmt) (any, error)
}

type Stmt interface {
	Accept(visitor VisitorStmt) (any, error)
}

type ExpressionStmt struct {
	Expression Expr
}

func (e ExpressionStmt) Accept(visitor VisitorStmt) (any, error) {
	return visitor.VisitExpressionStmt(e)
}

type PrintStmt struct {
	Expression Expr
}

func (p PrintStmt) Accept(visitor VisitorStmt) (any, error) {
	return visitor.VisitPrintStmt(p)
}

type VarStmt struct {
	Name        scanner.Token
	Initializer Expr
}

func (v VarStmt) Accept(visitor VisitorStmt) (any, error) {
	return visitor.VisitVarStmt(v)
}

type BlockStmt struct {
	Declarations []Stmt
}

func (b BlockStmt) Accept(visitor VisitorStmt) (any, error) {
	return visitor.VisitBlockStmt(b)
}

type IfStmt struct {
	Expression Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (i IfStmt) Accept(visitor VisitorStmt) (any, error) {
	return visitor.VisitIfStmt(i)
}

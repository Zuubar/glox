package parser

import "glox/scanner"

type VisitorStmt interface {
	VisitExpressionStmt(ExpressionStmt) (any, error)
	VisitPrintStmt(PrintStmt) (any, error)
	VisitVarStmt(VarStmt) (any, error)
	VisitBlockStmt(BlockStmt) (any, error)
	VisitIfStmt(IfStmt) (any, error)
	VisitWhileStmt(WhileStmt) (any, error)
	VisitForStmt(ForStmt) (any, error)
	VisitBreakStmt(BreakStmt) (any, error)
	VisitContinueStmt(ContinueStmt) (any, error)
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

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (w WhileStmt) Accept(visitor VisitorStmt) (any, error) {
	return visitor.VisitWhileStmt(w)
}

type ForStmt struct {
	Initializer Stmt
	Condition   Expr
	Increment   Stmt
	Body        Stmt
}

func (f ForStmt) Accept(visitor VisitorStmt) (any, error) {
	return visitor.VisitForStmt(f)
}

type BreakStmt struct {
	At scanner.Token
}

func (b BreakStmt) Accept(visitor VisitorStmt) (any, error) {
	return visitor.VisitBreakStmt(b)
}

type ContinueStmt struct {
	At scanner.Token
}

func (c ContinueStmt) Accept(visitor VisitorStmt) (any, error) {
	return visitor.VisitContinueStmt(c)
}

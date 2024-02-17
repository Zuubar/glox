package parser

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.WriteString("(" + name)
	for _, expr := range exprs {
		var value any = nil
		if expr != nil {
			value, _ = expr.Accept(a)
		}
		builder.WriteString(" " + fmt.Sprintf("%v", value))
	}
	builder.WriteString(")")
	return builder.String()
}

func (a *AstPrinter) parenthesizeStmt(name string, stmts ...Stmt) string {
	var builder strings.Builder

	builder.WriteString("(" + name)
	for _, expr := range stmts {
		var value any = nil
		if expr != nil {
			value, _ = expr.Accept(a)
		}
		builder.WriteString(" " + fmt.Sprintf("%v", value))
	}
	builder.WriteString(")")
	return builder.String()
}

func (a *AstPrinter) VisitLiteralExpr(literal LiteralExpr) (any, error) {
	if literal.Value == nil {
		return "", nil
	}

	return fmt.Sprintf("%v", literal.Value), nil
}

func (a *AstPrinter) VisitGroupingExpr(grouping GroupingExpr) (any, error) {
	return a.parenthesize("group", grouping.Expr), nil
}

func (a *AstPrinter) VisitUnaryExpr(unary UnaryExpr) (any, error) {
	return a.parenthesize(unary.Operator.Lexeme, unary.Right), nil
}

func (a *AstPrinter) VisitBinaryExpr(binary BinaryExpr) (any, error) {
	return a.parenthesize(binary.Operator.Lexeme, binary.Left, binary.Right), nil
}

func (a *AstPrinter) VisitLogicalExpr(expr LogicalExpr) (any, error) {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
}

func (a *AstPrinter) VisitAssignmentExpr(assignment AssignmentExpr) (any, error) {
	return a.parenthesize("assignment", assignment.Value), nil
}

func (a *AstPrinter) VisitTernaryExpr(ternary TernaryExpr) (any, error) {
	return a.parenthesize("?:", ternary.Condition, ternary.Left, ternary.Right), nil
}

func (a *AstPrinter) VisitVariableExpr(variableExpr VariableExpr) (any, error) {
	return a.parenthesize("variableExpr " + variableExpr.Name.Lexeme), nil
}

func (a *AstPrinter) VisitExpressionStmt(expressionStmt ExpressionStmt) (any, error) {
	return a.parenthesize("exprStmt", expressionStmt.Expression), nil
}

func (a *AstPrinter) VisitPrintStmt(printStmt PrintStmt) (any, error) {
	return a.parenthesize("printStmt", printStmt.Expression), nil
}

func (a *AstPrinter) VisitVarStmt(varStmt VarStmt) (any, error) {
	return a.parenthesize("varStmt "+varStmt.Name.Lexeme, varStmt.Initializer), nil
}

func (a *AstPrinter) VisitBlockStmt(stmt BlockStmt) (any, error) {
	return a.parenthesizeStmt("block", stmt.Declarations...), nil
}

func (a *AstPrinter) VisitIfStmt(stmt IfStmt) (any, error) {
	return a.parenthesizeStmt("ifStmt", stmt.ThenBranch, stmt.ElseBranch), nil
}

func (a *AstPrinter) VisitWhileStmt(stmt WhileStmt) (any, error) {
	return a.parenthesizeStmt("whileStmt", stmt.Body), nil
}

func (a *AstPrinter) VisitForStmt(stmt ForStmt) (any, error) {
	return a.parenthesizeStmt("forStmt", stmt.Body), nil
}

func (a *AstPrinter) VisitBreakStmt(_ BreakStmt) (any, error) {
	return a.parenthesizeStmt("breakStmt", nil), nil
}

func (a *AstPrinter) VisitContinueStmt(_ ContinueStmt) (any, error) {
	return a.parenthesizeStmt("continueStmt", nil), nil
}

func (a *AstPrinter) Print(statements []Stmt) any {
	results := make([]string, 0, 10)
	for _, statement := range statements {
		stmt, _ := statement.Accept(a)
		results = append(results, stmt.(string))
	}
	return results
}

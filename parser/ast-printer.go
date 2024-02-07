package parser

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
}

func (a *AstPrinter) VisitLiteralExpr(literal LiteralExpr) any {
	if literal.Value == nil {
		return ""
	}

	return fmt.Sprintf("%v", literal.Value)
}

func (a *AstPrinter) VisitGroupingExpr(grouping GroupingExpr) any {
	return a.parenthesize("group", grouping.Expr)
}

func (a *AstPrinter) VisitUnaryExpr(unary UnaryExpr) any {
	return a.parenthesize(unary.Operator.Lexeme, unary.Right)
}

func (a *AstPrinter) VisitBinaryExpr(binary BinaryExpr) any {
	return a.parenthesize(binary.Operator.Lexeme, binary.Left, binary.Right)
}

func (a *AstPrinter) VisitAssignmentExpr(assignment AssignmentExpr) any {
	return a.parenthesize("assignment", assignment.Value)
}

func (a *AstPrinter) VisitTernaryExpr(ternary TernaryExpr) any {
	return a.parenthesize("?:", ternary.Condition, ternary.Left, ternary.Right)
}

func (a *AstPrinter) VisitVariableExpr(variableExpr VariableExpr) any {
	return a.parenthesize("variableExpr " + variableExpr.Name.Lexeme)
}

func (a *AstPrinter) VisitExpressionStmt(expressionStmt ExpressionStmt) any {
	return a.parenthesize("exprStmt", expressionStmt.Expression)
}

func (a *AstPrinter) VisitPrintStmt(printStmt PrintStmt) any {
	return a.parenthesize("printStmt", printStmt.Expression)
}

func (a *AstPrinter) VisitVarStmt(varStmt VarStmt) any {
	return a.parenthesize("varStmt "+varStmt.Name.Lexeme, varStmt.Initializer)
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.WriteString("(" + name)
	for _, expr := range exprs {
		var value any = nil
		if expr != nil {
			value = expr.Accept(a)
		}
		builder.WriteString(" " + fmt.Sprintf("%v", value))
	}
	builder.WriteString(")")
	return builder.String()
}

func (a *AstPrinter) Print(statements []Stmt) any {
	results := make([]string, 0, 10)
	for _, statement := range statements {
		results = append(results, statement.Accept(a).(string))
	}
	return results
}

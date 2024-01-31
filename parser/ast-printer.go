package parser

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
}

func (a AstPrinter) VisitLiteralExpr(literal LiteralExpr) any {
	if literal.Value == nil {
		return ""
	}

	return fmt.Sprintf("%v", literal.Value)
}

func (a AstPrinter) VisitGroupingExpr(grouping GroupingExpr) any {
	return a.parenthesize("group", grouping.Expr)
}

func (a AstPrinter) VisitUnaryExpr(unary UnaryExpr) any {
	return a.parenthesize(unary.Operator.Lexeme, unary.Right)
}

func (a AstPrinter) VisitBinaryExpr(binary BinaryExpr) any {
	return a.parenthesize(binary.Operator.Lexeme, binary.Left, binary.Right)
}

func (a AstPrinter) VisitTernaryExpr(ternary TernaryExpr) any {
	return fmt.Sprintf("(?: %v %v %v)", ternary.Condition.Accept(a), ternary.Left.Accept(a), ternary.Right.Accept(a))
}

func (a AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.WriteString("(" + name)
	for _, expr := range exprs {
		builder.WriteString(" " + fmt.Sprintf("%v", expr.Accept(a)))
	}
	builder.WriteString(")")
	return builder.String()
}

func (a AstPrinter) Print(expr Expr) any {
	return expr.Accept(a)
}

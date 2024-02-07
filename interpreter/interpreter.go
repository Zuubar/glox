package interpreter

import (
	"fmt"
	"glox/parser"
	"glox/scanner"
	"reflect"
)

type Interpreter struct {
	environment *environment
}

func New() *Interpreter {
	return &Interpreter{environment: newEnvironment()}
}

func (i *Interpreter) isTruthy(obj any) bool {
	if obj == nil {
		return false
	}

	if t, ok := obj.(bool); ok {
		return t
	}

	return true
}

func (i *Interpreter) isType(obj any, targetType reflect.Kind) bool {
	return obj != nil && reflect.TypeOf(obj).Kind() == targetType
}

func (i *Interpreter) areNumberedOperands(obj1 any, obj2 any) bool {
	return i.isType(obj1, reflect.Float64) && i.isType(obj2, reflect.Float64)
}

func (i *Interpreter) areStringOperands(obj1 any, obj2 any) bool {
	return i.isType(obj1, reflect.String) && i.isType(obj2, reflect.String)
}

func (i *Interpreter) areEqual(obj1 any, obj2 any) bool {
	if obj1 == nil && obj2 == nil {
		return true
	}
	if obj1 == nil {
		return false
	}
	return obj1 == obj2
}

func (i *Interpreter) stringify(obj any) string {
	if obj == nil {
		return "nil"
	}

	return fmt.Sprintf("%v", obj)
}

func (i *Interpreter) VisitLiteralExpr(literal parser.LiteralExpr) any {
	return literal.Value
}

func (i *Interpreter) VisitGroupingExpr(grouping parser.GroupingExpr) any {
	return grouping.Expr.Accept(i)
}

func (i *Interpreter) VisitUnaryExpr(unary parser.UnaryExpr) any {
	token, obj := unary.Operator, unary.Right.Accept(i)

	switch unary.Operator.Type {
	case scanner.BANG:
		return !i.isTruthy(obj)
	case scanner.MINUS:
		if i.isType(obj, reflect.Float64) {
			return obj.(float64)
		}
	}

	return Error{Token: token, Message: "Operand must be a number."}
}

func (i *Interpreter) VisitBinaryExpr(binary parser.BinaryExpr) any {
	obj1, token, obj2 := binary.Left.Accept(i), binary.Operator, binary.Right.Accept(i)

	switch binary.Operator.Type {
	case scanner.PLUS:
		if i.areNumberedOperands(obj1, obj2) {
			return obj1.(float64) + obj2.(float64)
		}

		if i.areStringOperands(obj1, obj2) {
			return obj1.(string) + obj2.(string)
		}
		return Error{Token: token, Message: "Both operands should be numbers or strings."}
	case scanner.MINUS:
		if i.areNumberedOperands(obj1, obj2) {
			return obj1.(float64) - obj2.(float64)
		}
		return Error{Token: token, Message: "Both operands should be numbers."}
	case scanner.STAR:
		if i.areNumberedOperands(obj1, obj2) {
			return obj1.(float64) * obj2.(float64)
		}
		return Error{Token: token, Message: "Both operands should be numbers."}
	case scanner.SLASH:
		if i.areNumberedOperands(obj1, obj2) {
			left, _ := obj1.(float64)
			right, _ := obj2.(float64)

			if right == 0 {
				return Error{Token: token, Message: "Division by zero is prohibited."}
			}

			return left / right
		}
		return Error{Token: token, Message: "Both operands should be numbers."}
	case scanner.GREATER:
		if i.areNumberedOperands(obj1, obj2) {
			return obj1.(float64) > obj2.(float64)
		}
		return Error{Token: token, Message: "Both operands should be numbers."}
	case scanner.GREATER_EQUAL:
		if i.areNumberedOperands(obj1, obj2) {
			return obj1.(float64) >= obj2.(float64)
		}
		return Error{Token: token, Message: "Both operands should be numbers."}
	case scanner.LESS:
		if i.areNumberedOperands(obj1, obj2) {
			return obj1.(float64) < obj2.(float64)
		}
		return Error{Token: token, Message: "Both operands should be numbers."}
	case scanner.LESS_EQUAL:
		if i.areNumberedOperands(obj1, obj2) {
			return obj1.(float64) < obj2.(float64)
		}
		return Error{Token: token, Message: "Both operands should be numbers."}
	case scanner.EQUAL_EQUAL:
		return i.areEqual(obj1, obj2)
	case scanner.BANG_EQUAL:
		return !i.areEqual(obj1, obj2)
	}

	return nil
}

func (i *Interpreter) VisitAssignmentExpr(assignment parser.AssignmentExpr) any {
	token := assignment.Name
	if _, ok := i.environment.lookup(token.Lexeme); !ok {
		return Error{Token: token, Message: fmt.Sprintf("Undefined variable '%s'.", token.Lexeme)}
	}

	value := assignment.Value.Accept(i)
	if err, ok := value.(Error); ok {
		return err
	}

	i.environment.define(token.Lexeme, value)
	return value
}

func (i *Interpreter) VisitTernaryExpr(ternary parser.TernaryExpr) any {
	obj := ternary.Condition.Accept(i)
	if i.isTruthy(obj) {
		return ternary.Left.Accept(i)
	}

	return ternary.Right.Accept(i)
}

func (i *Interpreter) VisitVariableExpr(variableExpr parser.VariableExpr) any {
	value, ok := i.environment.lookup(variableExpr.Name.Lexeme)
	if !ok {
		return Error{Token: variableExpr.Name, Message: fmt.Sprintf("Undefined variable '%s'.", variableExpr.Name.Lexeme)}
	}
	return value
}

func (i *Interpreter) VisitExpressionStmt(expressionStmt parser.ExpressionStmt) any {
	result := expressionStmt.Expression.Accept(i)
	if err, ok := result.(Error); ok {
		return err
	}
	return nil
}

func (i *Interpreter) VisitPrintStmt(printStmt parser.PrintStmt) any {
	result := printStmt.Expression.Accept(i)
	if err, ok := result.(Error); ok {
		return err
	}
	fmt.Println(i.stringify(result))
	return nil
}

func (i *Interpreter) VisitVarStmt(varStmt parser.VarStmt) any {
	var value any = nil

	if varStmt.Initializer != nil {
		value = varStmt.Initializer.Accept(i)
		if err, ok := value.(Error); ok {
			return err
		}
	}

	i.environment.define(varStmt.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) Interpret(statements []parser.Stmt) error {
	for _, stmt := range statements {
		result := stmt.Accept(i)
		if err, ok := result.(Error); ok {
			return &err
		}
	}
	return nil
}

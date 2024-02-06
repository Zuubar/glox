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

func (i *Interpreter) isTruthy(object any) bool {
	if object == nil {
		return false
	}

	if t, ok := object.(bool); ok {
		return t
	}

	return true
}

func (i *Interpreter) isType(object any, targetType reflect.Kind) bool {
	return object != nil && reflect.TypeOf(object).Kind() == targetType
}

func (i *Interpreter) areNumberedOperands(object1 any, object2 any) bool {
	return i.isType(object1, reflect.Float64) && i.isType(object2, reflect.Float64)
}

func (i *Interpreter) areStringOperands(object1 any, object2 any) bool {
	return i.isType(object1, reflect.String) && i.isType(object2, reflect.String)
}

func (i *Interpreter) areEqual(object1 any, object2 any) bool {
	if object1 == nil && object2 == nil {
		return true
	}
	if object1 == nil {
		return false
	}
	return object1 == object2
}

func (i *Interpreter) VisitLiteralExpr(literal parser.LiteralExpr) any {
	return literal.Value
}

func (i *Interpreter) VisitGroupingExpr(grouping parser.GroupingExpr) any {
	return grouping.Expr.Accept(i)
}

func (i *Interpreter) VisitUnaryExpr(unary parser.UnaryExpr) any {
	token, object := unary.Operator, unary.Right.Accept(i)

	switch unary.Operator.Type {
	case scanner.BANG:
		return !i.isTruthy(object)
	case scanner.MINUS:
		if i.isType(object, reflect.Float64) {
			return object.(float64)
		}
	}

	return RuntimeError{Token: token, Message: "Operand must be a number."}
}

func (i *Interpreter) VisitBinaryExpr(binary parser.BinaryExpr) any {
	object1, token, object2 := binary.Left.Accept(i), binary.Operator, binary.Right.Accept(i)

	switch binary.Operator.Type {
	case scanner.PLUS:
		if i.areNumberedOperands(object1, object2) {
			return object1.(float64) + object2.(float64)
		}

		if i.areStringOperands(object1, object2) {
			return object1.(string) + object2.(string)
		}
		return RuntimeError{Token: token, Message: "Both operands should be numbers or strings."}
	case scanner.MINUS:
		if i.areNumberedOperands(object1, object2) {
			return object1.(float64) - object2.(float64)
		}
		return RuntimeError{Token: token, Message: "Both operands should be numbers."}
	case scanner.STAR:
		if i.areNumberedOperands(object1, object2) {
			return object1.(float64) * object2.(float64)
		}
		return RuntimeError{Token: token, Message: "Both operands should be numbers."}
	case scanner.SLASH:
		if i.areNumberedOperands(object1, object2) {
			left, _ := object1.(float64)
			right, _ := object2.(float64)

			if right == 0 {
				return RuntimeError{Token: token, Message: "Division by zero is prohibited."}
			}

			return left / right
		}
		return RuntimeError{Token: token, Message: "Both operands should be numbers."}
	case scanner.GREATER:
		if i.areNumberedOperands(object1, object2) {
			return object1.(float64) > object2.(float64)
		}
		return RuntimeError{Token: token, Message: "Both operands should be numbers."}
	case scanner.GREATER_EQUAL:
		if i.areNumberedOperands(object1, object2) {
			return object1.(float64) >= object2.(float64)
		}
		return RuntimeError{Token: token, Message: "Both operands should be numbers."}
	case scanner.LESS:
		if i.areNumberedOperands(object1, object2) {
			return object1.(float64) < object2.(float64)
		}
		return RuntimeError{Token: token, Message: "Both operands should be numbers."}
	case scanner.LESS_EQUAL:
		if i.areNumberedOperands(object1, object2) {
			return object1.(float64) < object2.(float64)
		}
		return RuntimeError{Token: token, Message: "Both operands should be numbers."}
	case scanner.EQUAL_EQUAL:
		return i.areEqual(object1, object2)
	case scanner.BANG_EQUAL:
		return !i.areEqual(object1, object2)
	}

	return nil
}

func (i *Interpreter) VisitTernaryExpr(ternary parser.TernaryExpr) any {
	object := ternary.Condition.Accept(i)
	if i.isTruthy(object) {
		return ternary.Left.Accept(i)
	}

	return ternary.Right.Accept(i)
}

func (i *Interpreter) VisitVariableExpr(v parser.VariableExpr) any {
	value, ok := i.environment.lookup(v.Name.Lexeme)
	if !ok {
		return RuntimeError{Token: v.Name, Message: fmt.Sprintf("Undefined variable '%s'.", v.Name.Lexeme)}
	}
	return value
}

func (i *Interpreter) VisitExpressionStmt(e parser.ExpressionStmt) any {
	result := e.Expression.Accept(i)
	if err, ok := result.(RuntimeError); ok {
		return err
	}
	return nil
}

func (i *Interpreter) VisitPrintStmt(p parser.PrintStmt) any {
	result := p.Expression.Accept(i)
	if err, ok := result.(RuntimeError); ok {
		return err
	}
	fmt.Println(result)
	return nil
}

func (i *Interpreter) VisitVarStmt(v parser.VarStmt) any {
	var value any = nil

	if v.Initializer != nil {
		value = v.Initializer.Accept(i)
		if err, ok := value.(RuntimeError); ok {
			return err
		}
	}

	i.environment.define(v.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) Interpret(statements []parser.Stmt) error {
	for _, stmt := range statements {
		result := stmt.Accept(i)
		if err, ok := result.(RuntimeError); ok {
			return &err
		}
	}
	return nil
}

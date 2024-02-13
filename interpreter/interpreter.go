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
	return &Interpreter{environment: newEnvironment(nil)}
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

func (i *Interpreter) VisitLiteralExpr(literal parser.LiteralExpr) (any, error) {
	return literal.Value, nil
}

func (i *Interpreter) VisitGroupingExpr(grouping parser.GroupingExpr) (any, error) {
	return grouping.Expr.Accept(i)
}

func (i *Interpreter) VisitUnaryExpr(unary parser.UnaryExpr) (any, error) {
	obj, err := unary.Right.Accept(i)
	if err != nil {
		return nil, err
	}

	switch unary.Operator.Type {
	case scanner.BANG:
		return !i.isTruthy(obj), nil
	case scanner.MINUS:
		if i.isType(obj, reflect.Float64) {
			return obj.(float64), nil
		}
	}

	return nil, &Error{Token: unary.Operator, Message: "Operand must be a number."}
}

func (i *Interpreter) VisitBinaryExpr(binary parser.BinaryExpr) (any, error) {
	obj1, err := binary.Left.Accept(i)
	if err != nil {
		return nil, err
	}

	obj2, err := binary.Right.Accept(i)
	if err != nil {
		return nil, err
	}

	token := binary.Operator

	switch binary.Operator.Type {
	case scanner.PLUS:
		if i.areNumberedOperands(obj1, obj2) {
			return obj1.(float64) + obj2.(float64), nil
		}

		if i.areStringOperands(obj1, obj2) {
			return obj1.(string) + obj2.(string), nil
		}
		return nil, &Error{Token: token, Message: "Both operands should be numbers or strings."}
	case scanner.MINUS:
		if i.areNumberedOperands(obj1, obj2) {
			return obj1.(float64) - obj2.(float64), nil
		}
		return nil, &Error{Token: token, Message: "Both operands should be numbers."}
	case scanner.STAR:
		if i.areNumberedOperands(obj1, obj2) {
			return obj1.(float64) * obj2.(float64), nil
		}
		return nil, &Error{Token: token, Message: "Both operands should be numbers."}
	case scanner.SLASH:
		if i.areNumberedOperands(obj1, obj2) {
			left, _ := obj1.(float64)
			right, _ := obj2.(float64)

			if right == 0 {
				return nil, &Error{Token: token, Message: "Division by zero is prohibited."}
			}

			return left / right, nil
		}
		return nil, &Error{Token: token, Message: "Both operands should be numbers."}
	case scanner.GREATER:
		if i.areNumberedOperands(obj1, obj2) {
			return obj1.(float64) > obj2.(float64), nil
		}
		return nil, &Error{Token: token, Message: "Both operands should be numbers."}
	case scanner.GREATER_EQUAL:
		if i.areNumberedOperands(obj1, obj2) {
			return obj1.(float64) >= obj2.(float64), nil
		}
		return nil, &Error{Token: token, Message: "Both operands should be numbers."}
	case scanner.LESS:
		if i.areNumberedOperands(obj1, obj2) {
			return obj1.(float64) < obj2.(float64), nil
		}
		return nil, &Error{Token: token, Message: "Both operands should be numbers."}
	case scanner.LESS_EQUAL:
		if i.areNumberedOperands(obj1, obj2) {
			return obj1.(float64) < obj2.(float64), nil
		}
		return nil, &Error{Token: token, Message: "Both operands should be numbers."}
	case scanner.EQUAL_EQUAL:
		return i.areEqual(obj1, obj2), nil
	case scanner.BANG_EQUAL:
		return !i.areEqual(obj1, obj2), nil
	}

	panic(&Error{Token: token, Message: "Unreachable."})
}

func (i *Interpreter) VisitLogicalExpr(expr parser.LogicalExpr) (any, error) {
	left, err := expr.Left.Accept(i)
	if err != nil {
		return nil, err
	}

	if expr.Operator.Type == scanner.OR {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if !i.isTruthy(left) {
			return left, nil
		}
	}

	return expr.Right.Accept(i)
}

func (i *Interpreter) VisitAssignmentExpr(assignment parser.AssignmentExpr) (any, error) {
	token := assignment.Name
	if _, ok := i.environment.get(token.Lexeme); !ok {
		return nil, &Error{Token: token, Message: fmt.Sprintf("Undefined variable '%s'.", token.Lexeme)}
	}

	value, err := assignment.Value.Accept(i)
	if err != nil {
		return nil, err
	}

	i.environment.define(token.Lexeme, value)
	return value, nil
}

func (i *Interpreter) VisitTernaryExpr(ternary parser.TernaryExpr) (any, error) {
	obj, err := ternary.Condition.Accept(i)
	if err != nil {
		return nil, err
	}

	if i.isTruthy(obj) {
		return ternary.Left.Accept(i)
	}

	return ternary.Right.Accept(i)
}

func (i *Interpreter) VisitVariableExpr(variableExpr parser.VariableExpr) (any, error) {
	value, ok := i.environment.get(variableExpr.Name.Lexeme)
	if !ok {
		return nil, &Error{Token: variableExpr.Name, Message: fmt.Sprintf("Undefined variable '%s'.", variableExpr.Name.Lexeme)}
	}
	return value, nil
}

func (i *Interpreter) VisitExpressionStmt(expressionStmt parser.ExpressionStmt) (any, error) {
	_, err := expressionStmt.Expression.Accept(i)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Interpreter) VisitPrintStmt(printStmt parser.PrintStmt) (any, error) {
	result, err := printStmt.Expression.Accept(i)
	if err != nil {
		return nil, err
	}
	fmt.Println(i.stringify(result))
	return nil, nil
}

func (i *Interpreter) VisitVarStmt(varStmt parser.VarStmt) (any, error) {
	var value any = nil
	var err error = nil

	if varStmt.Initializer != nil {
		value, err = varStmt.Initializer.Accept(i)
		if err != nil {
			return nil, err
		}
	}

	i.environment.define(varStmt.Name.Lexeme, value)
	return nil, nil
}

func (i *Interpreter) VisitBlockStmt(stmt parser.BlockStmt) (any, error) {
	previous := i.environment
	i.environment = newEnvironment(i.environment)
	for _, declaration := range stmt.Declarations {
		_, err := declaration.Accept(i)
		if err != nil {
			i.environment = previous
			return nil, err
		}
	}
	i.environment = previous
	return nil, nil
}

func (i *Interpreter) VisitIfStmt(stmt parser.IfStmt) (any, error) {
	value, err := stmt.Expression.Accept(i)
	if err != nil {
		return nil, err
	}

	if i.isTruthy(value) {
		value, err = stmt.ThenBranch.Accept(i)
	} else if stmt.ElseBranch != nil {
		value, err = stmt.ElseBranch.Accept(i)
	}

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *Interpreter) Interpret(statements []parser.Stmt) error {
	for _, stmt := range statements {
		_, err := stmt.Accept(i)
		if err != nil {
			return err
		}
	}
	return nil
}

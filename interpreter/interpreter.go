package interpreter

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"glox/parser"
	"glox/scanner"
	"math"
	"reflect"
)

type Interpreter struct {
	globalEnvironment *environment
	environment       *environment
	locals            map[string]int32
}

func New() *Interpreter {
	globalEnv := newEnvironment(nil)
	env := globalEnv

	globalEnv.define("clock", &nativeClock{})
	globalEnv.define("str", &nativeStringify{})

	return &Interpreter{globalEnvironment: globalEnv, environment: env, locals: make(map[string]int32)}
}

func (i *Interpreter) newError(token scanner.Token, message string) *Error {
	return &Error{Token: token, Message: message}
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

func (i *Interpreter) evaluate(expr parser.Expr) (any, error) {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt parser.Stmt) (any, error) {
	return stmt.Accept(i)
}

func (i *Interpreter) executeBlock(statements []parser.Stmt, env *environment) (any, error) {
	previous := i.environment

	i.environment = env
	for _, stmt := range statements {
		if _, err := i.execute(stmt); err != nil {
			i.environment = previous
			return nil, err
		}
	}
	i.environment = previous

	return nil, nil
}

func (i *Interpreter) encodeExpression(expr parser.Expr) string {
	serializedData, err := json.Marshal(expr)
	if err != nil {
		panic("Could not encode expression.")
	}
	hasher := md5.New()
	hasher.Write(serializedData)

	return hex.EncodeToString(hasher.Sum(nil))
}

func (i *Interpreter) lookupVariable(expr parser.Expr, name scanner.Token) (any, bool) {
	depth, ok := i.locals[i.encodeExpression(expr)]
	if !ok {
		return i.globalEnvironment.get(name.Lexeme)
	}

	return i.environment.getAt(name.Lexeme, depth)
}

func (i *Interpreter) Resolve(expr parser.Expr, depth int32) {
	i.locals[i.encodeExpression(expr)] = depth
}

func (i *Interpreter) VisitTernaryExpr(ternary parser.TernaryExpr) (any, error) {
	obj, err := i.evaluate(ternary.Condition)
	if err != nil {
		return nil, err
	}

	if i.isTruthy(obj) {
		return i.evaluate(ternary.Left)
	}

	return i.evaluate(ternary.Right)
}

func (i *Interpreter) VisitAssignmentExpr(assignment parser.AssignmentExpr) (any, error) {
	token := assignment.Name
	value, err := i.evaluate(assignment.Value)
	if err != nil {
		return nil, err
	}

	assigned := false
	depth, ok := i.locals[i.encodeExpression(assignment)]

	if !ok {
		assigned = i.globalEnvironment.assign(token.Lexeme, value)
	} else {
		assigned = i.environment.assignAt(token.Lexeme, value, depth)
	}

	if !assigned {
		return nil, i.newError(token, fmt.Sprintf("Undefined variable '%s'.", token.Lexeme))
	}

	return value, nil
}

func (i *Interpreter) VisitLogicalExpr(expr parser.LogicalExpr) (any, error) {
	left, err := i.evaluate(expr.Left)
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

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitSetExpr(expr parser.SetExpr) (any, error) {
	object, err := i.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	instance, ok := object.(*loxInstance)
	if !ok {
		return nil, i.newError(expr.Name, "Only instances have properties.")
	}

	instance.set(expr.Name, value)

	return value, nil
}

func (i *Interpreter) VisitBinaryExpr(binary parser.BinaryExpr) (any, error) {
	obj1, err := i.evaluate(binary.Left)
	if err != nil {
		return nil, err
	}

	obj2, err := i.evaluate(binary.Right)
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
		return nil, i.newError(token, "Both operands should be numbers or strings.")
	case scanner.MINUS:
		if i.areNumberedOperands(obj1, obj2) {
			return obj1.(float64) - obj2.(float64), nil
		}
		return nil, i.newError(token, "Both operands should be numbers.")
	case scanner.STAR:
		if i.areNumberedOperands(obj1, obj2) {
			return obj1.(float64) * obj2.(float64), nil
		}
		return nil, i.newError(token, "Both operands should be numbers.")
	case scanner.SLASH:
		if i.areNumberedOperands(obj1, obj2) {
			left, _ := obj1.(float64)
			right, _ := obj2.(float64)

			if right == 0 {
				return nil, i.newError(token, "Division by zero is prohibited.")
			}

			return left / right, nil
		}
		return nil, i.newError(token, "Both operands should be numbers.")
	case scanner.MODULO:
		if i.areNumberedOperands(obj1, obj2) {
			return math.Mod(obj1.(float64), obj2.(float64)), nil
		}
		return nil, i.newError(token, "Both operands should be numbers.")
	case scanner.GREATER:
		if i.areNumberedOperands(obj1, obj2) {
			return obj1.(float64) > obj2.(float64), nil
		}
		return nil, i.newError(token, "Both operands should be numbers.")
	case scanner.GREATER_EQUAL:
		if i.areNumberedOperands(obj1, obj2) {
			return obj1.(float64) >= obj2.(float64), nil
		}
		return nil, i.newError(token, "Both operands should be numbers.")
	case scanner.LESS:
		if i.areNumberedOperands(obj1, obj2) {
			return obj1.(float64) < obj2.(float64), nil
		}
		return nil, i.newError(token, "Both operands should be numbers.")
	case scanner.LESS_EQUAL:
		if i.areNumberedOperands(obj1, obj2) {
			return obj1.(float64) <= obj2.(float64), nil
		}
		return nil, i.newError(token, "Both operands should be numbers.")
	case scanner.EQUAL_EQUAL:
		return i.areEqual(obj1, obj2), nil
	case scanner.BANG_EQUAL:
		return !i.areEqual(obj1, obj2), nil
	}

	panic(i.newError(token, "Unreachable."))
}

func (i *Interpreter) VisitGroupingExpr(grouping parser.GroupingExpr) (any, error) {
	return i.evaluate(grouping.Expr)
}

func (i *Interpreter) VisitLiteralExpr(literal parser.LiteralExpr) (any, error) {
	return literal.Value, nil
}

func (i *Interpreter) VisitUnaryExpr(unary parser.UnaryExpr) (any, error) {
	obj, err := i.evaluate(unary.Right)
	if err != nil {
		return nil, err
	}

	switch unary.Operator.Type {
	case scanner.BANG:
		return !i.isTruthy(obj), nil
	case scanner.MINUS:
		if i.isType(obj, reflect.Float64) {
			return -obj.(float64), nil
		}
	}
	return nil, i.newError(unary.Operator, "Operand must be a number.")
}

func (i *Interpreter) VisitGetExpr(expr parser.GetExpr) (any, error) {
	object, err := i.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}
	instance, ok := object.(*loxInstance)
	if !ok {
		return nil, i.newError(expr.Name, "Only instances have properties.")
	}

	value, err := instance.get(expr.Name)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (i *Interpreter) VisitCallExpr(expr parser.CallExpr) (any, error) {
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	fun, ok := callee.(callable)

	if !ok {
		return nil, i.newError(expr.Parenthesis, "Non callable object, can only call functions and classes.")
	}

	if fun.arity() != int32(len(expr.Arguments)) {
		return nil, i.newError(expr.Parenthesis, fmt.Sprintf("Expected %d arguments, but got %d.", fun.arity(), len(expr.Arguments)))
	}

	arguments := make([]any, 0)

	for _, arg := range expr.Arguments {
		value, err := i.evaluate(arg)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, value)
	}

	return fun.call(i, arguments)
}

func (i *Interpreter) VisitLambdaExpr(expr parser.LambdaExpr) (any, error) {
	name := scanner.Token{Type: scanner.IDENTIFIER, Lexeme: "lambda", Literal: nil, Line: expr.Parenthesis.Line}
	return newLoxFunction(parser.FunctionStmt{Name: name, Parameters: expr.Parameters, Body: expr.Body}, i.environment), nil
}

func (i *Interpreter) VisitVariableExpr(variableExpr parser.VariableExpr) (any, error) {
	var lookupExpr parser.Expr = variableExpr
	value, ok := i.lookupVariable(lookupExpr, variableExpr.Name)
	if !ok {
		return nil, i.newError(variableExpr.Name, fmt.Sprintf("Undefined variable '%s'.", variableExpr.Name.Lexeme))
	}
	return value, nil
}

func (i *Interpreter) VisitExpressionStmt(expressionStmt parser.ExpressionStmt) (any, error) {
	_, err := i.evaluate(expressionStmt.Expression)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Interpreter) VisitPrintStmt(stmt parser.PrintStmt) (any, error) {
	value, err := i.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}
	fmt.Println(i.stringify(value))
	return nil, nil
}

func (i *Interpreter) VisitVarStmt(varStmt parser.VarStmt) (any, error) {
	var value any = nil
	var err error = nil

	if varStmt.Initializer != nil {
		value, err = i.evaluate(varStmt.Initializer)
		if err != nil {
			return nil, err
		}
	}

	i.environment.define(varStmt.Name.Lexeme, value)
	return nil, nil
}

func (i *Interpreter) VisitClassStmt(stmt parser.ClassStmt) (any, error) {
	i.environment.define(stmt.Name.Lexeme, newClass(stmt))
	return nil, nil
}

func (i *Interpreter) VisitFunctionStmt(stmt parser.FunctionStmt) (any, error) {
	i.environment.define(stmt.Name.Lexeme, newLoxFunction(stmt, i.environment))
	return nil, nil
}

func (i *Interpreter) VisitBlockStmt(stmt parser.BlockStmt) (any, error) {
	return i.executeBlock(stmt.Declarations, newEnvironment(i.environment))
}

func (i *Interpreter) VisitIfStmt(stmt parser.IfStmt) (any, error) {
	value, err := i.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}

	if i.isTruthy(value) {
		value, err = i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		value, err = i.execute(stmt.ElseBranch)
	}

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *Interpreter) VisitWhileStmt(stmt parser.WhileStmt) (any, error) {
	condition, err := i.evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}

	for i.isTruthy(condition) {
		if _, err := i.execute(stmt.Body); err != nil {
			if errors.Is(err, &parser.BreakInterrupt{}) {
				break
			}
			if errors.Is(err, &parser.ContinueInterrupt{}) {
				continue
			}
			return nil, err
		}

		condition, _ = i.evaluate(stmt.Condition)
	}
	return nil, nil
}

func (i *Interpreter) VisitForStmt(stmt parser.ForStmt) (any, error) {
	initializer, forCondition := stmt.Initializer, stmt.Condition

	if initializer != nil {
		if _, err := i.execute(initializer); err != nil {
			return nil, err
		}
	}

	if forCondition == nil {
		forCondition = parser.LiteralExpr{Value: true}
	}

	condition, err := i.evaluate(forCondition)
	if err != nil {
		return nil, err
	}

	evaluateLoop := func() error {
		if stmt.Increment != nil {
			if _, err := i.execute(stmt.Increment); err != nil {
				return err
			}
		}
		condition, _ = i.evaluate(forCondition)
		return nil
	}

	for i.isTruthy(condition) {
		if _, err := i.execute(stmt.Body); err != nil {
			if errors.Is(err, &parser.BreakInterrupt{}) {
				break
			}
			if errors.Is(err, &parser.ContinueInterrupt{}) {
				if err := evaluateLoop(); err != nil {
					return nil, err
				}
				continue
			}
			return nil, err
		}

		if err := evaluateLoop(); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Interpreter) VisitBreakStmt(_ parser.BreakStmt) (any, error) {
	return nil, &parser.BreakInterrupt{}
}

func (i *Interpreter) VisitContinueStmt(_ parser.ContinueStmt) (any, error) {
	return nil, &parser.ContinueInterrupt{}
}

func (i *Interpreter) VisitReturnStmt(stmt parser.ReturnStmt) (any, error) {
	var value any = nil
	var err error = nil

	if stmt.Expr != nil {
		value, err = i.evaluate(stmt.Expr)

		if err != nil {
			return nil, err
		}
	}

	return nil, &parser.ReturnInterrupt{Value: value}
}

func (i *Interpreter) Interpret(statements []parser.Stmt) error {
	for _, stmt := range statements {
		if _, err := i.execute(stmt); err != nil {
			return err
		}
	}
	return nil
}

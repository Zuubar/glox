package resolver

import (
	"fmt"
	"glox/interpreter"
	"glox/parser"
	"glox/scanner"
)

type variable struct {
	token scanner.Token
	state string
}

type Resolver struct {
	interpreter     *interpreter.Interpreter
	scopes          []map[string]*variable
	warnings        []*Warning
	currentFunction string
	currentClass    string
	loopLevel       int
}

func New(interpreter *interpreter.Interpreter) *Resolver {
	return &Resolver{interpreter: interpreter, scopes: make([]map[string]*variable, 0), warnings: make([]*Warning, 0), currentFunction: functionTypeNone, currentClass: classTypeNone, loopLevel: 0}
}

func (r *Resolver) Resolve(stmt []parser.Stmt) (any, error) {
	return r.resolveStmts(stmt)
}

func (r *Resolver) Warnings() []error {
	warnings := make([]error, 0)
	for _, warning := range r.warnings {
		warnings = append(warnings, warning)
	}
	return warnings
}

func (r *Resolver) newError(token scanner.Token, message string) *Error {
	return &Error{Token: token, Message: message}
}

func (r *Resolver) newWarning(token scanner.Token, message string) {
	r.warnings = append(r.warnings, &Warning{Token: token, Message: message})
}

func (r *Resolver) declare(name scanner.Token) error {
	if len(r.scopes) == 0 {
		return nil
	}

	scope := *r.peekScope()

	if _, ok := scope[name.Lexeme]; ok {
		return r.newError(name, fmt.Sprintf("Redeclared '%s' variable in this scope.", name.Lexeme))
	}

	scope[name.Lexeme] = &variable{token: name, state: variableStateDeclared}

	return nil
}

func (r *Resolver) define(name scanner.Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope := *r.peekScope()
	scope[name.Lexeme].state = variableStateDefined
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, map[string]*variable{})
}

func (r *Resolver) endScope() {
	for name, details := range *r.peekScope() {
		if details.state != variableStateRead {
			r.newWarning(details.token, fmt.Sprintf("Unused variable '%s'.", name))
		}
	}

	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) peekScope() *map[string]*variable {
	return &r.scopes[len(r.scopes)-1]
}

func (r *Resolver) beginLoop() int {
	r.loopLevel += 1
	return r.loopLevel
}

func (r *Resolver) endLoop() int {
	r.loopLevel -= 1
	return r.loopLevel
}

func (r *Resolver) insideLoop() bool {
	return r.loopLevel > 0
}

func (r *Resolver) resolveExpr(expr parser.Expr) (any, error) {
	return expr.Accept(r)
}

func (r *Resolver) resolveStmt(stmt parser.Stmt) (any, error) {
	return stmt.Accept(r)
}

func (r *Resolver) resolveStmts(stmt []parser.Stmt) (any, error) {
	for _, stmt := range stmt {
		if _, err := r.resolveStmt(stmt); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) resolveFunctions(function any, functionType string) (any, error) {
	var parameters []scanner.Token
	var body []parser.Stmt

	if stmt, ok := function.(parser.FunctionStmt); ok {
		parameters, body = stmt.Parameters, stmt.Body
	} else if expr, ok := function.(parser.LambdaExpr); ok {
		parameters, body = expr.Parameters, expr.Body
	} else {
		panic("Invalid AST function node type. resolveFunctions only receives FunctionStmt or LambdaExpr.")
	}

	previousFunction := r.currentFunction

	r.currentFunction = functionType
	r.beginScope()

	defer func() {
		r.currentFunction = previousFunction
		r.endScope()
	}()

	if parameters != nil {
		for _, parameter := range parameters {
			if err := r.declare(parameter); err != nil {
				return nil, err
			}
			r.define(parameter)
		}
	}

	return r.resolveStmts(body)
}

func (r *Resolver) resolveLocal(expr parser.Expr, name scanner.Token, isRead bool) (any, error) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		scope := r.scopes[i]
		if variable, ok := scope[name.Lexeme]; ok {
			if isRead {
				variable.state = variableStateRead
			}
			r.interpreter.Resolve(expr, int32(len(r.scopes)-1-i))
			break
		}
	}

	return nil, nil
}

func (r *Resolver) VisitTernaryExpr(expr parser.TernaryExpr) (any, error) {
	if _, err := r.resolveExpr(expr.Condition); err != nil {
		return nil, err
	}

	if _, err := r.resolveExpr(expr.Left); err != nil {
		return nil, err
	}

	return r.resolveExpr(expr.Right)
}

func (r *Resolver) VisitAssignmentExpr(expr parser.AssignmentExpr) (any, error) {
	if _, err := r.resolveExpr(expr.Value); err != nil {
		return nil, err
	}

	return r.resolveLocal(expr, expr.Name, false)
}

func (r *Resolver) VisitLogicalExpr(expr parser.LogicalExpr) (any, error) {
	if _, err := r.resolveExpr(expr.Left); err != nil {
		return nil, err
	}

	return r.resolveExpr(expr.Right)
}

func (r *Resolver) VisitSetExpr(expr parser.SetExpr) (any, error) {
	if _, err := r.resolveExpr(expr.Object); err != nil {
		return nil, err
	}

	return r.resolveExpr(expr.Value)
}

func (r *Resolver) VisitBinaryExpr(expr parser.BinaryExpr) (any, error) {
	if _, err := r.resolveExpr(expr.Left); err != nil {
		return nil, err
	}

	return r.resolveExpr(expr.Right)
}

func (r *Resolver) VisitGroupingExpr(expr parser.GroupingExpr) (any, error) {
	return r.resolveExpr(expr.Expr)
}

func (r *Resolver) VisitLiteralExpr(_ parser.LiteralExpr) (any, error) {
	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(expr parser.UnaryExpr) (any, error) {
	return r.resolveExpr(expr.Right)
}

func (r *Resolver) VisitGetExpr(expr parser.GetExpr) (any, error) {
	return r.resolveExpr(expr.Object)
}

func (r *Resolver) VisitCallExpr(expr parser.CallExpr) (any, error) {
	if _, err := r.resolveExpr(expr.Callee); err != nil {
		return nil, err
	}

	for _, argument := range expr.Arguments {
		if _, err := r.resolveExpr(argument); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) VisitLambdaExpr(expr parser.LambdaExpr) (any, error) {
	return r.resolveFunctions(expr, functionTypeFunction)
}

func (r *Resolver) VisitThisExpr(expr parser.ThisExpr) (any, error) {
	if r.currentClass == classTypeNone {
		return nil, r.newError(expr.Keyword, "Can't use 'this' outside of class.")
	}

	if r.currentFunction == functionTypeStaticMethod {
		return nil, r.newError(expr.Keyword, "Can't use 'this' inside static method, consider using 'className.property'.")
	}

	return r.resolveLocal(expr, expr.Keyword, true)
}

func (r *Resolver) VisitVariableExpr(expr parser.VariableExpr) (any, error) {
	varName := expr.Name.Lexeme
	if len(r.scopes) != 0 {
		v, ok := (*r.peekScope())[varName]

		if ok && v.state == variableStateDeclared {
			return nil, r.newError(expr.Name, fmt.Sprintf("Can't read local variable '%s' in it's own initializer", varName))
		}
	}

	return r.resolveLocal(expr, expr.Name, true)
}

func (r *Resolver) VisitExpressionStmt(stmt parser.ExpressionStmt) (any, error) {
	return r.resolveExpr(stmt.Expression)
}

func (r *Resolver) VisitPrintStmt(stmt parser.PrintStmt) (any, error) {
	return r.resolveExpr(stmt.Expression)
}

func (r *Resolver) VisitVarStmt(stmt parser.VarStmt) (any, error) {
	if err := r.declare(stmt.Name); err != nil {
		return nil, err
	}
	if stmt.Initializer != nil {
		if _, err := r.resolveExpr(stmt.Initializer); err != nil {
			return nil, err
		}
	}
	r.define(stmt.Name)

	return nil, nil
}

func (r *Resolver) VisitClassStmt(stmt parser.ClassStmt) (any, error) {
	currentClass := r.currentClass

	r.currentClass = classTypeClass
	if err := r.declare(stmt.Name); err != nil {
		return nil, err
	}
	r.define(stmt.Name)

	r.beginScope()
	defer func() {
		r.endScope()
		r.currentClass = currentClass
	}()

	lastScope := *r.peekScope()
	lastScope["this"] = &variable{state: variableStateRead}

	for _, method := range stmt.Methods {
		funcType := functionTypeMethod
		if method.Name.Lexeme == "init" {
			funcType = functionTypeClassInitializer
		}

		if _, err := r.resolveFunctions(method, funcType); err != nil {
			return nil, err
		}
	}

	for _, method := range stmt.StaticMethods {
		if _, err := r.resolveFunctions(method, functionTypeStaticMethod); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) VisitFunctionStmt(stmt parser.FunctionStmt) (any, error) {
	if err := r.declare(stmt.Name); err != nil {
		return nil, err
	}
	r.define(stmt.Name)

	return r.resolveFunctions(stmt, functionTypeFunction)
}

func (r *Resolver) VisitBlockStmt(stmt parser.BlockStmt) (any, error) {
	r.beginScope()
	defer r.endScope()
	if _, err := r.resolveStmts(stmt.Declarations); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) VisitIfStmt(stmt parser.IfStmt) (any, error) {
	if _, err := r.resolveExpr(stmt.Expression); err != nil {
		return nil, err
	}

	if _, err := r.resolveStmt(stmt.ThenBranch); err != nil {
		return nil, err
	}

	if stmt.ElseBranch != nil {
		return r.resolveStmt(stmt.ElseBranch)
	}

	return nil, nil
}

func (r *Resolver) VisitWhileStmt(stmt parser.WhileStmt) (any, error) {
	r.beginLoop()
	defer r.endLoop()

	if _, err := r.resolveExpr(stmt.Condition); err != nil {
		return nil, err
	}

	return r.resolveStmt(stmt.Body)
}

func (r *Resolver) VisitForStmt(stmt parser.ForStmt) (any, error) {
	r.beginLoop()
	defer r.endLoop()

	if stmt.Initializer != nil {
		if _, err := r.resolveStmt(stmt.Initializer); err != nil {
			return nil, err
		}
	}

	if stmt.Condition != nil {
		if _, err := r.resolveExpr(stmt.Condition); err != nil {
			return nil, err
		}
	}

	if stmt.Increment != nil {
		if _, err := r.resolveStmt(stmt.Increment); err != nil {
			return nil, err
		}
	}

	return r.resolveStmt(stmt.Body)
}

func (r *Resolver) VisitBreakStmt(stmt parser.BreakStmt) (any, error) {
	if !r.insideLoop() {
		return nil, r.newError(stmt.Keyword, "Unexpected 'break' outside of loop.")
	}
	return nil, nil
}

func (r *Resolver) VisitContinueStmt(stmt parser.ContinueStmt) (any, error) {
	if !r.insideLoop() {
		return nil, r.newError(stmt.Keyword, "Unexpected 'continue' outside of loop.")
	}

	return nil, nil
}

func (r *Resolver) VisitReturnStmt(stmt parser.ReturnStmt) (any, error) {
	if r.currentFunction == functionTypeNone {
		return nil, r.newError(stmt.Keyword, "Can't return from top-level code.")
	}

	if stmt.Expr != nil {
		if r.currentFunction == functionTypeClassInitializer {
			return nil, r.newError(stmt.Keyword, "Can't return a value from class initializer.")
		}

		return r.resolveExpr(stmt.Expr)
	}

	return nil, nil
}

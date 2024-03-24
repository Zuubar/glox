package parser

import (
	"fmt"
	"glox/scanner"
)

type Parser struct {
	tokens  []scanner.Token
	current int32
}

func New(tokens []scanner.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() ([]Stmt, []error) {
	statements, errors := make([]Stmt, 0, 10), make([]error, 0, 10)
	for !p.isAtEnd() {
		stmt, err := p.declaration()

		if err != nil {
			p.synchronize()
			errors = append(errors, err)
		}

		statements = append(statements, stmt)
	}
	return statements, errors
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == scanner.EOF
}

func (p *Parser) check(tokenType scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) advance() scanner.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.peekBehind()
}

func (p *Parser) peek() scanner.Token {
	return p.tokens[p.current]
}

func (p *Parser) peekBehind() scanner.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) match(types ...scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	for _, token := range types {
		if p.check(token) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) parseBinaryExpr(nonTerminal func() (Expr, error), types ...scanner.TokenType) (Expr, error) {
	expr, err := nonTerminal()
	if err != nil {
		return nil, err
	}

	for p.match(types...) {
		token := p.peekBehind()
		right, err := nonTerminal()

		if err != nil {
			return nil, err
		}

		expr = BinaryExpr{Left: expr, Operator: token, Right: right}
	}

	return expr, nil
}

func (p *Parser) parseLogicalExpr(nonTerminal func() (Expr, error), types ...scanner.TokenType) (Expr, error) {
	expr, err := nonTerminal()
	if err != nil {
		return nil, err
	}

	for p.match(types...) {
		token := p.peekBehind()
		right, err := nonTerminal()

		if err != nil {
			return nil, err
		}

		expr = LogicalExpr{Left: expr, Operator: token, Right: right}
	}

	return expr, nil
}

func (p *Parser) newError(token scanner.Token, message string) error {
	if token.Type == scanner.EOF {
		return &Error{Line: token.Line, Where: " at the end", Message: message}
	}
	return &Error{Line: token.Line, Where: fmt.Sprintf(" at '%s'", token.Lexeme), Message: message}
}

func (p *Parser) consume(tokenType scanner.TokenType, errorMsg string) (scanner.Token, error) {
	if !p.match(tokenType) {
		return scanner.Token{}, p.newError(p.peek(), errorMsg)
	}
	return p.peekBehind(), nil
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.peekBehind().Type == scanner.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case scanner.CLASS:
		case scanner.TRAIT:
		case scanner.FUN:
		case scanner.VAR:
		case scanner.FOR:
		case scanner.IF:
		case scanner.WHILE:
		case scanner.BREAK:
		case scanner.CONTINUE:
		case scanner.RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(scanner.VAR) {
		return p.varDecl()
	}
	if p.match(scanner.CLASS) {
		return p.classDecl()
	}
	if p.match(scanner.TRAIT) {
		return p.traitDecl()
	}
	if p.match(scanner.FUN) {
		return p.functionDecl("function")
	}

	return p.statement()
}

func (p *Parser) varDecl() (Stmt, error) {
	name, err := p.consume(scanner.IDENTIFIER, "Expected identifier after 'var'.")
	if err != nil {
		return nil, err
	}
	varDecl := VarStmt{Name: name, Initializer: nil}

	if p.match(scanner.EQUAL) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		varDecl.Initializer = expr
	}

	if _, err := p.consume(scanner.SEMICOLON, "Expected ';' after a variable declaration."); err != nil {
		return nil, err
	}
	return varDecl, nil
}

func (p *Parser) methodDecl() (Stmt, error) {
	name, err := p.consume(scanner.IDENTIFIER, fmt.Sprintf("Expteced method or getter name."))
	if err != nil {
		return nil, err
	}

	var parameters []scanner.Token = nil

	if p.match(scanner.LEFT_PAREN) {
		parameters = make([]scanner.Token, 0)
		if !p.check(scanner.RIGHT_PAREN) {
			parameters = append(parameters, p.advance())

			for p.match(scanner.COMMA) {
				if len(parameters) >= 255 {
					return nil, p.newError(p.peek(), "Can't have more than 255 parameters.")
				}

				parameters = append(parameters, p.advance())
			}
		}

		if _, err := p.consume(scanner.RIGHT_PAREN, "Expected ')' after parameters."); err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(scanner.LEFT_BRACE, fmt.Sprintf("Expected '{' before method or getter body.")); err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return FunctionStmt{Name: name, Parameters: parameters, Body: body.(BlockStmt).Declarations}, nil
}

func (p *Parser) classMethods() ([]FunctionStmt, []FunctionStmt, error) {
	methods, staticMethods := make([]FunctionStmt, 0), make([]FunctionStmt, 0)
	for !p.check(scanner.RIGHT_BRACE) && !p.isAtEnd() {
		parsingStaticMethod := p.match(scanner.CLASS)

		methodDecl, err := p.methodDecl()

		if err != nil {
			return nil, nil, err
		}

		method := methodDecl.(FunctionStmt)

		if parsingStaticMethod {
			staticMethods = append(staticMethods, method)
		} else {
			methods = append(methods, method)
		}
	}

	return methods, staticMethods, nil
}

func (p *Parser) classDecl() (Stmt, error) {
	name, err := p.consume(scanner.IDENTIFIER, "Expected class name.")

	if err != nil {
		return nil, err
	}

	var superclass VariableExpr

	if p.match(scanner.LESS) {
		name, err := p.consume(scanner.IDENTIFIER, "Excepted superclass name.")
		if err != nil {
			return nil, nil
		}

		superclass.Name = name
	}

	traits := make([]VariableExpr, 0)

	if p.match(scanner.USE_TRAIT) {
		traitName, err := p.consume(scanner.IDENTIFIER, "Excepted trait name.")
		if err != nil {
			return nil, nil
		}

		traits = append(traits, VariableExpr{Name: traitName})

		for p.match(scanner.COMMA) {
			traitName, err := p.consume(scanner.IDENTIFIER, "Excepted trait name.")
			if err != nil {
				return nil, err
			}

			traits = append(traits, VariableExpr{Name: traitName})
		}
	}

	if _, err := p.consume(scanner.LEFT_BRACE, "Expected '{' before class body."); err != nil {
		return nil, err
	}

	methods, staticMethods, err := p.classMethods()

	if err != nil {
		return nil, err
	}

	if _, err := p.consume(scanner.RIGHT_BRACE, "Expected '}' after class body."); err != nil {
		return nil, err
	}

	return ClassStmt{
		Name:          name,
		Superclass:    superclass,
		Traits:        traits,
		Methods:       methods,
		StaticMethods: staticMethods,
	}, nil
}

func (p *Parser) traitDecl() (Stmt, error) {
	name, err := p.consume(scanner.IDENTIFIER, "Expected trait name.")

	if err != nil {
		return nil, err
	}

	if _, err := p.consume(scanner.LEFT_BRACE, "Expected '{' before trait body."); err != nil {
		return nil, err
	}

	methods, staticMethods, err := p.classMethods()

	if _, err := p.consume(scanner.RIGHT_BRACE, "Expected '}' after trait body."); err != nil {
		return nil, err
	}

	return TraitStmt{
		Name:          name,
		Methods:       methods,
		StaticMethods: staticMethods,
	}, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(scanner.LEFT_BRACE) {
		return p.block()
	}
	if p.match(scanner.PRINT) {
		return p.printStmt()
	}
	if p.match(scanner.IF) {
		return p.ifStmt()
	}
	if p.match(scanner.WHILE) {
		return p.whileStmt()
	}
	if p.match(scanner.FOR) {
		return p.forStmt()
	}
	if p.match(scanner.BREAK, scanner.CONTINUE) {
		return p.loopInterruptStmts()
	}
	if p.match(scanner.RETURN) {
		return p.returnStmt()
	}

	return p.expressionStmt()
}

func (p *Parser) functionDecl(kind string) (Stmt, error) {
	name, err := p.consume(scanner.IDENTIFIER, fmt.Sprintf("Expteced %s name.", kind))
	if err != nil {
		return nil, err
	}
	parameters := make([]scanner.Token, 0)

	if _, err := p.consume(scanner.LEFT_PAREN, fmt.Sprintf("Expteced '(' after %s name.", kind)); err != nil {
		return nil, err
	}

	if !p.check(scanner.RIGHT_PAREN) {
		parameters = append(parameters, p.advance())

		for p.match(scanner.COMMA) {
			if len(parameters) >= 255 {
				return nil, p.newError(p.peek(), "Can't have more than 255 parameters.")
			}

			parameters = append(parameters, p.advance())
		}
	}

	if _, err := p.consume(scanner.RIGHT_PAREN, "Expected ')' after parameters."); err != nil {
		return nil, err
	}

	if _, err := p.consume(scanner.LEFT_BRACE, fmt.Sprintf("Expected '{' before %s body.", kind)); err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return FunctionStmt{Name: name, Parameters: parameters, Body: body.(BlockStmt).Declarations}, nil
}

func (p *Parser) expressionStmt() (Stmt, error) {
	expr, err := p.expression()

	if err != nil {
		return nil, err
	}

	if _, err := p.consume(scanner.SEMICOLON, "Expected ';' after a value."); err != nil {
		return nil, err
	}

	return ExpressionStmt{Expression: expr}, nil
}

func (p *Parser) printStmt() (Stmt, error) {
	expr, err := p.expression()

	if err != nil {
		return nil, err
	}

	if _, err := p.consume(scanner.SEMICOLON, "Expected ';' after a value."); err != nil {
		return nil, err
	}

	return PrintStmt{Expression: expr}, nil
}

func (p *Parser) block() (Stmt, error) {
	declarations := make([]Stmt, 0, 10)
	for !p.check(scanner.RIGHT_BRACE) && !p.isAtEnd() {
		declaration, err := p.declaration()
		if err != nil {
			return nil, err
		}
		declarations = append(declarations, declaration)
	}

	if _, err := p.consume(scanner.RIGHT_BRACE, "Expected '}' after a block."); err != nil {
		return nil, err
	}

	return BlockStmt{Declarations: declarations}, nil
}

func (p *Parser) ifStmt() (Stmt, error) {
	if _, err := p.consume(scanner.LEFT_PAREN, "Expected '(' after 'if'."); err != nil {
		return nil, err
	}

	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(scanner.RIGHT_PAREN, "Expected ')' at the end of 'if'."); err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch Stmt = nil
	if p.match(scanner.ELSE) {
		elseBranch, err = p.statement()

		if err != nil {
			return nil, err
		}
	}

	return IfStmt{Expression: expr, ThenBranch: thenBranch, ElseBranch: elseBranch}, nil
}

func (p *Parser) whileStmt() (Stmt, error) {
	if _, err := p.consume(scanner.LEFT_PAREN, "Expected '(' after 'while'."); err != nil {
		return nil, err
	}

	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(scanner.RIGHT_PAREN, "Expected ')' at the end of 'while'."); err != nil {
		return nil, err
	}

	stmt, err := p.statement()
	if err != nil {
		return nil, err
	}

	return WhileStmt{Condition: expr, Body: stmt}, nil
}

func (p *Parser) forStmt() (Stmt, error) {
	if _, err := p.consume(scanner.LEFT_PAREN, "Expected '(' after 'for'."); err != nil {
		return nil, err
	}

	var initializer Stmt
	var err error = nil

	if p.match(scanner.SEMICOLON) {
		initializer = nil
	} else if p.match(scanner.VAR) {
		initializer, err = p.varDecl()
	} else {
		initializer, err = p.expressionStmt()
	}

	if err != nil {
		return nil, err
	}

	var condition Expr = nil

	if !p.match(scanner.SEMICOLON) {
		condition, err = p.expression()

		if err != nil {
			return nil, err
		}

		if _, err := p.consume(scanner.SEMICOLON, "Expected ';' after condition."); err != nil {
			return nil, err
		}
	}

	var increment Stmt = nil

	if !p.match(scanner.RIGHT_PAREN) {
		incrementExpr, err := p.expression()

		if err != nil {
			return nil, err
		}

		increment = ExpressionStmt{Expression: incrementExpr}

		if _, err := p.consume(scanner.RIGHT_PAREN, "Expected ')' at the end of 'for'."); err != nil {
			return nil, err
		}
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return BlockStmt{Declarations: []Stmt{ForStmt{Initializer: initializer, Condition: condition, Increment: increment, Body: body}}}, nil

	// Desugar into a while.

	//whileBody := []Stmt{body}
	//if increment != nil {
	//	whileBody = append(whileBody, ExpressionStmt{Expression: increment})
	//}
	//
	//whileStmt := WhileStmt{Condition: condition, Body: BlockStmt{Declarations: whileBody}}
	//
	//whileBlock := BlockStmt{Declarations: make([]Stmt, 0)}
	//if initializer != nil {
	//	whileBlock.Declarations = append(whileBlock.Declarations, initializer)
	//}
	//whileBlock.Declarations = append(whileBlock.Declarations, whileStmt)
	//
	//return whileBlock, nil
}

func (p *Parser) loopInterruptStmts() (Stmt, error) {
	keyword := p.peekBehind()

	if _, err := p.consume(scanner.SEMICOLON, fmt.Sprintf("Expected ';' after a '%s'.", keyword.Lexeme)); err != nil {
		return nil, err
	}
	if keyword.Type == scanner.BREAK {
		return BreakStmt{Keyword: keyword}, nil
	}
	return ContinueStmt{Keyword: keyword}, nil
}

func (p *Parser) returnStmt() (Stmt, error) {
	keyword := p.peekBehind()

	var expr Expr = nil
	var err error = nil
	if !p.check(scanner.SEMICOLON) {
		if expr, err = p.expression(); err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(scanner.SEMICOLON, fmt.Sprintf("Expected ';' after a '%s'.", keyword.Lexeme)); err != nil {
		return nil, err
	}

	return ReturnStmt{Keyword: keyword, Expr: expr}, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.ternary()
}

func (p *Parser) ternary() (Expr, error) {
	condition, err := p.assignment()
	if err != nil {
		return nil, err
	}

	if !p.match(scanner.QUESTION) {
		return condition, err
	}

	left, err := p.ternary()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(scanner.COLON, "Expected ':' after '?'."); err != nil {
		return nil, err
	}

	right, err := p.ternary()
	if err != nil {
		return nil, err
	}

	return TernaryExpr{Condition: condition, Left: left, Right: right}, nil
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.logicalOr()
	if err != nil {
		return nil, err
	}

	if p.match(scanner.EQUAL) {
		equals := p.peekBehind()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		switch t := expr.(type) {
		case VariableExpr:
			return AssignmentExpr{Name: t.Name, Value: value}, nil
		case GetExpr:
			return SetExpr{Object: t.Object, Name: t.Name, Value: value}, nil
		case ArrayGetExpr:
			return ArraySetExpr{Array: t.Array, Bracket: t.Bracket, Index: t.Index, Value: value}, nil
		default:
			return nil, p.newError(equals, "Invalid assignment target.")
		}
	}

	return expr, nil
}

func (p *Parser) logicalOr() (Expr, error) {
	return p.parseLogicalExpr(p.logicalAnd, scanner.OR)
}

func (p *Parser) logicalAnd() (Expr, error) {
	return p.parseLogicalExpr(p.equality, scanner.AND)
}

func (p *Parser) equality() (Expr, error) {
	return p.parseBinaryExpr(p.comparison, scanner.BANG_EQUAL, scanner.EQUAL_EQUAL)
}

func (p *Parser) comparison() (Expr, error) {
	return p.parseBinaryExpr(p.modulo, scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL)
}

func (p *Parser) modulo() (Expr, error) {
	return p.parseBinaryExpr(p.term, scanner.MODULO)
}

func (p *Parser) term() (Expr, error) {
	return p.parseBinaryExpr(p.factor, scanner.PLUS, scanner.MINUS)
}

func (p *Parser) factor() (Expr, error) {
	return p.parseBinaryExpr(p.unary, scanner.STAR, scanner.SLASH)
}

func (p *Parser) unary() (Expr, error) {
	if !p.match(scanner.BANG, scanner.MINUS) {
		return p.call()
	}
	token := p.peekBehind()
	right, err := p.unary()

	if err != nil {
		return nil, err
	}

	return UnaryExpr{Operator: token, Right: right}, nil
}

func (p *Parser) finishCall(callee Expr) (Expr, error) {
	arguments := make([]Expr, 0)

	if p.match(scanner.RIGHT_PAREN) {
		return CallExpr{Callee: callee, Parenthesis: p.peekBehind(), Arguments: arguments}, nil
	}

	firstArg, err := p.expression()
	if err != nil {
		return nil, err
	}

	arguments = append(arguments, firstArg)

	for p.match(scanner.COMMA) {
		if len(arguments) >= 255 {
			return nil, p.newError(p.peek(), "Can't have more than 255 arguments.")
		}

		arg, err := p.expression()
		if err != nil {
			return nil, err
		}

		arguments = append(arguments, arg)
	}

	if _, err := p.consume(scanner.RIGHT_PAREN, "Expected ')' after arguments."); err != nil {
		return nil, err
	}

	return CallExpr{Callee: callee, Parenthesis: p.peekBehind(), Arguments: arguments}, nil
}

func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()

	if err != nil {
		return nil, err
	}

	for {
		if p.match(scanner.LEFT_BRACKET) {
			bracket := p.peekBehind()

			index, err := p.primary()
			if err != nil {
				return nil, err
			}

			if _, err := p.consume(scanner.RIGHT_BRACKET, "Excepted closing ']'."); err != nil {
				return nil, err
			}

			expr = ArrayGetExpr{Array: expr, Bracket: bracket, Index: index}

		} else if p.match(scanner.LEFT_PAREN) {
			if expr, err = p.finishCall(expr); err != nil {
				return nil, err
			}

		} else if p.match(scanner.DOT) {
			name, err := p.consume(scanner.IDENTIFIER, "Expected property name after '.'.")

			if err != nil {
				return nil, err
			}

			expr = GetExpr{Object: expr, Name: name}

		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) lambda() (Expr, error) {
	if _, err := p.consume(scanner.LEFT_PAREN, "Expected '(' before anonymous function parameters"); err != nil {
		return nil, err
	}

	parameters := make([]scanner.Token, 0)
	var parenthesisToken scanner.Token
	var err error

	if p.match(scanner.RIGHT_PAREN) {
		parenthesisToken = p.peekBehind()
	} else {
		firstParam := p.advance()
		parameters = append(parameters, firstParam)

		for p.match(scanner.COMMA) {
			parameters = append(parameters, p.advance())
		}

		parenthesisToken, err = p.consume(scanner.RIGHT_PAREN, "Expected ')' after anonymous function parameters.")
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(scanner.LEFT_BRACE, "Expected '{' before anonymous function body."); err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return LambdaExpr{Parenthesis: parenthesisToken, Parameters: parameters, Body: body.(BlockStmt).Declarations}, nil
}

func (p *Parser) array() (Expr, error) {
	elements := make([]Expr, 0)

	if p.match(scanner.RIGHT_BRACKET) {
		return ArrayExpr{Elements: elements, Bracket: p.peekBehind()}, nil
	}

	firstElem, err := p.expression()
	if err != nil {
		return nil, err
	}
	elements = append(elements, firstElem)

	for p.match(scanner.COMMA) {
		elem, err := p.expression()
		if err != nil {
			return nil, err
		}
		elements = append(elements, elem)
	}

	if _, err := p.consume(scanner.RIGHT_BRACKET, "Expected closing ']' for arrays."); err != nil {
		return nil, err
	}

	return ArrayExpr{Elements: elements, Bracket: p.peekBehind()}, nil
}

func (p *Parser) primary() (Expr, error) {
	if p.match(scanner.TRUE) {
		return LiteralExpr{Value: true}, nil
	}

	if p.match(scanner.FALSE) {
		return LiteralExpr{Value: false}, nil
	}

	if p.match(scanner.NIL) {
		return LiteralExpr{Value: nil}, nil
	}

	if p.match(scanner.NUMBER, scanner.STRING) {
		return LiteralExpr{Value: p.peekBehind().Literal}, nil
	}

	if p.match(scanner.THIS) {
		return ThisExpr{Keyword: p.peekBehind()}, nil
	}

	if p.match(scanner.IDENTIFIER) {
		return VariableExpr{Name: p.peekBehind()}, nil
	}

	if p.match(scanner.FUN) {
		return p.lambda()
	}

	if p.match(scanner.SUPER) {
		if _, err := p.consume(scanner.DOT, "Expected '.' after 'super'."); err != nil {
			return nil, err
		}

		method, err := p.consume(scanner.IDENTIFIER, "Expected superclass method name.")
		if err != nil {
			return nil, err
		}

		return SuperExpr{Keyword: p.peekBehind(), Method: method}, nil
	}

	if p.match(scanner.LEFT_BRACKET) {
		return p.array()
	}

	if p.match(scanner.LEFT_PAREN) {
		expr, err := p.expression()

		if err != nil {
			return nil, err
		}

		if _, err := p.consume(scanner.RIGHT_PAREN, "Expect ')' after expression."); err != nil {
			return nil, err
		}

		return GroupingExpr{Expr: expr}, nil
	}

	return nil, p.newError(p.peek(), "Expected an expression.")
}

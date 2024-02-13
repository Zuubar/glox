package parser

import (
	"fmt"
	"glox/scanner"
)

type Parser struct {
	tokens    []scanner.Token
	current   int32
	loopLevel int32
}

func New(tokens []scanner.Token) *Parser {
	return &Parser{tokens: tokens, loopLevel: 0}
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

func (p *Parser) generateError(token scanner.Token, message string) error {
	if token.Type == scanner.EOF {
		return &Error{Line: token.Line, Where: " at the end", Message: message}
	}
	return &Error{Line: token.Line, Where: fmt.Sprintf(" at '%s'", token.Lexeme), Message: message}
}

func (p *Parser) consume(tokenType scanner.TokenType, errorMsg string) error {
	if !p.match(tokenType) {
		return p.generateError(p.peek(), errorMsg)
	}
	return nil
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.peekBehind().Type == scanner.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case scanner.CLASS:
		case scanner.FUN:
		case scanner.VAR:
		case scanner.FOR:
		case scanner.IF:
		case scanner.WHILE:
		case scanner.BREAK:
		case scanner.PRINT:
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

	return p.statement()
}

func (p *Parser) varDecl() (Stmt, error) {
	if err := p.consume(scanner.IDENTIFIER, "Expected identifier after 'var'."); err != nil {
		return nil, err
	}
	varDecl := VarStmt{Name: p.peekBehind(), Initializer: nil}

	if p.match(scanner.EQUAL) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		varDecl.Initializer = expr
	}

	if err := p.consume(scanner.SEMICOLON, "Expected ';' after a variable declaration."); err != nil {
		return nil, err
	}
	return varDecl, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(scanner.PRINT) {
		return p.printStmt()
	}
	if p.match(scanner.LEFT_BRACE) {
		return p.block()
	}
	if p.match(scanner.IF) {
		return p.ifStmt()
	}
	if p.match(scanner.WHILE) {
		p.loopLevel += 1
		return p.whileStmt()
	}
	if p.match(scanner.FOR) {
		p.loopLevel += 1
		return p.forStmt()
	}
	if p.match(scanner.BREAK) {
		at := p.peekBehind()
		if p.loopLevel > 0 {
			if err := p.consume(scanner.SEMICOLON, "Expected ';' after a 'break'."); err != nil {
				return nil, err
			}
			return BreakStmt{At: at}, nil
		}
		return nil, p.generateError(p.peekBehind(), "Unexpected 'break' outside of while|for loop")
	}

	return p.expressionStmt()
}

func (p *Parser) expressionStmt() (Stmt, error) {
	expr, err := p.expression()

	if err != nil {
		return nil, err
	}

	if err := p.consume(scanner.SEMICOLON, "Expected ';' after a value."); err != nil {
		return nil, err
	}

	return ExpressionStmt{Expression: expr}, nil
}

func (p *Parser) printStmt() (Stmt, error) {
	expr, err := p.expression()

	if err != nil {
		return nil, err
	}

	if err := p.consume(scanner.SEMICOLON, "Expected ';' after a value."); err != nil {
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

	if err := p.consume(scanner.RIGHT_BRACE, "Expected '}' after a block."); err != nil {
		return nil, err
	}

	return BlockStmt{Declarations: declarations}, nil
}

func (p *Parser) ifStmt() (Stmt, error) {
	if err := p.consume(scanner.LEFT_PAREN, "Expected '(' after 'if'."); err != nil {
		return nil, err
	}

	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if err := p.consume(scanner.RIGHT_PAREN, "Expected ')' at the end of 'if'."); err != nil {
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
	defer func() { p.loopLevel -= 1 }()
	if err := p.consume(scanner.LEFT_PAREN, "Expected '(' after 'while'."); err != nil {
		return nil, err
	}

	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if err := p.consume(scanner.RIGHT_PAREN, "Expected ')' at the end of 'while'."); err != nil {
		return nil, err
	}

	stmt, err := p.statement()
	if err != nil {
		return nil, err
	}

	return WhileStmt{Condition: expr, Body: stmt}, nil
}

func (p *Parser) forStmt() (Stmt, error) {
	defer func() { p.loopLevel -= 1 }()
	if err := p.consume(scanner.LEFT_PAREN, "Expected '(' after 'for'."); err != nil {
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

	var condition Expr
	err = nil

	if p.match(scanner.SEMICOLON) {
		condition = LiteralExpr{Value: true}
	} else {
		condition, err = p.expression()

		if err := p.consume(scanner.SEMICOLON, "Expected ';' after condition."); err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	var increment Expr
	err = nil

	if p.match(scanner.RIGHT_PAREN) {
		increment = nil
	} else {
		increment, err = p.expression()

		if err := p.consume(scanner.RIGHT_PAREN, "Expected ')' at the end of 'for'."); err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	// Start of desugaring
	whileBody := []Stmt{body}
	if increment != nil {
		whileBody = append(whileBody, ExpressionStmt{Expression: increment})
	}

	whileStmt := WhileStmt{Condition: condition, Body: BlockStmt{Declarations: whileBody}}

	whileBlock := BlockStmt{Declarations: make([]Stmt, 0)}
	if initializer != nil {
		whileBlock.Declarations = append(whileBlock.Declarations, initializer)
	}
	whileBlock.Declarations = append(whileBlock.Declarations, whileStmt)

	return whileBlock, nil
}

func (p *Parser) expression() (Expr, error) {
	expr, err := p.ternary()

	if err != nil {
		return nil, err
	}

	for p.match(scanner.COMMA) {
		expr, err = p.ternary()

		if err != nil {
			return nil, err
		}
	}

	return expr, nil
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

	if err := p.consume(scanner.COLON, "Expected ':' after '?'."); err != nil {
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
		if t, ok := expr.(VariableExpr); ok {
			return AssignmentExpr{Name: t.Name, Value: value}, nil
		}
		return nil, p.generateError(equals, "Invalid assignment target.")
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
		return p.primary()
	}
	token := p.peekBehind()
	right, err := p.unary()

	if err != nil {
		return nil, err
	}

	return UnaryExpr{Operator: token, Right: right}, nil
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

	if p.match(scanner.IDENTIFIER) {
		return VariableExpr{Name: p.peekBehind()}, nil
	}

	if p.match(scanner.LEFT_PAREN) {
		expr, err := p.expression()

		if err != nil {
			return nil, err
		}

		if err := p.consume(scanner.RIGHT_PAREN, "Expect ')' after expression."); err != nil {
			return nil, err
		}

		return GroupingExpr{Expr: expr}, nil
	}

	return nil, p.generateError(p.peek(), "Expected an expression.")
}

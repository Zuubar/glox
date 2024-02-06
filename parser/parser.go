package parser

import (
	"fmt"
	"glox/scanner"
)

/*
Grammar:
	program -> declaration* EOF
	declaration -> varDecl | statement ;
	varDecl -> "var" IDENTIFIER ( "=" expression ";" )? ;
	statement -> expressionStmt | printStmt ;
	expressionStmt -> expression ";" ;
	printStmt -> "print" expression ";" ;

	expression -> ternary ("," ternary)* ;
	ternary -> equality "?" ternary ":" ternary | equality ;
	equality -> comparison ( ( "!=" | "==" ) comparison)* ;
	comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term)* ;
	term -> factor ( ( "+" | "-" ) factor)* ;
	factor -> unary ( ( "*" | "/" ) unary)* ;
	unary -> ( "!" | "-" ) unary | primary ;
	primary -> NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | IDENTIFIER ;
*/

type Parser struct {
	tokens  []scanner.Token
	current int32
}

func New(tokens []scanner.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() ([]Stmt, error) {
	statements := make([]Stmt, 0, 10)
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return statements, err
		}
		statements = append(statements, stmt)
	}
	return statements, nil
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

func (p *Parser) parseBinaryExprLeft(nonTerminal func() (Expr, error), types ...scanner.TokenType) (Expr, error) {
	expr, err := nonTerminal()

	for p.match(types...) {
		token := p.peekBehind()
		right, err := nonTerminal()

		if err != nil {
			return nil, err
		}

		expr = BinaryExpr{Left: expr, Operator: token, Right: right}
	}

	if err != nil {
		return nil, err
	}

	return expr, nil
}

func (p *Parser) generateError(token scanner.Token, message string) error {
	if token.Type == scanner.EOF {
		return &CompileTimeError{Line: token.Line, Where: " at the end", Message: message}
	}
	return &CompileTimeError{Line: token.Line, Where: fmt.Sprintf(" at '%s'", token.Lexeme), Message: message}
}

func (p *Parser) consume(tokenType scanner.TokenType, errorMsg string) error {
	if !p.match(tokenType) {
		return p.generateError(p.peek(), errorMsg)
	}
	return nil
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(scanner.VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() (Stmt, error) {
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
		return p.printStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() (Stmt, error) {
	expr, err := p.expression()

	if err != nil {
		return nil, err
	}

	if err := p.consume(scanner.SEMICOLON, "Expected ';' after a value."); err != nil {
		return nil, err
	}

	return PrintStmt{Expression: expr}, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()

	if err != nil {
		return nil, err
	}

	if err := p.consume(scanner.SEMICOLON, "Expected ';' after a value."); err != nil {
		return nil, err
	}

	return ExpressionStmt{Expression: expr}, nil
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
	condition, err := p.equality()
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

func (p *Parser) equality() (Expr, error) {
	return p.parseBinaryExprLeft(p.comparison, scanner.BANG_EQUAL, scanner.EQUAL_EQUAL)
}

func (p *Parser) comparison() (Expr, error) {
	return p.parseBinaryExprLeft(p.term, scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL)
}

func (p *Parser) term() (Expr, error) {
	return p.parseBinaryExprLeft(p.factor, scanner.PLUS, scanner.MINUS)
}

func (p *Parser) factor() (Expr, error) {
	return p.parseBinaryExprLeft(p.unary, scanner.STAR, scanner.SLASH)
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
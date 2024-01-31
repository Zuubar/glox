package parser

import (
	"fmt"
	"glox/scanner"
)

/*
Grammar:
	ternary -> expression "?" ternary ":" ternary
	expression -> equality ("," equality)*;
	equality -> comparison ( ( "!=" | "==" ) comparison)*
	comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term)*;
	term -> factor ( ( "+" | "-" ) factor)*;
	factor -> unary ( ( "*" | "/" ) unary)*;
	unary -> ( "!" | "-" ) unary | primary;
	primary -> NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")";
*/

type Parser struct {
	tokens  []scanner.Token
	current int32
}

func New(tokens []scanner.Token) Parser {
	return Parser{tokens: tokens}
}

func (p *Parser) Run() (Expr, error) {
	return p.ternary()
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
			return NilExpr{}, err
		}

		expr = BinaryExpr{Left: expr, Operator: token, Right: right}
	}

	if err != nil {
		return NilExpr{}, err
	}

	return expr, nil
}

func (p *Parser) error(token scanner.Token, message string) error {
	if token.Type == scanner.EOF {
		return &Error{Line: token.Line, Where: " at the end", Message: message}
	}
	return &Error{Line: token.Line, Where: fmt.Sprintf(" at '%s'", token.Lexeme), Message: message}
}

func (p *Parser) consume(tokenType scanner.TokenType, errorMsg string) error {
	if !p.match(tokenType) {
		return p.error(p.peekBehind(), errorMsg)
	}
	return nil
}

func (p *Parser) ternary() (Expr, error) {
	condition, err := p.expression()
	if err != nil {
		return NilExpr{}, err
	}

	if !p.match(scanner.QUESTION) {
		return condition, err
	}

	left, err := p.ternary()
	if err != nil {
		return NilExpr{}, err
	}

	if p.match(scanner.COLON) {
		right, err := p.ternary()
		if err != nil {
			return NilExpr{}, err
		}

		return TernaryExpr{Condition: condition, Left: left, Right: right}, nil
	}

	return NilExpr{}, p.error(p.peek(), "Expected ':' after '?'.")
}

func (p *Parser) expression() (Expr, error) {
	expr, err := p.equality()

	if err != nil {
		return NilExpr{}, err
	}

	for p.match(scanner.COMMA) {
		expr, err = p.equality()

		if err != nil {
			return NilExpr{}, err
		}
	}

	return expr, nil
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
		return NilExpr{}, err
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

	if p.match(scanner.LEFT_PAREN) {
		expr, err := p.expression()

		if err != nil {
			return NilExpr{}, err
		}

		if err := p.consume(scanner.RIGHT_PAREN, "Expect ')' after expression."); err != nil {
			return NilExpr{}, err
		}

		return GroupingExpr{Expr: expr}, nil
	}

	return NilExpr{}, p.error(p.peek(), "Expected an expression.")
}

package scanner

import (
	"errors"
	"strconv"
)

var Keywords = map[string]TokenType{
	"and":      AND,
	"class":    CLASS,
	"else":     ELSE,
	"false":    FALSE,
	"fun":      FUN,
	"for":      FOR,
	"if":       IF,
	"nil":      NIL,
	"or":       OR,
	"print":    PRINT,
	"return":   RETURN,
	"super":    SUPER,
	"this":     THIS,
	"true":     TRUE,
	"var":      VAR,
	"while":    WHILE,
	"break":    BREAK,
	"continue": CONTINUE,
}

type Scanner struct {
	source  string
	tokens  []Token
	start   int32
	current int32
	line    int32
	err     error
}

func New(source string) *Scanner {
	return &Scanner{source: source, tokens: make([]Token, 0, 100), line: 1, err: nil}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= int32(len(s.source))
}

func (s *Scanner) isDigit(char uint8) bool {
	return '0' <= char && char <= '9'
}

func (s *Scanner) isAlpha(char uint8) bool {
	return ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || char == '_'
}

func (s *Scanner) match(char uint8) bool {
	if s.isAtEnd() {
		return false
	}
	return s.peek() == char
}

func (s *Scanner) peek() uint8 {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() uint8 {
	if s.current+1 >= int32(len(s.source)) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *Scanner) advance() uint8 {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) addToken(tokenType TokenType, literal any) {
	lexeme := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{Type: tokenType, Lexeme: lexeme, Literal: literal, Line: s.line})
}

func (s *Scanner) blockComment() error {
	indents := 1

	s.advance()
	for !s.isAtEnd() && indents != 0 {
		switch s.peek() {
		case '/':
			if s.peekNext() == '*' {
				indents++
				s.advance()
			}
			break
		case '*':
			if s.peekNext() == '/' {
				indents--
				s.advance()
			}
			break
		case '\n':
			s.line++
			break
		}
		s.advance()
	}

	if indents != 0 {
		return &Error{Line: s.line, Message: "Unterminated block comment."}
	}

	return nil
}

func (s *Scanner) string() error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		return &Error{Line: s.line, Message: "Unterminated string."}
	}
	s.advance()
	s.addToken(STRING, s.source[s.start+1:s.current-1])

	return nil
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) && !s.isAtEnd() {
		s.advance()
	}
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) && !s.isAtEnd() {
			s.advance()
		}
	}
	num, _ := strconv.ParseFloat(s.source[s.start:s.current], 64)
	s.addToken(NUMBER, num)
}

func (s *Scanner) identifier() {
	for !s.isAtEnd() && (s.isAlpha(s.peek()) || s.isDigit(s.peek())) {
		s.advance()
	}
	tokenType, ok := Keywords[s.source[s.start:s.current]]

	if !ok {
		tokenType = IDENTIFIER
	}

	s.addToken(tokenType, nil)
}

func (s *Scanner) Run() ([]Token, error) {
	for !s.isAtEnd() {
		char := s.advance()

		switch char {
		case '(':
			s.addToken(LEFT_PAREN, nil)
			break
		case ')':
			s.addToken(RIGHT_PAREN, nil)
			break
		case '{':
			s.addToken(LEFT_BRACE, nil)
			break
		case '}':
			s.addToken(RIGHT_BRACE, nil)
			break
		case ',':
			s.addToken(COMMA, nil)
			break
		case '.':
			s.addToken(DOT, nil)
			break
		case '-':
			s.addToken(MINUS, nil)
			break
		case '+':
			s.addToken(PLUS, nil)
			break
		case ';':
			s.addToken(SEMICOLON, nil)
			break
		case '*':
			s.addToken(STAR, nil)
			break
		case '?':
			s.addToken(QUESTION, nil)
			break
		case ':':
			s.addToken(COLON, nil)
			break
		case '%':
			s.addToken(MODULO, nil)
			break
		case '/':
			if s.match('/') {
				for s.peek() != '\n' {
					s.advance()
				}
			} else if s.match('*') {
				errors.As(s.blockComment(), &s.err)
			} else {
				s.addToken(SLASH, nil)
			}
			break
		case '!':
			if s.match('=') {
				s.addToken(BANG_EQUAL, nil)
				s.advance()
			} else {
				s.addToken(BANG, nil)
			}
			break
		case '=':
			if s.match('=') {
				s.addToken(EQUAL_EQUAL, nil)
				s.advance()
			} else {
				s.addToken(EQUAL, nil)
			}
		case '>':
			if s.match('=') {
				s.addToken(GREATER_EQUAL, nil)
				s.advance()
			} else {
				s.addToken(GREATER, nil)
			}
		case '<':
			if s.match('=') {
				s.addToken(LESS_EQUAL, nil)
				s.advance()
			} else {
				s.addToken(LESS, nil)
			}
		case '"':
			errors.As(s.string(), &s.err)
			break
		case ' ':
		case '\r':
		case '\t':
			break
		case '\n':
			s.line++
			break
		default:
			if s.isDigit(char) {
				s.number()
			} else if s.isAlpha(char) {
				s.identifier()
			} else {
				s.err = &Error{Line: s.line, Message: "Unexpected character."}
			}
		}
		s.start = s.current
	}
	s.addToken(EOF, nil)

	return s.tokens, s.err
}

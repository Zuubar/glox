package scanner

type TokenType string

const (
	LEFT_PAREN    TokenType = "LEFT_PAREN"
	RIGHT_PAREN   TokenType = "RIGHT_PAREN"
	LEFT_BRACE    TokenType = "LEFT_BRACE"
	RIGHT_BRACE   TokenType = "RIGHT_BRACE"
	LEFT_BRACKET  TokenType = "LEFT_BRACKET"
	RIGHT_BRACKET TokenType = "RIGHT_BRACKET"
	COMMA         TokenType = "COMMA"
	DOT           TokenType = "DOT"
	MINUS         TokenType = "MINUS"
	PLUS          TokenType = "PLUS"
	SEMICOLON     TokenType = "SEMICOLON"
	SLASH         TokenType = "SLASH"
	STAR          TokenType = "STAR"
	MODULO        TokenType = "MODULO"
	// One or two character tokens.

	BANG          TokenType = "BANG"
	BANG_EQUAL    TokenType = "BANG_EQUAL"
	EQUAL         TokenType = "EQUAL"
	EQUAL_EQUAL   TokenType = "EQUAL_EQUAL"
	GREATER       TokenType = "GREATER"
	GREATER_EQUAL TokenType = "GREATER_EQUAL"
	LESS          TokenType = "LESS"
	LESS_EQUAL    TokenType = "LESS_EQUAL"
	QUESTION      TokenType = "QUESTION"
	COLON         TokenType = "COLON"
	// Literals.

	IDENTIFIER TokenType = "IDENTIFIER"
	STRING     TokenType = "STRING"
	NUMBER     TokenType = "NUMBER"
	// Keywords.

	AND       TokenType = "AND"
	CLASS     TokenType = "CLASS"
	TRAIT     TokenType = "TRAIT"
	USE_TRAIT TokenType = "USE_TRAIT"
	ELSE      TokenType = "ELSE"
	FALSE     TokenType = "FALSE"
	FUN       TokenType = "FUN"
	FOR       TokenType = "FOR"
	IF        TokenType = "IF"
	NIL       TokenType = "NIL"
	OR        TokenType = "OR"
	PRINT     TokenType = "PRINT"
	RETURN    TokenType = "RETURN"
	SUPER     TokenType = "SUPER"
	THIS      TokenType = "THIS"
	TRUE      TokenType = "TRUE"
	VAR       TokenType = "VAR"
	WHILE     TokenType = "WHILE"
	BREAK     TokenType = "BREAK"
	CONTINUE  TokenType = "CONTINUE"

	EOF TokenType = "EOF"
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int32
}

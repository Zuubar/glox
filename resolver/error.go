package resolver

import (
	"fmt"
	"glox/scanner"
)

type Error struct {
	Token   scanner.Token
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("[line %d] %s\n", e.Token.Line, e.Message)
}

type Warning struct {
	Token   scanner.Token
	Message string
}

func (e *Warning) Error() string {
	return fmt.Sprintf("[line %d] Warning: %s\n", e.Token.Line, e.Message)
}

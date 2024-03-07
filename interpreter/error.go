package interpreter

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

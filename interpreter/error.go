package interpreter

import (
	"fmt"
	"glox/scanner"
)

type RuntimeError struct {
	Token   scanner.Token
	Message string
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("%s \n[line %d]\n", e.Message, e.Token.Line)
}

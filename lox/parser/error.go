package parser

import "fmt"

type Error struct {
	Line    int32
	Where   string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("[line %d] Error%s: %s\n", e.Line, e.Where, e.Message)
}

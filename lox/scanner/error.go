package scanner

import (
	"fmt"
)

type Error struct {
	Line    int32
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("[line %d] Error: %s\n", e.Line, e.Message)
}

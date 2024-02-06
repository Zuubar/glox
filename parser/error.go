package parser

import "fmt"

type CompileTimeError struct {
	Line    int32
	Where   string
	Message string
}

func (e *CompileTimeError) Error() string {
	return fmt.Sprintf("[line %d] Error%s: %s\n", e.Line, e.Where, e.Message)
}

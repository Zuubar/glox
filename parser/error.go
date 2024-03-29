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

type BreakInterrupt struct {
}

func (e *BreakInterrupt) Error() string {
	return "Break interrupt"
}

type ContinueInterrupt struct {
}

func (e *ContinueInterrupt) Error() string {
	return "Continue interrupt"
}

type ReturnInterrupt struct {
	Value any
}

func (e *ReturnInterrupt) Error() string {
	return fmt.Sprintf("Return interrupt: %d", e.Value)
}

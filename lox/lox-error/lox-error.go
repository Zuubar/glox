package loxError

import (
	"fmt"
	"os"
)

var HadError = false

func Report(line int32, message string) {
	ReportAt(line, "", message)
}

func ReportAt(line int32, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	HadError = true
}

package lusherr

import (
	"fmt"
)

// Debug where an error was raised.
func Debug(err error) string {
	if err == nil {
		return fmt.Sprintf("unknown error")
	}
	frame, found := Locate(err)
	if !found {
		return fmt.Sprintf("%v (unknown caller frame)", err)
	}
	return fmt.Sprintf("%v (%s %s:%d)", err, frame.Function, frame.File, frame.Line)
}

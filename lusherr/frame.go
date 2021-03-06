package lusherr

import (
	"errors"
	"fmt"
	"runtime"
)

// originError is used to give any error an origin frame.
type originError struct {
	err   error
	frame runtime.Frame
}

func (e originError) Error() string {
	return fmt.Sprintf("%v", e.err)
}

// Unwrap the inner error.
func (e originError) Unwrap() error {
	return e.err
}

// Locate the frame of the error.
func (e originError) Locate() runtime.Frame {
	return e.frame
}

// Pin an error to a caller frame.
func (e originError) Pin(frame runtime.Frame) error {
	e.frame = frame
	return e
}

// Locator defines behavior for locating an error frame.
type Locator interface {
	Locate() runtime.Frame
}

// Locate where an error was raised.
func Locate(err error) (runtime.Frame, bool) {
	var frame runtime.Frame
	switch err := err.(type) {
	case Locator:
		return err.Locate(), true
	default:
		if err := errors.Unwrap(err); err != nil {
			return Locate(err)
		}
	}
	return frame, false
}

// Pinner defines behavior for defining an origin frame for an error.
type Pinner interface {
	Pin(runtime.Frame) error
}

// Pin an error to its caller frame to carry the location in code where the error occurred.
func Pin(err error) error {
	switch err := err.(type) {
	case Pinner:
		return err.Pin(frame(1))
	default:
		return originError{
			err:   err,
			frame: frame(1),
		}
	}
}

// frame of the caller, skipped from the caller
func frame(skip int) runtime.Frame {
	rpc := make([]uintptr, 1)
	runtime.Callers(skip+2, rpc)
	frame, _ := runtime.CallersFrames(rpc).Next()
	return frame
}

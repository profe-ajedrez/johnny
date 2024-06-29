package johnny

import (
	"fmt"
	"runtime"
)

// JohnnyError represents an error with additional information about the call stack.
type JohnnyError struct {
	// info is the error message with call stack information.
	info string
}

// NewJohnnyError returns a new JohnnyError with the given information and call stack.
func NewJohnnyError(info any) error {
	// Create a slice to store the call stack.
	stack := make([]uintptr, maxInfoCallstackSize)
	// Get the call stack.
	length := runtime.Callers(fromCaller, stack)
	// Create a new JohnnyError with the error message and call stack.
	return &JohnnyError{
		info: fmt.Sprintf("%s -- %v", info, stack[:length]),
	}
}

// Error returns the error message with call stack information.
func (e *JohnnyError) Error() string {
	return e.info
}

const (
	maxInfoCallstackSize = 12
	fromCaller           = 2
)

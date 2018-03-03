package main

import "fmt"

// We define a custom error type so that we can provide friendlier error messages
type FetchError struct {
	errorCode int    // error code
	details   string // underlying error message, if any
	err       error  // underlying golang error, if any
}

// Implement the golang Error interface
func (e *FetchError) Error() string {
	return fmt.Sprintf("%d - %s", e.errorCode, e.details)
}

func New(errorCode int, details string) error {
	return &FetchError{
		errorCode: errorCode,
		details:   details,
		err:       nil,
	}
}

func newError(errorCode int, details string) *FetchError {
	return &FetchError{
		errorCode: errorCode,
		details:   details,
		err:       nil,
	}
}

func wrapError(err error) *FetchError {
	if err == nil {
		return nil
	}
	return &FetchError{
		errorCode: -1,
		details:   err.Error(),
		err:       err,
	}
}

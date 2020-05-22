package processor

import "fmt"

// processorError is an processor error.
type processorError struct {
	Err    error
	TestID string
}

func (e processorError) Error() string {
	return fmt.Sprintf("TestID[%s]: %s", e.TestID, e.Err.Error())
}

// newError constructs new error with passed parameters.
func newError(err error, testID string) error {
	if err == nil {
		return nil
	}

	return &processorError{
		Err:    err,
		TestID: testID,
	}
}

// Package models describes models for json marshal-unmarshal.
package models

import (
	"fmt"
	"io"
)

// Report represents spamassasin report.
type Report struct {
	SpamAssassin SpamAssassin `json:"spamAssassin"`
}

// SpamAssassin is a root of report.
type SpamAssassin struct {
	Score   float64   `json:"score"`
	Headers []Headers `json:"headers"`
}

// Headers represents info for each header.
type Headers struct {
	Score       float64 `json:"score"`
	Tag         string  `json:"tag"`
	Description string  `json:"description"`
}

// ProcessorInput used for importing reports for processing.
type ProcessorInput struct {
	Data   io.ReadCloser
	TestID string
}

// NewProcessorInput constructs new ProcessorInput with passed parameters.
func NewProcessorInput(data io.ReadCloser, testID string) *ProcessorInput {
	return &ProcessorInput{Data: data, TestID: testID}
}

// ProcessorResponse contains processed input result.
type ProcessorResponse struct {
	TestID string
	Report Report
}

// NewProcessorResponse constructs new ProcessResponse with passed parameters.
func NewProcessorResponse(testID string, report Report) *ProcessorResponse {
	return &ProcessorResponse{TestID: testID, Report: report}
}

// Error is an processor error.
type Error struct {
	Err    error
	TestID string
}

func (e Error) Error() string {
	return fmt.Sprintf("TestID[%s]: %s", e.TestID, e.Err.Error())
}

// NewError constructs new error with passed parameters.
func NewError(err error, testID string) error {
	if err == nil {
		return nil
	}

	return &Error{
		Err:    err,
		TestID: testID,
	}
}

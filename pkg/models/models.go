// Package models describes models for json marshal-unmarshal.
package models

import (
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
	Data   io.Reader
	TestID string
}

// ProcessorResponse contains processed input result.
type ProcessorResponse struct {
	TestID string
	Report Report
	Error  error
}

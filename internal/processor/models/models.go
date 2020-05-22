// Package models describes models for json marshal-unmarshal.
package models

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

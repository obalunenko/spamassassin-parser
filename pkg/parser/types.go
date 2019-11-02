package parser

//go:generate stringer -type=reportType

// reportType
type reportType int

const (
	reportTypeUnknown reportType = iota
	reportType1
	reportType2

	reportTypeSentinel // should be always last.
)

func (i reportType) Valid() bool {
	return i > reportTypeUnknown && i < reportTypeSentinel
}

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

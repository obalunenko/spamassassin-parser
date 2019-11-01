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

type Report struct {
	SpamAssassin SpamAssassin `json:"spamAssassin"`
}

type Headers struct {
	Score       int    `json:"score"`
	Tag         string `json:"tag"`
	Description string `json:"description"`
}

type SpamAssassin struct {
	Score   int       `json:"score"`
	Headers []Headers `json:"headers"`
}

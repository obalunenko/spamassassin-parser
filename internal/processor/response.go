package processor

import (
	"github.com/obalunenko/spamassassin-parser/internal/processor/models"
)

// Response contains processed input result.
type Response struct {
	TestID string
	Report models.Report
}

// NewResponse constructs new ProcessResponse with passed parameters.
func NewResponse(testID string, report models.Report) *Response {
	return &Response{TestID: testID, Report: report}
}

// makeBufferedResponseChan creates buffered response channel.
func makeBufferedResponseChan(buf uint) chan *Response {
	return make(chan *Response, buf)
}

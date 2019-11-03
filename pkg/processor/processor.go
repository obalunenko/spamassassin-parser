// Package processor provides processing reports functionality with channels communication.
package processor

import (
	"context"
	"io"

	"github.com/oleg-balunenko/spamassassin-parser/pkg/models"
	"github.com/oleg-balunenko/spamassassin-parser/pkg/parser"
)

// Response contains processed input result.
type Response struct {
	TestID string
	Report models.Report
	Error  error
}

// MakeResponseChan creates not buffered response channel.
func MakeResponseChan() chan Response {
	return MakeBufferedResponseChan(0)
}

// MakeBufferedResponseChan creates buffered response channel.
func MakeBufferedResponseChan(buf uint) chan Response {
	return make(chan Response, buf)
}

// Input used for importing reports for processing.
type Input struct {
	Data       io.Reader
	TestID     string
	ResultChan chan Response
}

// MakeInputChan creates not buffered input channel.
func MakeInputChan() chan Input {
	return MakeBufferedInputChan(0)
}

// MakeBufferedInputChan creates buffered input channel.
func MakeBufferedInputChan(buf uint) chan Input {
	return make(chan Input, buf)
}

// ProcessReports handles imported reports and runs them through parser.
func ProcessReports(ctx context.Context, incomingReport <-chan Input) {
	for in := range incomingReport {
		if in.ResultChan != nil {
			if ctx.Err() != nil {
				in.ResultChan <- Response{
					TestID: in.TestID,
					Report: models.Report{},
					Error:  ctx.Err(),
				}
				return
			}

			report, err := parser.ParseReport(in.Data)

			in.ResultChan <- Response{
				TestID: in.TestID,
				Report: report,
				Error:  err,
			}
		}
	}
}

package processor

import (
	"context"
	"io"

	"github.com/oleg-balunenko/spamassassin-parser/pkg/parser"
)

// Response contains processed input result.
type Response struct {
	TestID string
	Report parser.Report
	Error  error
}

// InputReport used for importing reports for processing.
type InputReport struct {
	Data       io.Reader
	TestID     string
	ResultChan chan Response
}

// ProcessReports handles imported reports and runs them through parser.
func ProcessReports(ctx context.Context, incomingReport <-chan InputReport) {
	for in := range incomingReport {
		if in.ResultChan != nil {
			if ctx.Err() != nil {
				in.ResultChan <- Response{
					TestID: in.TestID,
					Report: parser.Report{},
					Error:  ctx.Err(),
				}
				return
			}

			report, err := parser.ProcessReport(in.Data)

			in.ResultChan <- Response{
				TestID: in.TestID,
				Report: report,
				Error:  err,
			}
		}
	}
}

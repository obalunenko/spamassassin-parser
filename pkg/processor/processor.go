// Package processor provides processing reports functionality with channels communication.
package processor

import (
	"context"
	"sync"

	"github.com/oleg-balunenko/spamassassin-parser/pkg/models"
	"github.com/oleg-balunenko/spamassassin-parser/pkg/parser"
)

// Processor manages spamassassin reports processing.
type Processor interface {
	// Process handles imported reports and runs them through parser.
	// User is responsible for canceling the context if process need to be stopped.
	Process(ctx context.Context)
	// Results returns the read channel for the messages that are returned by
	// the processor.
	Results() <-chan *models.ProcessorResponse
	// ProcessorInput is the output channel for the user to write messages to that they
	// wish to process.
	Input() chan<- *models.ProcessorInput
	// Close closes underlying input channel - means that no work expected.
	Close()
	closeResults()
}

type processor struct {
	closeOnce   sync.Once
	inChan      chan *models.ProcessorInput
	resultsChan chan *models.ProcessorResponse
}

// NewProcessor creates new instance of processor.
func NewProcessor() Processor {
	return &processor{
		inChan:      makeInputChan(),
		resultsChan: makeResponseChan(),
	}
}

// NewBuffered creates new instance of buffered processor.
func NewBuffered(buf uint) Processor {
	return &processor{
		inChan:      makeBufferedInputChan(buf),
		resultsChan: makeBufferedResponseChan(buf),
	}
}

// makeResponseChan creates not buffered response channel.
func makeResponseChan() chan *models.ProcessorResponse {
	return makeBufferedResponseChan(0)
}

// makeBufferedResponseChan creates buffered response channel.
func makeBufferedResponseChan(buf uint) chan *models.ProcessorResponse {
	return make(chan *models.ProcessorResponse, buf)
}

// makeInputChan creates not buffered input channel.
func makeInputChan() chan *models.ProcessorInput {
	return makeBufferedInputChan(0)
}

// makeBufferedInputChan creates buffered input channel.
func makeBufferedInputChan(buf uint) chan *models.ProcessorInput {
	return make(chan *models.ProcessorInput, buf)
}

func (p *processor) Process(ctx context.Context) {
	defer func() {
		p.closeResults()
	}()

	for in := range p.inChan {
		if p.resultsChan != nil {
			if ctx.Err() != nil {
				return
			}

			report, err := parser.ParseReport(in.Data)
			p.resultsChan <- &models.ProcessorResponse{
				TestID: in.TestID,
				Report: report,
				Error:  err,
			}
		}
	}
}

func (p *processor) Results() <-chan *models.ProcessorResponse {
	return p.resultsChan
}

func (p *processor) Input() chan<- *models.ProcessorInput {
	return p.inChan
}

func (p *processor) Close() {
	p.closeOnce.Do(func() {
		close(p.inChan)
	})
}

func (p *processor) closeResults() {
	p.closeOnce.Do(func() {
		close(p.resultsChan)
	})
}

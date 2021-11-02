// Package processor provides processing reports functionality with channels communication.
package processor

import (
	"context"
	"sync"

	log "github.com/obalunenko/logger"

	"github.com/obalunenko/spamassassin-parser/internal/processor/parser"
)

// Processor manages spamassassin reports processing.
type Processor interface {
	// Process handles imported reports and runs them through parser.
	// User is responsible for canceling the context if process need to be stopped.
	Process(ctx context.Context)
	// Results returns the read channel for the messages that are returned by
	// the processor.
	// Values from channel should be read or deadlock will be occurred if in config results channel is enabled.
	Results() <-chan *Response
	// Errors returns the read channel for the errors that are returned by processor.
	// Values from channel should be read or deadlock will be occurred if in config errors channel is enabled.
	Errors() <-chan error
	// Input is the output channel for the user to write messages to that they
	// wish to process.
	Input() chan<- *Input
	// Close closes underlying input channel - means that no work expected.
	Close()
}

type processor struct {
	closeOnce   sync.Once
	inChan      chan *Input
	resultsChan chan *Response
	errorsChan  chan error
}

// NewDefault creates new instance of processor with sane default config.
// Not buffered. Response is enabled. Errors are disabled.
func NewDefault() Processor {
	return New(NewConfig())
}

// New creates processor instance.
func New(cfg *Config) Processor {
	if cfg == nil {
		cfg = NewConfig()
	}

	var pr processor
	pr.inChan = makeBufferedInputChan(cfg.Buffer)

	if cfg.Receive.Response {
		pr.resultsChan = makeBufferedResponseChan(cfg.Buffer)
	}

	if cfg.Receive.Errors {
		pr.errorsChan = make(chan error)
	}

	return &pr
}

func (p *processor) Process(ctx context.Context) {
	defer func() {
		p.closeResults()
		p.closeErrors()
	}()

	for in := range p.inChan {
		if ctx.Err() != nil {
			return
		}

		p.processData(ctx, in)
	}
}

func (p *processor) processData(ctx context.Context, in *Input) {
	if in == nil {
		return
	}

	defer func() {
		if err := in.Data.Close(); err != nil {
			log.WithError(ctx, err).Error("processor: failed to close input reader")
		}
	}()

	report, err := parser.ParseReport(in.Data)
	if err != nil {
		err = newError(err, in.TestID)

		if p.errorsChan != nil {
			p.errorsChan <- err
		} else {
			log.WithError(ctx, err).Error("processor: failed to parse report")
		}

		return
	}

	resp := NewResponse(in.TestID, report)

	if p.resultsChan != nil {
		p.resultsChan <- resp
	} else {
		log.WithFields(ctx, log.Fields{
			"test_id": resp.TestID,
			"report":  resp.Report,
		}).Info("processor: processed")
	}
}

func (p *processor) Results() <-chan *Response {
	return p.resultsChan
}

func (p *processor) Input() chan<- *Input {
	return p.inChan
}

func (p *processor) Errors() <-chan error {
	return p.errorsChan
}

func (p *processor) Close() {
	p.closeOnce.Do(func() {
		if p.inChan != nil {
			close(p.inChan)
		}
	})
}

func (p *processor) closeResults() {
	p.closeOnce.Do(func() {
		if p.resultsChan != nil {
			close(p.resultsChan)
		}
	})
}

func (p *processor) closeErrors() {
	p.closeOnce.Do(func() {
		if p.errorsChan != nil {
			close(p.errorsChan)
		}
	})
}

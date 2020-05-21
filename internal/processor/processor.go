// Package processor provides processing reports functionality with channels communication.
package processor

import (
	"context"
	"io"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/oleg-balunenko/spamassassin-parser/internal/processor/models"
	"github.com/oleg-balunenko/spamassassin-parser/internal/processor/parser"
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
	// ProcessorInput is the output channel for the user to write messages to that they
	// wish to process.
	Input() chan<- *Input
	// Close closes underlying input channel - means that no work expected.
	Close()
}

// Input used for importing reports for processing.
type Input struct {
	Data   io.ReadCloser
	TestID string
}

// NewInput constructs new ProcessorInput with passed parameters.
func NewInput(data io.ReadCloser, testID string) *Input {
	return &Input{Data: data, TestID: testID}
}

// Response contains processed input result.
type Response struct {
	TestID string
	Report models.Report
}

// NewResponse constructs new ProcessResponse with passed parameters.
func NewResponse(testID string, report models.Report) *Response {
	return &Response{TestID: testID, Report: report}
}

type processor struct {
	closeOnce   sync.Once
	inChan      chan *Input
	resultsChan chan *Response
	errorsChan  chan error
}

// Config is a processor instance configuration.
type Config struct {
	Buffer  uint
	Receive struct {
		Response bool
		Errors   bool
	}
}

// NewConfig creates new config filled with sane default values.
func NewConfig() *Config {
	return &Config{
		Buffer: 0,
		Receive: struct {
			Response bool
			Errors   bool
		}{
			Response: true,
			Errors:   false,
		},
	}
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

// makeBufferedResponseChan creates buffered response channel.
func makeBufferedResponseChan(buf uint) chan *Response {
	return make(chan *Response, buf)
}

// makeBufferedInputChan creates buffered input channel.
func makeBufferedInputChan(buf uint) chan *Input {
	return make(chan *Input, buf)
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

		p.processData(in)
	}
}

func (p *processor) processData(in *Input) {
	if in == nil {
		return
	}

	defer func() {
		if err := in.Data.Close(); err != nil {
			log.Error(err)
		}
	}()

	report, err := parser.ParseReport(in.Data)
	if err != nil {
		err = newError(err, in.TestID)

		if p.errorsChan != nil {
			p.errorsChan <- err
		} else {
			log.Error(err)
		}

		return
	}

	resp := NewResponse(in.TestID, report)

	if p.resultsChan != nil {
		p.resultsChan <- resp
	} else {
		log.Infof("TestID[%s]: processed\n %+v \n", resp.TestID, resp.Report)
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

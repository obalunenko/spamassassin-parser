package processor

import "io"

// Input used for importing reports for processing.
type Input struct {
	Data   io.ReadCloser
	TestID string
}

// NewInput constructs new ProcessorInput with passed parameters.
func NewInput(data io.ReadCloser, testID string) *Input {
	return &Input{Data: data, TestID: testID}
}

// makeBufferedInputChan creates buffered input channel.
func makeBufferedInputChan(buf uint) chan *Input {
	return make(chan *Input, buf)
}

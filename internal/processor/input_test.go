package processor_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/obalunenko/spamassassin-parser/internal/processor"
)

func TestNewInput(t *testing.T) {
	type args struct {
		data   io.ReadCloser
		testID string
	}

	tests := []struct {
		name string
		args args
		want *processor.Input
	}{
		{
			name: "make processor input",
			args: args{
				data:   ioutil.NopCloser(bytes.NewReader([]byte("test reader"))),
				testID: "test 1",
			},
			want: &processor.Input{
				Data:   ioutil.NopCloser(bytes.NewReader([]byte("test reader"))),
				TestID: "test 1",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := processor.NewInput(tt.args.data, tt.args.testID)
			assert.Equal(t, tt.want, got)
		})
	}
}

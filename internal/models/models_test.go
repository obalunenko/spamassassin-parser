package models

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProcessorInput(t *testing.T) {
	type args struct {
		data   io.ReadCloser
		testID string
	}

	tests := []struct {
		name string
		args args
		want *ProcessorInput
	}{
		{
			name: "make processor input",
			args: args{
				data:   ioutil.NopCloser(bytes.NewReader([]byte("test reader"))),
				testID: "test 1",
			},
			want: &ProcessorInput{
				Data:   ioutil.NopCloser(bytes.NewReader([]byte("test reader"))),
				TestID: "test 1",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := NewProcessorInput(tt.args.data, tt.args.testID)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestError_Error(t *testing.T) {
	type fields struct {
		Err    error
		TestID string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "error message",
			fields: fields{
				Err:    errors.New("test error"),
				TestID: "testID1",
			},
			want: "TestID[testID1]: test error",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				Err:    tt.fields.Err,
				TestID: tt.fields.TestID,
			}
			got := e.Error()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewError(t *testing.T) {
	type args struct {
		err    error
		testID string
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "new error valid",
			args: args{
				err:    errors.New("test error"),
				testID: "testID1",
			},
			wantErr: &Error{
				Err:    errors.New("test error"),
				TestID: "testID1",
			},
		},
		{
			name: "new error nil",
			args: args{
				err:    nil,
				testID: "testID1",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := NewError(tt.args.err, tt.args.testID)
			assert.Equal(t, tt.wantErr, got)
		})
	}
}

func TestNewProcessorResponse(t *testing.T) {
	type args struct {
		testID string
		report Report
	}

	tests := []struct {
		name string
		args args
		want *ProcessorResponse
	}{
		{
			name: "make valid response",
			args: args{
				testID: "testID1",
				report: Report{
					SpamAssassin: SpamAssassin{
						Score: 1,
						Headers: []Headers{
							{
								Score:       1,
								Tag:         "TEST_TAG",
								Description: "descr",
							},
						},
					},
				},
			},
			want: &ProcessorResponse{
				TestID: "testID1",
				Report: Report{
					SpamAssassin: SpamAssassin{
						Score: 1,
						Headers: []Headers{
							{
								Score:       1,
								Tag:         "TEST_TAG",
								Description: "descr",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := NewProcessorResponse(tt.args.testID, tt.args.report)
			assert.Equal(t, tt.want, got)
		})
	}
}

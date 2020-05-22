package processor

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newError(t *testing.T) {
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
			wantErr: &processorError{
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
			got := newError(tt.args.err, tt.args.testID)
			assert.Equal(t, tt.wantErr, got)
		})
	}
}

func Test_processorError_Error(t *testing.T) {
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
			e := processorError{
				Err:    tt.fields.Err,
				TestID: tt.fields.TestID,
			}
			got := e.Error()
			assert.Equal(t, tt.want, got)
		})
	}
}

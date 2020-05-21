package processor

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oleg-balunenko/spamassassin-parser/internal/processor/models"
	"github.com/oleg-balunenko/spamassassin-parser/pkg/utils"
)

type want struct {
	filepath string
	wantErr  bool
}

type input struct {
	filepath string
	testID   string
}

type test struct {
	input input
	want  want
}

type expected struct {
	report  models.Report
	wantErr bool
}

func casesTestProcessReports(t testing.TB) ([]test, map[string]expected) {
	t.Helper()

	tests := []test{
		{
			input: input{
				filepath: filepath.Join("testdata", "report1.txt"),
				testID:   "report1.txt",
			},
			want: want{
				filepath: filepath.Join("testdata", "report1.golden.json"),
				wantErr:  false,
			},
		},
		{
			input: input{
				filepath: filepath.Join("testdata", "report2.txt"),
				testID:   "report2.txt",
			},
			want: want{
				filepath: filepath.Join("testdata", "report2.golden.json"),
				wantErr:  false,
			},
		},
		{
			input: input{
				filepath: filepath.Join("testdata", "report1.txt"),
				testID:   "report1.txt.repeat",
			},
			want: want{
				filepath: filepath.Join("testdata", "report1.golden.json"),
				wantErr:  false,
			},
		},
		{
			input: input{
				filepath: filepath.Join("testdata", "empty.json"),
				testID:   "empty",
			},
			want: want{
				filepath: filepath.Join("testdata", "empty.json"),
				wantErr:  true,
			},
		},
	}

	expResults := make(map[string]expected, len(tests))

	for _, tt := range tests {
		tt := tt
		report := models.GetReportFromFile(t, tt.want.filepath)

		expResults[tt.input.testID] = expected{
			report:  report,
			wantErr: tt.want.wantErr,
		}
	}

	return tests, expResults
}

func TestProcessReports(t *testing.T) {
	var secondsNum time.Duration = 5

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*secondsNum)

	defer cancel()

	cfg := NewConfig()
	cfg.Receive.Errors = true

	processor := New(cfg)

	go processor.Process(ctx)

	tests, expResults := casesTestProcessReports(t)

	go func() {
		for _, tt := range tests {
			tt := tt
			file := utils.GetReaderFromFile(t, tt.input.filepath)
			t.Logf("processing report: %s \n", tt.input.testID)
			processor.Input() <- &Input{
				Data:   file,
				TestID: tt.input.testID,
			}
		}

		processor.Close()
	}()

	// check all reports processed
	var processed int
LOOP:
	for {
		select {
		case res := <-processor.Results():
			if res != nil {
				processed++
				t.Logf("received result: %s\n", res.TestID)
				exp := expResults[res.TestID]

				assert.Equal(t, exp.report, res.Report)
			}
		case err := <-processor.Errors():
			require.IsType(t, &processorError{}, err, "unexpected error type")
			merr, ok := err.(*processorError)
			require.True(t, ok)

			exp := expResults[merr.TestID]

			if exp.wantErr {
				assert.Error(t, err)
				processed++
				continue
			}
			assert.NoError(t, err)

		case <-ctx.Done():
			assert.Equal(t, len(expResults), processed, "deadline reached, but not all results received")
			var secondsNum time.Duration = 2
			time.Sleep(time.Second * secondsNum)
			break LOOP
		}
	}
}

func TestNewConfig(t *testing.T) {
	expConfgig := &Config{
		Buffer: 0,
		Receive: struct {
			Response bool
			Errors   bool
		}{
			Response: true,
			Errors:   false,
		},
	}
	got := NewConfig()
	require.Equal(t, expConfgig, got)
}

func TestNewDefaultProcessor(t *testing.T) {
	got := NewDefault()
	assert.NotNil(t, got)
	assert.IsType(t, &processor{}, got)
	assert.NotNil(t, got.Results())
	assert.Nil(t, got.Errors())
	assert.NotNil(t, got.Input())
}

func TestNewProcessorInput(t *testing.T) {
	type args struct {
		data   io.ReadCloser
		testID string
	}

	tests := []struct {
		name string
		args args
		want *Input
	}{
		{
			name: "make processor input",
			args: args{
				data:   ioutil.NopCloser(bytes.NewReader([]byte("test reader"))),
				testID: "test 1",
			},
			want: &Input{
				Data:   ioutil.NopCloser(bytes.NewReader([]byte("test reader"))),
				TestID: "test 1",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := NewInput(tt.args.data, tt.args.testID)
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
			e := processorError{
				Err:    tt.fields.Err,
				TestID: tt.fields.TestID,
			}
			got := e.Error()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewProcessorResponse(t *testing.T) {
	type args struct {
		testID string
		report models.Report
	}

	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "make valid response",
			args: args{
				testID: "testID1",
				report: models.Report{
					SpamAssassin: models.SpamAssassin{
						Score: 1,
						Headers: []models.Headers{
							{
								Score:       1,
								Tag:         "TEST_TAG",
								Description: "descr",
							},
						},
					},
				},
			},
			want: &Response{
				TestID: "testID1",
				Report: models.Report{
					SpamAssassin: models.SpamAssassin{
						Score: 1,
						Headers: []models.Headers{
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
			got := NewResponse(tt.args.testID, tt.args.report)
			assert.Equal(t, tt.want, got)
		})
	}
}

package parser

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oleg-balunenko/spamassassin-parser/internal/processor/models"
	"github.com/oleg-balunenko/spamassassin-parser/pkg/utils"
)

const (
	testdata    = "testdata"
	goldenFile1 = "report1.golden.json"
	goldenFile2 = "report2.golden.json"
	testReport1 = "report1.txt"
	testReport2 = "report2.txt"
)

func Test_processReport(t *testing.T) {
	type args struct {
		filepath string
		rt       reportType
	}

	type expected struct {
		filepath string
		wantErr  bool
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "process report type 1",
			args: args{
				filepath: filepath.Join("..", testdata, testReport1),
				rt:       reportType1,
			},
			expected: expected{
				filepath: filepath.Join("..", testdata, goldenFile1),
				wantErr:  false,
			},
		},
		{
			name: "process report type 2",
			args: args{
				filepath: filepath.Join("..", testdata, testReport2),
				rt:       reportType2,
			},
			expected: expected{
				filepath: filepath.Join("..", testdata, goldenFile2),
				wantErr:  false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			data := utils.GetReaderFromFile(t, tt.args.filepath)
			parser, err := newParser(tt.args.rt)
			require.NoError(t, err)
			require.NotNil(t, parser)

			got, err := parser.Parse(data)
			if tt.expected.wantErr {
				assert.Error(t, err)
				return
			}
			wantReport := models.GetReportFromFile(t, tt.expected.filepath)
			assert.NoError(t, err)
			assert.Equal(t, wantReport, got)
		})
	}
}

func Test_getReportType(t *testing.T) {
	type args struct {
		filepath string
	}

	type expected struct {
		rt reportType
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "report type 1",
			args: args{
				filepath: filepath.Join("..", testdata, testReport1),
			},
			expected: expected{
				rt: reportType1,
			},
		},
		{
			name: "report type 2",
			args: args{
				filepath: filepath.Join("..", testdata, testReport2),
			},
			expected: expected{
				rt: reportType2,
			},
		},
		{
			name: "unknown type",
			args: args{
				filepath: filepath.Join("..", testdata, "empty.json"),
			},
			expected: expected{
				rt: reportTypeUnknown,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			data := utils.GetReaderFromFile(t, tt.args.filepath)
			got := getReportType(data)
			assert.Equal(t, tt.expected.rt.String(), got.String())
		})
	}
}

func Test_newParser(t *testing.T) {
	type args struct {
		rt reportType
	}

	tests := []struct {
		name    string
		args    args
		want    Parser
		wantErr bool
	}{
		{
			name: "report1 parser",
			args: args{
				rt: reportType1,
			},
			want:    report1Parser{},
			wantErr: false,
		},
		{
			name: "report2 parser",
			args: args{
				rt: reportType2,
			},
			want:    report2Parser{},
			wantErr: false,
		},
		{
			name: "unknown report parser",
			args: args{
				rt: reportTypeUnknown,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := newParser(tt.args.rt)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

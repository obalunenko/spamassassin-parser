package parser

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/obalunenko/spamassassin-parser/internal/processor/models"
	"github.com/obalunenko/spamassassin-parser/pkg/utils"
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
				filepath: filepath.FromSlash("../testdata/report1.txt"),
				rt:       reportType1,
			},
			expected: expected{
				filepath: filepath.FromSlash("../testdata/report1.golden.json"),
				wantErr:  false,
			},
		},
		{
			name: "process report type 2",
			args: args{
				filepath: filepath.FromSlash("../testdata/report2.txt"),
				rt:       reportType2,
			},
			expected: expected{
				filepath: filepath.FromSlash("../testdata/report2.golden.json"),
				wantErr:  false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			data := utils.GetReaderFromFile(t, tt.args.filepath)
			defer func() {
				require.NoError(t, data.Close(), "close reader")
			}()

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
	t.Parallel()

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
				filepath: filepath.FromSlash("../testdata/report1.txt"),
			},
			expected: expected{
				rt: reportType1,
			},
		},
		{
			name: "report type 2",
			args: args{
				filepath: filepath.FromSlash("../testdata/report2.txt"),
			},
			expected: expected{
				rt: reportType2,
			},
		},
		{
			name: "unknown type",
			args: args{
				filepath: filepath.FromSlash("../testdata/empty.json"),
			},
			expected: expected{
				rt: reportTypeUnknown,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data := utils.GetReaderFromFile(t, tt.args.filepath)
			got := getReportType(data)
			assert.Equal(t, tt.expected.rt.String(), got.String())
		})
	}
}

func Test_newParser(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

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

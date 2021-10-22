package parser_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/obalunenko/spamassassin-parser/internal/processor/models"
	"github.com/obalunenko/spamassassin-parser/internal/processor/parser"
	"github.com/obalunenko/spamassassin-parser/pkg/utils"
)

func TestProcessReport(t *testing.T) {
	type args struct {
		filepath string
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
			got, err := parser.ParseReport(data)
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

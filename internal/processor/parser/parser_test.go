package parser_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oleg-balunenko/spamassassin-parser/internal/processor/models"
	"github.com/oleg-balunenko/spamassassin-parser/internal/processor/parser"
	"github.com/oleg-balunenko/spamassassin-parser/pkg/utils"
)

const (
	testdata    = "testdata"
	goldenFile1 = "report1.golden.json"
	goldenFile2 = "report2.golden.json"
	testReport1 = "report1.txt"
	testReport2 = "report2.txt"
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
				filepath: filepath.Join("..", testdata, testReport1),
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

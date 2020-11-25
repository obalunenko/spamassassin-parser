package processor_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/obalunenko/spamassassin-parser/internal/processor"
	"github.com/obalunenko/spamassassin-parser/internal/processor/models"
)

func TestNewResponse(t *testing.T) {
	type args struct {
		testID string
		report models.Report
	}

	tests := []struct {
		name string
		args args
		want *processor.Response
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
			want: &processor.Response{
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
		t.Parallel()

		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := processor.NewResponse(tt.args.testID, tt.args.report)
			assert.Equal(t, tt.want, got)
		})
	}
}

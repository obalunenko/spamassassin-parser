package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oleg-balunenko/spamassassin-parser/pkg/models"
)

func Test_makeHeader(t *testing.T) {
	type args struct {
		score       string
		tag         string
		description string
	}

	tests := []struct {
		name    string
		args    args
		want    models.Headers
		wantErr bool
	}{
		{
			name: "make valid header",
			args: args{
				score:       "0.2",
				tag:         "TEST_TAG",
				description: "Description of header",
			},
			want: models.Headers{
				Score:       0.2,
				Tag:         "TEST_TAG",
				Description: "Description of header",
			},
			wantErr: false,
		},
		{
			name: "make invalid header",
			args: args{
				score:       "0.s2",
				tag:         "TEST_TAG",
				description: "Description of header",
			},
			want:    models.Headers{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := makeHeader(tt.args.score, tt.args.tag, tt.args.description)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

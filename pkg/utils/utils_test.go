package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/obalunenko/spamassassin-parser/pkg/utils"
)

func TestPrettyPrint(t *testing.T) {
	t.Parallel()
	tests := casesTestPrettyPrint(t)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := utils.PrettyPrint(tt.args.v, tt.args.prefix, tt.args.indent)
			if tt.wantErr {
				assert.Error(t, err)

				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

type args struct {
	v      interface{}
	prefix string
	indent string
}

type test struct {
	name    string
	args    args
	want    string
	wantErr bool
}

func casesTestPrettyPrint(t testing.TB) []test {
	t.Helper()

	return []test{
		{
			name: "print struct",
			args: args{
				v: struct {
					Field1 string
					Field2 string
				}{
					Field1: "testfield1",
					Field2: "testfield2",
				},
				prefix: "",
				indent: "\t",
			},
			want: `{
	"Field1": "testfield1",
	"Field2": "testfield2"
}
`,
			wantErr: false,
		},
		{
			name: "print map",
			args: args{
				v: map[string]int{
					"Field1": 1,
					"Field2": 2,
				},
				prefix: "",
				indent: "\t",
			},
			want: `{
	"Field1": 1,
	"Field2": 2
}
`,
			wantErr: false,
		},
		{
			name: "print nil",
			args: args{
				v:      nil,
				prefix: "",
				indent: "\t",
			},
			want:    "null\n",
			wantErr: false,
		},
	}
}

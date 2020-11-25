package processor_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/obalunenko/spamassassin-parser/internal/processor"
)

func TestNewConfig(t *testing.T) {
	expConfgig := &processor.Config{
		Buffer: 0,
		Receive: struct {
			Response bool
			Errors   bool
		}{
			Response: true,
			Errors:   false,
		},
	}
	got := processor.NewConfig()
	require.Equal(t, expConfgig, got)
}

package appconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oleg-balunenko/spamassassin-parser/pkg/env"
)

func TestLoad(t *testing.T) {
	t.Run("Load default", func(t *testing.T) {
		want := Config{
			InputDir:      "input",
			ResultDir:     "output",
			ArchiveDir:    "archive",
			ReceiveErrors: true,
		}

		got := Load()

		assert.Equal(t, want, got)
	})

	t.Run("Load with set env variables", func(t *testing.T) {
		inputDir := "datainput"
		reset := env.SetForTesting(t, "SPAMASSASSIN_INPUT", inputDir)
		defer reset()

		want := Config{
			InputDir:      inputDir,
			ResultDir:     "output",
			ArchiveDir:    "archive",
			ReceiveErrors: true,
		}

		got := Load()

		assert.Equal(t, want, got)
	})
}

package appconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oleg-balunenko/spamassassin-parser/internal/appconfig"
	"github.com/oleg-balunenko/spamassassin-parser/pkg/env"
)

func TestLoad(t *testing.T) {
	t.Run("Load default", func(t *testing.T) {
		want := appconfig.Config{
			InputDir:      "input",
			ResultDir:     "result",
			ArchiveDir:    "archive",
			ReceiveErrors: true,
		}

		got := appconfig.Load()

		assert.Equal(t, want, got)
	})

	t.Run("Load with set env variables", func(t *testing.T) {
		inputDir := "datainput"
		reset := env.SetForTesting(t, "SPAMASSASSIN_INPUT", inputDir)
		defer reset()

		want := appconfig.Config{
			InputDir:      inputDir,
			ResultDir:     "result",
			ArchiveDir:    "archive",
			ReceiveErrors: true,
		}

		got := appconfig.Load()

		assert.Equal(t, want, got)
	})
}

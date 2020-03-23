// Package appconfig provide application configuration.
package appconfig

import "github.com/oleg-balunenko/spamassassin-parser/pkg/env"

// Config stores application configuration.
type Config struct {
	InputDir      string
	ResultDir     string
	ArchiveDir    string
	ReceiveErrors bool
}

// Load loads application configuration.
func Load() Config {
	return Config{
		InputDir:      env.GetStringOrDefault("SPAMASSASSIN_INPUT", "input"),
		ResultDir:     env.GetStringOrDefault("SPAMASSASSIN_OUTPUT", "output"),
		ArchiveDir:    env.GetStringOrDefault("SPAMASSASSIN_ARCHIVE", "archive"),
		ReceiveErrors: env.GetBoolOrDefault("SPAMASSASSIN_RECEIVE_ERRORS", true),
	}
}

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

// Load loads application configuration with default values if environment variables not set.
// SPAMASSASSIN_INPUT - input
// SPAMASSASSIN_OUTPUT - output
// SPAMASSASSIN_ARCHIVE - archive
// SPAMASSASSIN_RECEIVE_ERRORS - true
func Load() Config {
	return Config{
		InputDir:      env.GetStringOrDefault("SPAMASSASSIN_INPUT", "input"),
		ResultDir:     env.GetStringOrDefault("SPAMASSASSIN_RESULT", "result"),
		ArchiveDir:    env.GetStringOrDefault("SPAMASSASSIN_ARCHIVE", "archive"),
		ReceiveErrors: env.GetBoolOrDefault("SPAMASSASSIN_RECEIVE_ERRORS", true),
	}
}

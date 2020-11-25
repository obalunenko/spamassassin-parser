// Package appconfig provide application configuration.
package appconfig

import "github.com/obalunenko/spamassassin-parser/pkg/getenv"

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
// SPAMASSASSIN_RECEIVE_ERRORS - true.
func Load() Config {
	return Config{
		InputDir:      getenv.StringOrDefault("SPAMASSASSIN_INPUT", "input"),
		ResultDir:     getenv.StringOrDefault("SPAMASSASSIN_RESULT", "result"),
		ArchiveDir:    getenv.StringOrDefault("SPAMASSASSIN_ARCHIVE", "archive"),
		ReceiveErrors: getenv.BoolOrDefault("SPAMASSASSIN_RECEIVE_ERRORS", true),
	}
}

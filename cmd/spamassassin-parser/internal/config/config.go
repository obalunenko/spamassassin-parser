// Package config provide application configuration.
package config

import (
	"flag"
)

var (
	inputDir    = flag.String("input", "input", "Input dir path")
	resultDir   = flag.String("result", "result", "Results dir path")
	archiveDir  = flag.String("archive", "archive", "Archive dir path")
	receiveErrs = flag.Bool("errors", true, "Receive parse errors")
	// Log related configs.
	logLevel       = flag.String("log_level", "INFO", "set log level of application")
	logFormat      = flag.String("log_format", "text", "Format of logs (supported values: text, json")
	logSentryDSN   = flag.String("log_sentry_dsn", "", "Sentry DSN")
	logSentryTrace = flag.Bool("log_sentry_trace", false,
		"Enables sending stacktrace to sentry")
	logSentryTraceLevel = flag.String("log_sentry_trace_level", "PANIC",
		"The level at which to start capturing stacktraces")
)

// ensureFlags panics if env is checked before flags are parsed.
// Ok to panic since this should be caught in dev or staging.
func ensureFlags() {
	if !flag.Parsed() {
		panic("flags not parsed yet")
	}
}

func init() {

}

// InputDir returns Input dir path.
func InputDir() string {
	ensureFlags()
	return *inputDir
}

// ResultDir returns Results dir path.
func ResultDir() string {
	ensureFlags()
	return *resultDir
}

// ArchiveDir returns Archive dir path.
func ArchiveDir() string {
	ensureFlags()
	return *archiveDir
}

// ReceiveErrors returns if Receive parse errors option enabled.
func ReceiveErrors() bool {
	ensureFlags()
	return *receiveErrs
}

// LogLevel config.
func LogLevel() string {
	ensureFlags()
	return *logLevel
}

// LogSentryDSN config.
func LogSentryDSN() string {
	ensureFlags()
	return *logSentryDSN
}

// LogSentryEnabled config.
func LogSentryEnabled() bool {
	ensureFlags()
	return LogSentryDSN() != ""
}

// LogSentryTraceEnabled config.
func LogSentryTraceEnabled() bool {
	ensureFlags()
	return *logSentryTrace
}

// LogSentryTraceLevel config.
func LogSentryTraceLevel() string {
	ensureFlags()
	return *logSentryTraceLevel
}

// LogFormat config.
func LogFormat() string {
	ensureFlags()
	return *logFormat
}

// Load loads application configuration.
func Load() {
	flag.Parse()
}

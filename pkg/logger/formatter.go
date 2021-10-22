package logger

import (
	"time"

	"github.com/sirupsen/logrus" //nolint:depguard // this is the only place where logrus should be imported.
)

const (
	jsonFmt = "json"
	textFmt = "text"
)

func makeFormatter(format string) logrus.Formatter {
	var f logrus.Formatter

	switch format {
	case jsonFmt:
		f = jsonFormatter()
	case textFmt:
		f = textFormatter()
	default:
		f = textFormatter()
	}

	return f
}

func jsonFormatter() logrus.Formatter {
	f := logrus.JSONFormatter{
		TimestampFormat:   time.RFC3339Nano,
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		DataKey:           "metadata",
		FieldMap:          nil,
		CallerPrettyfier:  nil,
		PrettyPrint:       false,
	}

	return &f
}

func textFormatter() logrus.Formatter {
	f := logrus.TextFormatter{
		ForceColors:               false,
		DisableColors:             false,
		ForceQuote:                false,
		DisableQuote:              false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             true,
		TimestampFormat:           "02-01-2006 15:04:05",
		DisableSorting:            false,
		SortingFunc:               nil,
		DisableLevelTruncation:    false,
		PadLevelText:              false,
		QuoteEmptyFields:          true,
		FieldMap:                  nil,
		CallerPrettyfier:          nil,
	}

	return &f
}

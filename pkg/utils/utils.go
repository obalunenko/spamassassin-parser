// Package utils provides common helper functions that used in code base.
package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/oleg-balunenko/spamassassin-parser/pkg/models"
)

// PrettyPrint appends to passed struct indents and returns a human readable form of struct.
// Each element of JSON object will start from indent with prefix.
func PrettyPrint(v interface{}, prefix string, indent string) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal")
	}

	var out bytes.Buffer
	if err := json.Indent(&out, b, prefix, indent); err != nil {
		return "", errors.Wrap(err, "failed to indent")
	}

	if _, err := out.WriteString("\n"); err != nil {
		return "", errors.Wrap(err, "failed to write string")
	}

	return out.String(), nil
}

// GetReaderFromFile is a test helper that opens passed filepath and returns reader.
func GetReaderFromFile(tb testing.TB, fPath string) io.ReadCloser {
	tb.Helper()

	file, err := os.Open(filepath.Clean(fPath))
	require.NoError(tb, err)

	return file
}

// GetReportFromFile is a test helper that unmarshal passed filepath into models.Report
func GetReportFromFile(tb testing.TB, fPath string) models.Report {
	tb.Helper()

	b, err := ioutil.ReadFile(filepath.Clean(fPath))
	require.NoError(tb, err)

	var rp models.Report

	err = json.Unmarshal(b, &rp)
	require.NoError(tb, err)

	return rp
}

// Package utils provides common helper functions that used in code base.
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// PrettyPrint appends to passed struct indents and returns a human-readable form of struct.
// Each element of JSON object will start from indent with prefix.
func PrettyPrint(v interface{}, prefix string, indent string) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("failed to marshal: %w", err)
	}

	var out bytes.Buffer
	if err := json.Indent(&out, b, prefix, indent); err != nil {
		return "", fmt.Errorf("failed to indent: %w", err)
	}

	if _, err := out.WriteString("\n"); err != nil {
		return "", fmt.Errorf("failed to write string: %w", err)
	}

	return out.String(), nil
}

// GetReaderFromFile is a test helper that opens passed filepath and returns reader.
// Caller of this function is responsible for closing io.ReadCloser.
func GetReaderFromFile(tb testing.TB, fPath string) io.ReadCloser {
	tb.Helper()

	file, err := os.Open(filepath.Clean(fPath))
	require.NoError(tb, err)

	return file
}

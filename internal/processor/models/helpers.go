package models

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// GetReportFromFile is a test helper that unmarshal passed filepath into models.Report.
func GetReportFromFile(tb testing.TB, fPath string) Report {
	tb.Helper()

	b, err := ioutil.ReadFile(filepath.Clean(fPath))
	require.NoError(tb, err)

	var rp Report

	err = json.Unmarshal(b, &rp)
	require.NoError(tb, err)

	return rp
}

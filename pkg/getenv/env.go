// Package getenv provide functionality for loading values from environment variables.
package getenv

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

// StringOrDefault returns string environment variable value or passed default.
func StringOrDefault(key, defVal string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return defVal
	}

	return val
}

// BoolOrDefault returns boolean environment variable value or passed default.
func BoolOrDefault(key string, defVal bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return defVal
	}

	v, err := strconv.ParseBool(val)
	if err != nil {
		return defVal
	}

	return v
}

// SetForTesting is a helper function for tests that sets and then resets env value to original.
// Defer should be called right after call of this function.
// Example:
// SetForTesting(t, "SOME_ENV", "new_value")
func SetForTesting(tb testing.TB, key string, value interface{}) {
	original := os.Getenv(key)

	if err := os.Setenv(key, fmt.Sprintf("%v", value)); err != nil {
		tb.Fatalf("failed to set flag: %v", err)
	}

	tb.Cleanup(func() {
		if err := os.Setenv(key, original); err != nil {
			tb.Fatalf("failed to revert flag: %v", err)
		}
	})
}

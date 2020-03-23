// Package env provide functionality for loading values from environment variables.
package env

import (
	"os"
	"strconv"
)

// GetStringOrDefault returns string environment variable value or passed default.
func GetStringOrDefault(key string, defVal string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return defVal
	}

	return val
}

// GetBoolOrDefault returns boolean environment variable value or passed default.
func GetBoolOrDefault(key string, defVal bool) bool {
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

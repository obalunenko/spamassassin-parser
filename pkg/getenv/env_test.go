package getenv_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oleg-balunenko/spamassassin-parser/pkg/getenv"
)

func TestGetBoolOrDefault(t *testing.T) {
	var (
		key    = "TestGetBoolOrDefault"
		defVal = true
	)

	t.Run("Env not set - default returned", func(t *testing.T) {
		got := getenv.BoolOrDefault(key, defVal)
		assert.Equal(t, defVal, got)
	})

	t.Run("Env set - value from env returned", func(t *testing.T) {
		want := false

		reset := getenv.SetForTesting(t, key, want)

		defer func() {
			reset()
		}()

		got := getenv.BoolOrDefault("TestGetBoolOrDefault", defVal)
		assert.Equal(t, want, got)
	})
}

func TestGetStringOrDefault(t *testing.T) {
	var (
		key    = "TestGetBoolOrDefault"
		defVal = "defVal"
	)

	t.Run("Env not set - default returned", func(t *testing.T) {
		got := getenv.StringOrDefault(key, defVal)
		assert.Equal(t, defVal, got)
	})

	t.Run("Env set - value from env returned", func(t *testing.T) {
		want := "MyVal"

		reset := getenv.SetForTesting(t, key, want)

		defer func() {
			reset()
		}()

		got := getenv.StringOrDefault("TestGetBoolOrDefault", defVal)
		assert.Equal(t, want, got)
	})
}

func TestSetForTesting(t *testing.T) {
	key := "TestSetForTesting"

	original := os.Getenv(key)
	assert.Equal(t, original, "", "Check that variable not set.")

	// Set new value for variable.
	reset := getenv.SetForTesting(t, key, "NEW_VAL")

	val := os.Getenv(key)
	assert.Equal(t, val, "NEW_VAL", "Check that variable changed value.")

	// Check that after calling reset - variable returns to original state
	reset()

	val = os.Getenv(key)

	assert.Equal(t, original, val, "Check that after calling reset - variable returns to original state")
}

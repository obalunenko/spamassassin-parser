package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBoolOrDefault(t *testing.T) {
	var (
		key    = "TestGetBoolOrDefault"
		defVal = true
	)

	t.Run("Env not set - default returned", func(t *testing.T) {
		got := GetBoolOrDefault(key, defVal)
		assert.Equal(t, defVal, got)
	})

	t.Run("Env set - value from env returned", func(t *testing.T) {
		want := false

		reset := SetForTesting(t, key, want)

		defer func() {
			reset()
		}()

		got := GetBoolOrDefault("TestGetBoolOrDefault", defVal)
		assert.Equal(t, want, got)
	})
}

func TestGetStringOrDefault(t *testing.T) {
	var (
		key    = "TestGetBoolOrDefault"
		defVal = "defVal"
	)

	t.Run("Env not set - default returned", func(t *testing.T) {
		got := GetStringOrDefault(key, defVal)
		assert.Equal(t, defVal, got)
	})

	t.Run("Env set - value from env returned", func(t *testing.T) {
		want := "MyVal"

		reset := SetForTesting(t, key, want)

		defer func() {
			reset()
		}()

		got := GetStringOrDefault("TestGetBoolOrDefault", defVal)
		assert.Equal(t, want, got)
	})
}

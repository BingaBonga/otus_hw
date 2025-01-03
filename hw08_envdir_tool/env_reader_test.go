package main

import (
	"testing"

	//nolint:depguard
	"github.com/stretchr/testify/require"
)

var testdataDir = "testdata/env"

func TestReadDir(t *testing.T) {
	t.Run("read test data", func(t *testing.T) {
		env, err := ReadDir(testdataDir)
		require.NoError(t, err)
		require.Equal(t, len(env), 5)
		require.Equal(t, env["BAR"].Value, "bar")
	})
}

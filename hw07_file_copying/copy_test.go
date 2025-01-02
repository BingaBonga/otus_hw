package main

import (
	"os"
	"testing"

	//nolint:depguard
	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("full copy case", func(t *testing.T) {
		file, err := os.CreateTemp("tmp", "out.txt")
		require.NoError(t, err)

		err = Copy("testdata/input.txt", file.Name(), 0, 0)
		require.NoError(t, err)
	})
}

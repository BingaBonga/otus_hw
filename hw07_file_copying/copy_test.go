package main

import (
	"os"
	"testing"

	//nolint:depguard
	"github.com/stretchr/testify/require"
)

var (
	fromFileName = "testdata/input.txt"
	toFileName   = "output_test.txt"
)

func TestCopy(t *testing.T) {
	t.Run("full copy", func(t *testing.T) {
		err := Copy(fromFileName, toFileName, 0, 0)
		require.NoError(t, err)

		fromFileStat, err := os.Stat(fromFileName)
		require.NoError(t, err)

		toFileStat, err := os.Stat(toFileName)
		require.NoError(t, err)
		require.Equal(t, fromFileStat.Size(), toFileStat.Size())

		err = os.Remove(toFileName)
		require.NoError(t, err)
	})

	t.Run("limit copy", func(t *testing.T) {
		err := Copy(fromFileName, toFileName, 0, 100)
		require.NoError(t, err)

		toFileStat, err := os.Stat(toFileName)
		require.NoError(t, err)

		require.Equal(t, int64(100), toFileStat.Size())

		err = os.Remove(toFileName)
		require.NoError(t, err)
	})

	t.Run("offset copy", func(t *testing.T) {
		err := Copy(fromFileName, toFileName, 100, 0)
		require.NoError(t, err)

		fromFileStat, err := os.Stat(fromFileName)
		require.NoError(t, err)

		toFileStat, err := os.Stat(toFileName)
		require.NoError(t, err)

		require.Equal(t, fromFileStat.Size()-100, toFileStat.Size())

		err = os.Remove(toFileName)
		require.NoError(t, err)
	})

	t.Run("complex copy", func(t *testing.T) {
		err := Copy(fromFileName, toFileName, 100, 1000)
		require.NoError(t, err)

		toFileStat, err := os.Stat(toFileName)
		require.NoError(t, err)

		require.Equal(t, int64(1000), toFileStat.Size())

		err = os.Remove(toFileName)
		require.NoError(t, err)
	})

	t.Run("offset eq file size copy", func(t *testing.T) {
		fromFileStat, err := os.Stat(fromFileName)
		require.NoError(t, err)

		err = Copy(fromFileName, toFileName, fromFileStat.Size(), 0)
		require.NoError(t, err)

		toFileStat, err := os.Stat(toFileName)
		require.NoError(t, err)

		require.Equal(t, int64(0), toFileStat.Size())

		err = os.Remove(toFileName)
		require.NoError(t, err)
	})

	t.Run("limit bigger then file size copy", func(t *testing.T) {
		fromFileStat, err := os.Stat(fromFileName)
		require.NoError(t, err)

		err = Copy(fromFileName, toFileName, 0, fromFileStat.Size()+1000)
		require.NoError(t, err)

		toFileStat, err := os.Stat(toFileName)
		require.NoError(t, err)

		require.Equal(t, fromFileStat.Size(), toFileStat.Size())

		err = os.Remove(toFileName)
		require.NoError(t, err)
	})

	t.Run("file not found from err", func(t *testing.T) {
		err := Copy("wrong_path", toFileName, 0, 0)
		require.Equal(t, err, ErrFileDoesNotExist)

		err = os.Remove(toFileName)
		require.Error(t, err)
	})

	t.Run("file read err", func(t *testing.T) {
		err := Copy("testdata", toFileName, 0, 0)
		require.Equal(t, err, ErrUnsupportedFile)

		err = os.Remove(toFileName)
		require.NoError(t, err)
	})

	t.Run("offset file err", func(t *testing.T) {
		fromFileStat, err := os.Stat(fromFileName)
		require.NoError(t, err)

		err = Copy(fromFileName, toFileName, fromFileStat.Size()+1000, 0)
		require.Equal(t, err, ErrOffsetExceedsFileSize)

		err = os.Remove(toFileName)
		require.NoError(t, err)
	})
}

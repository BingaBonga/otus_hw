package main

import (
	"os"
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
		require.Equal(t, env["BAR"], EnvValue{"bar", false})
		require.Equal(t, env["EMPTY"], EnvValue{"", false})
		require.Equal(t, env["FOO"], EnvValue{"   foo\nwith new line", false})
		require.Equal(t, env["HELLO"], EnvValue{"\"hello\"", false})
		require.Equal(t, env["UNSET"], EnvValue{"", true})
	})

	t.Run("test data dir not exists", func(t *testing.T) {
		env, err := ReadDir("notExistDir")
		require.Error(t, err)
		require.Nil(t, env)
	})

	t.Run("test data dir is empty", func(t *testing.T) {
		err := os.Mkdir("tmp", 0o777)
		require.NoError(t, err)

		env, err := ReadDir("tmp")
		require.NoError(t, err)
		require.Equal(t, len(env), 0)

		err = os.Remove("tmp")
		require.NoError(t, err)
	})
}

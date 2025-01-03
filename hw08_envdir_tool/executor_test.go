package main

import (
	"testing"

	//nolint:depguard
	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("test cmd", func(t *testing.T) {
		env := map[string]EnvValue{
			"FOO": {"foo", false},
			"BAR": {"bar", false},
		}
		cmd := []string{"/bin/bash", "-c", "echo $FOO$BAR"}
		code := RunCmd(cmd, env)

		require.Equal(t, 0, code)
	})
}

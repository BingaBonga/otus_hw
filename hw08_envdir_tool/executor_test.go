package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunCmd(t *testing.T) {
	t.Run("test cmd", func(t *testing.T) {
		env := map[string]EnvValue{
			"FOO": {"foo", false},
			"BAR": {"bar", false},
		}
		cmd := []string{"/bin/bash", "echo -e \"${FOO}${BAR}\""}
		code := RunCmd(cmd, env)

		assert.Equal(t, 0, code)
	})
}

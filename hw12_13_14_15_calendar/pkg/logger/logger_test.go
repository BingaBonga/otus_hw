package logger

import (
	"os"
	"testing"

	//nolint:depguard
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestLogger(t *testing.T) {
	t.Run("logger created", func(t *testing.T) {
		var level zapcore.Level
		logger, err := New(level, os.TempDir()+"/test.log")
		require.NoError(t, err)
		require.NotNil(t, logger)
	})
}

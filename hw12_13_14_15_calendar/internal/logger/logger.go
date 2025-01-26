package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level zapcore.Level, filePath string) (*zap.Logger, error) {
	err := makeLogFile(filePath)
	if err != nil {
		return nil, err
	}

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.OutputPaths = []string{filePath, "stdout"}
	loggerConfig.ErrorOutputPaths = []string{filePath, "stderr"}
	loggerConfig.Level = zap.NewAtomicLevelAt(level)
	loggerConfig.EncoderConfig.EncodeTime = timeEncoder

	return loggerConfig.Build()
}

func makeLogFile(path string) error {
	logFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if os.IsNotExist(err) {
		logFile, err = os.Create(path)
	}
	if err != nil {
		return err
	}
	defer logFile.Close()

	return nil
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.UTC().Format("Jan 02 03:04:05"))
}

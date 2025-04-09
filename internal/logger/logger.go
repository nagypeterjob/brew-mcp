package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewStructuredLogger(level string) (*zap.SugaredLogger, error) {
	atom := zap.NewAtomicLevel()
	if err := atom.UnmarshalText([]byte(level)); err != nil {
		return nil, fmt.Errorf("could parse log level: %w", err)
	}

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level = atom

	loggerConfig.OutputPaths = []string{"stdout"}
	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("new zap logger: %w", err)
	}

	return logger.Sugar(), nil
}

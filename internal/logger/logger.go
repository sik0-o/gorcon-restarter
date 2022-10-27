package logger

import (
	"fmt"
	"log"

	"github.com/sik0-o/gorcon-restarter/v2/internal/environment"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New returns new zap logger for the given service that is launched
// in the given environment. Optionally, non-empty log level will be set.
func New(
	version string,
	env environment.Env,
	level string,
) (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	if env.IsProduction() {
		config = zap.NewProductionConfig()
	} else {
		config.OutputPaths = append(config.OutputPaths, `tmp/output.log`)
		config.ErrorOutputPaths = append(config.ErrorOutputPaths, `tmp/error.log`)
	}

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if err := config.Level.UnmarshalText([]byte(level)); err != nil && level != "" {
		log.Printf("failed to set log level %q: %v", level, err)
	}

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	logger = logger.With(
		zap.String("version", version),
		zap.String("environment", env.String()),
	)

	return logger, nil
}

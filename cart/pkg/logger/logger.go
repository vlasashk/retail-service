package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New(lvl, serviceName string) (zerolog.Logger, error) {
	logLevel, err := zerolog.ParseLevel(lvl)
	if err != nil {
		return zerolog.Logger{}, err
	}

	logger := zerolog.New(os.Stderr).
		Level(logLevel).
		With().
		Str("name", serviceName).
		Timestamp().
		Caller().
		Logger()

	zerolog.DefaultContextLogger = &logger

	return logger, nil
}

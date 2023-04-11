package infrustructure

import (
	"os"

	"github.com/rs/zerolog"
)

func NewZerologLogger(serviceName string) *zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	skipFrameCount := 3
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("service", serviceName).CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).Logger()
	return &logger
}

package config

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ConfigureLogging() {
	var writers []io.Writer

	level, err := zerolog.ParseLevel(Get().Logging.Level)
	conf := Get()
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	if conf.Logging.Console {
		consoleWriter := zerolog.NewConsoleWriter()
		consoleWriter.NoColor = !conf.Logging.Color
		writers = append(writers, consoleWriter)
	}

	if len(writers) == 0 {
		writers = append(writers, os.Stderr)
	}

	log.Logger = zerolog.New(io.MultiWriter(writers...)).With().Timestamp().Logger()
	log.Logger = log.Logger.Level(level)
	zerolog.SetGlobalLevel(level)
	zerolog.DefaultContextLogger = &log.Logger
}

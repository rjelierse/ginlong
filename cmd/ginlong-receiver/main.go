package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"

	"github.com/rjelierse/ginlong/internal/server"
)

type Config struct {
	LogLevel        string `default:"info" split_words:"true"`
	ReceiverAddress string `default:":10000" split_words:"true"`
}

func main() {
	var config Config
	envconfig.MustProcess("GINLONG", &config)

	level, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		panic(err)
	}
	log := zerolog.New(os.Stdout).With().Timestamp().Logger().Level(level)

	receiver := server.New(log)
	if err := receiver.Listen(config.ReceiverAddress); err != nil {
		log.Fatal().Err(err).Msg("Failed to start TCP server")
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-done:
		log.Info().Str("signal", sig.String()).Msg("Terminating after receiving signal")
	}

	receiver.Shutdown()
}

package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	server2 "github.com/rjelierse/ginlong/internal/server"

	"github.com/rs/zerolog"
)

var logLevel = flag.String("log", "info", "Set the log output level")

func main() {
	flag.Parse()

	level, err := zerolog.ParseLevel(*logLevel)
	if err != nil {
		panic(err)
	}
	log := zerolog.New(os.Stdout).With().Logger().Level(level)

	server := server2.NewServer(log)
	if err := server.Listen(":10000"); err != nil {
		log.Fatal().Err(err).Msg("Failed to start TCP server")
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-done:
		log.Info().Str("signal", sig.String()).Msg("Terminating after receiving signal")
	}

	server.Shutdown()
}

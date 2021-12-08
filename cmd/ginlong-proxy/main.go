package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/rjelierse/ginlong/internal/proxy"
)

type config struct {
	ListenAddress   string `envconfig:"LISTEN_ADDR" default:"0.0.0.0:10000"`
	UpstreamAddress string `envconfig:"UPSTREAM_ADDR" default:"47.88.8.200:10000"`
}

func main() {
	var conf config
	if err := envconfig.Process("", &conf); err != nil {
		log.Fatal("Failed to parse config:", err)
	}

	log.Printf("%s - running on %s (arch: %s os: %s)", os.Args[0], runtime.Version(), runtime.GOARCH, runtime.GOOS)

	ctx, done := context.WithCancel(context.Background())

	// Run proxy server.
	server := proxy.New(conf.ListenAddress, conf.UpstreamAddress)
	go func() {
		if err := server.Listen(ctx); err != nil {
			log.Println("Error:", err)
		}
		done()
	}()

	// Wait for process termination.
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt)
	go func() {
		<-sigC
		log.Println("Received signal to shut down.")

		server.Shutdown()

		// Wait for server to shut down or terminate after 10 seconds.
		select {
		case <-time.After(time.Second * 10):
			os.Exit(1)
		case <-ctx.Done():
			// no-op
		}
	}()

	// Process measurements.
	go func() {
		for m := range server.Measurements() {
			fmt.Printf("Measurement%+v\n", m)
		}
	}()

	// Wait for graceful shutdown.
	<-ctx.Done()

	log.Println("Shutdown completed.")
}

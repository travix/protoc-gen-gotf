package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const shutdownDelay = 5 * time.Second

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// create service servers
	shutdown := make(chan bool)
	done := make(chan bool)
	// start servers
	NewServer(&Server{}, shutdown, done)
	handleShutdown(shutdown, done)
	log.Info().Msg("shutdown complete")
}

// handleShutdown will wait for syscall.SIGINT, syscall.SIGTERM
// Once received will gracefully shutdown server in 10 seconds else will.
func handleShutdown(closed chan<- bool, done <-chan bool) {
	infoLogger := log.Level(zerolog.InfoLevel)
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	// received shutdown
	infoLogger.Info().Msgf("received %s(%#v)! shutting down", sig.String(), sig)
	close(closed)
	infoLogger.Info().Msg("waining for grpc to stop")
	select {
	case <-time.After(shutdownDelay):
		return
	case <-done:
	}
}

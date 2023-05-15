package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// TestMain runs before all tests in this package.
func TestMain(m *testing.M) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, PartsExclude: []string{zerolog.TimestampFieldName}})
	cmd := exec.Command("/bin/bash", "-c", "./gen.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal().Err(err).Msg("failed to run ./gen.sh")
	}
	os.Exit(m.Run())
}

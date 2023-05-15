package plugin

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// TestMain runs before all tests in this package.
func TestMain(m *testing.M) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, PartsExclude: []string{zerolog.TimestampFieldName}})
	os.Exit(m.Run())
}

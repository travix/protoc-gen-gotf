package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/travix/protoc-gen-goterraform/plugin"
)

//go:generate ./gen.sh
func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, FormatMessage: formatter})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	plugin.Run()
}

func formatter(msg any) string {
	if msg == nil {
		return ""
	}
	return fmt.Sprintf("[%s] %s", plugin.Name, msg)
}

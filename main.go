package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-goterraform/plugin"
)

//go:generate ./gen.sh
func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, FormatMessage: formatter, PartsExclude: []string{zerolog.TimestampFieldName}})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	protogen.Options{}.Run(plugin.Run)
}

func formatter(msg any) string {
	if msg == nil {
		return ""
	}
	return fmt.Sprintf("[%s] %s", plugin.Name, msg)
}

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-gotf/plugin"
)

//go:generate ./gen.sh
func main() {
	versionFlag := flag.Bool("version", false, "Print version and exit")
	flag.Parse()
	if *versionFlag {
		fmt.Println(plugin.Version())
		os.Exit(0)
	}
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

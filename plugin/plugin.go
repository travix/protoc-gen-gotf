package plugin

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-goterraform/internal/terraform"
)

const Name = "protoc-gen-goterraform"

// opt is the plugin options. set only once in plugin run.
var opt = &options{logLevel: zerolog.WarnLevel, suffix: ".tf.pb.go"}

// options that can be passed to plugin.
type options struct {
	logLevel zerolog.Level // logLevel mode
	module   string        // module name of go tf code being generated
	prefix   string        // prefix for go tf files
	suffix   string        // suffix for go tf files
}

// Plugin is the interface for the goterraform plugin.
type Plugin interface {
	Run(Input) ([]*protogen.GeneratedFile, error)
}

type plugin struct {
	*protogen.Plugin
}

// Run creates new goterraform plugin and runs it.
func Run(gen *protogen.Plugin) error {
	var err error
	SetOptions(gen.Request.GetParameter())
	p := &plugin{Plugin: gen}
	var in Input
	if in, err = NewInput(gen, terraform.NewSynthesizer(protogen.GoImportPath(opt.module))); err != nil {
		return err
	} else if in == nil {
		log.Debug().Msg("no files to generate")
		return nil
	}
	log.Debug().Msg("generating files")
	_, err = p.Run(in)
	return err
}

// SetOptions sets the plugin options. Parameters passed should be a comma-separated list example:
//
//	prefix=tf_
//
// or
//
//	prefix=tf_,suffix=.pb.go
//
// Available options:
//
//	log_level= for plugin, available values trace, debug, info, warn, error, fatal, panic, disable. Default is warn.
//	module= module name of go tf code being generated
//	prefix= prefix for go tf files
//	suffix= suffix for go tf files
func SetOptions(params string) {
	for _, param := range strings.Split(params, ",") {
		var value string
		if i := strings.Index(param, "="); i >= 0 {
			value = param[i+1:]
			param = param[0:i]
		}
		switch param {
		case "module":
			opt.module = strings.TrimSpace(value)
		case "prefix":
			opt.prefix = strings.TrimSpace(value)
		case "suffix":
			opt.suffix = strings.TrimSpace(value)
		case "log_level":
			var err error
			if opt.logLevel, err = zerolog.ParseLevel(value); err != nil {
				panic(fmt.Errorf("invalid log_level %s: %w", value, err))
			}
			zerolog.SetGlobalLevel(opt.logLevel)
		case "", "paths", "annotate_code":
			// Ignore go plugin options.
		default:
			if param[0] == 'M' {
				// Ignore go plugin options.
				continue
			}
			log.Warn().Msgf("ignoring %s, unknown option", param)
		}
	}
}

// Run executes the plugin and generates the go tf code files.
func (p *plugin) Run(Input) ([]*protogen.GeneratedFile, error) {
	// TODO: generate all files
	return nil, nil
}

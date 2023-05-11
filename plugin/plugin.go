package plugin

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/compiler/protogen"
)

const Name = "protoc-gen-goterraform"

var (
	// opt is the plugin options. set only once in plugin run.
	opt = &options{lock: &atomic.Bool{}, suffix: ".tf.pb.go"}
)

// Plugin is the interface for the goterraform plugin.
type Plugin interface {
	Run() error
}

// options that can be passed to plugin.
type options struct {
	debug  bool         // debug mode
	lock   *atomic.Bool // lock prevents opt from being set more than once which prevents Run from being called more than once in single plugin run
	module string       // module name of go tf code being generated
	prefix string       // prefix for go tf files
	suffix string       // suffix for go tf files
}

type plugin struct {
	*protogen.Plugin
}

// NewPlugin returns goterraform plugin.
func NewPlugin(gen *protogen.Plugin) (Plugin, error) {
	p := &plugin{Plugin: gen}
	for _, file := range p.Plugin.Files {
		if !file.Generate {
			log.Debug().Msgf("skipped %s not in requested files", file.Proto.GetName())
			continue
		}
		log.Debug().Msgf("parsing %s files", file.Proto.GetName())
		// TODO: parse all files
	}
	return p, nil
}

// Run creates new goterraform plugin and runs it.
func Run() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		if err := SetOptions(gen.Request.GetParameter()); err != nil {
			return err
		}
		p, err := NewPlugin(gen)
		if err != nil {
			return err
		}
		log.Debug().Msg("generating fieles")
		return p.Run()
	})
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
//	debug=true enable debug logging
//	module= module name of go tf code being generated
//	prefix= prefix for go tf files
//	suffix= suffix for go tf files
func SetOptions(params string) error {
	if opt.lock.Load() {
		return fmt.Errorf("can't set plugin options again, this should not happen")
	}
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
		case "debug":
			opt.debug = true
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
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
	opt.lock.Store(true)
	return nil
}

// Run executes the plugin and generates the go tf code files.
func (p *plugin) Run() error {
	// TODO: generate all files
	return nil
}

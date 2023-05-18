package plugin

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-gotf/extension"
	"github.com/travix/protoc-gen-gotf/extensionimpl"
	"github.com/travix/protoc-gen-gotf/gocode"
	"github.com/travix/protoc-gen-gotf/pb"
)

const Name = "protoc-gen-gotf"

// opt is the plugin options. set only once in plugin run.
var opt = &options{logLevel: zerolog.WarnLevel}

// options that can be passed to plugin.
type options struct {
	logLevel zerolog.Level // logLevel mode
	module   string        // module name of go tf code being generated
}

// Plugin is the interface for the gotf plugin.
type Plugin interface {
	Run(Input) ([]*protogen.GeneratedFile, error)
}

type plugin struct {
	*protogen.Plugin
	gocode.Writer
}

// Run creates new gotf plugin and runs it.
func Run(gen *protogen.Plugin) error {
	var err error
	SetOptions(gen.Request.GetParameter())
	log.Debug().Str("module", opt.module).Msg("plugin options")
	p := &plugin{Plugin: gen}
	var in Input
	if in, err = NewInput(extensionimpl.NewSynthesizer(opt.module), gen); err != nil {
		return err
	}
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
func (p *plugin) Run(in Input) ([]*protogen.GeneratedFile, error) {
	if in == nil {
		log.Warn().Msg("no input to plugin")
		return nil, nil
	}
	if in.Provider() == nil {
		log.Warn().Msgf("no provider found: %s option not set in any of the proto files", pb.E_Provider.TypeDescriptor().FullName())
		return nil, nil
	}
	if len(in.Resources())+len(in.Datasources()) == 0 {
		log.Warn().Msgf("no resources or datasources found: %s or %s option not set in any of the proto files", pb.E_Resource.TypeDescriptor().FullName(), pb.E_Datasource.TypeDescriptor().FullName())
		return nil, nil
	}
	log.Debug().Msg("generating files")
	provider := in.Provider()
	importPath := provider.ImportPath()
	p.Writer = gocode.NewWriter(provider.PbImportPath(), importPath, provider.PbPackageName(), provider.PackageName())
	var files []*protogen.GeneratedFile
	var err error
	generatedFiles := make([]*protogen.GeneratedFile, 0)
	if files, err = p.genBlocks(in, importPath); err != nil {
		return nil, err
	}
	generatedFiles = append(generatedFiles, files...)
	if files, err = p.genDependencies(in, importPath); err != nil {
		return nil, err
	}
	generatedFiles = append(generatedFiles, files...)
	hasServiceClient := p.HasServiceClient(in.AllBlocks())
	if files, err = p.genProvider(importPath, provider, hasServiceClient); err != nil {
		return nil, err
	}
	return append(generatedFiles, files...), nil
}

func (p *plugin) genProvider(importPath protogen.GoImportPath, provider extension.Provider, hasServiceClient bool) ([]*protogen.GeneratedFile, error) {
	filename := filepath.Join(string(importPath), "provider.go")
	file := p.NewGeneratedFile(filename, importPath)
	if err := p.WriteProvider(filename, file, provider, hasServiceClient); err != nil {
		return nil, err
	}
	return []*protogen.GeneratedFile{file}, nil // just one file but... keep the same pattern
}

func (p *plugin) genBlocks(in Input, importPath protogen.GoImportPath) ([]*protogen.GeneratedFile, error) {
	generatedFiles := make([]*protogen.GeneratedFile, 0)
	for _, resource := range in.Resources() {
		filename := filepath.Join(string(importPath), resource.Filename())
		file := p.NewGeneratedFile(filename, importPath)
		if err := p.WriteResource(filename, file, resource); err != nil {
			return nil, err
		}
		generatedFiles = append(generatedFiles, file)
	}
	for _, datasource := range in.Datasources() {
		filename := filepath.Join(string(importPath), datasource.Filename())
		file := p.NewGeneratedFile(filename, importPath)
		if err := p.WriteDatasource(filename, file, datasource); err != nil {
			return nil, err
		}
		generatedFiles = append(generatedFiles, file)
	}
	return generatedFiles, nil
}

func (p *plugin) genDependencies(in Input, importPath protogen.GoImportPath) ([]*protogen.GeneratedFile, error) {
	generatedFiles := make([]*protogen.GeneratedFile, 0)
	for fullName, models := range in.Dependencies() {
		var generatedFilename string
		for _, f := range p.Files {
			if f.Desc.FullName() == fullName {
				generatedFilename = f.GeneratedFilenamePrefix + ".tf.pb.go"
				break
			}
		}
		if generatedFilename == "" {
			return nil, fmt.Errorf("could not find proto file for %s, this shouldn't have happened", fullName)
		}
		file := p.NewGeneratedFile(generatedFilename, importPath)
		if err := p.WriteDependency(generatedFilename, file, models...); err != nil {
			return nil, err
		}
		generatedFiles = append(generatedFiles, file)
	}
	return generatedFiles, nil
}

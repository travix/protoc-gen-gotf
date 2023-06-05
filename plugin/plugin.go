package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/travix/protoc-gen-gotf/extension"
	"github.com/travix/protoc-gen-gotf/extensionimpl"
	"github.com/travix/protoc-gen-gotf/gocode"
	"github.com/travix/protoc-gen-gotf/pb"
	"github.com/travix/protoc-gen-gotf/structtag"
)

const Name = "protoc-gen-gotf"

// opt is the plugin options. set only once in plugin run.
var opt = &options{logLevel: zerolog.WarnLevel}

// version in format "{version} {buildDate} {commit} {commitDate}", populated at compile time.
var version = "local 1970-01-01T00:00:00Z (e83c516 1970-01-01T00:00:00Z)"

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
	structtag.Tagger
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
	gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	return err
}

func Version() string {
	return version
}

// SetOptions sets the plugin options. Parameters passed should be a comma-separated list example:
// Available options:
//
//	log_level= for plugin, available values trace, debug, info, warn, error, fatal, panic, disable. Default is warn.
//	module= module name of go tf code being generated
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
	packageData := provider.PackageData()
	var err error
	// Available options:
	p.Writer, err = gocode.NewWriter(opt.module, packageData, strings.Split(version, " ")[0])
	if err != nil {
		return nil, err
	}
	p.Tagger = structtag.NewTagger("tfsdk")
	var files []*protogen.GeneratedFile
	generatedFiles := make([]*protogen.GeneratedFile, 0)
	if files, err = p.genBlocks(in, packageData); err != nil {
		return nil, err
	}
	generatedFiles = append(generatedFiles, files...)
	if files, err = p.genDependencies(in, packageData); err != nil {
		return nil, err
	}
	generatedFiles = append(generatedFiles, files...)
	hasServiceClient := false
	for _, block := range in.AllBlocks() {
		if block.HasServiceClient() {
			hasServiceClient = true
			break
		}
	}
	if files, err = p.genProvider(packageData, provider, hasServiceClient); err != nil {
		return nil, err
	}
	return append(generatedFiles, files...), nil
}

func (p *plugin) genProvider(packageData extension.PackageData, provider extension.Provider, hasServiceClient bool) ([]*protogen.GeneratedFile, error) {
	filename := filepath.Join(string(packageData.ProviderImportPath), provider.Filename())
	file := p.NewGeneratedFile(filename, packageData.ProviderImportPath)
	if err := p.WriteProvider(filename, file, provider, hasServiceClient); err != nil {
		return nil, err
	}
	return []*protogen.GeneratedFile{file}, nil // just one file but... keep the same pattern
}

func (p *plugin) genBlocks(in Input, packageData extension.PackageData) ([]*protogen.GeneratedFile, error) {
	generatedFiles := make([]*protogen.GeneratedFile, 0)
	for _, block := range in.AllBlocks() {
		filename := filepath.Join(string(packageData.ProviderImportPath), block.Filename())
		file := p.NewGeneratedFile(filename, packageData.ProviderImportPath)
		var writeBlock gocode.Write
		var writeExec gocode.Write
		if block.IsResource() {
			writeBlock = p.WriteResource
			writeExec = p.WriteResourceExec
		} else {
			writeBlock = p.WriteDatasource
			writeExec = p.WriteDatasourceExec
		}
		if err := writeBlock(filename, file, block); err != nil {
			return nil, err
		}
		generatedFiles = append(generatedFiles, file)
		if packageData.ExecImportPath != "" {
			filename = filepath.Join(string(packageData.ExecImportPath), block.ExecFilename())
			existingFile := strings.TrimPrefix(filename, opt.module+"/")
			if _, err := os.Stat(existingFile); err == nil {
				log.Debug().Msgf("skipping %s, Exec file already exists", existingFile)
				continue
			}
			file = p.NewGeneratedFile(filename, packageData.ExecImportPath)
			if err := writeExec(filename, file, block); err != nil {
				return nil, err
			}
			generatedFiles = append(generatedFiles, file)
		}
	}
	return generatedFiles, nil
}

func (p *plugin) genDependencies(in Input, packageData extension.PackageData) ([]*protogen.GeneratedFile, error) {
	generatedFiles := make([]*protogen.GeneratedFile, 0)
	for fullName, models := range in.Dependencies() {
		var generatedFilename string
		for _, f := range p.Files {
			if *f.Proto.Name == fullName {
				generatedFilename = f.GeneratedFilenamePrefix + "_tf.pb.go"
				break
			}
		}
		if generatedFilename == "" {
			return nil, fmt.Errorf("could not find proto file for %s, this shouldn't have happened", fullName)
		}
		file := p.NewGeneratedFile(generatedFilename, packageData.ProviderImportPath)
		if err := p.WriteDependency(generatedFilename, file, models...); err != nil {
			return nil, err
		}
		protoFilePath := strings.ReplaceAll(generatedFilename, opt.module+"/", "")
		protoFilePath = strings.ReplaceAll(protoFilePath, "_tf.pb.go", ".pb.go")
		if err := p.Tag(protoFilePath, models); err != nil {
			return nil, err
		}
		generatedFiles = append(generatedFiles, file)
	}
	return generatedFiles, nil
}

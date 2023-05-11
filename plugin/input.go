package plugin

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/runtime/protoimpl"

	"github.com/travix/protoc-gen-goterraform/internal/terraform"
	"github.com/travix/protoc-gen-goterraform/pb"
)

var _ Input = &input{}

type Input interface {
	Provider() terraform.Provider
	Resources() []terraform.Block
	Datasources() []terraform.Block
	Dependencies() []terraform.TypeValue
}

type input struct {
	providerProtoFile string
	provider          terraform.Provider
	resources         []terraform.Block
	datasources       []terraform.Block
	dependencies      []terraform.TypeValue
}

func (in *input) Provider() terraform.Provider {
	return in.provider
}

func (in *input) Resources() []terraform.Block {
	return in.resources
}

func (in *input) Datasources() []terraform.Block {
	return in.datasources
}

func (in *input) Dependencies() []terraform.TypeValue {
	return in.dependencies
}

// set provider if not already set, and error if multiple providers found.
func (in *input) funcName1(file *protogen.File, synthesizer terraform.Synthesizer) error {
	provider, err := synthesizer.Provider(file.Desc)
	if err != nil {
		return err
	}
	if in.provider == nil {
		in.provider = provider
		in.providerProtoFile = file.Desc.Path()
	} else if provider != nil {
		return fmt.Errorf("error multiple providers: %s options found in %s and %s", pb.E_Provider.TypeDescriptor().FullName(), in.providerProtoFile, file.Desc.Path())
	}
	return nil
}

func (in *input) funcName(file *protogen.File, synthesizer terraform.Synthesizer) error {
	for _, message := range file.Messages {
		for _, blockType := range []*protoimpl.ExtensionInfo{pb.E_Resource, pb.E_Datasource} {
			block, dependencies, err := synthesizer.Block(message, blockType)
			if err != nil {
				return err
			}
			if block != nil {
				in.addBlock(block, blockType)
				in.setDependencies(dependencies)
			}
		}
	}
	return nil
}

func (in *input) setDependencies(dependencies []terraform.TypeValue) {
	// TODO: check for duplicates
	in.dependencies = append(in.dependencies, dependencies...)
}

func (in *input) addBlock(block terraform.Block, blockType *protoimpl.ExtensionInfo) {
	switch blockType {
	case pb.E_Resource:
		in.resources = append(in.resources, block)
	case pb.E_Datasource:
		in.datasources = append(in.datasources, block)
	}
}

// NewInput returns goterraform *input.
func NewInput(gen *protogen.Plugin, synthesizer terraform.Synthesizer) (Input, error) {
	in := &input{}
	for _, file := range gen.Files {
		if !file.Generate {
			log.Debug().Msgf("skipped %s not in requested files", file.Proto.GetName())
			continue
		}
		log.Debug().Msgf("parsing %s files", file.Proto.GetName())
		if err := in.funcName1(file, synthesizer); err != nil {
			return nil, err
		}
		if err := in.funcName(file, synthesizer); err != nil {
			return nil, err
		}
	}
	if in.provider == nil {
		log.Warn().Msgf("no provider found: %s option not set in any of the proto files", pb.E_Provider.TypeDescriptor().FullName())
		return nil, nil
	}
	if len(in.resources) == 0 && len(in.datasources) == 0 {
		log.Warn().Msgf("no resources or datasources found: %s or %s option not set in any of the proto files", pb.E_Resource.TypeDescriptor().FullName(), pb.E_Datasource.TypeDescriptor().FullName())
		return nil, nil
	}
	return in, nil
}

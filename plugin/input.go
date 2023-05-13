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
	Dependencies() []*protogen.Message
}

type input struct {
	providerProtoFile string
	provider          terraform.Provider
	resources         []terraform.Block
	datasources       []terraform.Block
	dependencies      []*protogen.Message
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

func (in *input) Dependencies() []*protogen.Message {
	return in.dependencies
}

// setProvider if not already set, and error if multiple providers found.
func (in *input) setProvider(file *protogen.File, synthesizer terraform.Synthesizer) error {
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

// setBlocks sets resources and datasources as input.
func (in *input) setBlocks(file *protogen.File, synthesizer terraform.Synthesizer) error {
	for _, message := range file.Messages {
		for _, blockType := range []*protoimpl.ExtensionInfo{pb.E_Resource, pb.E_Datasource} {
			block, err := synthesizer.Block(message, blockType)
			if err != nil {
				return err
			}
			if block == nil {
				continue
			}
			in.addBlock(block, blockType)
			for _, attr := range block.Attributes() {
				err = in.setDependencies(attr, message)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// setDependencies sets proto messages from non native terraform.TypeValue as dependencies.
func (in *input) setDependencies(attr terraform.Attribute, message *protogen.Message) error {
	dependencies, err := in.dependencyChain(attr, message)
	if err != nil {
		return err
	}
	for _, dependency := range dependencies {
		for _, existing := range in.dependencies {
			// make sure we don't add the same dependency twice
			if existing.GoIdent.String() == dependency.GoIdent.String() {
				continue
			}
			in.dependencies = append(in.dependencies, dependency)
		}
	}
	return nil
}

func (in *input) dependencyChain(dependency terraform.Attribute, parent *protogen.Message) ([]*protogen.Message, error) {
	if dependency.TypeValue().TerraformNative() {
		return nil, nil
	}
	message := dependency.TypeValue().Message()
	if message == nil {
		return nil, fmt.Errorf("error dependency not a native terraform type value %s#%s.%s: expected proto message",
			parent.Location.SourceFile, parent.Location.Path, dependency.Name())
	}
	return in.messageChain(message), nil
}

func (in *input) messageChain(message *protogen.Message) []*protogen.Message {
	dependencies := []*protogen.Message{message}
	for _, field := range message.Fields {
		if field.Message != nil {
			deps := in.messageChain(field.Message)
			dependencies = append(dependencies, deps...)
		}
	}
	return dependencies
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
		if err := in.setProvider(file, synthesizer); err != nil {
			return nil, err
		}
		if err := in.setBlocks(file, synthesizer); err != nil {
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

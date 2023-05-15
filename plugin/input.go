package plugin

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/runtime/protoimpl"

	"github.com/travix/protoc-gen-goterraform/extension"
	"github.com/travix/protoc-gen-goterraform/pb"
)

var _ Input = &input{}

type Input interface {
	Datasources() []extension.Block
	Dependencies() []extension.Model
	Provider() extension.Provider
	Resources() []extension.Block
}

type input struct {
	providerProtoFile string
	provider          extension.Provider
	blocks            []extension.Block
	dependencies      []extension.Model
}

// NewInput returns goterraform *input.
func NewInput(gen *protogen.Plugin, synthesizer extension.Synthesizer) (Input, error) {
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
		if err := in.addBlocks(file, synthesizer); err != nil {
			return nil, err
		}
	}
	if in.provider == nil {
		log.Warn().Msgf("no provider found: %s option not set in any of the proto files", pb.E_Provider.TypeDescriptor().FullName())
		return nil, nil
	}
	if len(in.blocks) == 0 {
		log.Warn().Msgf("no resources or datasources found: %s or %s option not set in any of the proto files", pb.E_Resource.TypeDescriptor().FullName(), pb.E_Datasource.TypeDescriptor().FullName())
		return nil, nil
	}
	return in, nil
}

func (in *input) Datasources() []extension.Block {
	resource := make([]extension.Block, 0)
	for _, block := range in.blocks {
		if block.Type() == pb.E_Datasource {
			resource = append(resource, block)
		}
	}
	return resource
}

func (in *input) Dependencies() []extension.Model {
	return in.dependencies
}

func (in *input) Provider() extension.Provider {
	return in.provider
}

func (in *input) Resources() []extension.Block {
	resource := make([]extension.Block, 0)
	for _, block := range in.blocks {
		if block.Type() == pb.E_Resource {
			resource = append(resource, block)
		}
	}
	return resource
}

// addBlocks sets resources and datasources as input.
func (in *input) addBlocks(file *protogen.File, synthesizer extension.Synthesizer) error {
	blocks := make([]extension.Block, 0)
	for _, message := range file.Messages {
		for _, blockType := range []*protoimpl.ExtensionInfo{pb.E_Resource, pb.E_Datasource} {
			block, err := synthesizer.Block(message, blockType)
			if err != nil {
				return err
			}
			if block == nil {
				continue
			}
			err = in.addDependencies(synthesizer, block.Model())
			if err != nil {
				return err
			}
			blocks = append(blocks, block)
		}
	}
	for _, block := range blocks {
		err := in.setBlock(block)
		if err != nil {
			return err
		}
	}
	return nil
}

func (in *input) messageChain(message *protogen.Message) []*protogen.Message {
	dependencies := make([]*protogen.Message, 0)
	for _, field := range message.Fields {
		if field.Message != nil {
			dependencies = append(dependencies, field.Message)
			dependencies = append(dependencies, in.messageChain(field.Message)...)
		}
	}
	return dependencies
}

func (in *input) setBlock(block extension.Block) error {
	for _, b := range in.blocks {
		name := b.Name()
		typeName := b.TypeName()
		if name == block.Name() && typeName == block.TypeName() {
			return fmt.Errorf("error duplicate terraform blocks: name: %s, type: %s", name, typeName)
		}
	}
	in.blocks = append(in.blocks, block)
	return nil
}

// addDependencies sets proto messages from non native terraform.TypeValue as dependencies.
func (in *input) addDependencies(synthesizer extension.Synthesizer, model extension.Model) error {
	messages := in.messageChain(model.Message())
	newDependencies := []extension.Model{model}
	for _, message := range messages {
		dependency, err := synthesizer.Model(message, false)
		if err != nil {
			return err
		}
		if dependency == nil {
			continue
		}
		newDependencies = append(newDependencies, dependency)
	}
	// add message as dependency since it's not a native terraform type
DEPENDENCIES:
	for _, dependency := range newDependencies {
		for _, existing := range in.dependencies {
			// make sure we don't add the same dependency twice
			if existing.Message().GoIdent.String() == dependency.Message().GoIdent.String() {
				continue DEPENDENCIES
			}
		}
		in.dependencies = append(in.dependencies, dependency)
	}
	return nil
}

// setProvider if not already set, and error if multiple providers found.
func (in *input) setProvider(file *protogen.File, synthesizer extension.Synthesizer) error {
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

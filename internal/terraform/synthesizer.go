package terraform

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/travix/protoc-gen-goterraform/pb"
)

var _ Synthesizer = &synthesizer{}

// Synthesizer for synthesizing Provider or resource, datasource Block from proto files, messages.
//
// //go:generate mockery --name Synthesizer --output ../../mocks.
type Synthesizer interface {
	// Provider synthesis from proto file
	Provider(desc protoreflect.FileDescriptor) (Provider, error)
	Resource(msg *protogen.Message) (Block, []TypeValue, error)
	Datasource(msg *protogen.Message) (Block, []TypeValue, error)
	Block(msg *protogen.Message, blockType protoreflect.ExtensionType) (Block, []TypeValue, error)
}

type synthesizer struct {
	module protogen.GoImportPath
}

func NewSynthesizer(module protogen.GoImportPath) Synthesizer {
	return &synthesizer{module}
}

func (s synthesizer) Provider(desc protoreflect.FileDescriptor) (Provider, error) {
	p := &provider{}
	p.option = getFileOption(desc)
	if p.option == nil {
		return nil, nil
	}
	// TODO: parse option attributes
	return p, nil
}

func (s synthesizer) Resource(msg *protogen.Message) (Block, []TypeValue, error) {
	return s.Block(msg, pb.E_Resource)
}

func (s synthesizer) Datasource(msg *protogen.Message) (Block, []TypeValue, error) {
	return s.Block(msg, pb.E_Datasource)
}

func (s synthesizer) Block(msg *protogen.Message, blockType protoreflect.ExtensionType) (Block, []TypeValue, error) {
	if msg.Desc.IsMapEntry() {
		return nil, nil, nil
	}
	dep := make([]TypeValue, 0)
	b := &block{}
	b.option = getMessageOption(msg.Desc, blockType)
	if b.option == nil {
		return nil, nil, nil
	}
	for _, field := range msg.Fields {
		_, _, _ = NewAttributeFromProto(field, b.option.ExplicitFields)
	}
	// TODO: parse message attributes
	return b, dep, nil
}

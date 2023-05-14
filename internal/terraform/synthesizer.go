package terraform

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/travix/protoc-gen-goterraform/extensions"
	"github.com/travix/protoc-gen-goterraform/pb"
)

var _ extensions.Synthesizer = &synthesizer{}

type synthesizer struct {
	module protogen.GoImportPath
}

func NewSynthesizer(module protogen.GoImportPath) extensions.Synthesizer {
	return &synthesizer{module}
}

func (s synthesizer) Block(msg *protogen.Message, blockType protoreflect.ExtensionType) (extensions.Block, error) {
	return NewBlock(s, msg, blockType)
}

func (s synthesizer) BlockAttribute(field *protogen.Field, explicit bool) (extensions.Attribute, error) {
	return NewBlockAttribute(s, field, explicit)
}

func (s synthesizer) Datasource(msg *protogen.Message) (extensions.Block, error) {
	return NewBlock(s, msg, pb.E_Datasource)
}

func (s synthesizer) FieldOption(desc protoreflect.FieldDescriptor) *pb.Attribute {
	if mo, ok := getOptions[*descriptorpb.FieldOptions](desc); ok {
		option, _ := proto.GetExtension(mo, pb.E_Attribute).(*pb.Attribute)
		return option
	}
	return nil
}

func (s synthesizer) FileOption(desc protoreflect.FileDescriptor) *pb.Option {
	if fo, ok := getOptions[*descriptorpb.FileOptions](desc); ok {
		option, _ := proto.GetExtension(fo, pb.E_Provider).(*pb.Option)
		log.Debug().Interface("asd", proto.GetExtension(fo, pb.E_Provider)).Msgf("getFileOption: %T", proto.GetExtension(fo, pb.E_Provider))
		return option
	}
	return nil
}

func (s synthesizer) MessageOption(desc protoreflect.MessageDescriptor, extType protoreflect.ExtensionType) *pb.Block {
	if mo, ok := getOptions[*descriptorpb.MessageOptions](desc); ok {
		option, _ := proto.GetExtension(mo, extType).(*pb.Block)
		return option
	}
	return nil
}

func (s synthesizer) Provider(desc protoreflect.FileDescriptor) (extensions.Provider, error) {
	p := &provider{}
	p.option = s.FileOption(desc)
	if p.option == nil {
		return nil, nil
	}
	for index, optionAttribute := range p.option.Attributes {
		attr, err := s.ProviderAttribute(optionAttribute)
		if err != nil {
			return nil, fmt.Errorf("error failed to parse %s.attributes[%d] from %s: %w", pb.E_Provider.TypeDescriptor().FullName(), index, desc.Path(), err)
		}
		if attr == nil {
			continue
		}
		p.attributes = append(p.attributes, attr)
	}
	return p, nil
}

func (s synthesizer) ProviderAttribute(option *pb.Attribute) (extensions.Attribute, error) {
	return NewAttribute(option)
}

func (s synthesizer) Resource(msg *protogen.Message) (extensions.Block, error) {
	return NewBlock(s, msg, pb.E_Resource)
}

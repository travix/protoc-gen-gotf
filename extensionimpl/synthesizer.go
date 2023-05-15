package extensionimpl

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/travix/protoc-gen-goterraform/extension"
	"github.com/travix/protoc-gen-goterraform/pb"
)

var _ extension.Synthesizer = &synthesizer{}

type synthesizer struct {
	module string
}

func (s synthesizer) Module() string {
	return s.module
}

func NewSynthesizer(module string) extension.Synthesizer {
	s := &synthesizer{module: module}
	// , getOptions[*pb.Block], getOptions[*pb.Option],
	return s
}

func (s synthesizer) Model(msg *protogen.Message, explicit bool) (extension.Model, error) {
	return NewModel(s, msg, explicit)
}

func (s synthesizer) Block(msg *protogen.Message, blockType protoreflect.ExtensionType) (extension.Block, error) {
	return NewBlock(s, msg, blockType)
}

func (s synthesizer) FieldAttribute(field *protogen.Field, explicit bool) (extension.Attribute, error) {
	option := s.FieldOption(field.Desc)
	if option == nil {
		if explicit {
			return nil, nil
		}
		option = &pb.Attribute{}
	}
	if option.Name == nil {
		option.Name = proto.String(field.GoName)
	}
	// ignore attr on block attributes
	option.Attr = nil
	option.Description = deferToComment(option.Description, field.Comments)
	option.MdDescription = deferToComment(option.MdDescription, field.Comments)
	a, err := NewAttribute(option)
	if err != nil {
		return nil, err
	}
	aTyped, _ := a.(*attribute)
	if aTyped.typeValue, err = inferTypeValue(field); err != nil {
		return nil, err
	}
	if aTyped.typeValue == nil {
		return nil, fmt.Errorf("error failed to infer type value for %s", field.Desc.FullName())
	}
	return a, setSchema(aTyped, field)
}

func (s synthesizer) Datasource(msg *protogen.Message) (extension.Block, error) {
	return NewBlock(s, msg, pb.E_Datasource)
}

func (s synthesizer) FieldOption(desc protoreflect.FieldDescriptor) *pb.Attribute {
	if option, ok := getOptions[*pb.Attribute](desc, pb.E_Attribute); ok {
		return option
	}
	return nil
}

func (s synthesizer) FileOption(desc protoreflect.FileDescriptor) *pb.Option {
	if option, ok := getOptions[*pb.Option](desc, pb.E_Provider); ok {
		log.Debug().Msgf("getFileOption: %T", option)
		return option
	}
	return nil
}

func (s synthesizer) MessageOption(desc protoreflect.MessageDescriptor, extType protoreflect.ExtensionType) *pb.Block {
	if option, ok := getOptions[*pb.Block](desc, extType); ok {
		return option
	}
	return nil
}

func (s synthesizer) Provider(desc protoreflect.FileDescriptor) (extension.Provider, error) {
	return NewProvider(s, desc)
}

func (s synthesizer) Attribute(option *pb.Attribute) (extension.Attribute, error) {
	return NewAttribute(option)
}

func (s synthesizer) Resource(msg *protogen.Message) (extension.Block, error) {
	return NewBlock(s, msg, pb.E_Resource)
}

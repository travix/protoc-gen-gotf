package extensionimpl

import (
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/travix/protoc-gen-gotf/extension"
	"github.com/travix/protoc-gen-gotf/pb"
)

var _ extension.Synthesizer = &synthesizer{}

type synthesizer struct {
	module string
}

func (s synthesizer) MessagePackageName(msg *protogen.Message) protogen.GoPackageName {
	return protogen.GoPackageName(getPkgName(msg.Desc.ParentFile().Options()))
}

func (s synthesizer) MessageImportPath(msg *protogen.Message) protogen.GoImportPath {
	return protogen.GoImportPath(getImportPath(msg.Desc.ParentFile().Options()))
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

func (s synthesizer) Attribute(field *protogen.Field, explicit bool) (extension.Attribute, error) {
	return NewAttribute(s, field, explicit)
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

func (s synthesizer) ProviderOption(desc protoreflect.MessageDescriptor) *pb.Provider {
	if option, ok := getOptions[*pb.Provider](desc, pb.E_Provider); ok {
		log.Debug().Msgf("ProviderOption: %T", option)
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

func (s synthesizer) Provider(msg *protogen.Message) (extension.Provider, error) {
	return NewProvider(s, msg)
}

func (s synthesizer) Resource(msg *protogen.Message) (extension.Block, error) {
	return NewBlock(s, msg, pb.E_Resource)
}

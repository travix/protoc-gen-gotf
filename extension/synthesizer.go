package extension

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/travix/protoc-gen-gotf/pb"
)

// Synthesizer for synthesizing Provider or resource, datasource Block from proto files, messages.
//
// //go:gocode mockery --name Synthesizer --output ../../mocks.
type Synthesizer interface {
	Block(msg *protogen.Message, blockType protoreflect.ExtensionType) (Block, error)
	Attribute(field *protogen.Field, explicit bool) (Attribute, error)
	Datasource(msg *protogen.Message) (Block, error)
	FieldOption(desc protoreflect.FieldDescriptor) *pb.Attribute
	ProviderOption(desc protoreflect.MessageDescriptor) *pb.Provider
	MessageOption(desc protoreflect.MessageDescriptor, extType protoreflect.ExtensionType) *pb.Block
	Provider(msg *protogen.Message) (Provider, error)
	Resource(msg *protogen.Message) (Block, error)
	Model(msg *protogen.Message, explicit bool) (Model, error)
	Module() string
}

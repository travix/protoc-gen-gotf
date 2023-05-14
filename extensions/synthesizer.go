package extensions

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/travix/protoc-gen-goterraform/pb"
)

// Synthesizer for synthesizing Provider or resource, datasource Block from proto files, messages.
//
// //go:generate mockery --name Synthesizer --output ../../mocks.
type Synthesizer interface {
	Block(msg *protogen.Message, blockType protoreflect.ExtensionType) (Block, error)
	BlockAttribute(field *protogen.Field, explicit bool) (Attribute, error)
	Datasource(msg *protogen.Message) (Block, error)
	FieldOption(desc protoreflect.FieldDescriptor) *pb.Attribute
	FileOption(desc protoreflect.FileDescriptor) *pb.Option
	MessageOption(desc protoreflect.MessageDescriptor, extType protoreflect.ExtensionType) *pb.Block
	Provider(desc protoreflect.FileDescriptor) (Provider, error)
	ProviderAttribute(option *pb.Attribute) (Attribute, error)
	Resource(msg *protogen.Message) (Block, error)
}

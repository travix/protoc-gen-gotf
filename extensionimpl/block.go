package extensionimpl

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/travix/protoc-gen-gotf/extension"
	"github.com/travix/protoc-gen-gotf/pb"
)

var _ extension.Block = &block{}

type block struct {
	_type   protoreflect.ExtensionType
	members map[string]*pb.GoType
	option  *pb.Block
	model   extension.Model
}

func NewBlock(synth extension.Synthesizer, msg *protogen.Message, blockType protoreflect.ExtensionType) (extension.Block, error) {
	b := &block{_type: blockType}
	b.option = synth.MessageOption(msg.Desc, blockType)
	if b.option == nil {
		return nil, nil
	}
	b.setName(msg)
	b.option.Description = *deferToComment(&b.option.Description, msg.Comments)
	var err error
	b.model, err = synth.Model(msg, b.option.ExplicitFields)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (b *block) Description() string {
	return b.option.Description
}

func (b *block) Filename() string {
	name := b.TfName()
	if b._type == pb.E_Resource {
		name += "_resource"
	} else {
		name += "_datasource"
	}
	return name + ".pb.go"
}

func (b *block) Members() map[string]*pb.GoType {
	return b.option.Members
}

func (b *block) Model() extension.Model {
	return b.model
}

func (b *block) GoName() string {
	return *b.option.Name
}

func (b *block) Option() *pb.Block {
	return b.option
}

func (b *block) TfName() string {
	return toSnakeCase(b.GoName())
}

func (b *block) Type() protoreflect.ExtensionType {
	return b._type
}

func (b *block) TypeName() string {
	return string(b._type.TypeDescriptor().FullName())
}

func (b *block) setName(msg *protogen.Message) {
	if b.option.Name == nil {
		b.option.Name = proto.String(msg.GoIdent.GoName)
	}
}
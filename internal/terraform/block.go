package terraform

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/travix/protoc-gen-goterraform/extensions"
	"github.com/travix/protoc-gen-goterraform/pb"
)

var _ extensions.Block = &block{}

type block struct {
	_type      protoreflect.ExtensionType
	attributes []extensions.Attribute
	members    map[string]*pb.GoType
	option     *pb.Block
}

func NewBlock(synth extensions.Synthesizer, msg *protogen.Message, blockType protoreflect.ExtensionType) (extensions.Block, error) {
	b := &block{_type: blockType}
	b.option = synth.MessageOption(msg.Desc, blockType)
	if b.option == nil {
		return nil, nil
	}
	for _, field := range msg.Fields {
		attr, err := synth.BlockAttribute(field, b.option.ExplicitFields)
		if err != nil {
			return nil, err
		}
		if attr == nil {
			continue
		}
		b.attributes = append(b.attributes, attr)
	}
	b.setName(msg)
	return b, nil
}

func (b *block) Attributes() []extensions.Attribute {
	return b.attributes
}

func (b *block) Members() map[string]*pb.GoType {
	return b.option.Members
}

func (b *block) Name() string {
	return *b.option.Name
}

func (b *block) Option() *pb.Block {
	return b.option
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

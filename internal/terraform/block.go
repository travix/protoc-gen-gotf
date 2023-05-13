package terraform

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/travix/protoc-gen-goterraform/pb"
)

var _ Block = &block{}

// Block is helper to generate a terraform block.
//
// //go:generate mockery --name Block --output ../../mocks.
type Block interface {
	Attributes() []Attribute
	Name() string
	Option() *pb.Block
	Members() map[string]protogen.GoIdent
}

type block struct {
	members    map[string]protogen.GoIdent
	attributes []Attribute
	option     *pb.Block
}

func (b *block) Attributes() []Attribute {
	return b.attributes
}

func (b *block) Name() string {
	return *b.option.Name
}

func (b *block) Option() *pb.Block {
	return b.option
}

func (b *block) Members() map[string]protogen.GoIdent {
	// TODO: set from synthesizer
	return b.members
}

func (b *block) setName(msg *protogen.Message) {
	if b.option.Name == nil {
		b.option.Name = proto.String(msg.GoIdent.GoName)
	}
}

func NewBlock(synth Synthesizer, msg *protogen.Message, blockType protoreflect.ExtensionType) (Block, error) {
	b := &block{}
	b.option = synth.getMessageOption(msg.Desc, blockType)
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

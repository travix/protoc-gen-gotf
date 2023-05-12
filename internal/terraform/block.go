package terraform

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-goterraform/pb"
)

var _ Block = &block{}

// Block is helper to generate a terraform block.
//
// //go:generate mockery --name Block --output ../../mocks.
type Block interface {
	Attributes() []*Attribute
	Name() string
	Option() *pb.Block
	StructMembers() map[string]protogen.GoIdent
}

type block struct {
	name          string
	structMembers map[string]protogen.GoIdent
	attributes    []*Attribute
	option        *pb.Block
}

func (b *block) Attributes() []*Attribute {
	return b.attributes
}

func (b *block) Name() string {
	return b.name
}

func (b *block) Option() *pb.Block {
	return b.option
}

func (b *block) StructMembers() map[string]protogen.GoIdent {
	return b.structMembers
}

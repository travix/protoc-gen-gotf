package extensions

import (
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/travix/protoc-gen-goterraform/pb"
)

// Block is helper to generate a terraform block.
//
// //go:generate mockery --name Block --output ../../mocks.
type Block interface {
	Attributes() []Attribute
	Members() map[string]*pb.GoType
	Name() string
	Option() *pb.Block
	Type() protoreflect.ExtensionType
	TypeName() string
}

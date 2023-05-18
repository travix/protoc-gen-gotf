package extension

import (
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/travix/protoc-gen-gotf/pb"
)

// Block is helper to gocode a terraform block.
//
// //go:gocode mockery --name Block --output ../../mocks.
type Block interface {
	Description() string
	Filename() string
	Members() map[string]*pb.GoType
	Model() Model
	GoName() string
	Option() *pb.Block
	TfName() string
	Type() protoreflect.ExtensionType
	TypeName() string
}

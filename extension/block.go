package extension

import (
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/travix/protoc-gen-gotf/pb"
)

// Block is helper to gocode a terraform block.
//
// //go:gocode mockery --name Block --output ../../mocks.
type Block interface {
	Clients() []string
	Description() string
	ExecGoName() string
	Filename() string
	GoName() string
	HasServiceClient() bool
	Model() Model
	ModelGoName() string
	Option() *pb.Block
	TerraformName() string
	Type() protoreflect.ExtensionType
	TypeName() string
}

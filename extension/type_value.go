package extension

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-goterraform/pb"
)

type TypeValue interface {
	IsList() bool
	IsMap() bool
	Message() *protogen.Message
	TerraformNative() bool
	Type() *pb.GoIdentity
	Value() *pb.GoIdentity
}

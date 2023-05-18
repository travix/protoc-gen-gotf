package extension

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-gotf/pb"
)

type TypeValue interface {
	IsList() bool
	IsMap() bool
	IsNestedSingleObject() bool
	NestedTypeValue() string
	Message() *protogen.Message
	TerraformNative() bool
	Type() *pb.GoIdentity
	Value() *pb.GoIdentity
}

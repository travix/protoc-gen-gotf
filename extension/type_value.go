package extension

import (
	"github.com/travix/protoc-gen-gotf/pb"
)

type TypeValue interface {
	IsList() bool
	IsMap() bool
	IsNestedSingleObject() bool
	NestedTypeValue() string
	TerraformNative() bool
	Type() *pb.GoIdentity
	Value() *pb.GoIdentity
}

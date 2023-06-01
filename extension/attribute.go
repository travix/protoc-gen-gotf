package extension

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-gotf/pb"
)

type Attribute interface {
	Computed() bool
	DefaultValue() string
	Deprecation() string
	Description() string
	ElementType() string
	Field() *protogen.Field
	HasNestedType() bool
	IsPointer() bool
	MdDescription() string
	Name() string
	NeedsDefaultValue() bool
	Optional() bool
	Required() bool
	Schema() *pb.GoIdentity
	Sensitive() bool
	TypeValue() TypeValue
	// TODO: support CustomType and Validators
}

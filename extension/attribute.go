package extension

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-gotf/pb"
)

type Attribute interface {
	Computed() bool
	Deprecation() string
	Description() string
	ElementType() string
	Field() *protogen.Field
	MdDescription() string
	Name() string
	Optional() bool
	Required() bool
	Schema() *pb.GoIdentity
	Sensitive() bool
	TypeValue() TypeValue
	HasNestedType() bool
	// TODO: support CustomType and Validators
}

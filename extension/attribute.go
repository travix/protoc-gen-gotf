package extension

import (
	"github.com/travix/protoc-gen-gotf/pb"
)

type Attribute interface {
	Computed() bool
	Deprecation() string
	Description() string
	ElementType() string
	GoName() string
	MdDescription() string
	Name() string
	Optional() bool
	Required() bool
	Schema() *pb.GoIdentity
	Sensitive() bool
	TypeValue() TypeValue
	// TODO: support CustomType and Validators
}

package extensions

import (
	"github.com/travix/protoc-gen-goterraform/pb"
)

type Attribute interface {
	Computed() bool
	Deprecation() string
	Description() string
	MdDescription() string
	Name() string
	Optional() bool
	Required() bool
	Schema() *pb.GoIdentity
	Sensitive() bool
	TypeValue() TypeValue
	// TODO: support CustomType and Validators
}

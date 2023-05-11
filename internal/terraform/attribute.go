package terraform

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-goterraform/pb"
)

var _ Attribute = &attribute{}

type Attribute interface {
	Computed() bool
	Deprecation() string
	Description() string
	GoIdent() protogen.GoIdent
	MdDescription() string
	Name() string
	NestedType() *protogen.GoIdent
	Option() *pb.Attribute
	Optional() bool
	Required() bool
	Sensitive() bool
	// TODO: support CustomType and Validators
}

type attribute struct {
	computed      bool
	deprecation   string
	description   string
	goIdent       protogen.GoIdent
	mdDescription string
	name          string
	nestedType    *protogen.GoIdent
	option        *pb.Attribute
	optional      bool
	required      bool
	sensitive     bool
}

func NewAttributeFromProto(field *protogen.Field, explicit bool) (Attribute, TypeValue, error) {
	a := &attribute{}
	a.option = getFieldOption(field.Desc)
	if a.option == nil {
		if explicit {
			return nil, nil, nil
		}
		a.option = &pb.Attribute{}
		// TODO: set defaults
	}
	return nil, nil, nil
}

func (a *attribute) Computed() bool {
	return a.computed
}

func (a *attribute) Deprecation() string {
	return a.deprecation
}

func (a *attribute) Description() string {
	return a.description
}

func (a *attribute) GoIdent() protogen.GoIdent {
	return a.goIdent
}

func (a *attribute) MdDescription() string {
	return a.mdDescription
}

func (a *attribute) Name() string {
	return a.name
}

func (a *attribute) NestedType() *protogen.GoIdent {
	return a.nestedType
}

func (a *attribute) Option() *pb.Attribute {
	return a.option
}

func (a *attribute) Optional() bool {
	return a.optional
}

func (a *attribute) Required() bool {
	return a.required
}

func (a *attribute) Sensitive() bool {
	return a.sensitive
}

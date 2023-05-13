package terraform

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/travix/protoc-gen-goterraform/pb"
)

var _ Attribute = &attribute{}

type Attribute interface {
	Computed() bool
	Deprecation() string
	Description() string
	MdDescription() string
	Name() string
	Optional() bool
	Required() bool
	Sensitive() bool
	TypeValue() TypeValue
	Schema() *pb.GoIdentity
	// TODO: support CustomType and Validators
}

type attribute struct {
	*pb.Attribute
	typeValue TypeValue
	schema    *pb.GoIdentity
}

func NewBlockAttribute(synth Synthesizer, field *protogen.Field, explicit bool) (Attribute, error) {
	option := synth.getFieldOption(field.Desc)
	if option == nil {
		if explicit {
			return nil, nil
		}
		option = &pb.Attribute{}
	}
	if option.Name == nil {
		option.Name = proto.String(toSnakeCase(field.GoName))
	}
	// ignore attr on block attributes
	option.Attr = nil
	option.Description = getString(option.Description, field.Comments)
	option.MdDescription = getString(option.MdDescription, field.Comments)
	a, err := NewAttribute(option)
	if err != nil {
		return nil, err
	}
	aTyped, _ := a.(*attribute)
	if aTyped.typeValue, err = inferTypeValue(field); err != nil {
		return nil, err
	}
	if aTyped.typeValue == nil {
		return nil, fmt.Errorf("error failed to infer type value for %s", field.Desc.FullName())
	}
	return a, aTyped.setSchema(field)
}

func NewAttribute(option *pb.Attribute) (Attribute, error) {
	if option == nil {
		return nil, nil
	}
	a := &attribute{}
	a.Attribute = option
	if a.Attribute.Skip {
		return nil, nil
	}
	a.setDeprecation(option, protogen.CommentSet{})
	a.setDescription(option, protogen.CommentSet{})
	a.setMdDescription(option, protogen.CommentSet{})
	if option.Name == nil {
		return nil, fmt.Errorf("attribute name is required")
	}
	if strings.TrimSpace(*option.Name) == "" {
		return nil, fmt.Errorf("attribute name can't be empty string")
	}
	var err error
	if option.Attr != nil {
		if a.typeValue, err = explicitTypeValue(option.Attr); err != nil {
			return nil, err
		}
		if err = a.setSchemaFromOption(*option.Attr); err != nil {
			return nil, err
		}
	}
	return a, nil
}

func (a *attribute) Computed() bool {
	return a.MustBe == pb.MustBe_Computed || a.MustBe == pb.MustBe_OptionalAndComputed
}

func (a *attribute) Deprecation() string {
	return *a.Attribute.Description
}

func (a *attribute) Description() string {
	return *a.Attribute.Description
}

func (a *attribute) MdDescription() string {
	return *a.Attribute.MdDescription
}

func (a *attribute) Name() string {
	return *a.Attribute.Name
}

func (a *attribute) Optional() bool {
	return a.MustBe == pb.MustBe_Optional || a.MustBe == pb.MustBe_OptionalAndComputed
}

func (a *attribute) Required() bool {
	return a.MustBe == pb.MustBe_Required
}

func (a *attribute) Schema() *pb.GoIdentity {
	return a.schema
}

func (a *attribute) Sensitive() bool {
	return a.Attribute.Sensitive != nil && *a.Attribute.Sensitive
}

func (a *attribute) TypeValue() TypeValue {
	return a.typeValue
}

func (a *attribute) setDeprecation(option *pb.Attribute, comments protogen.CommentSet) {
	a.Attribute.Deprecation = getString(option.Deprecation, comments)
}

func (a *attribute) setDescription(option *pb.Attribute, comments protogen.CommentSet) {
	a.Attribute.Description = getString(option.Description, comments)
}

func (a *attribute) setMdDescription(option *pb.Attribute, comments protogen.CommentSet) {
	a.Attribute.MdDescription = getString(option.MdDescription, comments)
}

func (a *attribute) setSchemaFromOption(attrType pb.AttrType) error {
	switch attrType {
	case pb.AttrType_bool_attr:
		a.schema = &pb.GoIdentity{Name: "BoolAttribute", ImportPath: "github.com/hashicorp/terraform-plugin-framework/datasource/schema"}
	case pb.AttrType_float64_attr:
		a.schema = &pb.GoIdentity{Name: "Float64Attribute", ImportPath: "github.com/hashicorp/terraform-plugin-framework/datasource/schema"}
	case pb.AttrType_int64_attr:
		a.schema = &pb.GoIdentity{Name: "Int64Attribute", ImportPath: "github.com/hashicorp/terraform-plugin-framework/datasource/schema"}
	case pb.AttrType_string_attr:
		a.schema = &pb.GoIdentity{Name: "StringAttribute", ImportPath: "github.com/hashicorp/terraform-plugin-framework/datasource/schema"}
	default:
		return fmt.Errorf("%w: %s", ErrUnknownAttrType, attrType)
	}
	return nil
}

func (a *attribute) setSchema(field *protogen.Field) error {
	if field.Desc.IsList() {
		a.schema = SchemaList()
		return nil
	}
	if field.Desc.IsMap() {
		a.schema = SchemaMap()
		return nil
	}
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		a.schema = SchemaBool()
	case protoreflect.DoubleKind, protoreflect.FloatKind:
		a.schema = SchemaFloat64()
	case protoreflect.Int32Kind, protoreflect.Int64Kind:
		a.schema = SchemaInt64()
	case protoreflect.MessageKind:
		a.schema = SchemaSingleNested()
	case protoreflect.StringKind:
		a.schema = SchemaString()
	default:
		return fmt.Errorf("error at %s#%s: %w", field.Location.SourceFile, field.Location.Path, ErrUnsupportedKind)
	}
	return nil
}

func SchemaList() *pb.GoIdentity {
	return &pb.GoIdentity{Name: "ListAttribute", ImportPath: "github.com/hashicorp/terraform-plugin-framework/datasource/schema"}
}

func SchemaMap() *pb.GoIdentity {
	return &pb.GoIdentity{Name: "MapAttribute", ImportPath: "github.com/hashicorp/terraform-plugin-framework/datasource/schema"}
}

func SchemaBool() *pb.GoIdentity {
	return &pb.GoIdentity{Name: "BoolAttribute", ImportPath: "github.com/hashicorp/terraform-plugin-framework/datasource/schema"}
}

func SchemaFloat64() *pb.GoIdentity {
	return &pb.GoIdentity{Name: "Float64Attribute", ImportPath: "github.com/hashicorp/terraform-plugin-framework/datasource/schema"}
}

func SchemaInt64() *pb.GoIdentity {
	return &pb.GoIdentity{Name: "Int64Attribute", ImportPath: "github.com/hashicorp/terraform-plugin-framework/datasource/schema"}
}

func SchemaSingleNested() *pb.GoIdentity {
	return &pb.GoIdentity{Name: "SingleNestedAttribute", ImportPath: "github.com/hashicorp/terraform-plugin-framework/datasource/schema"}
}

func SchemaString() *pb.GoIdentity {
	return &pb.GoIdentity{Name: "StringAttribute", ImportPath: "github.com/hashicorp/terraform-plugin-framework/datasource/schema"}
}

func getString(direct *string, comments protogen.CommentSet) *string {
	if direct != nil {
		return direct
	}
	var str string
	for index, c := range comments.LeadingDetached {
		if index > 0 {
			str += "\n"
		}
		str += c.String()
	}
	str += string(comments.Leading)
	str += string(comments.Trailing)
	str = strings.TrimSpace(str)
	return &str
}

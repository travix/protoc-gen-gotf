package extensionimpl

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"

	"github.com/travix/protoc-gen-gotf/extension"
	"github.com/travix/protoc-gen-gotf/pb"
)

var _ extension.Attribute = &attribute{}

type attribute struct {
	defaultValue string
	elementType  string
	field        *protogen.Field
	option       *pb.Attribute
	schema       *pb.GoIdentity
	typeValue    extension.TypeValue
}

func NewAttribute(synth extension.Synthesizer, field *protogen.Field, explicit bool) (extension.Attribute, error) {
	a := &attribute{}
	a.option = synth.FieldOption(field.Desc)
	if a.option == nil {
		if explicit {
			return nil, nil
		}
		a.option = &pb.Attribute{}
	}
	if a.option.Skip {
		return nil, nil
	}
	a.field = field
	if a.option.Name == nil {
		a.option.Name = proto.String(field.GoName)
	}
	if strings.TrimSpace(*a.option.Name) == "" {
		return nil, fmt.Errorf("attribute name can't be empty string")
	}
	var err error
	a.defaultValue, err = defaultValue(a.field)
	if err != nil {
		return nil, err
	}
	a.option.Description = deferToComment(a.option.Description, field.Comments)
	a.option.MdDescription = deferToComment(a.option.MdDescription, field.Comments)
	a.option.Deprecation = deferToComment(a.option.Deprecation, protogen.CommentSet{})
	if a.typeValue, err = inferTypeValue(field); err != nil {
		return nil, err
	}
	if a.schema, err = Schema(field); err != nil {
		return nil, err
	}
	if field.Message != nil {
		// note message will implement terraform typable through extension.Model
		a.elementType = fmt.Sprintf("&%s{}", field.Message.GoIdent.GoName)
	}
	return a, nil
}

func (a *attribute) DefaultValue() string {
	return a.defaultValue
}

func (a *attribute) NeedsDefaultValue() bool {
	return !a.IsPointer() && !a.TypeValue().IsList() && !a.TypeValue().IsMap() && !a.TypeValue().IsNestedSingleObject()
}

func (a *attribute) Computed() bool {
	return a.option.MustBe == pb.MustBe_Computed || a.option.MustBe == pb.MustBe_OptionalAndComputed
}

func (a *attribute) Deprecation() string {
	return *a.option.Deprecation
}

func (a *attribute) Description() string {
	return *a.option.Description
}

func (a *attribute) ElementType() string {
	return a.elementType
}

func (a *attribute) Field() *protogen.Field {
	return a.field
}

func (a *attribute) IsPointer() bool {
	return a.field.Desc.HasOptionalKeyword()
}

func (a *attribute) MdDescription() string {
	return *a.option.MdDescription
}

func (a *attribute) Name() string {
	return toSnakeCase(*a.option.Name)
}

func (a *attribute) Optional() bool {
	return a.option.MustBe == pb.MustBe_Optional || a.option.MustBe == pb.MustBe_OptionalAndComputed
}

func (a *attribute) Required() bool {
	return a.option.MustBe == pb.MustBe_Required
}

func (a *attribute) Schema() *pb.GoIdentity {
	return a.schema
}

func (a *attribute) Sensitive() bool {
	return a.option.Sensitive != nil && *a.option.Sensitive
}

func (a *attribute) TypeValue() extension.TypeValue {
	return a.typeValue
}

func (a *attribute) HasNestedType() bool {
	return a.typeValue.IsList() || a.typeValue.IsMap() || a.typeValue.IsNestedSingleObject()
}

func (a *attribute) NestedType() string {
	if strings.HasPrefix(a.typeValue.NestedTypeValue(), "types.") && strings.HasSuffix(a.typeValue.NestedTypeValue(), "Type") {
		return a.typeValue.NestedTypeValue()
	}
	return fmt.Sprintf("types.ObjectType{ AttrTypes: (%s).AttributeTypes() }", a.TypeValue().NestedTypeValue())
}

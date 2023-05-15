package extensionimpl

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-goterraform/extension"
	"github.com/travix/protoc-gen-goterraform/pb"
)

var _ extension.Attribute = &attribute{}

type attribute struct {
	*pb.Attribute
	elementType string
	schema      *pb.GoIdentity
	typeValue   extension.TypeValue
}

func NewAttribute(option *pb.Attribute) (extension.Attribute, error) {
	if option == nil {
		return nil, nil
	}
	a := &attribute{}
	a.Attribute = option
	if a.Attribute.Skip {
		return nil, nil
	}
	a.Attribute.Deprecation = deferToComment(option.Deprecation, protogen.CommentSet{})
	a.Attribute.Description = deferToComment(option.Description, protogen.CommentSet{})
	a.Attribute.MdDescription = deferToComment(option.MdDescription, protogen.CommentSet{})
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
		if a.schema, err = schemaForAttrType(*option.Attr); err != nil {
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

func (a *attribute) ElementType() string {
	return a.elementType
}

func (a *attribute) GoName() string {
	return *a.Attribute.Name
}

func (a *attribute) MdDescription() string {
	return *a.Attribute.MdDescription
}

func (a *attribute) Name() string {
	return toSnakeCase(*a.Attribute.Name)
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

func (a *attribute) TypeValue() extension.TypeValue {
	return a.typeValue
}

func deferToComment(direct *string, comments protogen.CommentSet) *string {
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

package extensionimpl

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/travix/protoc-gen-gotf/extension"
	"github.com/travix/protoc-gen-gotf/pb"
)

var (
	_                  extension.TypeValue = &typeValue{}
	ErrUnsupportedKind                     = errors.New("unsupported protoreflect.Kind")
)

type typeValue struct {
	_type           *pb.GoIdentity
	isList          bool
	isMap           bool
	terraformNative bool
	value           *pb.GoIdentity
	nestedType      string
}

func (t typeValue) NestedTypeValue() string {
	return t.nestedType
}

func (t typeValue) IsNestedSingleObject() bool {
	return t.nestedType != "" && !t.isList && !t.isMap
}

func NewMapTypeValue(message *protogen.Message) extension.TypeValue {
	return newTypeValue(message, false, true)
}

func NewListTypeValue(message *protogen.Message) extension.TypeValue {
	return newTypeValue(message, true, false)
}

func NewNestedSingleObjectTypeValue(message *protogen.Message) extension.TypeValue {
	return newTypeValue(message, false, false)
}

func TypeValueBool() extension.TypeValue {
	return &typeValue{
		_type: &pb.GoIdentity{
			Name:       "BoolType",
			ImportPath: "github.com/hashicorp/terraform-plugin-framework/types",
		},
		value: &pb.GoIdentity{
			Name:       "Bool",
			ImportPath: "github.com/hashicorp/terraform-plugin-framework/types",
		},
		terraformNative: true,
	}
}

func TypeValueString() extension.TypeValue {
	return &typeValue{
		_type: &pb.GoIdentity{
			Name:       "StringType",
			ImportPath: "github.com/hashicorp/terraform-plugin-framework/types",
		},
		value: &pb.GoIdentity{
			Name:       "String",
			ImportPath: "github.com/hashicorp/terraform-plugin-framework/types",
		},
		terraformNative: true,
	}
}

func TypeValueInt64() extension.TypeValue {
	return &typeValue{
		_type: &pb.GoIdentity{
			Name:       "Int64Type",
			ImportPath: "github.com/hashicorp/terraform-plugin-framework/types",
		},
		value: &pb.GoIdentity{
			Name:       "Int64",
			ImportPath: "github.com/hashicorp/terraform-plugin-framework/types",
		},
		terraformNative: true,
	}
}

func TypeValueFloat64() extension.TypeValue {
	return &typeValue{
		_type: &pb.GoIdentity{
			Name:       "Float64Type",
			ImportPath: "github.com/hashicorp/terraform-plugin-framework/types",
		},
		value: &pb.GoIdentity{
			Name:       "Float64",
			ImportPath: "github.com/hashicorp/terraform-plugin-framework/types",
		},
		terraformNative: true,
	}
}

func (t typeValue) IsList() bool {
	return t.isList
}

func (t typeValue) IsMap() bool {
	return t.isMap
}

// func (t typeValue) Message() *protogen.Message {
//	return t.message
//}

func (t typeValue) TerraformNative() bool {
	return t.terraformNative
}

func (t typeValue) Value() *pb.GoIdentity {
	return t.value
}

func (t typeValue) Type() *pb.GoIdentity {
	return t._type
}

func inferTypeValue(field *protogen.Field) (extension.TypeValue, error) {
	kind := field.Desc.Kind()
	var tv extension.TypeValue
	var err error
	switch kind {
	case protoreflect.BoolKind:
		tv = TypeValueBool()
	case protoreflect.DoubleKind, protoreflect.FloatKind:
		tv = TypeValueFloat64()
	case protoreflect.Int32Kind, protoreflect.Int64Kind:
		tv = TypeValueInt64()
	case protoreflect.MessageKind:
		if field.Desc.IsList() {
			return NewListTypeValue(field.Message), nil
		}
		if field.Desc.IsMap() {
			if field.Message.Fields[0].Desc.Kind() != protoreflect.StringKind {
				err = fmt.Errorf("unsupported map key type: %s", field.Message.Fields[0].Desc.Kind())
				break
			}
			return NewMapTypeValue(field.Message.Fields[2].Message), nil
		}
		return NewNestedSingleObjectTypeValue(field.Message), nil
	case protoreflect.StringKind:
		tv = TypeValueString()
	default:
		err = ErrUnsupportedKind
	}
	if err != nil {
		return nil, fmt.Errorf("error at %s#%s: %w", field.Location.SourceFile, field.Location.Path, err)
	}
	if tv == nil {
		return nil, fmt.Errorf("error failed to infer type value for %s", field.Desc.FullName())
	}
	if field.Desc.IsList() {
		tvTyped, _ := tv.(*typeValue)
		tvTyped.isList = true
		tvTyped.nestedType = fmt.Sprintf("types.%s", tv.Type().Name)
	}
	return tv, nil
}

func newTypeValue(message *protogen.Message, isList, isMap bool) *typeValue {
	return &typeValue{
		_type: &pb.GoIdentity{
			Name:       fmt.Sprintf("%sTfType", message.GoIdent.GoName),
			ImportPath: string(message.GoIdent.GoImportPath),
		},
		value: &pb.GoIdentity{
			// TODO: push implementation of github.com/hashicorp/terraform-plugin-framework/types/basetypes.ObjectValuable to pb
			Name:       message.GoIdent.GoName,
			ImportPath: string(message.GoIdent.GoImportPath),
		},
		isList:     isList,
		isMap:      isMap,
		nestedType: fmt.Sprintf("&%s{}", message.GoIdent.GoName),
	}
}

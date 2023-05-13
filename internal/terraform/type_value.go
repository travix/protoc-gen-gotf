package terraform

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/travix/protoc-gen-goterraform/pb"
)

var (
	_                  TypeValue = &typeValue{}
	ErrUnsupportedKind           = errors.New("unsupported protoreflect.Kind")
	ErrUnknownAttrType           = errors.New("error unknown value")
)

type TypeValue interface {
	TerraformNative() bool
	Type() *pb.GoIdentity
	Value() *pb.GoIdentity
	Message() *protogen.Message
	IsList() bool
	IsMap() bool
}

type typeValue struct {
	_type           *pb.GoIdentity
	isList          bool
	isMap           bool
	message         *protogen.Message
	terraformNative bool
	value           *pb.GoIdentity
}

func TypeValueBool() TypeValue {
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

func TypeValueString() TypeValue {
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

func TypeValueInt64() TypeValue {
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

func TypeValueFloat64() TypeValue {
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

func (t typeValue) Message() *protogen.Message {
	return t.message
}

func (t typeValue) TerraformNative() bool {
	return t.terraformNative
}

func (t typeValue) Value() *pb.GoIdentity {
	return t.value
}

func (t typeValue) Type() *pb.GoIdentity {
	return t._type
}

func inferTypeValue(field *protogen.Field) (TypeValue, error) {
	kind := field.Desc.Kind()
	var tv TypeValue
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
			tv = NewListTypeValue(field.Message)
		} else if field.Desc.IsMap() {
			if field.Message.Fields[0].Desc.Kind() != protoreflect.StringKind {
				err = fmt.Errorf("unsupported map key type: %s", field.Message.Fields[0].Desc.Kind())
				break
			}
			tv = NewMapTypeValue(field.Message.Fields[2].Message)
		}

		return tv, nil
	case protoreflect.StringKind:
		tv, _ = TypeValueString().(*typeValue)
	default:
		err = ErrUnsupportedKind
	}
	if err != nil {
		return nil, fmt.Errorf("error at %s#%s: %w", field.Location.SourceFile, field.Location.Path, err)
	}
	return tv, nil
}

func NewMapTypeValue(message *protogen.Message) TypeValue {
	return newTypeValue(message, false, true)
}

func NewListTypeValue(message *protogen.Message) TypeValue {
	return newTypeValue(message, true, false)
}

func newTypeValue(message *protogen.Message, isList, isMap bool) *typeValue {
	return &typeValue{
		_type: &pb.GoIdentity{
			// TODO: push implementation of github.com/hashicorp/terraform-plugin-framework/types/basetypes.ObjectTypable to pb
			Name:       fmt.Sprintf("%sTfType", message.GoIdent.GoName),
			ImportPath: string(message.GoIdent.GoImportPath),
		},
		value: &pb.GoIdentity{
			// TODO: push implementation of github.com/hashicorp/terraform-plugin-framework/types/basetypes.ObjectValuable to pb
			Name:       message.GoIdent.GoName,
			ImportPath: string(message.GoIdent.GoImportPath),
		},
		isList:  isList,
		isMap:   isMap,
		message: message,
	}
}

func explicitTypeValue(attrType *pb.AttrType) (TypeValue, error) {
	if attrType == nil {
		return nil, nil
	}
	switch *attrType {
	case pb.AttrType_bool_attr:
		return TypeValueBool(), nil
	case pb.AttrType_string_attr:
		return TypeValueString(), nil
	case pb.AttrType_int64_attr:
		return TypeValueInt64(), nil
	case pb.AttrType_float64_attr:
		return TypeValueFloat64(), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownAttrType, attrType)
	}
}

package extensionimpl

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/travix/protoc-gen-gotf/pb"
)

func SchemaBool() *pb.GoIdentity {
	return &pb.GoIdentity{Name: "BoolAttribute"}
}

func SchemaFloat64() *pb.GoIdentity {
	return &pb.GoIdentity{Name: "Float64Attribute"}
}

func SchemaInt64() *pb.GoIdentity {
	return &pb.GoIdentity{Name: "Int64Attribute"}
}

func SchemaList() *pb.GoIdentity {
	return &pb.GoIdentity{Name: "ListAttribute"}
}

func SchemaMap() *pb.GoIdentity {
	return &pb.GoIdentity{Name: "MapAttribute"}
}

func SchemaSingleNested() *pb.GoIdentity {
	return &pb.GoIdentity{Name: "SingleNestedAttribute"}
}

func SchemaString() *pb.GoIdentity {
	return &pb.GoIdentity{Name: "StringAttribute"}
}

func Schema(field *protogen.Field) (*pb.GoIdentity, error) {
	var schema *pb.GoIdentity
	if field.Desc.IsList() {
		return SchemaList(), nil
	}
	if field.Desc.IsMap() {
		return SchemaMap(), nil
	}
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		schema = SchemaBool()
	case protoreflect.DoubleKind, protoreflect.FloatKind:
		schema = SchemaFloat64()
	case protoreflect.Int32Kind, protoreflect.Int64Kind:
		schema = SchemaInt64()
	case protoreflect.MessageKind:
		schema = SchemaSingleNested()
	case protoreflect.StringKind:
		schema = SchemaString()
	default:
		return nil, fmt.Errorf("error at %s#%s: %w", field.Location.SourceFile, field.Location.Path, ErrUnsupportedKind)
	}
	return schema, nil
}

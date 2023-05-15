package extensionimpl

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/travix/protoc-gen-goterraform/pb"
)

func SchemaBool() *pb.GoIdentity {
	return &pb.GoIdentity{Name: "BoolAttribute", ImportPath: "github.com/hashicorp/terraform-plugin-framework/datasource/schema"}
}

func SchemaFloat64() *pb.GoIdentity {
	return &pb.GoIdentity{Name: "Float64Attribute", ImportPath: "github.com/hashicorp/terraform-plugin-framework/datasource/schema"}
}

func SchemaInt64() *pb.GoIdentity {
	return &pb.GoIdentity{Name: "Int64Attribute", ImportPath: "github.com/hashicorp/terraform-plugin-framework/datasource/schema"}
}

func SchemaList() *pb.GoIdentity {
	return &pb.GoIdentity{Name: "ListAttribute", ImportPath: "github.com/hashicorp/terraform-plugin-framework/datasource/schema"}
}

func SchemaMap() *pb.GoIdentity {
	return &pb.GoIdentity{Name: "MapAttribute", ImportPath: "github.com/hashicorp/terraform-plugin-framework/datasource/schema"}
}

func SchemaSingleNested() *pb.GoIdentity {
	return &pb.GoIdentity{Name: "SingleNestedAttribute", ImportPath: "github.com/hashicorp/terraform-plugin-framework/datasource/schema"}
}

func SchemaString() *pb.GoIdentity {
	return &pb.GoIdentity{Name: "StringAttribute", ImportPath: "github.com/hashicorp/terraform-plugin-framework/datasource/schema"}
}

func setSchema(a *attribute, field *protogen.Field) error {
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
		// TODO: assuming message implements terraform typable interface
		a.elementType = fmt.Sprintf("&%s{}", field.Message.GoIdent.GoName)
	case protoreflect.StringKind:
		a.schema = SchemaString()
	default:
		return fmt.Errorf("error at %s#%s: %w", field.Location.SourceFile, field.Location.Path, ErrUnsupportedKind)
	}
	return nil
}

func schemaForAttrType(attrType pb.AttrType) (*pb.GoIdentity, error) {
	var schema *pb.GoIdentity
	switch attrType {
	case pb.AttrType_bool_attr:
		schema = &pb.GoIdentity{Name: "BoolAttribute", ImportPath: "github.com/hashicorp/terraform-plugin-framework/datasource/schema"}
	case pb.AttrType_float64_attr:
		schema = &pb.GoIdentity{Name: "Float64Attribute", ImportPath: "github.com/hashicorp/terraform-plugin-framework/datasource/schema"}
	case pb.AttrType_int64_attr:
		schema = &pb.GoIdentity{Name: "Int64Attribute", ImportPath: "github.com/hashicorp/terraform-plugin-framework/datasource/schema"}
	case pb.AttrType_string_attr:
		schema = &pb.GoIdentity{Name: "StringAttribute", ImportPath: "github.com/hashicorp/terraform-plugin-framework/datasource/schema"}
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownAttrType, attrType)
	}
	return schema, nil
}

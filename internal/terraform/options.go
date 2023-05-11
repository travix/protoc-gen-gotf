package terraform

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/travix/protoc-gen-goterraform/pb"
)

func getFieldOption(desc protoreflect.FieldDescriptor) *pb.Attribute {
	if mo, ok := desc.Options().(*descriptorpb.FieldOptions); ok && mo != nil {
		option, _ := proto.GetExtension(mo, pb.E_Attribute).(*pb.Attribute)
		return option
	}
	return nil
}

func getFileOption(desc protoreflect.FileDescriptor) *pb.Option {
	if fo, ok := desc.Options().(*descriptorpb.FileOptions); ok && fo != nil {
		option, _ := proto.GetExtension(fo, pb.E_Provider).(*pb.Option)
		return option
	}
	return nil
}

func getMessageOption(desc protoreflect.MessageDescriptor, datasource protoreflect.ExtensionType) *pb.Block {
	if mo, ok := desc.Options().(*descriptorpb.MessageOptions); ok && mo != nil {
		option, _ := proto.GetExtension(mo, datasource).(*pb.Block)
		return option
	}
	return nil
}

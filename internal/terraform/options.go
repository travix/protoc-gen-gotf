package terraform

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/travix/protoc-gen-goterraform/pb"
)

func getBlockOption(desc protoreflect.MessageDescriptor, datasource *protoimpl.ExtensionInfo) *pb.Block {
	var block *pb.Block
	if mo, ok := desc.Options().(*descriptorpb.MessageOptions); ok && mo != nil {
		block, _ = proto.GetExtension(mo, datasource).(*pb.Block)
		return block
	}
	return nil
}

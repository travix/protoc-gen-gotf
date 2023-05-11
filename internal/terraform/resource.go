package terraform

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-goterraform/pb"
)

type Resource struct {
	*pb.Block
}

func ResourcesFromProto(file *protogen.File) ([]*Resource, []*protogen.Message, error) {
	resources := make([]*Resource, 0)
	dep := make([]*protogen.Message, 0)
	for _, msg := range file.Messages {
		if msg.Desc.IsMapEntry() {
			continue
		}
		resource := getBlockOption(msg.Desc, pb.E_Resource)
		if resource == nil {
			continue
		}
		resources = append(resources, &Resource{resource})
		// TODO: parse message attributes
	}
	return resources, dep, nil
}

package terraform

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-goterraform/pb"
)

type DataSource struct {
	*pb.Block
}

func DataSourcesFromProto(file *protogen.File) ([]*DataSource, []*protogen.Message, error) {
	datasources := make([]*DataSource, 0)
	dep := make([]*protogen.Message, 0)
	for _, msg := range file.Messages {
		if msg.Desc.IsMapEntry() {
			continue
		}
		dataSource := getBlockOption(msg.Desc, pb.E_Datasource)
		if dataSource == nil {
			continue
		}
		datasources = append(datasources, &DataSource{dataSource})
	}
	return datasources, dep, nil
}

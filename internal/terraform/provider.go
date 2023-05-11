package terraform

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-goterraform/pb"
)

type Provider struct {
	*pb.Option
	Resources    []*Resource
	DataSources  []*DataSource
	Dependencies []*protogen.Message
}

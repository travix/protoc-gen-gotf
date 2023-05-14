package extensions

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-goterraform/pb"
)

// Provider is helper to generate a terraform provider.
//
// //go:generate mockery --name Provider --output ../../mocks.
type Provider interface {
	Attributes() []Attribute
	GoImportPath() protogen.GoImportPath
	Members() map[string]*pb.GoType
	Name() string
	Option() *pb.Option
	Package() protogen.GoPackageName
}

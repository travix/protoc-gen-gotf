package extension

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-gotf/pb"
)

// Provider is helper to gocode a terraform provider.
//
// //go:gocode mockery --name Provider --output ../../mocks.
type Provider interface {
	Description() string
	ExecGoName() string
	Filename() string
	GoName() string
	ImportPath() protogen.GoImportPath
	Model() Model
	ModelGoName() string
	Option() *pb.Provider
	PackageName() protogen.GoPackageName
	PbImportPath() protogen.GoImportPath
	PbPackageName() protogen.GoPackageName
	TerraformName() string
}

package extension

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-gotf/pb"
)

// Provider is helper to gocode a terraform provider.
//
// //go:gocode mockery --name Provider --output ../../mocks.
type Provider interface {
	AddDatasource(string, string)
	AddResource(string, string)
	Datasources() map[string]string
	Description() string
	ExecFilename() string
	ExecGoName() string
	Filename() string
	GoName() string
	HasServiceClient() bool
	Model() Model
	ModelGoName() string
	Option() *pb.Provider
	PackageData() PackageData
	Resources() map[string]string
	SetHasServiceClient(bool)
	TerraformName() string
}

type PackageData struct {
	ExecImportPath      protogen.GoImportPath
	ExecPackageName     protogen.GoPackageName
	PbImportPath        protogen.GoImportPath
	PbPackageName       protogen.GoPackageName
	ProviderImportPath  protogen.GoImportPath
	ProviderPackageName protogen.GoPackageName
}

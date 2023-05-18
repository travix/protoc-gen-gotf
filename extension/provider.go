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
	Filename() string
	ImportPath() protogen.GoImportPath // old
	Members() map[string]*pb.GoType    // old
	Model() Model
	GoName() string                        // old
	Option() *pb.Provider                  // old
	PackageName() protogen.GoPackageName   // old
	PbImportPath() protogen.GoImportPath   // old
	PbPackageName() protogen.GoPackageName // old
	TfName() string
}

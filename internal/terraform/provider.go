package terraform

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-goterraform/extensions"
	"github.com/travix/protoc-gen-goterraform/pb"
)

var _ extensions.Provider = &provider{}

type provider struct {
	attributes   []extensions.Attribute
	goImportPath protogen.GoImportPath
	goPackage    protogen.GoPackageName
	name         string
	option       *pb.Option
}

func (p *provider) Attributes() []extensions.Attribute {
	return p.attributes
}

func (p *provider) GoImportPath() protogen.GoImportPath {
	return p.goImportPath
}

func (p *provider) Members() map[string]*pb.GoType {
	return p.option.Members
}

func (p *provider) Name() string {
	return p.name
}

func (p *provider) Option() *pb.Option {
	return p.option
}

func (p *provider) Package() protogen.GoPackageName {
	return p.goPackage
}

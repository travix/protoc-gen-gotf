package terraform

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-goterraform/pb"
)

var _ Provider = &provider{}

// Provider is helper to generate a terraform provider.
//
//go:generate mockery --name Provider --output ../../mocks
type Provider interface {
	Attributes() []*Attribute
	GoImportPath() protogen.GoImportPath
	Name() string
	Option() *pb.Option
	Package() protogen.GoPackageName
}

type provider struct {
	attributes   []*Attribute
	goImportPath protogen.GoImportPath
	goPackage    protogen.GoPackageName
	name         string
	option       *pb.Option
}

func (p *provider) Name() string {
	return p.name
}

func (p *provider) Attributes() []*Attribute {
	return p.attributes
}

func (p *provider) Package() protogen.GoPackageName {
	return p.goPackage
}

func (p *provider) GoImportPath() protogen.GoImportPath {
	return p.goImportPath
}

func (p *provider) Option() *pb.Option {
	return p.option
}

package extensionimpl

import (
	"fmt"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-gotf/extension"
	"github.com/travix/protoc-gen-gotf/pb"
)

var _ extension.Provider = &provider{}

type provider struct {
	importPath    protogen.GoImportPath
	model         extension.Model
	module        string
	option        *pb.Provider
	packageName   protogen.GoPackageName
	pbImportPath  protogen.GoImportPath
	pbPackageName protogen.GoPackageName
}

func NewProvider(synth extension.Synthesizer, msg *protogen.Message) (extension.Provider, error) {
	option := synth.ProviderOption(msg.Desc)
	if option == nil {
		return nil, nil
	}
	p := &provider{option: option}
	if p.option.Name == "" {
		p.option.Name = msg.GoIdent.GoName
	}
	p.option.Description = *deferToComment(&p.option.Description, msg.Comments)
	p.pbPackageName = synth.MessagePackageName(msg)
	p.pbImportPath = synth.MessageImportPath(msg)
	if !strings.HasPrefix(string(p.pbImportPath), p.module) {
		p.pbImportPath = protogen.GoImportPath(filepath.Join(p.module, string(p.pbImportPath)))
	}
	p.module = synth.Module()
	if p.option.ProviderPackage == "" {
		p.option.ProviderPackage = filepath.Join(p.module, "providerpb")
	}
	p.packageName = protogen.GoPackageName(filepath.Base(p.option.ProviderPackage))
	p.importPath = protogen.GoImportPath(p.option.ProviderPackage)
	if !strings.HasPrefix(string(p.importPath), p.module) {
		p.importPath = protogen.GoImportPath(filepath.Join(p.module, string(p.importPath)))
	}
	var err error
	p.model, err = synth.Model(msg, false)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *provider) GoName() string {
	return toCamelCase(p.option.Name)
}

func (p *provider) Description() string {
	return p.option.Description
}

func (p *provider) Filename() string {
	return "provider.go"
}

func (p *provider) ImportPath() protogen.GoImportPath {
	return p.importPath
}

func (p *provider) Model() extension.Model {
	return p.model
}

func (p *provider) ModelGoName() string {
	return p.model.GoName()
}

func (p *provider) Name() string {
	return p.option.Name
}

func (p *provider) Option() *pb.Provider {
	return p.option
}

func (p *provider) PackageName() protogen.GoPackageName {
	return p.packageName
}

func (p *provider) PbImportPath() protogen.GoImportPath {
	return p.pbImportPath
}

func (p *provider) PbPackageName() protogen.GoPackageName {
	return p.pbPackageName
}

func (p *provider) ExecGoName() string {
	return fmt.Sprintf("%sExec", p.GoName())
}

func (p *provider) TerraformName() string {
	return toSnakeCase(p.option.Name)
}

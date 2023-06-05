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
	model       extension.Model
	module      string
	option      *pb.Provider
	packageData extension.PackageData
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
	p.packageData.PbPackageName = synth.MessagePackageName(msg)
	p.packageData.PbImportPath = synth.MessageImportPath(msg)
	if !strings.HasPrefix(string(p.packageData.PbImportPath), p.module) {
		p.packageData.PbImportPath = protogen.GoImportPath(filepath.Join(p.module, string(p.packageData.PbImportPath)))
	}
	p.module = synth.Module()
	if p.option.ProviderPackage == "" {
		p.option.ProviderPackage = filepath.Join(p.module, "providerpb")
	}
	p.packageData.ProviderPackageName = protogen.GoPackageName(filepath.Base(p.option.ProviderPackage))
	p.packageData.ProviderImportPath = protogen.GoImportPath(p.option.ProviderPackage)
	if !strings.HasPrefix(string(p.packageData.ProviderImportPath), p.module) {
		p.packageData.ProviderImportPath = protogen.GoImportPath(filepath.Join(p.module, string(p.packageData.ProviderImportPath)))
	}
	if p.option.ExecPackage != nil {
		p.packageData.ExecImportPath = protogen.GoImportPath(*p.option.ExecPackage)
		p.packageData.ExecPackageName = protogen.GoPackageName(filepath.Base(*p.option.ExecPackage))
	}
	var err error
	p.model, err = synth.Model(msg, false)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *provider) Description() string {
	return p.option.Description
}

func (p *provider) Filename() string {
	return "provider.go"
}

func (p *provider) GoName() string {
	return toCamelCase(p.option.Name)
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

func (p *provider) PackageData() extension.PackageData {
	return p.packageData
}

func (p *provider) ExecGoName() string {
	return fmt.Sprintf("%sExec", p.GoName())
}

func (p *provider) TerraformName() string {
	return toSnakeCase(p.option.Name)
}

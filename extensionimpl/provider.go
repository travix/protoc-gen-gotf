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
	datasources      map[string]string
	hasServiceClient bool
	model            extension.Model
	module           string
	option           *pb.Provider
	packageData      extension.PackageData
	resources        map[string]string
}

func NewProvider(synth extension.Synthesizer, msg *protogen.Message) (extension.Provider, error) {
	option := synth.ProviderOption(msg.Desc)
	if option == nil {
		return nil, nil
	}
	p := &provider{option: option, datasources: map[string]string{}, resources: map[string]string{}}
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
	p.packageData.ProviderPackageName = protogen.GoPackageName(getPkgName(p.option.ProviderPackage))
	p.packageData.ProviderImportPath = protogen.GoImportPath(getImportPath(p.option.ProviderPackage))
	if !strings.HasPrefix(string(p.packageData.ProviderImportPath), p.module) {
		p.packageData.ProviderImportPath = protogen.GoImportPath(filepath.Join(p.module, string(p.packageData.ProviderImportPath)))
	}
	if p.option.ExecPackage != nil {
		p.packageData.ExecPackageName = protogen.GoPackageName(getPkgName(*p.option.ExecPackage))
		p.packageData.ExecImportPath = protogen.GoImportPath(getImportPath(*p.option.ExecPackage))
	}
	var err error
	p.model, err = synth.Model(msg, false)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *provider) AddDatasource(name string, exec string) {
	p.datasources[name] = exec
}

func (p *provider) AddResource(name string, exec string) {
	p.resources[name] = exec
}

func (p *provider) Datasources() map[string]string {
	return p.datasources
}

func (p *provider) Description() string {
	return p.option.Description
}

func (p *provider) ExecFilename() string {
	return "provider_exec.go"
}

func (p *provider) ExecGoName() string {
	return fmt.Sprintf("%sExec", p.GoName())
}

func (p *provider) Filename() string {
	return "provider.go"
}

func (p *provider) HasServiceClient() bool {
	return p.hasServiceClient
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

func (p *provider) Resources() map[string]string {
	return p.resources
}

func (p *provider) SetHasServiceClient(has bool) {
	p.hasServiceClient = has
}

func (p *provider) TerraformName() string {
	return toSnakeCase(p.option.Name)
}

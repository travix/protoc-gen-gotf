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
	if p.option.PbPackage == "" {
		return nil, fmt.Errorf("error gotf.provider.package option is required in %s#%s", msg.Location.SourceFile, msg.Location.Path)
	}
	p.setPb(synth, option)
	p.setPkgAndPath()
	var err error
	p.model, err = synth.Model(msg, false)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *provider) setPkgAndPath() {
	seg := strings.Split(p.option.ProviderPackage, ";")
	p.packageName = protogen.GoPackageName(filepath.Base(seg[len(seg)-1]))
	p.importPath = protogen.GoImportPath(seg[0])
	if !strings.HasPrefix(string(p.importPath), p.module) {
		p.importPath = protogen.GoImportPath(filepath.Join(p.module, string(p.importPath)))
	}
}

func (p *provider) setPb(synth extension.Synthesizer, option *pb.Provider) {
	seg := strings.Split(p.option.PbPackage, ";")
	pkg := filepath.Base(seg[len(seg)-1])
	p.pbPackageName = protogen.GoPackageName(pkg)
	seg = strings.Split(p.option.PbPackage, ";")
	p.pbImportPath = protogen.GoImportPath(seg[0])
	p.module = ""
	if p.option.ProviderPackage == "" {
		p.option.ProviderPackage = "providerpb"
		p.module = synth.Module()
	}
	if option.Module != nil {
		p.module = *option.Module
	}
	if !strings.HasPrefix(string(p.pbImportPath), p.module) {
		p.pbImportPath = protogen.GoImportPath(filepath.Join(p.module, string(p.pbImportPath)))
	}
}

func (p *provider) Members() map[string]*pb.GoType {
	return p.option.Members
}

func (p *provider) Name() string {
	return p.option.Name
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

func (p *provider) Model() extension.Model {
	return p.model
}

func (p *provider) TfName() string {
	return toSnakeCase(p.GoName())
}

func (p *provider) Option() *pb.Provider {
	return p.option
}

func (p *provider) PackageName() protogen.GoPackageName {
	return p.packageName
}

func (p *provider) ImportPath() protogen.GoImportPath {
	return p.importPath
}

func (p *provider) PbImportPath() protogen.GoImportPath {
	return p.pbImportPath
}

func (p *provider) PbPackageName() protogen.GoPackageName {
	return p.pbPackageName
}

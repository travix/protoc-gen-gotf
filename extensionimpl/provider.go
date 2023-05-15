package extensionimpl

import (
	"fmt"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/travix/protoc-gen-goterraform/extension"
	"github.com/travix/protoc-gen-goterraform/pb"
)

var _ extension.Provider = &provider{}

type provider struct {
	attributes    []extension.Attribute
	option        *pb.Option
	importPath    protogen.GoImportPath
	packageName   protogen.GoPackageName
	pbPackageName protogen.GoPackageName
	pbImportPath  protogen.GoImportPath
	module        string
}

func NewProvider(synth extension.Synthesizer, desc protoreflect.FileDescriptor) (extension.Provider, error) {
	option := synth.FileOption(desc)
	if option == nil {
		return nil, nil
	}
	p := &provider{option: option}
	if p.option.Name == "" {
		return nil, fmt.Errorf("error goterraform.provider.name option not found in %s", desc.Path())
	}
	if p.option.Package == "" {
		return nil, fmt.Errorf("error goterraform.provider.package option is required in %s", desc.Path())
	}
	p.setPb(synth, option)
	p.setPkgAndPath()
	for index, optionAttribute := range p.option.Attributes {
		attr, err := synth.Attribute(optionAttribute)
		if err != nil {
			return nil, fmt.Errorf("error failed to parse %s.attributes[%d] from %s: %w", pb.E_Provider.TypeDescriptor().FullName(), index, desc.Path(), err)
		}
		if attr == nil {
			continue
		}
		p.attributes = append(p.attributes, attr)
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

func (p *provider) setPb(synth extension.Synthesizer, option *pb.Option) {
	seg := strings.Split(p.option.Package, ";")
	pkg := filepath.Base(seg[len(seg)-1])
	p.pbPackageName = protogen.GoPackageName(pkg)
	seg = strings.Split(p.option.Package, ";")
	p.pbImportPath = protogen.GoImportPath(seg[0])
	p.module = ""
	if p.option.ProviderPackage == "" {
		p.option.ProviderPackage = p.option.Package
		p.module = synth.Module()
	}
	if option.Module != nil {
		p.module = *option.Module
	}
	if !strings.HasPrefix(string(p.pbImportPath), p.module) {
		p.pbImportPath = protogen.GoImportPath(filepath.Join(p.module, string(p.pbImportPath)))
	}
}

func (p *provider) Attributes() []extension.Attribute {
	return p.attributes
}

func (p *provider) GoImportPath() protogen.GoImportPath {
	return protogen.GoImportPath(*p.option.Module)
}

func (p *provider) Members() map[string]*pb.GoType {
	return p.option.Members
}

func (p *provider) Name() string {
	return p.option.Name
}

func (p *provider) Option() *pb.Option {
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

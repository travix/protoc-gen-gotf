package generate

import (
	_ "embed"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-goterraform/extension"
)

var _ Generator = &generator{}
var (
	//go:embed data_source.go.tmpl
	dataSourceTmpl string
	//go:embed dependency.go.tmpl
	dependencyTmpl string
	//go:embed resource.go.tmpl
	resourceTmpl string
	//go:embed provider.go.tmpl
	providerTmpl string
)

type Generator interface {
	Datasource(*protogen.GeneratedFile, extension.Block) error
	Dependency(*protogen.GeneratedFile, ...*protogen.Message) error
	Provider(*protogen.GeneratedFile, extension.Provider) error
	Resource(*protogen.GeneratedFile, extension.Block) error
}

type generator struct {
	pbImportPath, providerImportPath protogen.GoImportPath
	pbPackage, providerPackage       protogen.GoPackageName
	dataSourceTmpl                   *template.Template
	dependencyTmpl                   *template.Template
	providerTmpl                     *template.Template
	resourceTmpl                     *template.Template
}

func NewGenerator(pbImportPath, providerImportPath protogen.GoImportPath, pbPackage, providerPackage protogen.GoPackageName) Generator {
	g := &generator{pbImportPath: pbImportPath, providerImportPath: providerImportPath, pbPackage: pbPackage, providerPackage: providerPackage}
	g.dataSourceTmpl = template.Must(template.New("data_source.tmpl").Funcs(sprig.TxtFuncMap()).Parse(dataSourceTmpl))
	g.dependencyTmpl = template.Must(template.New("dependency.tmpl").Funcs(sprig.TxtFuncMap()).Parse(dependencyTmpl))
	g.providerTmpl = template.Must(template.New("resource.tmpl").Funcs(sprig.TxtFuncMap()).Parse(providerTmpl))
	g.resourceTmpl = template.Must(template.New("provider.tmpl").Funcs(sprig.TxtFuncMap()).Parse(resourceTmpl))
	return g
}

type entry struct {
	string
	any
}

func (g *generator) data(entries ...entry) map[string]any {
	data := map[string]any{
		"PbImportPath":       g.pbImportPath,
		"ProviderImportPath": g.providerImportPath,
		"PbPackage":          g.pbPackage,
		"ProviderPackage":    g.providerPackage,
	}
	for _, e := range entries {
		data[e.string] = e.any
	}
	return data
}

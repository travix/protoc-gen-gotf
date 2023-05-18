package gocode

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-gotf/extension"
	"github.com/travix/protoc-gen-gotf/pb"
)

// ensure writer implements Writer.
var _ Writer = &writer{}

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

type Writer interface {
	WriteDatasource(string, *protogen.GeneratedFile, extension.Block) error
	WriteDependency(string, *protogen.GeneratedFile, ...extension.Model) error
	WriteProvider(string, *protogen.GeneratedFile, extension.Provider, bool) error
	WriteResource(string, *protogen.GeneratedFile, extension.Block) error
	HasServiceClient(blocks []extension.Block) bool
}

type writer struct {
	pbImportPath, providerImportPath   protogen.GoImportPath
	pbPackageName, providerPackageName protogen.GoPackageName
	dataSourceTmpl                     *template.Template
	dependencyTmpl                     *template.Template
	providerTmpl                       *template.Template
	resourceTmpl                       *template.Template
}

func NewWriter(pbImportPath, providerImportPath protogen.GoImportPath, pbPackageName, providerPackageName protogen.GoPackageName) Writer {
	g := &writer{pbImportPath: pbImportPath, providerImportPath: providerImportPath, pbPackageName: pbPackageName, providerPackageName: providerPackageName}
	funcs := sprig.TxtFuncMap()
	funcs["MemberType"] = g.memberType
	funcs["NewServiceClient"] = g.newServiceClient
	funcs["IsServiceClient"] = g.isServiceClient
	g.dataSourceTmpl = template.Must(template.New("data_source.tmpl").Funcs(funcs).Parse(dataSourceTmpl))
	g.dependencyTmpl = template.Must(template.New("dependency.tmpl").Funcs(funcs).Parse(dependencyTmpl))
	g.providerTmpl = template.Must(template.New("provider.tmpl").Funcs(funcs).Parse(providerTmpl))
	g.resourceTmpl = template.Must(template.New("resource.tmpl").Funcs(funcs).Parse(resourceTmpl))
	return g
}

func (w *writer) memberType(goType *pb.GoType) string {
TYPE:
	switch goType.Type.(type) {
	case *pb.GoType_Builtin:
		return w.typeFromGoType(goType.GetSlice().GetBuiltin().String())
	case *pb.GoType_GoIdentity:
		return w.typeFromGoIdentity(goType.GetGoIdentity())
	case *pb.GoType_Slice:
		switch goType.GetSlice().Type.(type) {
		case *pb.Slice_Builtin:
			return fmt.Sprintf("[]%s", w.typeFromGoType(goType.GetSlice().GetBuiltin().String()))
		case *pb.Slice_Identity:
			return fmt.Sprintf("[]%s", w.typeFromGoIdentity(goType.GetSlice().GetIdentity()))
		default:
			break TYPE
		}
	case *pb.GoType_Map:
		mKey := ""
		switch goType.GetMap().Key.(type) {
		case *pb.Map_BuiltinKey:
			mKey += w.typeFromGoType(goType.GetMap().GetBuiltinKey().String())
		case *pb.Map_IdentityKey:
			mKey += w.typeFromGoIdentity(goType.GetMap().GetIdentityKey())
		default:
			break TYPE
		}
		switch goType.GetMap().Value.(type) {
		case *pb.Map_BuiltinValue:
			return fmt.Sprintf("map[%s]%s", mKey, w.typeFromGoType(goType.GetMap().GetBuiltinValue().String()))
		case *pb.Map_IdentityValue:
			return fmt.Sprintf("map[%s]%s", mKey, w.typeFromGoIdentity(goType.GetMap().GetIdentityValue()))
		default:
			break TYPE
		}
	}
	panic("not a valid member type")
}

func (w *writer) newServiceClient(goType *pb.GoType) string {
	if _, ok := goType.Type.(*pb.GoType_GoIdentity); ok {
		return fmt.Sprintf("%s.New%s", w.pbPackageName, goType.GetGoIdentity().Name)
	}
	panic("not a valid service client")
}

func (w *writer) typeFromGoIdentity(identity *pb.GoIdentity) string {
	pkg := string(w.pbPackageName)
	if identity.ImportPath != "" {
		pkg = filepath.Base(identity.ImportPath)
	}
	ptr := ""
	if identity.Ptr {
		ptr = "*"
	}
	return fmt.Sprintf("%s%s.%s", ptr, pkg, identity.Name)
}

func (w *writer) HasServiceClient(blocks []extension.Block) bool {
	for _, b := range blocks {
		for _, m := range b.Members() {
			if w.isServiceClient(m) {
				return true
			}
		}
	}
	return false
}

func (w *writer) isServiceClient(goType *pb.GoType) bool {
	if _, ok := goType.Type.(*pb.GoType_GoIdentity); ok {
		return strings.HasSuffix(goType.GetGoIdentity().Name, "ServiceClient")
	}
	return false
}

func (w *writer) typeFromGoType(tp string) string {
	if strings.HasSuffix(tp, "_ptr") {
		tp = "*" + strings.TrimSuffix(tp, "_ptr")
	}
	return tp
}

func (w *writer) formatAndWrite(filename string, file *protogen.GeneratedFile, src []byte) error {
	var err error
	log.Debug().Str("filename", filename).Bool("src", src == nil).Msg("formatting")
	if src, err = w.Format(src, filename); err != nil {
		return err
	}
	if _, err = file.Write(src); err != nil {
		return fmt.Errorf("failed to write %s: %w", filename, err)
	}
	return nil
}

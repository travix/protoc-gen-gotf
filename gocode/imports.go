package gocode

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/travix/protoc-gen-gotf/pb"
)

type _import struct {
	path string
	string
}

func (w *writer) importForGoType(goType *pb.GoType) []_import {
TYPE:
	switch goType.Type.(type) {
	case *pb.GoType_GoIdentity:
		return w.newImport(goType.GetGoIdentity())
	case *pb.GoType_Slice:
		switch goType.GetSlice().Type.(type) {
		case *pb.Slice_Identity:
			return w.newImport(goType.GetSlice().GetIdentity())
		default:
			break TYPE
		}
	case *pb.GoType_Map:
		imports := make([]_import, 0)
		if _, ok := goType.GetMap().Key.(*pb.Map_IdentityKey); ok {
			imports = append(imports, w.newImport(goType.GetMap().GetIdentityKey())...)
		}
		if _, ok := goType.GetMap().Value.(*pb.Map_IdentityValue); ok {
			imports = append(imports, w.newImport(goType.GetMap().GetIdentityValue())...)
		}
		if len(imports) > 0 {
			return imports
		}
	}
	return nil
}

func (w *writer) newImport(goIdentity *pb.GoIdentity) []_import {
	importPath := goIdentity.GetImportPath()
	if importPath == "" {
		if strings.HasSuffix(goIdentity.Name, "ServiceClient") {
			log.Debug().Msgf("Importing google.golang.org/grpc for %s", goIdentity.Name)
			return []_import{
				{path: "fmt"},
				{path: "google.golang.org/grpc"},
				{path: "github.com/travix/gotf/cntxt"},
			}
		}
		return nil
	}
	return []_import{{importPath, filepath.Base(importPath)}}
}

func (w *writer) importStrings(imports []_import) []string {
	importStrings := make([]string, 0)
	imports = uniqImports(imports)
	for _, imprt := range imports {
		importStrings = append(importStrings, strings.TrimSpace(fmt.Sprintf(`%s "%s"`, imprt.string, imprt.path)))
	}
	return importStrings
}

func uniqImports(imports []_import) []_import {
	sort.Slice(imports, func(i, j int) bool {
		return imports[i].path < imports[j].path
	})
	dest := make([]_import, 0)
UNIQUE:
	for i := 0; i < len(imports); i++ {
		newImport := imports[i]
		for _, existing := range dest {
			if existing.path == newImport.path {
				continue UNIQUE
			}
		}
		dest = append(dest, newImport)
	}
	return dest
}

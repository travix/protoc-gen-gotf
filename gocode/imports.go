package gocode

import (
	"fmt"
	"sort"
	"strings"
)

type _import struct {
	path string
	string
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

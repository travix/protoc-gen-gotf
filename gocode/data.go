package gocode

import "github.com/travix/protoc-gen-gotf/extension"

type entry struct {
	string
	any
}

func (w *writer) data(entries ...entry) map[string]any {
	data := map[string]any{
		"PbImportPath":        w.pbImportPath,
		"ProviderImportPath":  w.providerImportPath,
		"PbPackageName":       w.pbPackageName,
		"ProviderPackageName": w.providerPackageName,
		"Version":             "local",
		"ProtocVersion":       "local",
		"Location":            "local",
	}
	for _, e := range entries {
		data[e.string] = e.any
	}
	return data
}

func (w *writer) blockData(block extension.Block, defaultImports []_import) map[string]any {
	hasServiceClient := false
	imports := append([]_import{}, defaultImports...)
	for _, m := range block.Members() {
		imports = append(imports, w.importForGoType(m)...)
		if w.isServiceClient(m) {
			hasServiceClient = true
		}
	}
	importStrings := w.importStrings(imports)
	return w.data(
		entry{"Block", block},
		entry{"Imports", importStrings},
		entry{"HasServiceClient", hasServiceClient})
}

func (w *writer) providerData(provider extension.Provider, hasServiceClient bool, defaultImports []_import) map[string]any {
	imports := append([]_import{}, defaultImports...)
	for _, m := range provider.Members() {
		imports = append(imports, w.importForGoType(m)...)
	}
	importStrings := w.importStrings(imports)
	return w.data(
		entry{"Provider", provider},
		entry{"Imports", importStrings},
		entry{"HasServiceClient", hasServiceClient})
}

func (w *writer) dependencyData(models []extension.Model, defaultImports []_import) map[string]any {
	imports := append([]_import{}, defaultImports...)
	importStrings := w.importStrings(imports)
	return w.data(
		entry{"Models", models},
		entry{"Imports", importStrings})
}

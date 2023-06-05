package gocode

import "github.com/travix/protoc-gen-gotf/extension"

type entry struct {
	string
	any
}

func (w *writer) data(entries ...entry) map[string]any {
	data := map[string]any{
		"ExecImportPath":      w.ExecImportPath,
		"ExecPackageName":     w.ExecPackageName,
		"PbImportPath":        w.PbImportPath,
		"PbPackageName":       w.PbPackageName,
		"ProviderImportPath":  w.ProviderImportPath,
		"ProviderPackageName": w.ProviderPackageName,
		"Version":             w.version,
	}
	for _, e := range entries {
		data[e.string] = e.any
	}
	return data
}

func (w *writer) blockData(block extension.Block, importsArg []_import) map[string]any {
	imports := make([]_import, len(importsArg))
	copy(imports, importsArg)
	// nolint:makezero // https://github.com/ashanbrown/makezero/issues/12
	imports = append(imports, _import{path: string(w.PbImportPath), string: string(w.PbPackageName)})
	if block.HasServiceClient() {
		// nolint:makezero
		imports = append(imports, _import{path: "google.golang.org/grpc"}, _import{path: "fmt"})
	}
	return w.data(
		entry{"Block", block},
		entry{"Imports", w.importStrings(imports)})
}

func (w *writer) execData(block extension.Block, importsArg []_import) map[string]any {
	imports := make([]_import, len(importsArg))
	copy(imports, importsArg)
	// nolint:makezero // https://github.com/ashanbrown/makezero/issues/12
	imports = append(imports,
		_import{path: string(w.PbImportPath), string: string(w.PbPackageName)},
		_import{path: string(w.ProviderImportPath), string: string(w.ProviderPackageName)},
	)
	return w.data(
		entry{"Block", block},
		entry{"Imports", w.importStrings(imports)})
}

func (w *writer) providerData(provider extension.Provider, hasServiceClient bool, importsArg []_import) map[string]any {
	imports := make([]_import, len(importsArg))
	copy(imports, importsArg)
	// nolint:makezero // https://github.com/ashanbrown/makezero/issues/12
	imports = append(imports, _import{path: string(w.PbImportPath), string: string(w.PbPackageName)})
	return w.data(
		entry{"Provider", provider},
		entry{"Imports", w.importStrings(imports)},
		entry{"HasServiceClient", hasServiceClient})
}

func (w *writer) dependencyData(models []extension.Model, defaultImports []_import) map[string]any {
	return w.data(
		entry{"Models", models},
		entry{"Imports", w.importStrings(defaultImports)})
}

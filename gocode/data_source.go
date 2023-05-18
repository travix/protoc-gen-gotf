package gocode

import (
	"bytes"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-gotf/extension"
)

var defaultDatasourceImports = []_import{
	{path: "context"},
	{path: "github.com/hashicorp/terraform-plugin-framework/datasource"},
	{path: "github.com/hashicorp/terraform-plugin-framework/datasource/schema"},
	{path: "github.com/hashicorp/terraform-plugin-log/tflog"},
	{path: "github.com/travix/gotf/dtsrc"},
}

func (w *writer) WriteDatasource(filename string, file *protogen.GeneratedFile, block extension.Block) error {
	imports := make([]_import, len(defaultDatasourceImports))
	copy(imports, defaultDatasourceImports)
	// nolint:makezero // https://github.com/ashanbrown/makezero/issues/12
	imports = append(imports, _import{path: string(w.pbImportPath), string: string(w.pbPackageName)})
	data := w.blockData(block, imports)
	code := &bytes.Buffer{}
	if err := w.dataSourceTmpl.Execute(code, data); err != nil {
		return fmt.Errorf("failed to execute %s template: %w", w.dataSourceTmpl.Name(), err)
	}
	return w.formatAndWrite(filename, file, code.Bytes())
}

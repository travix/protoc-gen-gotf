package gocode

import (
	"bytes"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-gotf/extension"
)

var defaultResourceImports = []_import{
	{path: "context"},
	{path: "github.com/hashicorp/terraform-plugin-framework/resource"},
	{path: "github.com/hashicorp/terraform-plugin-framework/resource/schema"},
	{path: "github.com/hashicorp/terraform-plugin-log/tflog"},
	{path: "github.com/travix/gotf/rsrc"},
}

func (w *writer) WriteResource(filename string, file *protogen.GeneratedFile, block extension.Block) error {
	imports := make([]_import, len(defaultResourceImports))
	copy(imports, defaultResourceImports)
	// nolint:makezero // https://github.com/ashanbrown/makezero/issues/12
	imports = append(imports, _import{path: string(w.pbImportPath), string: string(w.pbPackageName)})
	data := w.blockData(block, imports)
	code := &bytes.Buffer{}
	if err := w.resourceTmpl.Execute(code, data); err != nil {
		return fmt.Errorf("failed to execute %s template: %w", w.resourceTmpl.Name(), err)
	}
	return w.formatAndWrite(filename, file, code.Bytes())
}

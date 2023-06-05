package gocode

import (
	"bytes"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-gotf/extension"
)

var defaultDatasourceExecImports = []_import{
	{path: "context"},
	{path: "github.com/hashicorp/terraform-plugin-framework/datasource"},
	{path: "github.com/hashicorp/terraform-plugin-framework/diag"},
}

func (w *writer) WriteDatasourceExec(filename string, file *protogen.GeneratedFile, block extension.Block) error {
	data := w.execData(block, defaultDatasourceExecImports)
	code := &bytes.Buffer{}
	if err := w.templates.ExecuteTemplate(code, dataSourceExecTemplate, data); err != nil {
		return fmt.Errorf("failed to execute %s template: %w", dataSourceExecTemplate, err)
	}
	return w.formatAndWrite(filename, file, code.Bytes())
}

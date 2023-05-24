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
	data := w.blockData(block, defaultDatasourceImports)
	code := &bytes.Buffer{}
	if err := w.templates.ExecuteTemplate(code, dataSourceTemplate, data); err != nil {
		return fmt.Errorf("failed to execute %s template: %w", dataSourceTemplate, err)
	}
	return w.formatAndWrite(filename, file, code.Bytes())
}

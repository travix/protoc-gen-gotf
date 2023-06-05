package gocode

import (
	"bytes"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-gotf/extension"
)

var defaultResourceExecImports = []_import{
	{path: "context"},
	{path: "github.com/hashicorp/terraform-plugin-framework/resource"},
	{path: "github.com/hashicorp/terraform-plugin-framework/diag"},
}

func (w *writer) WriteResourceExec(filename string, file *protogen.GeneratedFile, block extension.Block) error {
	data := w.execData(block, defaultResourceExecImports)
	code := &bytes.Buffer{}
	if err := w.templates.ExecuteTemplate(code, resourceExecTemplate, data); err != nil {
		return fmt.Errorf("failed to execute %s template: %w", resourceExecTemplate, err)
	}
	return w.formatAndWrite(filename, file, code.Bytes())
}

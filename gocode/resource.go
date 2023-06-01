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
	{path: "github.com/travix/gotf"},
	{path: "github.com/travix/gotf/rsrc"},
}

func (w *writer) WriteResource(filename string, file *protogen.GeneratedFile, block extension.Block) error {
	data := w.blockData(block, defaultResourceImports)
	code := &bytes.Buffer{}
	if err := w.templates.ExecuteTemplate(code, resourceTemplate, data); err != nil {
		return fmt.Errorf("failed to execute %s template: %w", resourceTemplate, err)
	}
	return w.formatAndWrite(filename, file, code.Bytes())
}

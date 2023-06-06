package gocode

import (
	"bytes"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-gotf/extension"
)

var defaultProviderExecImports = []_import{
	{path: "context"},
	{path: "github.com/hashicorp/terraform-plugin-framework/datasource"},
	{path: "github.com/hashicorp/terraform-plugin-framework/resource"},
	{path: "github.com/hashicorp/terraform-plugin-framework/diag"},
	{path: "google.golang.org/grpc"},
}

func (w *writer) WriteProviderExec(filename string, file *protogen.GeneratedFile, provider extension.Provider) error {
	data := w.providerExecData(provider, defaultProviderExecImports)
	code := &bytes.Buffer{}
	if err := w.templates.ExecuteTemplate(code, providerExecTemplate, data); err != nil {
		return fmt.Errorf("failed to execute %s template: %w", providerExecTemplate, err)
	}
	return w.formatAndWrite(filename, file, code.Bytes())
}

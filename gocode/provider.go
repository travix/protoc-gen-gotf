package gocode

import (
	"bytes"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-gotf/extension"
)

var defaultProviderImports = []_import{
	{path: "context"},
	{path: "github.com/hashicorp/terraform-plugin-framework/datasource"},
	{path: "github.com/hashicorp/terraform-plugin-framework/provider"},
	{path: "github.com/hashicorp/terraform-plugin-framework/provider/schema"},
	{path: "github.com/hashicorp/terraform-plugin-framework/resource"},
	{path: "github.com/hashicorp/terraform-plugin-log/tflog"},
	{path: "github.com/travix/gotf"},
	{path: "github.com/travix/gotf/prvdr"},
}

func (w *writer) WriteProvider(filename string, file *protogen.GeneratedFile, provider extension.Provider, hasServiceClient bool) error {
	data := w.providerData(provider, hasServiceClient, defaultProviderImports)
	code := &bytes.Buffer{}
	if err := w.templates.ExecuteTemplate(code, providerTemplate, data); err != nil {
		return fmt.Errorf("failed to execute %s template: %w", providerTemplate, err)
	}
	return w.formatAndWrite(filename, file, code.Bytes())
}

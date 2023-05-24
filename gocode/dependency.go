package gocode

import (
	"bytes"
	"fmt"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-gotf/extension"
)

var defaultDependencyImports = []_import{
	{path: "context"},
	{"github.com/hashicorp/terraform-plugin-framework/datasource/schema", "dschema"},
	{"github.com/hashicorp/terraform-plugin-framework/provider/schema", "pschema"},
	{"github.com/hashicorp/terraform-plugin-framework/resource/schema", "rschema"},
	{path: "github.com/hashicorp/terraform-plugin-framework/attr"},
	{path: "github.com/hashicorp/terraform-plugin-framework/diag"},
	{path: "github.com/hashicorp/terraform-plugin-framework/types"},
	{path: "github.com/hashicorp/terraform-plugin-framework/types/basetypes"},
	{path: "github.com/hashicorp/terraform-plugin-go/tftypes"},
}

func (w *writer) WriteDependency(filename string, file *protogen.GeneratedFile, models ...extension.Model) error {
	data := w.dependencyData(models, defaultDependencyImports)
	log.Debug().Int("models", len(models)).Msg("generating dependency")
	code := &bytes.Buffer{}
	if err := w.templates.ExecuteTemplate(code, dependencyTemplate, data); err != nil {
		return fmt.Errorf("failed to execute %s template: %w", dependencyTemplate, err)
	}
	return w.formatAndWrite(filename, file, code.Bytes())
}

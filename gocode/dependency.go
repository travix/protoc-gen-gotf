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
	{path: "github.com/hashicorp/terraform-plugin-framework/attr"},
	{"github.com/hashicorp/terraform-plugin-framework/datasource/schema", "dschema"},
	{path: "github.com/hashicorp/terraform-plugin-framework/diag"},
	{"github.com/hashicorp/terraform-plugin-framework/resource/schema", "rschema"},
	{path: "github.com/hashicorp/terraform-plugin-framework/types"},
	{path: "github.com/hashicorp/terraform-plugin-framework/types/basetypes"},
	{path: "github.com/hashicorp/terraform-plugin-go/tftypes"},
}

func (w *writer) WriteDependency(filename string, file *protogen.GeneratedFile, models ...extension.Model) error {
	data := w.dependencyData(models, defaultDependencyImports)
	log.Debug().Int("models", len(models)).Msg("generating dependency")
	code := &bytes.Buffer{}
	if err := w.dependencyTmpl.Execute(code, data); err != nil {
		return fmt.Errorf("failed to execute %s template: %w", w.dependencyTmpl.Name(), err)
	}
	return w.formatAndWrite(filename, file, code.Bytes())
}

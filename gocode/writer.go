package gocode

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-gotf/extension"
)

const (
	dataSourceTemplate = "data_source.tmpl"
	dependencyTemplate = "proto_tf.tmpl"
	resourceTemplate   = "resource.tmpl"
	providerTemplate   = "provider.tmpl"
)

// ensure writer implements Writer.
var _ Writer = &writer{}

//go:embed tmpls
var templates embed.FS

type Writer interface {
	WriteDatasource(string, *protogen.GeneratedFile, extension.Block) error
	WriteDependency(string, *protogen.GeneratedFile, ...extension.Model) error
	WriteProvider(string, *protogen.GeneratedFile, extension.Provider, bool) error
	WriteResource(string, *protogen.GeneratedFile, extension.Block) error
}

type writer struct {
	pbImportPath, providerImportPath   protogen.GoImportPath
	pbPackageName, providerPackageName protogen.GoPackageName
	templates                          *template.Template
	version                            string
}

func NewWriter(pbImportPath, providerImportPath protogen.GoImportPath, pbPackageName, providerPackageName protogen.GoPackageName, version string) (Writer, error) {
	w := &writer{pbImportPath: pbImportPath, providerImportPath: providerImportPath, pbPackageName: pbPackageName, providerPackageName: providerPackageName, version: version}
	err := w.addTemplates()
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *writer) formatAndWrite(filename string, file *protogen.GeneratedFile, src []byte) error {
	var err error
	if src, err = w.Format(src, filename); err != nil {
		return err
	}
	if _, err = file.Write(src); err != nil {
		return fmt.Errorf("failed to write %s: %w", filename, err)
	}
	return nil
}

func (w *writer) addTemplates() error {
	funcs := sprig.TxtFuncMap()
	w.templates = template.New("tf_templates").Funcs(funcs)
	err := fs.WalkDir(templates, "tmpls", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk templates: %w", err)
		}
		if d == nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".tmpl") {
			log.Debug().Msgf("gocode templates: file %s skipped", path)
			return nil
		}
		var templateData []byte
		if templateData, err = templates.ReadFile(path); err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}
		_, err = w.templates.New(filepath.Base(path)).Parse(string(templateData))
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}
		return nil
	})
	if err != nil {
		return err //nolint:wrapcheck
	}
	return nil
}

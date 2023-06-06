package gocode

import (
	"fmt"
	"go/format"

	"github.com/daixiang0/gci/pkg/config"
	"github.com/daixiang0/gci/pkg/gci"
	"github.com/daixiang0/gci/pkg/log"
	"github.com/daixiang0/gci/pkg/section"
	zlog "github.com/rs/zerolog/log"
)

type srcFile struct {
	path string
	src  []byte
}

func (f srcFile) Load() ([]byte, error) {
	return f.src, nil
}
func (f srcFile) Path() string {
	return f.path
}

func (w *writer) Format(src []byte, path string) ([]byte, error) {
	dst, err := format.Source(src)
	if err != nil {
		zlog.Trace().Msgf("source:\n%s", src)
		return nil, fmt.Errorf("failed to gofmt file %s: %w", path, err)
	}
	log.InitLogger()
	cfg := config.Config{
		BoolConfig: config.BoolConfig{
			CustomOrder: true,
		},
		Sections: []section.Section{
			section.Standard{},
			section.Default{},
			section.Dot{},
			section.Custom{Prefix: w.module},
		},
		SectionSeparators: section.DefaultSectionSeparators(),
	}
	_, dst, err = gci.LoadFormatGoFile(srcFile{path, dst}, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to format imports in file %s: %w", path, err)
	}
	return dst, nil
}

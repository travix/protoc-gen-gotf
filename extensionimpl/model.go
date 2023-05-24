package extensionimpl

import (
	"fmt"
	"go/token"
	"strings"
	"unicode"
	"unicode/utf8"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-gotf/extension"
	"github.com/travix/protoc-gen-gotf/pb"
)

var _ extension.Model = &model{}

type model struct {
	message    *protogen.Message
	pkgName    string
	attributes []extension.Attribute
}

func NewModel(synth synthesizer, msg *protogen.Message, explicit bool) (extension.Model, error) {
	m := &model{message: msg}
	m.pkgName = getPkgName(msg.Desc.ParentFile().Options())
	for _, field := range msg.Fields {
		attr, err := synth.Attribute(field, explicit)
		if err != nil {
			return nil, err
		}
		if attr == nil {
			continue
		}
		m.attributes = append(m.attributes, attr)
	}
	return m, nil
}

func (m *model) FindAttribute(name string) (extension.Attribute, bool) {
	for _, attr := range m.attributes {
		if attr.Field().GoName == name {
			return attr, true
		}
	}
	return nil, false
}

func (m *model) GoName() string {
	return fmt.Sprintf("%s.%s", m.pkgName, m.message.GoIdent.GoName)
}

func (m *model) Attributes() []extension.Attribute {
	return m.attributes
}

func (m *model) PackageName() string {
	return m.pkgName
}

func (m *model) Message() *protogen.Message {
	return m.message
}

func (m *model) IsProvider() bool {
	_, ok := getOptions[*pb.Provider](m.message.Desc, pb.E_Provider)
	return ok
}

func (m *model) IsDatasource() bool {
	_, ok := getOptions[*pb.Block](m.message.Desc, pb.E_Datasource)
	return ok
}

func (m *model) IsResource() bool {
	_, ok := getOptions[*pb.Block](m.message.Desc, pb.E_Resource)
	return ok
}

func GoSanitized(s string) string {
	s = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return '_'
	}, s)
	r, _ := utf8.DecodeRuneInString(s)
	if token.Lookup(s).IsKeyword() || !unicode.IsLetter(r) {
		return "_" + s
	}
	return s
}

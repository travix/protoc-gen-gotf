package extensionimpl

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-goterraform/extension"
)

var _ extension.Model = &model{}

type model struct {
	message    *protogen.Message
	attributes []extension.Attribute
}

func NewModel(synth synthesizer, msg *protogen.Message, explicit bool) (extension.Model, error) {
	m := &model{message: msg}
	for _, field := range msg.Fields {
		attr, err := synth.FieldAttribute(field, explicit)
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

func (m *model) Attributes() []extension.Attribute {
	return m.attributes
}

func (m *model) Message() *protogen.Message {
	return m.message
}

package extension

import (
	"google.golang.org/protobuf/compiler/protogen"
)

type Model interface {
	Attributes() []Attribute
	Message() *protogen.Message
}

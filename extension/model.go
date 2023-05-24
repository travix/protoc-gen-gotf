package extension

import (
	"google.golang.org/protobuf/compiler/protogen"
)

type Model interface {
	Attributes() []Attribute
	Message() *protogen.Message
	GoName() string
	PackageName() string
	FindAttribute(name string) (Attribute, bool)
	IsProvider() bool
	IsDatasource() bool
	IsResource() bool
}

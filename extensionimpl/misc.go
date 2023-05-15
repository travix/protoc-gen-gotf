package extensionimpl

import (
	"regexp"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/travix/protoc-gen-goterraform/pb"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func getOptions[T *pb.Attribute | *pb.Block | *pb.Option](desc protoreflect.Descriptor, extType protoreflect.ExtensionType) (T, bool) {
	if desc == nil {
		return nil, false
	}
	optMayBe := desc.Options()
	if optMayBe == nil {
		return nil, false
	}
	option, ok := proto.GetExtension(optMayBe, extType).(T)
	return option, ok && option != nil
}

func toSnakeCase(name string) string {
	snake := matchFirstCap.ReplaceAllString(name, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

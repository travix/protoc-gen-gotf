package terraform

import (
	"regexp"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func toSnakeCase(name string) string {
	snake := matchFirstCap.ReplaceAllString(name, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func getOptions[T *descriptorpb.FileOptions | *descriptorpb.MessageOptions | *descriptorpb.FieldOptions](desc protoreflect.Descriptor) (T, bool) {
	if desc == nil {
		return nil, false
	}
	optMayBe := desc.Options()
	if optMayBe == nil {
		return nil, false
	}
	option, ok := optMayBe.(T)
	return option, ok && option != nil
}

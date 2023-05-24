package extensionimpl

import (
	"path"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/travix/protoc-gen-gotf/pb"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func getOptions[T *pb.Attribute | *pb.Block | *pb.Provider](desc protoreflect.Descriptor, extType protoreflect.ExtensionType) (T, bool) {
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

func toCamelCase(name string) string {
	strs := strings.Split(name, "_")
	title := cases.Title(language.Und).String
	for index, str := range strs {
		strs[index] = title(str)
	}
	return strings.Join(strs, "")
}

func deferToComment(direct *string, comments protogen.CommentSet) *string {
	if direct != nil && *direct != "" {
		return direct
	}
	var str string
	for index, c := range comments.LeadingDetached {
		if index > 0 {
			str += "\n"
		}
		str += c.String()
	}
	str += string(comments.Leading)
	str += string(comments.Trailing)
	str = strings.TrimSpace(str)
	return &str
}

func getPkgName(options protoreflect.ProtoMessage) string {
	fileOpt, _ := options.(*descriptorpb.FileOptions)
	goPkg := fileOpt.GetGoPackage()
	if i := strings.Index(goPkg, ";"); i >= 0 {
		return goPkg[i+1:]
	}
	return GoSanitized(path.Base(goPkg))
}

func getImportPath(options protoreflect.ProtoMessage) string {
	fileOpt, _ := options.(*descriptorpb.FileOptions)
	goPkg := fileOpt.GetGoPackage()
	if i := strings.Index(goPkg, ";"); i >= 0 {
		return goPkg[:i]
	}
	return goPkg
}

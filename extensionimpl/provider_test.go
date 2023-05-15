package extensionimpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/travix/protoc-gen-goterraform/extension"
	"github.com/travix/protoc-gen-goterraform/pb"
)

type MockedFileDescriptor struct {
	protoreflect.FileDescriptor
	mock.Mock
}

func TestNewProvider(t *testing.T) {
	t.Run("Returns nil if option is nil", func(t *testing.T) {
		mockedSynth := &extension.MockedSynthesizer{}
		mockedSynth.On("FileOption", mock.Anything).Return(nil)
		mocked := &MockedFileDescriptor{}
		got, err := NewProvider(mockedSynth, mocked)
		assert.NoError(t, err)
		assert.Nil(t, got)
		mockedSynth.AssertExpectations(t)
	})
	t.Run("Returns from field and option", func(t *testing.T) {
		mockedSynth := &extension.MockedSynthesizer{}
		mockedSynth.On("Attribute", mock.Anything).Return(&extension.MockedAttribute{}, nil)
		mockedSynth.On("Module").Return(protogen.GoImportPath(""))
		mocked := &MockedFileDescriptor{}
		mockedSynth.On("FileOption", mock.Anything).Return(&pb.Option{
			Name:            "p1",
			Package:         "pb",
			ProviderPackage: "provider",
			Module:          proto.String("module"),
			Attributes: []*pb.Attribute{
				{
					Name: proto.String("name"),
					Attr: pb.AttrType_string_attr.Enum(),
				},
			},
			Members: map[string]*pb.GoType{
				"access_key": {
					Type: &pb.GoType_Builtin{
						Builtin: pb.Builtin_string,
					},
				},
			},
		})
		got, err := NewProvider(mockedSynth, mocked)
		assert.NoError(t, err)
		assert.Equal(t, "p1", got.Name())
		assert.Equal(t, protogen.GoPackageName("pb"), got.PbPackageName())
		assert.Equal(t, protogen.GoPackageName("provider"), got.PackageName())
		assert.Equal(t, protogen.GoImportPath("module/pb"), got.PbImportPath())
		assert.Equal(t, protogen.GoImportPath("module/provider"), got.ImportPath())
		assert.Len(t, got.Members(), 1, "should have one member")
		assert.Len(t, got.Attributes(), 1, "should have one attribute")
		mocked.AssertExpectations(t)
	})
}

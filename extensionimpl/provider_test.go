package extensionimpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"

	"github.com/travix/protoc-gen-gotf/extension"
	"github.com/travix/protoc-gen-gotf/pb"
)

func TestNewProvider(t *testing.T) {
	t.Run("Returns nil if option is nil", func(t *testing.T) {
		mockedSynth := &extension.MockedSynthesizer{}
		mockedSynth.On("ProviderOption", mock.Anything).Return(nil)
		arg := &protogen.Message{}
		got, err := NewProvider(mockedSynth, arg)
		assert.NoError(t, err)
		assert.Nil(t, got)
		mockedSynth.AssertExpectations(t)
	})
	t.Run("Returns from field and option", func(t *testing.T) {
		mockedSynth := &extension.MockedSynthesizer{}
		mockedSynth.On("Model", mock.Anything, false).Return(&extension.MockedModel{}, nil)
		mockedSynth.On("ProviderOption", mock.Anything).Return(&pb.Provider{
			Name:            "p1",
			PbPackage:       "pb",
			ProviderPackage: "provider",
			Module:          proto.String("module"),
			Members: map[string]*pb.GoType{
				"access_key": {
					Type: &pb.GoType_Builtin{
						Builtin: pb.Builtin_string,
					},
				},
			},
		})
		arg := &protogen.Message{
			GoIdent: protogen.GoIdent{
				GoName: "test",
			},
			Fields: []*protogen.Field{
				{
					GoIdent: protogen.GoIdent{
						GoName: "test",
					},
				},
			},
		}
		got, err := NewProvider(mockedSynth, arg)
		assert.NoError(t, err)
		assert.Equal(t, "p1", got.TfName())
		assert.Equal(t, protogen.GoPackageName("pb"), got.PbPackageName())
		assert.Equal(t, protogen.GoPackageName("provider"), got.PackageName())
		assert.Equal(t, protogen.GoImportPath("module/pb"), got.PbImportPath())
		assert.Equal(t, protogen.GoImportPath("module/provider"), got.ImportPath())
		assert.Len(t, got.Members(), 1, "should have one member")
		assert.NotNil(t, got.Model(), "should have model")
		mockedSynth.AssertExpectations(t)
	})
}

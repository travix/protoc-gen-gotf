package extensionimpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/compiler/protogen"

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
		mockedSynth.On("MessagePackageName", mock.Anything).Return(protogen.GoPackageName("pb"))
		mockedSynth.On("MessageImportPath", mock.Anything).Return(protogen.GoImportPath("./pb"))
		mockedSynth.On("Module").Return("mod-name")
		mockedSynth.On("ProviderOption", mock.Anything).Return(&pb.Provider{
			Name:            "p1",
			ProviderPackage: "provider",
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
		assert.Equal(t, "p1", got.TerraformName())
		assert.Equal(t, protogen.GoPackageName("pb"), got.PackageData().PbPackageName)
		assert.Equal(t, protogen.GoPackageName("provider"), got.PackageData().ProviderPackageName)
		assert.Equal(t, protogen.GoImportPath("./pb"), got.PackageData().PbImportPath)
		assert.Equal(t, protogen.GoImportPath("mod-name/provider"), got.PackageData().ProviderImportPath)
		assert.NotNil(t, got.Model(), "should have model")
		mockedSynth.AssertExpectations(t)
	})
}

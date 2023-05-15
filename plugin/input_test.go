package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/travix/protoc-gen-goterraform/extension"
	"github.com/travix/protoc-gen-goterraform/pb"
	"github.com/travix/protoc-gen-goterraform/testdata"
)

func TestNewInput(t *testing.T) {
	stub := testdata.NewStub(t, "../testdata/valid-01/code_generator_request.pb.bin")
	// setup mocks
	mockedModel1 := &extension.MockedModel{}
	mockedModel2 := &extension.MockedModel{}
	mockedModel1.On("Message").Return(&protogen.Message{GoIdent: protogen.GoIdent{GoName: "User"}})
	mockedModel2.On("Message").Return(&protogen.Message{GoIdent: protogen.GoIdent{GoName: "UserData"}})
	mockedBlock1 := &extension.MockedBlock{}
	mockedBlock2 := &extension.MockedBlock{}
	mockedBlock1.On("Model").Return(mockedModel1)
	mockedBlock1.On("Name").Return("User")
	mockedBlock1.On("TypeName").Return("resource")
	mockedBlock1.On("Type").Return(pb.E_Resource)
	mockedBlock2.On("Model").Return(mockedModel2)
	mockedBlock2.On("Name").Return("User")
	mockedBlock2.On("TypeName").Return("datasource")
	mockedBlock2.On("Type").Return(pb.E_Datasource)
	mockedSynthesizer := &extension.MockedSynthesizer{}
	mockedSynthesizer.On("Provider", mock.Anything).Return(&extension.MockedProvider{}, nil)
	mockedSynthesizer.On("Block", mock.AnythingOfType("*protogen.Message"), mock.AnythingOfType("*impl.ExtensionInfo")).Return(
		func(msg *protogen.Message, blockType protoreflect.ExtensionType) (extension.Block, error) {
			if string(msg.Desc.Name()) == "User" && blockType == pb.E_Resource {
				return mockedBlock1, nil
			}
			if string(msg.Desc.Name()) == "UserData" && blockType == pb.E_Datasource {
				return mockedBlock2, nil
			}
			return nil, nil
		})
	in, err := NewInput(stub.Plugin(), mockedSynthesizer)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, in) {
		return
	}
	assert.NotNil(t, in.Provider())
	assert.Len(t, in.Resources(), 1, "should have 1 resource")
	assert.Len(t, in.Datasources(), 1, "should have 1 datasource")
	assert.Len(t, in.Dependencies(), 2, "should have 2 dependencies")
	mockedSynthesizer.AssertExpectations(t)
	mockedBlock1.AssertExpectations(t)
	mockedBlock2.AssertExpectations(t)
}

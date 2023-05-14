package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/travix/protoc-gen-goterraform/internal/terraform"
	"github.com/travix/protoc-gen-goterraform/pb"
	"github.com/travix/protoc-gen-goterraform/testdata"
)

func TestNewInput(t *testing.T) {
	stub := testdata.NewStub(t, "../testdata/min-valid/code_generator_request.pb.bin")
	// setup mocks
	mockedTypeValue1 := &terraform.MockedTypeValue{}
	mockedTypeValue2 := &terraform.MockedTypeValue{}
	mockedTypeValue1.On("TerraformNative").Return(false)
	mockedTypeValue1.On("Message").Return(stub.Message("User"))
	mockedTypeValue2.On("TerraformNative").Return(false)
	mockedTypeValue2.On("Message").Return(stub.Message("UserData"))
	mockedAttribute1 := &terraform.MockedAttribute{}
	mockedAttribute2 := &terraform.MockedAttribute{}
	mockedAttribute1.On("TypeValue").Return(mockedTypeValue1)
	mockedAttribute2.On("TypeValue").Return(mockedTypeValue2)
	mockedBlock1 := &terraform.MockedBlock{}
	mockedBlock2 := &terraform.MockedBlock{}
	mockedBlock1.On("Attributes").Return([]terraform.Attribute{mockedAttribute1})
	mockedBlock1.On("Name").Return("User")
	mockedBlock1.On("TypeName").Return("resource")
	mockedBlock1.On("Type").Return(pb.E_Resource)
	mockedBlock2.On("Attributes").Return([]terraform.Attribute{mockedAttribute2})
	mockedBlock2.On("Name").Return("User")
	mockedBlock2.On("TypeName").Return("datasource")
	mockedBlock2.On("Type").Return(pb.E_Datasource)
	mockedSynthesizer := &terraform.MockedSynthesizer{}
	mockedSynthesizer.On("Provider", mock.Anything).Return(&terraform.MockedProvider{}, nil)
	mockedSynthesizer.On("Block", mock.AnythingOfType("*protogen.Message"), mock.AnythingOfType("*impl.ExtensionInfo")).Return(
		func(msg *protogen.Message, blockType protoreflect.ExtensionType) (terraform.Block, error) {
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
	mockedAttribute1.AssertExpectations(t)
	mockedAttribute2.AssertExpectations(t)
	mockedTypeValue1.AssertExpectations(t)
	mockedTypeValue2.AssertExpectations(t)
}

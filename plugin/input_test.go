package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/travix/protoc-gen-goterraform/internal/terraform"
	"github.com/travix/protoc-gen-goterraform/pb"
	"github.com/travix/protoc-gen-goterraform/testdata"
)

func TestNewInput(t *testing.T) {
	mocked := &terraform.MockedSynthesizer{}
	mockedBlock := &terraform.MockedBlock{}
	mocked.On("Provider", mock.Anything).Return(&terraform.MockedProvider{}, nil)
	mocked.On("Block", mock.Anything, pb.E_Resource).Return(mockedBlock, nil, nil)
	mocked.On("Block", mock.Anything, pb.E_Datasource).Return(mockedBlock, nil, nil)
	mockedBlock.On("Attributes").Return(nil)
	_, err := NewInput(testdata.NewStub("../testdata/min-valid/code_generator_request.pb.bin").Plugin(), mocked)
	assert.Nil(t, err)
	mocked.AssertExpectations(t)
	// TODO: test dependencies
}

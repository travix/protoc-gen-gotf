package plugin

import (
	"os"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/travix/protoc-gen-goterraform/internal/terraform/mocks"
	"github.com/travix/protoc-gen-goterraform/pb"
)

func TestNewInput(t *testing.T) {
	mocked := &mocks.Synthesizer{}
	mocked.On("Provider", mock.Anything).Return(&mocks.Provider{}, nil)
	mocked.On("Block", mock.Anything, pb.E_Resource).Return(&mocks.Block{}, nil, nil)
	mocked.On("Block", mock.Anything, pb.E_Datasource).Return(&mocks.Block{}, nil, nil)
	_, err := NewInput(testProtoGen("testdata/minimum-valid/code_generator_request.pb.bin"), mocked)
	assert.Nil(t, err)
	mocked.AssertExpectations(t)
	// TODO: test dependencies
}

func testProtoGen(dataFile string, extraParam ...string) *protogen.Plugin {
	reqData, err := os.ReadFile(dataFile)
	if err != nil {
		log.Fatal().Err(err).Msgf("read request data file failed: %s", dataFile)
	}
	req := &pluginpb.CodeGeneratorRequest{}
	if err = proto.Unmarshal(reqData, req); err != nil {
		log.Fatal().Err(err).Msg("unable to unmarshal request")
	}
	params := ",paths=source_relative"
	for _, p := range extraParam {
		params += "," + p
	}
	req.Parameter = proto.String(*req.Parameter + params)
	gen, err := protogen.Options{}.New(req)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create proto generator")
	}
	return gen
}

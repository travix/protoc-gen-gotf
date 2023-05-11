package plugin

import (
	"os"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/travix/protoc-gen-goterraform/pb"
)

func TestNewPlugin(t *testing.T) {
	var plug *plugin
	previous := t.Run("Plugin", func(t *testing.T) {
		_assert := assert.New(t)
		_plug, err := NewPlugin(testProtoGen("testdata/minimum-valid/code_generator_request.pb.bin"))
		if !_assert.Nil(err) || !_assert.NotNil(_plug, "returned plugin is nil") {
			return
		}
		var ok bool
		{
			if plug, ok = _plug.(*plugin); !_assert.True(ok, "p is not a *plugin") {
				return
			}
		}
	})
	previous = t.Run("Provider", func(t *testing.T) {
		if !previous {
			t.SkipNow()
		}
		_assert := assert.New(t)
		if !_assert.NotNil(plug.provider) {
			return
		}
		if !_assert.NotNil(plug.provider.Option) {
			return
		}
		expected := &pb.Option{
			Attributes: nil,
			Module:     nil,
			Name:       "valid",
			Package:    nil,
		}
		_assert.True(proto.Equal(plug.provider.Option, expected), "provider option is not equal to expected")
	})
	previous = t.Run("Resources", func(t *testing.T) {
		if !previous {
			t.SkipNow()
		}
		_assert := assert.New(t)
		p := plug.provider
		if !_assert.Len(p.Resources, 1) {
			return
		}
	})
	previous = t.Run("DataSources", func(t *testing.T) {
		if !previous {
			t.SkipNow()
		}
		_assert := assert.New(t)
		p := plug.provider
		if !_assert.Len(p.DataSources, 1) {
			return
		}
	})
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

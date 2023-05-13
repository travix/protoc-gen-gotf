package testdata

import (
	"os"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

type Stub struct {
	dataFile   string
	extraParam []string
}

func NewStub(dataFile string, extraParam ...string) *Stub {
	return &Stub{dataFile: dataFile, extraParam: extraParam}
}

func (s *Stub) Plugin() *protogen.Plugin {
	gen, err := protogen.Options{}.New(s.readBin())
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create proto generator")
	}
	return gen
}

func (s *Stub) readBin() *pluginpb.CodeGeneratorRequest {
	reqData, err := os.ReadFile(s.dataFile)
	if err != nil {
		log.Fatal().Err(err).Msgf("read request data file failed: %s", s.dataFile)
	}
	req := &pluginpb.CodeGeneratorRequest{}
	if err = proto.Unmarshal(reqData, req); err != nil {
		log.Fatal().Err(err).Msg("unable to unmarshal request")
	}
	params := ",paths=source_relative"
	for _, p := range s.extraParam {
		params += "," + p
	}
	req.Parameter = proto.String(*req.Parameter + params)
	return req
}

func (s *Stub) Messages() []*protogen.Message {
	messages := make([]*protogen.Message, 0)
	req := s.Plugin()
	for _, f := range req.Files {
		for _, m := range f.Messages {
			messages = append(messages, m)
		}
	}
	return messages
}

func (s *Stub) File(file string) *protogen.File {
	req := s.Plugin()
	for _, f := range req.Files {
		if f.Desc.Path() == file {
			return f
		}
	}
	log.Fatal().Msgf("file %s not found", file)
	return nil
}

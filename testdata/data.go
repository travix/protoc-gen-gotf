package testdata

import (
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

type Stub struct {
	*testing.T
	dataFile   string
	extraParam []string
}

func NewStub(t *testing.T, dataFile string, extraParam ...string) *Stub {
	return &Stub{dataFile: dataFile, extraParam: extraParam, T: t}
}

func (s *Stub) readBin() *pluginpb.CodeGeneratorRequest {
	s.Helper()
	reqData, err := os.ReadFile(s.dataFile)
	if err != nil {
		s.Fatalf("read request %s failed: %q", s.dataFile, err)
	}
	req := &pluginpb.CodeGeneratorRequest{}
	if err = proto.Unmarshal(reqData, req); err != nil {
		s.Fatalf("unable to unmarshal request: %q", err)
	}
	params := ",paths=source_relative"
	for _, p := range s.extraParam {
		params += "," + p
	}
	req.Parameter = proto.String(*req.Parameter + params)
	return req
}

func (s *Stub) Field(name string) *protogen.Field {
	s.Helper()
	req := s.Plugin()
	for _, f := range req.Files {
		for _, m := range f.Messages {
			for _, fld := range m.Fields {
				if strings.HasSuffix(string(fld.Desc.FullName()), name) {
					return fld
				}
			}
		}
	}
	s.Fatalf("field %s not found", name)
	return nil
}

func (s *Stub) File(file string) *protogen.File {
	s.Helper()
	req := s.Plugin()
	for _, f := range req.Files {
		if f.Desc.Path() == file {
			return f
		}
	}
	log.Fatal().Msgf("file %s not found", file)
	return nil
}

func (s *Stub) Files() []*protogen.File {
	s.Helper()
	req := s.Plugin()
	return req.Files
}

func (s *Stub) Message(name string) *protogen.Message {
	s.Helper()
	req := s.Plugin()
	for _, f := range req.Files {
		for _, m := range f.Messages {
			if string(m.Desc.Name()) == name {
				return m
			}
		}
	}
	log.Fatal().Msgf("message %s not found", name)
	return nil
}

func (s *Stub) Messages() []*protogen.Message {
	s.Helper()
	messages := make([]*protogen.Message, 0)
	req := s.Plugin()
	for _, f := range req.Files {
		messages = append(messages, f.Messages...)
	}
	return messages
}

func (s *Stub) Plugin() *protogen.Plugin {
	s.Helper()
	gen, err := protogen.Options{}.New(s.readBin())
	if err != nil {
		s.Fatalf("unable to create proto generator: %q", err)
	}
	return gen
}

func (s *Stub) Service(name string) *protogen.Service {
	s.Helper()
	req := s.Plugin()
	for _, f := range req.Files {
		for _, svc := range f.Services {
			if string(svc.Desc.Name()) == name {
				return svc
			}
		}
	}
	log.Fatal().Msgf("service %s not found", name)
	return nil
}

func (s *Stub) Services() []*protogen.Service {
	s.Helper()
	services := make([]*protogen.Service, 0)
	req := s.Plugin()
	for _, f := range req.Files {
		services = append(services, f.Services...)
	}
	return services
}

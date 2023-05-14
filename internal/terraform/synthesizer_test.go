package terraform

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/travix/protoc-gen-goterraform/pb"
	"github.com/travix/protoc-gen-goterraform/testdata"
)

func Test_synthesizer_Provider(t *testing.T) {
	s := &synthesizer{}
	sub := testdata.NewStub(t, "../../testdata/valid-01/code_generator_request.pb.bin")
	got, err := s.Provider(sub.File("test.proto").Desc)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, got) {
		return
	}
}

func Test_synthesizer_FieldOption(t *testing.T) {
	s := &synthesizer{}
	sub := testdata.NewStub(t, "../../testdata/valid-01/code_generator_request.pb.bin")
	got := s.FieldOption(sub.Field("UserData.email").Desc)
	if !assert.NotNil(t, got) {
		return
	}
}

func Test_synthesizer_MessageOption(t *testing.T) {
	s := &synthesizer{}
	sub := testdata.NewStub(t, "../../testdata/valid-01/code_generator_request.pb.bin")
	got := s.MessageOption(sub.Message("User").Desc, pb.E_Resource)
	if !assert.NotNil(t, got) {
		return
	}
}

package terraform

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/travix/protoc-gen-goterraform/testdata"
)

func Test_synthesizer_Provider(t *testing.T) {
	s := &synthesizer{}
	sub := testdata.NewStub("../../testdata/min-valid/code_generator_request.pb.bin")
	got, err := s.Provider(sub.File("test.proto").Desc)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, got) {
		return
	}
}

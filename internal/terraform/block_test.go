package terraform

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/travix/protoc-gen-goterraform/extensions"
	"github.com/travix/protoc-gen-goterraform/pb"
)

func TestNewBlock(t *testing.T) {
	t.Run("Returns nil when no option is found", func(t *testing.T) {
		mocked := &extensions.MockedSynthesizer{}
		arg := &protogen.Message{}
		mocked.On("MessageOption", mock.Anything, pb.E_Resource).Once().Return(nil, nil)
		got, err := NewBlock(mocked, arg, pb.E_Resource)
		if !assert.Nil(t, err) {
			return
		}
		if !assert.Nil(t, got) {
			return
		}
		mocked.AssertExpectations(t)
	})
	t.Run("Returns Block when option is found", func(t *testing.T) {
		mocked := &extensions.MockedSynthesizer{}
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
		mocked.On("MessageOption", mock.Anything, pb.E_Resource).Once().Return(&pb.Block{
			Members: map[string]*pb.GoType{
				"m1": {
					Type: &pb.GoType_Builtin{
						Builtin: pb.Builtin_int,
					},
				},
			},
		}, nil)
		mocked.On("BlockAttribute", arg.Fields[0], false).Once().Return(&extensions.MockedAttribute{}, nil)
		got, err := NewBlock(mocked, arg, pb.E_Resource)
		if !assert.Nil(t, err) {
			return
		}
		if !assert.NotNil(t, got) {
			return
		}
		assert.Equal(t, "test", got.Name())
		assert.Len(t, got.Attributes(), 1, "len(Attributes()) = 1")
		assert.Len(t, got.Members(), 1, "len(Members()) = 1")
		mocked.AssertExpectations(t)
	})
}

func Test_block(t *testing.T) {
	b := &block{
		members:    nil,
		attributes: []extensions.Attribute{&extensions.MockedAttribute{}},
		option:     &pb.Block{},
	}
	assert.Equal(t, []extensions.Attribute{&extensions.MockedAttribute{}}, b.Attributes())
	assert.Equal(t, &pb.Block{}, b.Option())
	b.setName(&protogen.Message{
		GoIdent: protogen.GoIdent{
			GoName: "test",
		},
	})
	assert.Equal(t, "test", b.Name())
}

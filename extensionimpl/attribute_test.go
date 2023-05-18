package extensionimpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/travix/protoc-gen-gotf/extension"
	"github.com/travix/protoc-gen-gotf/pb"
)

type MockedFieldDescriptor struct {
	protoreflect.FieldDescriptor
	mock.Mock
}

func (m *MockedFieldDescriptor) Kind() protoreflect.Kind {
	args := m.Called()
	return args.Get(0).(protoreflect.Kind) // nolint:forcetypeassert
}

func (m *MockedFieldDescriptor) IsList() bool {
	args := m.Called()
	return args.Get(0).(bool) // nolint:forcetypeassert
}

func (m *MockedFieldDescriptor) IsMap() bool {
	args := m.Called()
	return args.Get(0).(bool) // nolint:forcetypeassert
}

func TestNewAttribute(t *testing.T) {
	t.Run("returns nil when no option is found", func(t *testing.T) {
		mockedSynth := &extension.MockedSynthesizer{}
		mockedSynth.On("FieldOption", mock.Anything).Once().Return(nil, nil)
		field := &protogen.Field{}
		got, err := NewAttribute(mockedSynth, field, true)
		if !assert.NoError(t, err) {
			return
		}
		assert.Nil(t, got)
		mockedSynth.AssertExpectations(t)
	})
	t.Run("returns Attribute when no option is found and explicit is false", func(t *testing.T) {
		mockedSynth := &extension.MockedSynthesizer{}
		mockedSynth.On("FieldOption", mock.Anything).Once().Return(&pb.Attribute{}, nil)
		mockedDesc := &MockedFieldDescriptor{}
		mockedDesc.On("Kind").Twice().Return(protoreflect.StringKind)
		mockedDesc.On("IsList").Return(false)
		mockedDesc.On("IsMap").Return(false)
		field := &protogen.Field{GoName: "test", Desc: mockedDesc}
		got, err := NewAttribute(mockedSynth, field, false)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.NotNil(t, got) {
			return
		}
		mockedSynth.AssertExpectations(t)
		mockedDesc.AssertExpectations(t)
	})
	t.Run("skips when option.skip true", func(t *testing.T) {
		mockedSynth := &extension.MockedSynthesizer{}
		mockedSynth.On("FieldOption", mock.Anything).Once().Return(&pb.Attribute{Skip: true}, nil)
		field := &protogen.Field{}
		got, err := NewAttribute(mockedSynth, field, false)
		if !assert.NoError(t, err) {
			return
		}
		assert.Nil(t, got)
		mockedSynth.AssertExpectations(t)
	})
	t.Run("returns Attribute", func(t *testing.T) {
		mockedSynth := &extension.MockedSynthesizer{}
		mockedSynth.On("FieldOption", mock.Anything).Once().Return(&pb.Attribute{}, nil)
		mockedDesc := &MockedFieldDescriptor{}
		mockedDesc.On("Kind").Twice().Return(protoreflect.StringKind)
		mockedDesc.On("IsList").Return(false)
		mockedDesc.On("IsMap").Return(false)
		field := &protogen.Field{GoName: "test", Desc: mockedDesc}
		got, err := NewAttribute(mockedSynth, field, false)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.NotNil(t, got) {
			return
		}
		mockedSynth.AssertExpectations(t)
		mockedDesc.AssertExpectations(t)
		assert.Equal(t, SchemaString(), got.Schema())
		assert.Equal(t, TypeValueString(), got.TypeValue())
	})
}

func Test_attribute(t *testing.T) {
	t.Run("Computed", func(t *testing.T) {
		a := &attribute{option: &pb.Attribute{MustBe: pb.MustBe_Computed}}
		assert.True(t, a.Computed())
		a = &attribute{option: &pb.Attribute{MustBe: pb.MustBe_OptionalAndComputed}}
		assert.True(t, a.Computed())
		a = &attribute{option: &pb.Attribute{MustBe: pb.MustBe_Required}}
		assert.False(t, a.Computed())
		a = &attribute{option: &pb.Attribute{MustBe: pb.MustBe_Optional}}
		assert.False(t, a.Computed())
	})
	t.Run("Required", func(t *testing.T) {
		a := &attribute{option: &pb.Attribute{MustBe: pb.MustBe_Computed}}
		assert.False(t, a.Required())
		a = &attribute{option: &pb.Attribute{MustBe: pb.MustBe_OptionalAndComputed}}
		assert.False(t, a.Required())
		a = &attribute{option: &pb.Attribute{MustBe: pb.MustBe_Required}}
		assert.True(t, a.Required())
		a = &attribute{option: &pb.Attribute{MustBe: pb.MustBe_Optional}}
		assert.False(t, a.Required())
	})
	t.Run("Optional", func(t *testing.T) {
		a := &attribute{option: &pb.Attribute{MustBe: pb.MustBe_Computed}}
		assert.False(t, a.Optional())
		a = &attribute{option: &pb.Attribute{MustBe: pb.MustBe_OptionalAndComputed}}
		assert.True(t, a.Optional())
		a = &attribute{option: &pb.Attribute{MustBe: pb.MustBe_Required}}
		assert.False(t, a.Optional())
		a = &attribute{option: &pb.Attribute{MustBe: pb.MustBe_Optional}}
		assert.True(t, a.Optional())
	})
	t.Run("Sensitive", func(t *testing.T) {
		a := &attribute{option: &pb.Attribute{}}
		assert.False(t, a.Sensitive())
		a = &attribute{option: &pb.Attribute{Sensitive: proto.Bool(false)}}
		assert.False(t, a.Sensitive())
		a = &attribute{option: &pb.Attribute{Sensitive: proto.Bool(true)}}
		assert.True(t, a.Sensitive())
	})
}

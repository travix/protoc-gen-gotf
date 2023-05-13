package terraform

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/travix/protoc-gen-goterraform/mocks"
	"github.com/travix/protoc-gen-goterraform/pb"
)

func TestNewAttribute(t *testing.T) {
	type args struct {
		option *pb.Attribute
	}
	opt1 := &pb.Attribute{
		Name:          proto.String("name"),
		MustBe:        pb.MustBe_Required,
		Sensitive:     proto.Bool(true),
		Description:   proto.String("description"),
		MdDescription: proto.String("md_description"),
		Deprecation:   proto.String("deprecation"),
	}
	opt2 := &pb.Attribute{
		Name:          proto.String("name"),
		MustBe:        pb.MustBe_Required,
		Sensitive:     proto.Bool(true),
		Description:   proto.String("description"),
		MdDescription: proto.String("md_description"),
		Deprecation:   proto.String("deprecation"),
		Attr:          pb.AttrType_string_attr.Enum(),
	}
	tv := TypeValueString()
	schema := SchemaString()
	tests := []struct {
		name    string
		args    args
		want    Attribute
		wantErr assert.ErrorAssertionFunc
	}{
		{name: "returns nil", wantErr: assert.NoError},
		{name: "skips", args: args{&pb.Attribute{Skip: true}}, wantErr: assert.NoError},
		{name: "returns from option", args: args{opt1}, want: &attribute{Attribute: opt1}, wantErr: assert.NoError},
		{name: "returns with TypeValue", args: args{opt2}, want: &attribute{Attribute: opt2, typeValue: tv, schema: schema}, wantErr: assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAttribute(tt.args.option)
			if !tt.wantErr(t, err, fmt.Sprintf("NewAttribute(%v)", tt.args.option)) {
				return
			}
			assert.Equalf(t, tt.want, got, "NewAttribute(%v)", tt.args.option)
		})
	}
}

func TestNewBlockAttribute(t *testing.T) {
	t.Run("returns nil if explicit fields", func(t *testing.T) {
		mocked := &MockedSynthesizer{}
		mocked.On("getFieldOption", mock.Anything).Return(nil)
		got, err := NewBlockAttribute(mocked, &protogen.Field{}, true)
		assert.NoError(t, err)
		assert.Nil(t, got)
		mocked.AssertExpectations(t)
	})
	t.Run("returns from field and option", func(t *testing.T) {
		mockedSynth := &MockedSynthesizer{}
		mocked := &mocks.MockedFieldDescriptor{}
		mockedSynth.On("getFieldOption", mock.Anything).Return(&pb.Attribute{Name: proto.String("name")})
		mocked.On("Kind").Return(protoreflect.BoolKind)
		mocked.On("IsList").Return(false)
		mocked.On("IsMap").Return(false)
		got, err := NewBlockAttribute(mockedSynth, &protogen.Field{
			Desc: mocked,
			Comments: protogen.CommentSet{
				Leading: "description",
			},
		}, false)
		assert.NoError(t, err)
		assert.Equal(t, &attribute{
			Attribute: &pb.Attribute{
				Name:          proto.String("name"),
				Description:   proto.String("description"),
				MdDescription: proto.String("description"),
				Deprecation:   proto.String(""),
			},
			typeValue: TypeValueBool(),
			schema:    SchemaBool(),
		}, got)
		mocked.AssertExpectations(t)
	})
	t.Run("override's the type defined in option", func(t *testing.T) {
		mockedSynth := &MockedSynthesizer{}
		mocked := &mocks.MockedFieldDescriptor{}
		option := &pb.Attribute{
			Name: proto.String("name"),
			Attr: pb.AttrType_int64_attr.Enum(),
		}
		mockedSynth.On("getFieldOption", mock.Anything).Return(option)
		mocked.On("Kind").Return(protoreflect.FloatKind)
		mocked.On("IsList").Return(false)
		mocked.On("IsMap").Return(false)
		got, err := NewBlockAttribute(mockedSynth, &protogen.Field{
			Desc: mocked,
			Comments: protogen.CommentSet{
				Leading: "description",
			},
		}, false)
		assert.NoError(t, err)
		assert.Equal(t, &attribute{
			Attribute: &pb.Attribute{
				Name:          proto.String("name"),
				Description:   proto.String("description"),
				MdDescription: proto.String("description"),
				Deprecation:   proto.String(""),
			},
			typeValue: TypeValueFloat64(),
			schema:    SchemaFloat64(),
		}, got)
		mocked.AssertExpectations(t)
	})
}

func Test_attribute(t *testing.T) {
	t.Run("Computed", func(t *testing.T) {
		a := &attribute{Attribute: &pb.Attribute{MustBe: pb.MustBe_Computed}}
		assert.True(t, a.Computed())
		a = &attribute{Attribute: &pb.Attribute{MustBe: pb.MustBe_OptionalAndComputed}}
		assert.True(t, a.Computed())
		a = &attribute{Attribute: &pb.Attribute{MustBe: pb.MustBe_Required}}
		assert.False(t, a.Computed())
		a = &attribute{Attribute: &pb.Attribute{MustBe: pb.MustBe_Optional}}
		assert.False(t, a.Computed())
	})
	t.Run("Required", func(t *testing.T) {
		a := &attribute{Attribute: &pb.Attribute{MustBe: pb.MustBe_Computed}}
		assert.False(t, a.Required())
		a = &attribute{Attribute: &pb.Attribute{MustBe: pb.MustBe_OptionalAndComputed}}
		assert.False(t, a.Required())
		a = &attribute{Attribute: &pb.Attribute{MustBe: pb.MustBe_Required}}
		assert.True(t, a.Required())
		a = &attribute{Attribute: &pb.Attribute{MustBe: pb.MustBe_Optional}}
		assert.False(t, a.Required())
	})
	t.Run("Optional", func(t *testing.T) {
		a := &attribute{Attribute: &pb.Attribute{MustBe: pb.MustBe_Computed}}
		assert.False(t, a.Optional())
		a = &attribute{Attribute: &pb.Attribute{MustBe: pb.MustBe_OptionalAndComputed}}
		assert.True(t, a.Optional())
		a = &attribute{Attribute: &pb.Attribute{MustBe: pb.MustBe_Required}}
		assert.False(t, a.Optional())
		a = &attribute{Attribute: &pb.Attribute{MustBe: pb.MustBe_Optional}}
		assert.True(t, a.Optional())
	})
	t.Run("Sensitive", func(t *testing.T) {
		a := &attribute{Attribute: &pb.Attribute{}}
		assert.False(t, a.Sensitive())
		a = &attribute{Attribute: &pb.Attribute{Sensitive: proto.Bool(false)}}
		assert.False(t, a.Sensitive())
		a = &attribute{Attribute: &pb.Attribute{Sensitive: proto.Bool(true)}}
		assert.True(t, a.Sensitive())
	})
}

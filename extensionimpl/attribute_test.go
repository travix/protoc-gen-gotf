package extensionimpl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	"github.com/travix/protoc-gen-goterraform/extension"
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
		want    extension.Attribute
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

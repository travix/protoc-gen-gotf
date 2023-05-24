// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: gotf.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Must is used to set terraform block attribute's required, optional or computed
// values
type MustBe int32

const (
	MustBe_Computed            MustBe = 0 // field will be computed by provider, user can't set it in terraform script (default).
	MustBe_Optional            MustBe = 1 // field is optional, user can set it in terraform script.
	MustBe_OptionalAndComputed MustBe = 2 // field is optional, but will be computed or overriden by provider.
	MustBe_Required            MustBe = 3 // field is required, user must set it in terraform script.
)

// Enum value maps for MustBe.
var (
	MustBe_name = map[int32]string{
		0: "Computed",
		1: "Optional",
		2: "OptionalAndComputed",
		3: "Required",
	}
	MustBe_value = map[string]int32{
		"Computed":            0,
		"Optional":            1,
		"OptionalAndComputed": 2,
		"Required":            3,
	}
)

func (x MustBe) Enum() *MustBe {
	p := new(MustBe)
	*p = x
	return p
}

func (x MustBe) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MustBe) Descriptor() protoreflect.EnumDescriptor {
	return file_gotf_proto_enumTypes[0].Descriptor()
}

func (MustBe) Type() protoreflect.EnumType {
	return &file_gotf_proto_enumTypes[0]
}

func (x MustBe) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MustBe.Descriptor instead.
func (MustBe) EnumDescriptor() ([]byte, []int) {
	return file_gotf_proto_rawDescGZIP(), []int{0}
}

// GoIdentity to indicate user defined go types such as structs, interfaces, etc.
// for terraform provider and block attribute type.
type GoIdentity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// name of the user defined go type.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// import_path (go mod name) of the user defined go type.
	// defaults to current proto files go_package name.
	ImportPath string `protobuf:"bytes,2,opt,name=import_path,json=importPath,proto3" json:"import_path,omitempty"`
	// if set to true, this type will be generated as pointer type.
	// *should not be set for interfaces.*
	Ptr bool `protobuf:"varint,3,opt,name=ptr,proto3" json:"ptr,omitempty"`
}

func (x *GoIdentity) Reset() {
	*x = GoIdentity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gotf_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GoIdentity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GoIdentity) ProtoMessage() {}

func (x *GoIdentity) ProtoReflect() protoreflect.Message {
	mi := &file_gotf_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GoIdentity.ProtoReflect.Descriptor instead.
func (*GoIdentity) Descriptor() ([]byte, []int) {
	return file_gotf_proto_rawDescGZIP(), []int{0}
}

func (x *GoIdentity) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GoIdentity) GetImportPath() string {
	if x != nil {
		return x.ImportPath
	}
	return ""
}

func (x *GoIdentity) GetPtr() bool {
	if x != nil {
		return x.Ptr
	}
	return false
}

// Block allows to generate terraform resource or datasource from proto message
type Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// name for terraform resource or datasource, defaults to message name.
	// defaults to snake_case of message name.
	Name *string `protobuf:"bytes,1,opt,name=name,proto3,oneof" json:"name,omitempty"`
	// explicit_fields if set to `false` (default), all fields will be generated as
	// attributes in this terraform block.
	//
	// If set to `true`, only fields with `gotf.Attribute` will be generated
	// as an attributes in this terraform block.
	//
	// Check gotf.Attribute for default attribute properties.
	ExplicitFields bool `protobuf:"varint,2,opt,name=explicit_fields,json=explicitFields,proto3" json:"explicit_fields,omitempty"`
	// client (names) used by this block. Generated terraform go code expects these
	// clients to be generated by grpc in same go_package as other protobufs.
	// Should be complete go type name
	//
	// This calls exec with `Set<ClientName>(client)`
	// The provider exec needs to implement github.com/travix/gotf/prvdr.CanConfigureGrpc interface.
	Client []string `protobuf:"bytes,3,rep,name=client,proto3" json:"client,omitempty"`
	// description is used to set terraform block's description.
	Description string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *Block) Reset() {
	*x = Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gotf_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_gotf_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_gotf_proto_rawDescGZIP(), []int{1}
}

func (x *Block) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

func (x *Block) GetExplicitFields() bool {
	if x != nil {
		return x.ExplicitFields
	}
	return false
}

func (x *Block) GetClient() []string {
	if x != nil {
		return x.Client
	}
	return nil
}

func (x *Block) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

// Attribute to set kind properties of terraform attribute & schema for terraform block.
type Attribute struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// skip this field from from this terraform block.
	Skip bool `protobuf:"varint,1,opt,name=skip,proto3" json:"skip,omitempty"`
	// attribute must_be computed, optional, required or optional and computed.
	//
	// If attribute is optional make sure `optional` is set to true on fields as terraform
	// will return nil value for this proto field.
	MustBe MustBe `protobuf:"varint,2,opt,name=must_be,json=mustBe,proto3,enum=gotf.MustBe" json:"must_be,omitempty"`
	// name of terraform attribute, check [attribute-names] documentation for more details.
	// defaults to snake_case of message field name.
	//
	// [attribute-names]: https://developer.hashicorp.com/terraform/plugin/best-practices/naming#attribute-names
	Name *string `protobuf:"bytes,3,opt,name=name,proto3,oneof" json:"name,omitempty"`
	// sensitive if set to true field will be marked as sensitive.
	Sensitive *bool `protobuf:"varint,4,opt,name=sensitive,proto3,oneof" json:"sensitive,omitempty"`
	// description is used to set terraform attribute's description. Defaults to
	// field comment.
	Description *string `protobuf:"bytes,5,opt,name=description,proto3,oneof" json:"description,omitempty"`
	// md_description is used to set terraform attribute's markdown description of
	// field. Defaults to field comment.
	MdDescription *string `protobuf:"bytes,6,opt,name=md_description,json=mdDescription,proto3,oneof" json:"md_description,omitempty"`
	// deprecation is used to set terraform attribute's deprecation message.
	Deprecation *string `protobuf:"bytes,7,opt,name=deprecation,proto3,oneof" json:"deprecation,omitempty"`
}

func (x *Attribute) Reset() {
	*x = Attribute{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gotf_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Attribute) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Attribute) ProtoMessage() {}

func (x *Attribute) ProtoReflect() protoreflect.Message {
	mi := &file_gotf_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Attribute.ProtoReflect.Descriptor instead.
func (*Attribute) Descriptor() ([]byte, []int) {
	return file_gotf_proto_rawDescGZIP(), []int{2}
}

func (x *Attribute) GetSkip() bool {
	if x != nil {
		return x.Skip
	}
	return false
}

func (x *Attribute) GetMustBe() MustBe {
	if x != nil {
		return x.MustBe
	}
	return MustBe_Computed
}

func (x *Attribute) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

func (x *Attribute) GetSensitive() bool {
	if x != nil && x.Sensitive != nil {
		return *x.Sensitive
	}
	return false
}

func (x *Attribute) GetDescription() string {
	if x != nil && x.Description != nil {
		return *x.Description
	}
	return ""
}

func (x *Attribute) GetMdDescription() string {
	if x != nil && x.MdDescription != nil {
		return *x.MdDescription
	}
	return ""
}

func (x *Attribute) GetDeprecation() string {
	if x != nil && x.Deprecation != nil {
		return *x.Deprecation
	}
	return ""
}

// Option to create provider for this package.
type Provider struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// name is used to set terraform provider type name, all resources and datasources
	// will be generated under this provider.
	//
	// Example:
	//
	//	name = "xyz";
	//
	// This will produce terraform blocks will be
	//
	//	resource "xyz_resource_name" { ... }
	//	data "xyz_datasource_name" { ... }
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// package where generated terraform go code should be placed
	//
	// If not set files will be placed in providerpb package.
	ProviderPackage string `protobuf:"bytes,3,opt,name=provider_package,json=providerPackage,proto3" json:"provider_package,omitempty"`
	// description is used to set terraform provider's description.
	Description string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *Provider) Reset() {
	*x = Provider{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gotf_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Provider) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Provider) ProtoMessage() {}

func (x *Provider) ProtoReflect() protoreflect.Message {
	mi := &file_gotf_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Provider.ProtoReflect.Descriptor instead.
func (*Provider) Descriptor() ([]byte, []int) {
	return file_gotf_proto_rawDescGZIP(), []int{3}
}

func (x *Provider) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Provider) GetProviderPackage() string {
	if x != nil {
		return x.ProviderPackage
	}
	return ""
}

func (x *Provider) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

var file_gotf_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*Provider)(nil),
		Field:         82848,
		Name:          "gotf.provider",
		Tag:           "bytes,82848,opt,name=provider",
		Filename:      "gotf.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*Block)(nil),
		Field:         82849,
		Name:          "gotf.resource",
		Tag:           "bytes,82849,opt,name=resource",
		Filename:      "gotf.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*Block)(nil),
		Field:         82850,
		Name:          "gotf.datasource",
		Tag:           "bytes,82850,opt,name=datasource",
		Filename:      "gotf.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*Attribute)(nil),
		Field:         82851,
		Name:          "gotf.attribute",
		Tag:           "bytes,82851,opt,name=attribute",
		Filename:      "gotf.proto",
	},
}

// Extension fields to descriptorpb.MessageOptions.
var (
	// define provider for this package, should be set only once per package.
	//
	// optional gotf.Provider provider = 82848;
	E_Provider = &file_gotf_proto_extTypes[0]
	// define this message as a resource
	//
	// optional gotf.Block resource = 82849;
	E_Resource = &file_gotf_proto_extTypes[1]
	// define this message as a datasource
	//
	// optional gotf.Block datasource = 82850;
	E_Datasource = &file_gotf_proto_extTypes[2]
)

// Extension fields to descriptorpb.FieldOptions.
var (
	// define this field's terraform attribute on this resource and/ or datasource.
	//
	// optional gotf.Attribute attribute = 82851;
	E_Attribute = &file_gotf_proto_extTypes[3]
)

var File_gotf_proto protoreflect.FileDescriptor

var file_gotf_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x67, 0x6f, 0x74, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x67, 0x6f,
	0x74, 0x66, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x53, 0x0a, 0x0a, 0x47, 0x6f, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74,
	0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x69, 0x6d, 0x70,
	0x6f, 0x72, 0x74, 0x50, 0x61, 0x74, 0x68, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x74, 0x72, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x03, 0x70, 0x74, 0x72, 0x22, 0x8c, 0x01, 0x0a, 0x05, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x12, 0x17, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x00, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x27, 0x0a, 0x0f,
	0x65, 0x78, 0x70, 0x6c, 0x69, 0x63, 0x69, 0x74, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0e, 0x65, 0x78, 0x70, 0x6c, 0x69, 0x63, 0x69, 0x74, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x20, 0x0a,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x42,
	0x07, 0x0a, 0x05, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0xc6, 0x02, 0x0a, 0x09, 0x41, 0x74, 0x74,
	0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6b, 0x69, 0x70, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x73, 0x6b, 0x69, 0x70, 0x12, 0x25, 0x0a, 0x07, 0x6d, 0x75,
	0x73, 0x74, 0x5f, 0x62, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x67, 0x6f,
	0x74, 0x66, 0x2e, 0x4d, 0x75, 0x73, 0x74, 0x42, 0x65, 0x52, 0x06, 0x6d, 0x75, 0x73, 0x74, 0x42,
	0x65, 0x12, 0x17, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x00, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x09, 0x73, 0x65,
	0x6e, 0x73, 0x69, 0x74, 0x69, 0x76, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x48, 0x01, 0x52,
	0x09, 0x73, 0x65, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x76, 0x65, 0x88, 0x01, 0x01, 0x12, 0x25, 0x0a,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x02, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x88, 0x01, 0x01, 0x12, 0x2a, 0x0a, 0x0e, 0x6d, 0x64, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x48, 0x03, 0x52, 0x0d,
	0x6d, 0x64, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01,
	0x12, 0x25, 0x0a, 0x0b, 0x64, 0x65, 0x70, 0x72, 0x65, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x09, 0x48, 0x04, 0x52, 0x0b, 0x64, 0x65, 0x70, 0x72, 0x65, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x73, 0x65, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x76, 0x65, 0x42, 0x0e,
	0x0a, 0x0c, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x11,
	0x0a, 0x0f, 0x5f, 0x6d, 0x64, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x64, 0x65, 0x70, 0x72, 0x65, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x22, 0x6b, 0x0a, 0x08, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x29, 0x0a, 0x10, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x5f, 0x70, 0x61,
	0x63, 0x6b, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x70, 0x72, 0x6f,
	0x76, 0x69, 0x64, 0x65, 0x72, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x12, 0x20, 0x0a, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2a, 0x4b,
	0x0a, 0x06, 0x4d, 0x75, 0x73, 0x74, 0x42, 0x65, 0x12, 0x0c, 0x0a, 0x08, 0x43, 0x6f, 0x6d, 0x70,
	0x75, 0x74, 0x65, 0x64, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x61, 0x6c, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c,
	0x41, 0x6e, 0x64, 0x43, 0x6f, 0x6d, 0x70, 0x75, 0x74, 0x65, 0x64, 0x10, 0x02, 0x12, 0x0c, 0x0a,
	0x08, 0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x10, 0x03, 0x3a, 0x50, 0x0a, 0x08, 0x70,
	0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa0, 0x87, 0x05, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0e, 0x2e, 0x67, 0x6f, 0x74, 0x66, 0x2e, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72,
	0x52, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x88, 0x01, 0x01, 0x3a, 0x4d, 0x0a,
	0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa1, 0x87, 0x05, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x67, 0x6f, 0x74, 0x66, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52,
	0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x88, 0x01, 0x01, 0x3a, 0x51, 0x0a, 0x0a,
	0x64, 0x61, 0x74, 0x61, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa2, 0x87, 0x05, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x67, 0x6f, 0x74, 0x66, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b,
	0x52, 0x0a, 0x64, 0x61, 0x74, 0x61, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x88, 0x01, 0x01, 0x3a,
	0x51, 0x0a, 0x09, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x12, 0x1d, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa3, 0x87, 0x05, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x67, 0x6f, 0x74, 0x66, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x69,
	0x62, 0x75, 0x74, 0x65, 0x52, 0x09, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x88,
	0x01, 0x01, 0x42, 0x18, 0x5a, 0x16, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x74, 0x72, 0x61, 0x76, 0x69, 0x78, 0x2f, 0x67, 0x6f, 0x74, 0x66, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gotf_proto_rawDescOnce sync.Once
	file_gotf_proto_rawDescData = file_gotf_proto_rawDesc
)

func file_gotf_proto_rawDescGZIP() []byte {
	file_gotf_proto_rawDescOnce.Do(func() {
		file_gotf_proto_rawDescData = protoimpl.X.CompressGZIP(file_gotf_proto_rawDescData)
	})
	return file_gotf_proto_rawDescData
}

var file_gotf_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_gotf_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_gotf_proto_goTypes = []interface{}{
	(MustBe)(0),                         // 0: gotf.MustBe
	(*GoIdentity)(nil),                  // 1: gotf.GoIdentity
	(*Block)(nil),                       // 2: gotf.Block
	(*Attribute)(nil),                   // 3: gotf.Attribute
	(*Provider)(nil),                    // 4: gotf.Provider
	(*descriptorpb.MessageOptions)(nil), // 5: google.protobuf.MessageOptions
	(*descriptorpb.FieldOptions)(nil),   // 6: google.protobuf.FieldOptions
}
var file_gotf_proto_depIdxs = []int32{
	0, // 0: gotf.Attribute.must_be:type_name -> gotf.MustBe
	5, // 1: gotf.provider:extendee -> google.protobuf.MessageOptions
	5, // 2: gotf.resource:extendee -> google.protobuf.MessageOptions
	5, // 3: gotf.datasource:extendee -> google.protobuf.MessageOptions
	6, // 4: gotf.attribute:extendee -> google.protobuf.FieldOptions
	4, // 5: gotf.provider:type_name -> gotf.Provider
	2, // 6: gotf.resource:type_name -> gotf.Block
	2, // 7: gotf.datasource:type_name -> gotf.Block
	3, // 8: gotf.attribute:type_name -> gotf.Attribute
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	5, // [5:9] is the sub-list for extension type_name
	1, // [1:5] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_gotf_proto_init() }
func file_gotf_proto_init() {
	if File_gotf_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gotf_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GoIdentity); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gotf_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Block); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gotf_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Attribute); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gotf_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Provider); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_gotf_proto_msgTypes[1].OneofWrappers = []interface{}{}
	file_gotf_proto_msgTypes[2].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_gotf_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 4,
			NumServices:   0,
		},
		GoTypes:           file_gotf_proto_goTypes,
		DependencyIndexes: file_gotf_proto_depIdxs,
		EnumInfos:         file_gotf_proto_enumTypes,
		MessageInfos:      file_gotf_proto_msgTypes,
		ExtensionInfos:    file_gotf_proto_extTypes,
	}.Build()
	File_gotf_proto = out.File
	file_gotf_proto_rawDesc = nil
	file_gotf_proto_goTypes = nil
	file_gotf_proto_depIdxs = nil
}

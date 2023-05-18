// Code generated by protoc-gen-terraform. DO NOT EDIT.
// versions:
//   protoc-gen-gotf local
// 	 protoc          local
// source: local

package providerpb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/travix/gotf/cntxt"
	"github.com/travix/gotf/dtsrc"
	"google.golang.org/grpc"

	pb "github.com/travix/protoc-gen-gotf/example/pb"
)

// Ensure *GroupsDataSource fully satisfy terraform framework interfaces.
var _ datasource.DataSource = &GroupsDataSource{}

type GroupsDataSource struct {
	proxy  dtsrc.Datasource
	client pb.GroupServiceClient // defined in proto file
}

func NewGroupsDataSource(proxy dtsrc.Datasource) func() datasource.DataSource {
	if proxy == nil {
		panic("github.com/travix/gotf/dtsrc.Datasource is required")
	}
	return func() datasource.DataSource {
		return &GroupsDataSource{proxy: proxy}
	}
}

func (d *GroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_groups"
	if _proxy, ok := d.proxy.(dtsrc.CanMetadata); ok {
		ctx = d.setupContext(ctx)
		_proxy.Metadata(ctx, req, resp)
	}
}

func (d *GroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "terraform datasource",
		Attributes:          (&pb.Groups{}).DatasourceSchema(),
	}
	if _proxy, ok := d.proxy.(dtsrc.CanSchema); ok {
		ctx = d.setupContext(ctx)
		_proxy.Schema(ctx, req, resp)
	}
}

func (d *GroupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Get the gRPC client connection from the ProviderData
	if req.ProviderData == nil {
		resp.Diagnostics.AddError("Expected ProviderData to be not nil", "req.ProviderData is nil")
		return
	}
	conn, ok := req.ProviderData.(grpc.ClientConnInterface)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ProviderData Type",
			fmt.Sprintf("Expected grpc.ClientConnInterface, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	// Create the service clients
	d.client = pb.NewGroupServiceClient(conn)
	if _proxy, ok := d.proxy.(dtsrc.CanConfigure); ok {
		ctx = d.setupContext(ctx)
		_proxy.Configure(ctx, req, resp)
		return
	}
}

func (d *GroupsDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	if _proxy, ok := d.proxy.(dtsrc.CanConfigValidators); ok {
		ctx = d.setupContext(ctx)
		return _proxy.ConfigValidators(ctx)
	}
	tflog.Warn(ctx, "ConfigValidators method not implemented. Make sure argument to NewGroupsDataSource() implements dtsrc.CanConfigValidators interface")
	return nil
}

func (d *GroupsDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	if _proxy, ok := d.proxy.(dtsrc.CanValidateConfig); ok {
		ctx = d.setupContext(ctx)
		_proxy.ValidateConfig(ctx, req, resp)
		return
	}
	tflog.Warn(ctx, "ValidateConfig method not implemented. Make sure argument to NewGroupsDataSource() implements dtsrc.CanValidateConfig interface")
}

func (d *GroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data pb.Groups
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Set the service clients into the context
	ctx = cntxt.WithValue(ctx, "client", d.client)

	ctx = d.setupContext(ctx)
	diagnostics := d.proxy.Read(ctx, req, resp, &data)
	if diagnostics.HasError() {
		resp.Diagnostics.Append(diagnostics...)
		return
	}
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *GroupsDataSource) setupContext(ctx context.Context) context.Context {
	// Pass members via context
	// also passes the service clients
	ctx = cntxt.WithValue(ctx, "client", d.client)
	return ctx
}
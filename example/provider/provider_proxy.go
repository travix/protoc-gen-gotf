package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/travix/gotf/prvdr"
	"github.com/travix/protoc-gen-gotf/example/pb"
	"github.com/travix/protoc-gen-gotf/example/providerpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ prvdr.Provider = &ProviderProxy{}
var _ prvdr.CanConfigureGrpc = &ProviderProxy{}

type ProviderProxy struct {
	group  groupProxy
	groups groupsProxy
	user   userProxy
	users  usersProxy
}

func (p *ProviderProxy) ConfigureGrpc(ctx context.Context, data any) (conn grpc.ClientConnInterface, diagnostics diag.Diagnostics) {
	var model *pb.ProviderModel
	var ok bool
	if model, ok = data.(*pb.ProviderModel); ok {
		diagnostics.AddError("Failed to cast data to *pb.ProviderModel", "expected data arg to be *pb.ProviderModel")
		return
	}
	// credentials and serverAddr can be fetched from req.Config by setting
	opts := []grpc.DialOption{
		// credentials or tokens can be fetched from model by setting fields on ProviderModel protobuf
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	var err error
	tflog.Info(ctx, fmt.Sprintf("dialing grpc connection with example grcp %s", model.Endpoint))
	conn, err = grpc.Dial(model.Endpoint, opts...)
	if err != nil {
		diagnostics.AddError("Failed connecting to example grcp", fmt.Sprintf("eror in grpc connection with %s: %v", model.Endpoint, err))
		return nil, diagnostics
	}
	return
}

func (p *ProviderProxy) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		providerpb.NewUsersDataSource(p.users),
		providerpb.NewGroupsDataSource(p.groups),
	}
}

func (p *ProviderProxy) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		providerpb.NewUserResource(p.user),
		providerpb.NewGroupResource(p.group),
	}
}

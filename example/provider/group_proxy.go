package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/travix/gotf/cntxt"
	"github.com/travix/gotf/rsrc"
)

var _ rsrc.Resource = &groupProxy{}

type groupProxy struct{}

func (g groupProxy) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, data any) diag.Diagnostics {
	cntxt.Value(ctx, "client")
	//TODO implement me
	panic("implement me")
}

func (g groupProxy) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, data any) diag.Diagnostics {
	//TODO implement me
	panic("implement me")
}

func (g groupProxy) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, data any) diag.Diagnostics {
	//TODO implement me
	panic("implement me")
}

func (g groupProxy) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, data any) diag.Diagnostics {
	//TODO implement me
	panic("implement me")
}

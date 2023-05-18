package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/travix/gotf/rsrc"
)

var _ rsrc.Resource = &userProxy{}

type userProxy struct{}

func (g userProxy) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, data any) diag.Diagnostics {
	//TODO implement me
	panic("implement me")
}

func (g userProxy) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, data any) diag.Diagnostics {
	//TODO implement me
	panic("implement me")
}

func (g userProxy) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, data any) diag.Diagnostics {
	//TODO implement me
	panic("implement me")
}

func (g userProxy) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, data any) diag.Diagnostics {
	//TODO implement me
	panic("implement me")
}

package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/travix/gotf/dtsrc"
)

var _ dtsrc.Datasource = &usersProxy{}

type usersProxy struct{}

func (u usersProxy) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse, data any) diag.Diagnostics {
	//TODO implement me
	panic("implement me")
}

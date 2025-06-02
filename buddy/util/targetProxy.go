package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type targetProxyModel struct {
	Name types.String `tfsdk:"name"`
	Host types.String `tfsdk:"host"`
	Port types.String `tfsdk:"port"`
	Auth types.Set    `tfsdk:"auth"`
}

func TargetProxyModelToApi(ctx context.Context, s *types.Set) (*buddy.TargetProxy, diag.Diagnostics) {
	var t []targetProxyModel
	diags := s.ElementsAs(ctx, &t, false)
	if len(t) == 0 {
		return nil, diags
	}
	if len(t) != 1 {
		diags.Append(diag.NewErrorDiagnostic("Wrong target proxy settings", "There should be only one target proxy entry"))
		return nil, diags
	}
	tt := t[0]
	var result buddy.TargetProxy
	result.Name = tt.Name.ValueString()
	if !tt.Port.IsNull() && !tt.Port.IsUnknown() {
		result.Port = tt.Port.ValueString()
	}
	if !tt.Host.IsNull() && !tt.Host.IsUnknown() {
		result.Host = tt.Host.ValueString()
	}
	if !tt.Auth.IsNull() && !tt.Auth.IsUnknown() {
		auth, d := TargetAuthModelToApi(ctx, &tt.Auth)
		diags.Append(d...)
		result.Auth = auth
	}
	return &result, diags
}

func TargetProxyModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Required: true,
		},
		"host": schema.StringAttribute{
			Optional: true,
		},
		"port": schema.StringAttribute{
			Optional: true,
		},
	}
}

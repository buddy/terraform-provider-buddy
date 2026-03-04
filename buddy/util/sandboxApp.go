package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type sandboxAppModel struct {
	Id      types.String `tfsdk:"id"`
	Command types.String `tfsdk:"command"`
	Status  types.String `tfsdk:"status"`
}

func sandboxAppModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"id":      types.StringType,
		"command": types.StringType,
		"status":  types.StringType,
	}
}

func (s *sandboxAppModel) loadAPI(app *buddy.SandboxApp) {
	s.Id = types.StringValue(app.Id)
	s.Status = types.StringValue(app.AppStatus)
	s.Command = types.StringValue(app.Command)
}

func SandboxAppsFromApi(ctx context.Context, apps *[]*buddy.SandboxApp) (basetypes.SetValue, diag.Diagnostics) {
	s := make([]*sandboxAppModel, len(*apps))
	for i, v := range *apps {
		s[i] = &sandboxAppModel{}
		s[i].loadAPI(v)
	}
	return types.SetValueFrom(ctx, types.ObjectType{AttrTypes: sandboxAppModelAttrs()}, &s)
}

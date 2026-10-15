package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type sandboxAppModel struct {
	Id         types.String `tfsdk:"id"`
	Identifier types.String `tfsdk:"identifier"`
	Command    types.String `tfsdk:"command"`
	Path       types.String `tfsdk:"path"`
	Status     types.String `tfsdk:"status"`
}

func sandboxAppModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"id":         types.StringType,
		"identifier": types.StringType,
		"command":    types.StringType,
		"path":       types.StringType,
		"status":     types.StringType,
	}
}

func (s *sandboxAppModel) loadAPI(app *buddy.SandboxApp) {
	s.Id = types.StringValue(app.Id)
	s.Identifier = types.StringValue(app.Identifier)
	s.Status = types.StringValue(app.AppStatus)
	s.Command = types.StringValue(app.Command)
	s.Path = types.StringValue(app.Path)
}

func SandboxAppsFromApi(ctx context.Context, apps *[]*buddy.SandboxApp) (basetypes.SetValue, diag.Diagnostics) {
	s := make([]*sandboxAppModel, len(*apps))
	for i, v := range *apps {
		s[i] = &sandboxAppModel{}
		s[i].loadAPI(v)
	}
	return types.SetValueFrom(ctx, types.ObjectType{AttrTypes: sandboxAppModelAttrs()}, &s)
}

type sandboxAppConfigModel struct {
	Command    types.String `tfsdk:"command"`
	Identifier types.String `tfsdk:"identifier"`
	Path       types.String `tfsdk:"path"`
}

func ResourceSandboxAppConfigModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"command": schema.StringAttribute{
			MarkdownDescription: "The app's run command",
			Required:            true,
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The app's identifier. Generated from the run command when not set",
			Optional:            true,
		},
		"path": schema.StringAttribute{
			MarkdownDescription: "The directory the app is run in. Defaults to the sandbox's `app_dir`",
			Optional:            true,
		},
	}
}

func SandboxAppCommandsToApi(ctx context.Context, s *types.Set) (*[]*buddy.SandboxAppOps, diag.Diagnostics) {
	commands, diags := StringSetToApi(ctx, s)
	apps := make([]*buddy.SandboxAppOps, len(*commands))
	for i := range *commands {
		apps[i] = &buddy.SandboxAppOps{
			Command: &(*commands)[i],
		}
	}
	return &apps, diags
}

func SandboxAppConfigsToApi(ctx context.Context, l *types.List) (*[]*buddy.SandboxAppOps, diag.Diagnostics) {
	var m []sandboxAppConfigModel
	diags := l.ElementsAs(ctx, &m, false)
	apps := make([]*buddy.SandboxAppOps, len(m))
	for i, v := range m {
		app := buddy.SandboxAppOps{
			Command: v.Command.ValueStringPointer(),
		}
		if !v.Identifier.IsNull() && !v.Identifier.IsUnknown() {
			app.Identifier = v.Identifier.ValueStringPointer()
		}
		if !v.Path.IsNull() && !v.Path.IsUnknown() {
			app.Path = v.Path.ValueStringPointer()
		}
		apps[i] = &app
	}
	return &apps, diags
}

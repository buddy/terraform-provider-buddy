package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type workspaceModel struct {
	HtmlUrl     types.String `tfsdk:"html_url"`
	WorkspaceId types.Int64  `tfsdk:"workspace_id"`
	Name        types.String `tfsdk:"name"`
	Domain      types.String `tfsdk:"domain"`
}

func workspaceModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"html_url":     types.StringType,
		"workspace_id": types.Int64Type,
		"name":         types.StringType,
		"domain":       types.StringType,
	}
}

func (w *workspaceModel) loadAPI(workspace *buddy.Workspace) {
	w.HtmlUrl = types.StringValue(workspace.HtmlUrl)
	w.WorkspaceId = types.Int64Value(int64(workspace.Id))
	w.Name = types.StringValue(workspace.Name)
	w.Domain = types.StringValue(workspace.Domain)
}

func SourceWorkspaceModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"html_url": schema.StringAttribute{
			Computed: true,
		},
		"workspace_id": schema.Int64Attribute{
			Computed: true,
		},
		"name": schema.StringAttribute{
			Computed: true,
		},
		"domain": schema.StringAttribute{
			Computed: true,
		},
	}
}

func WorkspacesModelFromApi(ctx context.Context, workspaces *[]*buddy.Workspace) (basetypes.SetValue, diag.Diagnostics) {
	r := make([]*workspaceModel, len(*workspaces))
	for i, v := range *workspaces {
		r[i] = &workspaceModel{}
		r[i].loadAPI(v)
	}
	return types.SetValueFrom(ctx, types.ObjectType{AttrTypes: workspaceModelAttrs()}, &r)
}

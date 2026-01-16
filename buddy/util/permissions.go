package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	sourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type permissionModel struct {
	HtmlUrl                types.String `tfsdk:"html_url"`
	PermissionId           types.Int64  `tfsdk:"permission_id"`
	Name                   types.String `tfsdk:"name"`
	Type                   types.String `tfsdk:"type"`
	PipelineAccessLevel    types.String `tfsdk:"pipeline_access_level"`
	ProjectTeamAccessLevel types.String `tfsdk:"project_team_access_level"`
	RepositoryAccessLevel  types.String `tfsdk:"repository_access_level"`
	SandboxAccessLevel     types.String `tfsdk:"sandbox_access_level"`
	TargetAccessLevel      types.String `tfsdk:"target_access_level"`
	EnvironmentAccessLevel types.String `tfsdk:"environment_access_level"`
}

func permissionModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"html_url":                  types.StringType,
		"permission_id":             types.Int64Type,
		"name":                      types.StringType,
		"type":                      types.StringType,
		"pipeline_access_level":     types.StringType,
		"project_team_access_level": types.StringType,
		"repository_access_level":   types.StringType,
		"sandbox_access_level":      types.StringType,
		"target_access_level":       types.StringType,
		"environment_access_level":  types.StringType,
	}
}

func (r *permissionModel) loadAPI(permission *buddy.Permission) {
	r.HtmlUrl = types.StringValue(permission.HtmlUrl)
	r.PermissionId = types.Int64Value(int64(permission.Id))
	r.Name = types.StringValue(permission.Name)
	r.Type = types.StringValue(permission.Type)
	r.PipelineAccessLevel = types.StringValue(permission.PipelineAccessLevel)
	r.ProjectTeamAccessLevel = types.StringValue(permission.ProjectTeamAccessLevel)
	r.RepositoryAccessLevel = types.StringValue(permission.RepositoryAccessLevel)
	r.SandboxAccessLevel = types.StringValue(permission.SandboxAccessLevel)
	r.EnvironmentAccessLevel = types.StringValue(permission.EnvironmentAccessLevel)
	r.TargetAccessLevel = types.StringValue(permission.TargetAccessLevel)
}

func SourcePermissionModelAttributes() map[string]sourceschema.Attribute {
	return map[string]sourceschema.Attribute{
		"html_url": sourceschema.StringAttribute{
			Computed: true,
		},
		"permission_id": sourceschema.Int64Attribute{
			Computed: true,
		},
		"name": sourceschema.StringAttribute{
			Computed: true,
		},
		"type": sourceschema.StringAttribute{
			Computed: true,
		},
		"pipeline_access_level": sourceschema.StringAttribute{
			Computed: true,
		},
		"project_team_access_level": sourceschema.StringAttribute{
			Computed: true,
		},
		"repository_access_level": sourceschema.StringAttribute{
			Computed: true,
		},
		"sandbox_access_level": sourceschema.StringAttribute{
			Computed: true,
		},
		"target_access_level": sourceschema.StringAttribute{
			Computed: true,
		},
		"environment_access_level": sourceschema.StringAttribute{
			Computed: true,
		},
	}
}

func PermissionModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"html_url": schema.StringAttribute{
			Computed: true,
		},
		"permission_id": schema.Int64Attribute{
			Computed: true,
		},
		"name": schema.StringAttribute{
			Computed: true,
		},
		"type": schema.StringAttribute{
			Computed: true,
		},
		"pipeline_access_level": schema.StringAttribute{
			Computed: true,
		},
		"project_team_access_level": schema.StringAttribute{
			Computed: true,
		},
		"repository_access_level": schema.StringAttribute{
			Computed: true,
		},
		"sandbox_access_level": schema.StringAttribute{
			Computed: true,
		},
		"target_access_level": schema.StringAttribute{
			Computed: true,
		},
		"environment_access_level": schema.StringAttribute{
			Computed: true,
		},
	}
}

func PermissionsModelFromApi(ctx context.Context, permissions *[]*buddy.Permission) (basetypes.SetValue, diag.Diagnostics) {
	p := make([]*permissionModel, len(*permissions))
	for i, v := range *permissions {
		p[i] = &permissionModel{}
		p[i].loadAPI(v)
	}
	return types.SetValueFrom(ctx, types.ObjectType{AttrTypes: permissionModelAttrs()}, &p)
}

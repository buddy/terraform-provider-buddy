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

type projectModel struct {
	HtmlUrl     types.String `tfsdk:"html_url"`
	Name        types.String `tfsdk:"name"`
	DisplayName types.String `tfsdk:"display_name"`
	Status      types.String `tfsdk:"status"`
}

func projectModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"html_url":     types.StringType,
		"name":         types.StringType,
		"display_name": types.StringType,
		"status":       types.StringType,
	}
}

func (r *projectModel) loadAPI(project *buddy.Project) {
	r.HtmlUrl = types.StringValue(project.HtmlUrl)
	r.Name = types.StringValue(project.Name)
	r.DisplayName = types.StringValue(project.DisplayName)
	r.Status = types.StringValue(project.Status)
}

func SourceProjectsModelAttributes() map[string]sourceschema.Attribute {
	return map[string]sourceschema.Attribute{
		"html_url": schema.StringAttribute{
			Computed: true,
		},
		"name": schema.StringAttribute{
			Computed: true,
		},
		"display_name": schema.StringAttribute{
			Computed: true,
		},
		"status": schema.StringAttribute{
			Computed: true,
		},
	}
}

func ProjectsModelFromApi(ctx context.Context, projects *[]*buddy.Project) (basetypes.SetValue, diag.Diagnostics) {
	p := make([]*projectModel, len(*projects))
	for i, v := range *projects {
		p[i] = &projectModel{}
		p[i].loadAPI(v)
	}
	return types.SetValueFrom(ctx, types.ObjectType{AttrTypes: projectModelAttrs()}, &p)
}

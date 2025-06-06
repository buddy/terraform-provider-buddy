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

type environmentModel struct {
	HtmlUrl       types.String `tfsdk:"html_url"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	Name          types.String `tfsdk:"name"`
	Identifier    types.String `tfsdk:"identifier"`
	Tags          types.Set    `tfsdk:"tags"`
	PublicUrl     types.String `tfsdk:"public_url"`
}

func environmentModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"html_url":       types.StringType,
		"environment_id": types.StringType,
		"name":           types.StringType,
		"identifier":     types.StringType,
		"tags":           types.SetType{ElemType: types.StringType},
		"public_url":     types.StringType,
	}
}

func (e *environmentModel) loadAPI(ctx context.Context, environment *buddy.Environment) diag.Diagnostics {
	e.HtmlUrl = types.StringValue(environment.HtmlUrl)
	e.EnvironmentId = types.StringValue(environment.Id)
	e.Name = types.StringValue(environment.Name)
	e.Identifier = types.StringValue(environment.Identifier)
	e.PublicUrl = types.StringValue(environment.PublicUrl)
	t, d := types.SetValueFrom(ctx, types.StringType, &environment.Tags)
	e.Tags = t
	return d
}

func SourceEnvironmentModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"html_url": schema.StringAttribute{
			Computed: true,
		},
		"environment_id": schema.StringAttribute{
			Computed: true,
		},
		"name": schema.StringAttribute{
			Computed: true,
		},
		"identifier": schema.StringAttribute{
			Computed: true,
		},
		"public_url": schema.StringAttribute{
			Computed: true,
		},
		"tags": schema.SetAttribute{
			Computed:    true,
			ElementType: types.StringType,
		},
	}
}

func EnvironmentsModelFromApi(ctx context.Context, environments *[]*buddy.Environment) (basetypes.SetValue, diag.Diagnostics) {
	l := make([]*environmentModel, len(*environments))
	for i, v := range *environments {
		l[i] = &environmentModel{}
		l[i].loadAPI(ctx, v)
	}
	return types.SetValueFrom(ctx, types.ObjectType{AttrTypes: environmentModelAttrs()}, &l)
}

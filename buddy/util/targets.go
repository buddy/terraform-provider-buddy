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

type TargetModel struct {
	HtmlUrl    types.String `tfsdk:"html_url"`
	TargetId   types.String `tfsdk:"target_id"`
	Identifier types.String `tfsdk:"identifier"`
	Name       types.String `tfsdk:"name"`
	Type       types.String `tfsdk:"type"`
	Tags       types.List   `tfsdk:"tags"`
	Host       types.String `tfsdk:"host"`
	Scope      types.String `tfsdk:"scope"`
	Repository types.String `tfsdk:"repository"`
	Port       types.String `tfsdk:"port"`
	Path       types.String `tfsdk:"path"`
	Secure     types.Bool   `tfsdk:"secure"`
	Disabled   types.Bool   `tfsdk:"disabled"`
}

func TargetModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"html_url":   types.StringType,
		"target_id":  types.StringType,
		"identifier": types.StringType,
		"name":       types.StringType,
		"type":       types.StringType,
		"tags":       types.ListType{ElemType: types.StringType},
		"host":       types.StringType,
		"scope":      types.StringType,
		"repository": types.StringType,
		"port":       types.StringType,
		"path":       types.StringType,
		"secure":     types.BoolType,
		"disabled":   types.BoolType,
	}
}

func (t *TargetModel) LoadAPI(ctx context.Context, target *buddy.Target) {
	t.HtmlUrl = types.StringValue(target.HtmlUrl)
	t.TargetId = types.StringValue(target.Id)
	t.Identifier = types.StringValue(target.Identifier)
	t.Name = types.StringValue(target.Name)
	t.Type = types.StringValue(target.Type)
	t.Scope = types.StringValue(target.Scope)

	if target.Tags != nil && len(target.Tags) > 0 {
		tags, _ := types.ListValueFrom(ctx, types.StringType, target.Tags)
		t.Tags = tags
	} else {
		t.Tags = types.ListNull(types.StringType)
	}

	t.Host = types.StringValue(target.Host)
	t.Repository = types.StringValue(target.Repository)
	t.Port = types.StringValue(target.Port)
	t.Path = types.StringValue(target.Path)
	t.Secure = types.BoolValue(target.Secure)
	t.Disabled = types.BoolValue(target.Disabled)
}

func SourceTargetModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"html_url": schema.StringAttribute{
			Computed: true,
		},
		"target_id": schema.StringAttribute{
			Computed: true,
		},
		"identifier": schema.StringAttribute{
			Computed: true,
		},
		"name": schema.StringAttribute{
			Computed: true,
		},
		"type": schema.StringAttribute{
			Computed: true,
		},
		"tags": schema.ListAttribute{
			Computed:    true,
			ElementType: types.StringType,
		},
		"host": schema.StringAttribute{
			Computed: true,
		},
		"scope": schema.StringAttribute{
			Computed: true,
		},
		"repository": schema.StringAttribute{
			Computed: true,
		},
		"port": schema.StringAttribute{
			Computed: true,
		},
		"path": schema.StringAttribute{
			Computed: true,
		},
		"secure": schema.BoolAttribute{
			Computed: true,
		},
		"disabled": schema.BoolAttribute{
			Computed: true,
		},
	}
}

func TargetsModelFromApi(ctx context.Context, targets *[]*buddy.Target) (basetypes.SetValue, diag.Diagnostics) {
	l := make([]*TargetModel, len(*targets))
	for i, v := range *targets {
		l[i] = &TargetModel{}
		l[i].LoadAPI(ctx, v)
	}
	return types.SetValueFrom(ctx, types.ObjectType{AttrTypes: TargetModelAttrs()}, &l)
}

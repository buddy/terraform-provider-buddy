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

type groupModel struct {
	HtmlUrl types.String `tfsdk:"html_url"`
	GroupId types.Int64  `tfsdk:"group_id"`
	Name    types.String `tfsdk:"name"`
}

func groupModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"html_url": types.StringType,
		"group_id": types.Int64Type,
		"name":     types.StringType,
	}
}

func (g *groupModel) loadAPI(group *buddy.Group) {
	g.HtmlUrl = types.StringValue(group.HtmlUrl)
	g.Name = types.StringValue(group.Name)
	g.GroupId = types.Int64Value(int64(group.Id))
}

func SourceGroupModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"html_url": schema.StringAttribute{
			Computed: true,
		},
		"group_id": schema.Int64Attribute{
			Computed: true,
		},
		"name": schema.StringAttribute{
			Computed: true,
		},
	}
}

func GroupsModelFromApi(ctx context.Context, groups *[]*buddy.Group) (basetypes.SetValue, diag.Diagnostics) {
	r := make([]*groupModel, len(*groups))
	for i, v := range *groups {
		r[i] = &groupModel{}
		r[i].loadAPI(v)
	}
	return types.SetValueFrom(ctx, types.ObjectType{AttrTypes: groupModelAttrs()}, &r)
}

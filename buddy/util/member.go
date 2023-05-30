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

type memberModel struct {
	HtmlUrl        types.String `tfsdk:"html_url"`
	MemberId       types.Int64  `tfsdk:"member_id"`
	Name           types.String `tfsdk:"name"`
	Email          types.String `tfsdk:"email"`
	Admin          types.Bool   `tfsdk:"admin"`
	WorkspaceOwner types.Bool   `tfsdk:"workspace_owner"`
	AvatarUrl      types.String `tfsdk:"avatar_url"`
}

func memberModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"html_url":        types.StringType,
		"member_id":       types.Int64Type,
		"name":            types.StringType,
		"email":           types.StringType,
		"admin":           types.BoolType,
		"workspace_owner": types.BoolType,
		"avatar_url":      types.StringType,
	}
}

func (r *memberModel) loadAPI(member *buddy.Member) {
	r.HtmlUrl = types.StringValue(member.HtmlUrl)
	r.MemberId = types.Int64Value(int64(member.Id))
	r.Name = types.StringValue(member.Name)
	r.Email = types.StringValue(member.Email)
	r.Admin = types.BoolValue(member.Admin)
	r.WorkspaceOwner = types.BoolValue(member.WorkspaceOwner)
	r.AvatarUrl = types.StringValue(member.AvatarUrl)
}

func MemberModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"html_url": schema.StringAttribute{
			Computed: true,
		},
		"member_id": schema.Int64Attribute{
			Computed: true,
		},
		"name": schema.StringAttribute{
			Computed: true,
		},
		"email": schema.StringAttribute{
			Computed: true,
		},
		"admin": schema.BoolAttribute{
			Computed: true,
		},
		"workspace_owner": schema.BoolAttribute{
			Computed: true,
		},
		"avatar_url": schema.StringAttribute{
			Computed: true,
		},
	}
}

func MembersModelFromApi(ctx context.Context, members *[]*buddy.Member) (basetypes.SetValue, diag.Diagnostics) {
	r := make([]*memberModel, len(*members))
	for i, v := range *members {
		r[i] = &memberModel{}
		r[i].loadAPI(v)
	}
	return types.SetValueFrom(ctx, types.ObjectType{AttrTypes: memberModelAttrs()}, &r)
}

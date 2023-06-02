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

type memberModel struct {
	HtmlUrl        types.String `tfsdk:"html_url"`
	MemberId       types.Int64  `tfsdk:"member_id"`
	Name           types.String `tfsdk:"name"`
	Email          types.String `tfsdk:"email"`
	Admin          types.Bool   `tfsdk:"admin"`
	WorkspaceOwner types.Bool   `tfsdk:"workspace_owner"`
	AvatarUrl      types.String `tfsdk:"avatar_url"`
	Status         types.String `tfsdk:"status"`
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
		"status":          types.StringType,
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
	r.Status = types.StringValue(member.Status)
}

func SourceMemberModelAttributes() map[string]sourceschema.Attribute {
	return map[string]sourceschema.Attribute{
		"html_url": schema.StringAttribute{
			Computed: true,
		},
		"member_id": sourceschema.Int64Attribute{
			Computed: true,
		},
		"name": sourceschema.StringAttribute{
			Computed: true,
		},
		"email": sourceschema.StringAttribute{
			Computed: true,
		},
		"admin": sourceschema.BoolAttribute{
			Computed: true,
		},
		"status": sourceschema.StringAttribute{
			Computed: true,
		},
		"workspace_owner": sourceschema.BoolAttribute{
			Computed: true,
		},
		"avatar_url": sourceschema.StringAttribute{
			Computed: true,
		},
	}
}

func ResourceMemberModelAttributes() map[string]schema.Attribute {
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
		"status": schema.StringAttribute{
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

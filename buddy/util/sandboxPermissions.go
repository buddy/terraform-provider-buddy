package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type sandboxPermissionsAccessModel struct {
	Id          types.Int64  `tfsdk:"id"`
	AccessLevel types.String `tfsdk:"access_level"`
}

type sandboxPermissionsModel struct {
	Others types.String `tfsdk:"others"`
	Users  types.Set    `tfsdk:"user"`
	Groups types.Set    `tfsdk:"group"`
}

func sandboxPermissionsAccessModelToApi(ctx context.Context, s *types.Set) (*[]*buddy.SandboxResourcePermission, diag.Diagnostics) {
	var spa []sandboxPermissionsAccessModel
	diags := s.ElementsAs(ctx, &spa, false)
	sandboxResourcePermissions := make([]*buddy.SandboxResourcePermission, len(spa))
	for i, v := range spa {
		srp := &buddy.SandboxResourcePermission{}
		if !v.Id.IsNull() && !v.Id.IsUnknown() {
			srp.Id = int(v.Id.ValueInt64())
		}
		if !v.AccessLevel.IsNull() && !v.AccessLevel.IsUnknown() {
			srp.AccessLevel = v.AccessLevel.ValueString()
		}
		sandboxResourcePermissions[i] = srp
	}
	return &sandboxResourcePermissions, diags
}

func SandboxPermissionsModelToApi(ctx context.Context, s *types.Set) (*buddy.SandboxPermissions, diag.Diagnostics) {
	var p []sandboxPermissionsModel
	diags := s.ElementsAs(ctx, &p, false)
	if len(p) == 0 {
		return nil, diags
	}
	if len(p) != 1 {
		diags.Append(diag.NewErrorDiagnostic("Wrong sandbox permissions settings", "There should be only one sandbox permissions entry"))
		return nil, diags
	}
	pp := p[0]
	var result buddy.SandboxPermissions
	result.Others = pp.Others.ValueString()
	if result.Others == "" {
		result.Others = buddy.SandboxPermissionDefault
	}
	if !pp.Users.IsNull() && !pp.Users.IsUnknown() {
		users, d := sandboxPermissionsAccessModelToApi(ctx, &pp.Users)
		diags.Append(d...)
		result.Users = *users
	} else {
		result.Users = []*buddy.SandboxResourcePermission{}
	}
	if !pp.Groups.IsNull() && !pp.Groups.IsUnknown() {
		groups, d := sandboxPermissionsAccessModelToApi(ctx, &pp.Groups)
		diags.Append(d...)
		result.Groups = *groups
	} else {
		result.Groups = []*buddy.SandboxResourcePermission{}
	}
	return &result, diags
}

func SandboxPermissionsAccessModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Required: true,
		},
		"access_level": schema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					buddy.SandboxPermissionDefault,
					buddy.SandboxPermissionDenied,
					buddy.SandboxPermissionReadOnly,
					buddy.SandboxPermissionManage,
				),
			},
		},
	}
}

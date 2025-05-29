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

type targetPermissionsAccessModel struct {
	Id          types.Int64  `tfsdk:"id"`
	AccessLevel types.String `tfsdk:"access_level"`
}

type targetPermissionsModel struct {
	Others types.String `tfsdk:"others"`
	Users  types.Set    `tfsdk:"user"`
	Groups types.Set    `tfsdk:"group"`
}

func targetPermissionsAccessModelToApi(ctx context.Context, s *types.Set) (*[]*buddy.TargetResourcePermission, diag.Diagnostics) {
	var tpam []targetPermissionsAccessModel
	diags := s.ElementsAs(ctx, &tpam, false)
	targetResourcePermissions := make([]*buddy.TargetResourcePermission, len(tpam))
	for i, v := range tpam {
		trp := &buddy.TargetResourcePermission{}
		if !v.Id.IsNull() && !v.Id.IsUnknown() {
			trp.Id = int(v.Id.ValueInt64())
		}
		if !v.AccessLevel.IsNull() && !v.AccessLevel.IsUnknown() {
			trp.AccessLevel = v.AccessLevel.ValueString()
		}
		targetResourcePermissions[i] = trp
	}
	return &targetResourcePermissions, diags
}

func TargetPermissionsModelToApi(ctx context.Context, s *types.Set) (*buddy.TargetPermissions, diag.Diagnostics) {
	var p []targetPermissionsModel
	diags := s.ElementsAs(ctx, &p, false)
	if len(p) == 0 {
		return nil, diags
	}
	if len(p) != 1 {
		diags.Append(diag.NewErrorDiagnostic("Wrong target permissions settings", "There should be only one target permissions entry"))
		return nil, diags
	}
	pp := p[0]
	var result buddy.TargetPermissions
	result.Others = pp.Others.ValueString()
	if result.Others == "" {
		result.Others = buddy.TargetPermissionUseOnly
	}
	if !pp.Users.IsNull() && !pp.Users.IsUnknown() {
		users, d := targetPermissionsAccessModelToApi(ctx, &pp.Users)
		diags.Append(d...)
		result.Users = *users
	} else {
		result.Users = []*buddy.TargetResourcePermission{}
	}
	if !pp.Groups.IsNull() && !pp.Groups.IsUnknown() {
		groups, d := targetPermissionsAccessModelToApi(ctx, &pp.Groups)
		diags.Append(d...)
		result.Groups = *groups
	} else {
		result.Groups = []*buddy.TargetResourcePermission{}
	}
	return &result, diags
}

func TargetPermissionsAccessModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Required: true,
		},
		"access_level": schema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					buddy.TargetPermissionManage,
					buddy.TargetPermissionUseOnly,
				),
			},
		},
	}
}

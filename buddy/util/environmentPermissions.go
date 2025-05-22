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

type environmentPermissionsAccessModel struct {
	Id          types.Int64  `tfsdk:"id"`
	AccessLevel types.String `tfsdk:"access_level"`
}

type environmentPermissionsModel struct {
	Others types.String `tfsdk:"others"`
	Users  types.Set    `tfsdk:"user"`
	Groups types.Set    `tfsdk:"group"`
}

func environmentPermissionsAccessModelToApi(ctx context.Context, s *types.Set) (*[]*buddy.EnvironmentResourcePermissions, diag.Diagnostics) {
	var epam []environmentPermissionsAccessModel
	diags := s.ElementsAs(ctx, &epam, false)
	environmentResourcePermissions := make([]*buddy.EnvironmentResourcePermissions, len(epam))
	for i, v := range epam {
		erp := &buddy.EnvironmentResourcePermissions{}
		if !v.Id.IsNull() && !v.Id.IsUnknown() {
			erp.Id = int(v.Id.ValueInt64())
		}
		if !v.AccessLevel.IsNull() && !v.AccessLevel.IsUnknown() {
			erp.AccessLevel = v.AccessLevel.ValueString()
		}
		environmentResourcePermissions[i] = erp
	}
	return &environmentResourcePermissions, diags
}

func EnvironmentPermissionsModelToApi(ctx context.Context, s *types.Set) (*buddy.EnvironmentPermissions, diag.Diagnostics) {
	var p []environmentPermissionsModel
	diags := s.ElementsAs(ctx, &p, false)
	if len(p) == 0 {
		return nil, diags
	}
	if len(p) != 1 {
		diags.Append(diag.NewErrorDiagnostic("Wrong environment permissions settings", "There should be only one environment permissions entry"))
		return nil, diags
	}
	pp := p[0]
	var result buddy.EnvironmentPermissions
	result.Others = pp.Others.ValueString()
	if result.Others == "" {
		result.Others = buddy.EnvironmentPermissionAccessLevelUseOnly
	}
	if !pp.Users.IsNull() && !pp.Users.IsUnknown() {
		users, d := environmentPermissionsAccessModelToApi(ctx, &pp.Users)
		diags.Append(d...)
		result.Users = *users
	} else {
		result.Users = []*buddy.EnvironmentResourcePermissions{}
	}
	if !pp.Groups.IsNull() && !pp.Groups.IsUnknown() {
		groups, d := environmentPermissionsAccessModelToApi(ctx, &pp.Groups)
		diags.Append(d...)
		result.Groups = *groups
	} else {
		result.Groups = []*buddy.EnvironmentResourcePermissions{}
	}
	return &result, diags
}

func EnvironmentPermissionsAccessModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Required: true,
		},
		"access_level": schema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					buddy.EnvironmentPermissionAccessLevelManage,
					buddy.EnvironmentPermissionAccessLevelUseOnly,
					buddy.EnvironmentPermissionAccessLevelDefault,
					buddy.EnvironmentPermissionAccessLevelDenied,
				),
			},
		},
	}
}

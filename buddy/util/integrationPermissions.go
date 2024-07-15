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

type integrationPermissionsAccessModel struct {
	Id          types.Int64  `tfsdk:"id"`
	AccessLevel types.String `tfsdk:"access_level"`
}

type integrationPermissionsModel struct {
	Admins types.String `tfsdk:"admins"`
	Others types.String `tfsdk:"others"`
	Users  types.Set    `tfsdk:"user"`
	Groups types.Set    `tfsdk:"group"`
}

func integrationPermissionsAccessModelToApi(ctx context.Context, s *types.Set) (*[]*buddy.IntegrationResourcePermission, diag.Diagnostics) {
	var ipam []integrationPermissionsAccessModel
	diags := s.ElementsAs(ctx, &ipam, false)
	integrationResourcePermissions := make([]*buddy.IntegrationResourcePermission, len(ipam))
	for i, v := range ipam {
		irp := &buddy.IntegrationResourcePermission{}
		if !v.Id.IsNull() && !v.Id.IsUnknown() {
			irp.Id = int(v.Id.ValueInt64())
		}
		if !v.AccessLevel.IsNull() && !v.AccessLevel.IsUnknown() {
			irp.AccessLevel = v.AccessLevel.ValueString()
		}
		integrationResourcePermissions[i] = irp
	}
	return &integrationResourcePermissions, diags
}

func IntegrationPermissionsModelToApi(ctx context.Context, s *types.Set) (*buddy.IntegrationPermissions, diag.Diagnostics) {
	var p []integrationPermissionsModel
	diags := s.ElementsAs(ctx, &p, false)
	if len(p) == 0 {
		return nil, diags
	}
	if len(p) != 1 {
		diags.Append(diag.NewErrorDiagnostic("Wrong integration permissions settings", "There should be only one integration permissions entry"))
		return nil, diags
	}
	pp := p[0]
	var result buddy.IntegrationPermissions
	result.Others = pp.Others.ValueString()
	if result.Others == "" {
		result.Others = buddy.IntegrationPermissionDenied
	}
	result.Admins = pp.Admins.ValueString()
	if result.Admins == "" {
		result.Admins = buddy.IntegrationPermissionManage
	}
	if !pp.Users.IsNull() && !pp.Users.IsUnknown() {
		users, d := integrationPermissionsAccessModelToApi(ctx, &pp.Users)
		diags.Append(d...)
		result.Users = *users
	} else {
		result.Users = []*buddy.IntegrationResourcePermission{}
	}
	if !pp.Groups.IsNull() && !pp.Groups.IsUnknown() {
		groups, d := integrationPermissionsAccessModelToApi(ctx, &pp.Groups)
		diags.Append(d...)
		result.Groups = *groups
	} else {
		result.Groups = []*buddy.IntegrationResourcePermission{}
	}
	return &result, diags
}

func IntegrationPermissionsAccessModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Required: true,
		},
		"access_level": schema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					buddy.IntegrationPermissionManage,
					buddy.IntegrationPermissionUseOnly,
					buddy.IntegrationPermissionDenied,
				),
			},
		},
	}
}

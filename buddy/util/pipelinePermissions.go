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

type pipelinePermissionsAccessModel struct {
	Id          types.Int64  `tfsdk:"id"`
	AccessLevel types.String `tfsdk:"access_level"`
}

type pipelinePermissionsModel struct {
	Others types.String `tfsdk:"others"`
	Users  types.Set    `tfsdk:"user"`
	Groups types.Set    `tfsdk:"group"`
}

func pipelinePermissionsAccessModelToApi(ctx context.Context, s *types.Set) (*[]*buddy.PipelineResourcePermission, diag.Diagnostics) {
	var ppam []pipelinePermissionsAccessModel
	diags := s.ElementsAs(ctx, &ppam, false)
	pipelineResourcePermissions := make([]*buddy.PipelineResourcePermission, len(ppam))
	for i, v := range ppam {
		prp := &buddy.PipelineResourcePermission{}
		if !v.Id.IsNull() && !v.Id.IsUnknown() {
			prp.Id = int(v.Id.ValueInt64())
		}
		if !v.AccessLevel.IsNull() && !v.AccessLevel.IsUnknown() {
			prp.AccessLevel = v.AccessLevel.ValueString()
		}
		pipelineResourcePermissions[i] = prp
	}
	return &pipelineResourcePermissions, diags
}

func PipelinePermissionsModelToApi(ctx context.Context, s *types.Set) (*buddy.PipelinePermissions, diag.Diagnostics) {
	var p []pipelinePermissionsModel
	diags := s.ElementsAs(ctx, &p, false)
	if len(p) == 0 {
		return nil, diags
	}
	if len(p) != 1 {
		diags.Append(diag.NewErrorDiagnostic("Wrong pipeline permissions settings", "There should be only one pipeline permissions entry"))
		return nil, diags
	}
	pp := p[0]
	var result buddy.PipelinePermissions
	result.Others = pp.Others.ValueString()
	if result.Others == "" {
		result.Others = buddy.PipelinePermissionDefault
	}
	if !pp.Users.IsNull() && !pp.Users.IsUnknown() {
		users, d := pipelinePermissionsAccessModelToApi(ctx, &pp.Users)
		diags.Append(d...)
		result.Users = *users
	} else {
		result.Users = []*buddy.PipelineResourcePermission{}
	}
	if !pp.Groups.IsNull() && !pp.Groups.IsUnknown() {
		groups, d := pipelinePermissionsAccessModelToApi(ctx, &pp.Groups)
		diags.Append(d...)
		result.Groups = *groups
	} else {
		result.Groups = []*buddy.PipelineResourcePermission{}
	}
	return &result, diags
}

func PipelinePermissionsAccessModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Required: true,
		},
		"access_level": schema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					buddy.PipelinePermissionDefault,
					buddy.PipelinePermissionDenied,
					buddy.PipelinePermissionReadOnly,
					buddy.PipelinePermissionRunOnly,
					buddy.PipelinePermissionReadWrite,
				),
			},
		},
	}
}

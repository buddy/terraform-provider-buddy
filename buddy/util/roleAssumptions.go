package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type roleAssumptionModel struct {
	Arn        types.String `tfsdk:"arn"`
	ExternalId types.String `tfsdk:"external_id"`
	Duration   types.Int64  `tfsdk:"duration"`
}

func RoleAssumptionsModelToAPI(ctx context.Context, s *types.List) (*[]*buddy.RoleAssumption, diag.Diagnostics) {
	var ram []roleAssumptionModel
	diags := s.ElementsAs(ctx, &ram, false)
	roleAssumptions := make([]*buddy.RoleAssumption, len(ram))
	for i, v := range ram {
		ra := &buddy.RoleAssumption{}
		if !v.Arn.IsNull() && !v.Arn.IsUnknown() {
			ra.Arn = v.Arn.ValueString()
		}
		if !v.ExternalId.IsNull() && !v.ExternalId.IsUnknown() {
			ra.ExternalId = v.ExternalId.ValueString()
		}
		if !v.Duration.IsNull() && !v.Duration.IsUnknown() {
			ra.Duration = int(v.Duration.ValueInt64())
		}
		roleAssumptions[i] = ra
	}
	return &roleAssumptions, diags
}

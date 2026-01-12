package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type environmentEnvironmentModel struct {
	Project     types.String `tfsdk:"project"`
	Environment types.String `tfsdk:"environment"`
}

func EnvironmentEnvironmentsModelToApi(ctx context.Context, s *types.Set) (*[]*buddy.EnvironmentAllowedEnvironment, diag.Diagnostics) {
	var ee []environmentEnvironmentModel
	diags := s.ElementsAs(ctx, &ee, false)
	result := make([]*buddy.EnvironmentAllowedEnvironment, len(ee))
	for i, e := range ee {
		env := buddy.EnvironmentAllowedEnvironment{}
		if !e.Project.IsNull() && !e.Project.IsUnknown() {
			env.Project = e.Project.ValueString()
		}
		if !e.Environment.IsNull() && !e.Environment.IsUnknown() {
			env.Environment = e.Environment.ValueString()
		}
		result[i] = &env
	}
	return &result, diags
}

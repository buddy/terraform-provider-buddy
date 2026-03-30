package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type targetSandboxModel struct {
	Project     types.String `tfsdk:"project"`
	Sandbox     types.String `tfsdk:"sandbox"`
	AccessLevel types.String `tfsdk:"access_level"`
}

func TargetSandboxesModelToApi(ctx context.Context, s *types.Set) (*[]*buddy.TargetAllowedSandbox, diag.Diagnostics) {
	var ss []targetSandboxModel
	diags := s.ElementsAs(ctx, &ss, false)
	result := make([]*buddy.TargetAllowedSandbox, len(ss))
	for i, s := range ss {
		sandbox := buddy.TargetAllowedSandbox{}
		if !s.Project.IsNull() && !s.Project.IsUnknown() {
			sandbox.Project = s.Project.ValueString()
		}
		if !s.Sandbox.IsNull() && !s.Sandbox.IsUnknown() {
			sandbox.Sandbox = s.Sandbox.ValueString()
		}
		if !s.AccessLevel.IsNull() && !s.AccessLevel.IsUnknown() {
			sandbox.AccessLevel = s.AccessLevel.ValueString()
		}
		result[i] = &sandbox
	}
	return &result, diags
}

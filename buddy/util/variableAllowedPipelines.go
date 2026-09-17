package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type variableAllowedPipelineModel struct {
	Project     types.String `tfsdk:"project"`
	Pipeline    types.String `tfsdk:"pipeline"`
	Action      types.String `tfsdk:"action"`
	AccessLevel types.String `tfsdk:"access_level"`
}

func VariableAllowedPipelinesModelToApi(ctx context.Context, s *types.Set) (*[]*buddy.VariableAllowedPipeline, diag.Diagnostics) {
	var pp []variableAllowedPipelineModel
	diags := s.ElementsAs(ctx, &pp, false)
	result := make([]*buddy.VariableAllowedPipeline, len(pp))
	for i, p := range pp {
		pipeline := buddy.VariableAllowedPipeline{}
		if !p.Project.IsNull() && !p.Project.IsUnknown() {
			pipeline.Project = p.Project.ValueString()
		}
		if !p.Pipeline.IsNull() && !p.Pipeline.IsUnknown() {
			pipeline.Pipeline = p.Pipeline.ValueString()
		}
		if !p.Action.IsNull() && !p.Action.IsUnknown() {
			pipeline.Action = p.Action.ValueString()
		}
		if !p.AccessLevel.IsNull() && !p.AccessLevel.IsUnknown() {
			pipeline.AccessLevel = p.AccessLevel.ValueString()
		}
		result[i] = &pipeline
	}
	return &result, diags
}

package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type environmentPipelineModel struct {
	Project     types.String `tfsdk:"project"`
	Pipeline    types.String `tfsdk:"pipeline"`
	AccessLevel types.String `tfsdk:"access_level"`
}

func EnvironmentPipelinesModelToApi(ctx context.Context, s *types.Set) (*[]*buddy.EnvironmentAllowedPipeline, diag.Diagnostics) {
	var pp []environmentPipelineModel
	diags := s.ElementsAs(ctx, &pp, false)
	result := make([]*buddy.EnvironmentAllowedPipeline, len(pp))
	for i, p := range pp {
		pipeline := buddy.EnvironmentAllowedPipeline{}
		if !p.Project.IsNull() && !p.Project.IsUnknown() {
			pipeline.Project = p.Project.ValueString()
		}
		if !p.Pipeline.IsNull() && !p.Pipeline.IsUnknown() {
			pipeline.Pipeline = p.Pipeline.ValueString()
		}
		if !p.AccessLevel.IsNull() && !p.AccessLevel.IsUnknown() {
			pipeline.AccessLevel = p.AccessLevel.ValueString()
		}
		result[i] = &pipeline
	}
	return &result, diags
}

package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	sourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type sandboxPipelineModel struct {
	Project     types.String `tfsdk:"project"`
	Pipeline    types.String `tfsdk:"pipeline"`
	AccessLevel types.String `tfsdk:"access_level"`
}

func sandboxPipelineModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"project":      types.StringType,
		"pipeline":     types.StringType,
		"access_level": types.StringType,
	}
}

func SandboxPipelinesModelToApi(ctx context.Context, s *types.Set) (*[]*buddy.SandboxAllowedPipeline, diag.Diagnostics) {
	var pp []sandboxPipelineModel
	diags := s.ElementsAs(ctx, &pp, false)
	result := make([]*buddy.SandboxAllowedPipeline, len(pp))
	for i, p := range pp {
		pipeline := buddy.SandboxAllowedPipeline{}
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

func SandboxPipelinesFromApi(ctx context.Context, pipelines *[]*buddy.SandboxAllowedPipeline) (basetypes.SetValue, diag.Diagnostics) {
	r := make([]*sandboxPipelineModel, len(*pipelines))
	for i, p := range *pipelines {
		r[i] = &sandboxPipelineModel{
			Project:     types.StringValue(p.Project),
			Pipeline:    types.StringValue(p.Pipeline),
			AccessLevel: types.StringValue(p.AccessLevel),
		}
	}
	return types.SetValueFrom(ctx, types.ObjectType{AttrTypes: sandboxPipelineModelAttrs()}, &r)
}

func ResourceSandboxPipelineModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"project": schema.StringAttribute{
			MarkdownDescription: "The pipeline's project name",
			Required:            true,
		},
		"pipeline": schema.StringAttribute{
			MarkdownDescription: "The pipeline's identifier",
			Required:            true,
		},
		"access_level": schema.StringAttribute{
			MarkdownDescription: "The pipeline's access level. Allowed: `USE_ONLY`, `DENIED`",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.OneOf(buddy.SandboxPipelineAccessLevelUseOnly, buddy.SandboxPipelineAccessLevelDenied),
			},
		},
	}
}

func SourceSandboxPipelineModelAttributes() map[string]sourceschema.Attribute {
	return map[string]sourceschema.Attribute{
		"project": sourceschema.StringAttribute{
			MarkdownDescription: "The pipeline's project name",
			Computed:            true,
		},
		"pipeline": sourceschema.StringAttribute{
			MarkdownDescription: "The pipeline's identifier",
			Computed:            true,
		},
		"access_level": sourceschema.StringAttribute{
			MarkdownDescription: "The pipeline's access level",
			Computed:            true,
		},
	}
}

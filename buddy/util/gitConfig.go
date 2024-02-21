package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	sourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type GitConfigModel struct {
	Project types.String `tfsdk:"project"`
	Branch  types.String `tfsdk:"branch"`
	Path    types.String `tfsdk:"path"`
}

func GitConfigModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"project": types.StringType,
		"branch":  types.StringType,
		"path":    types.StringType,
	}
}

func (gc *GitConfigModel) loadAPI(_ context.Context, gitConfig *buddy.PipelineGitConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	gc.Project = types.StringValue(gitConfig.Project)
	gc.Branch = types.StringValue(gitConfig.Branch)
	gc.Path = types.StringValue(gitConfig.Path)
	return diags
}

func SourceGitConfigModelAttributes() map[string]sourceschema.Attribute {
	return map[string]sourceschema.Attribute{
		"project": sourceschema.StringAttribute{
			Computed: true,
		},
		"branch": sourceschema.StringAttribute{
			Computed: true,
		},
		"path": sourceschema.StringAttribute{
			Computed: true,
		},
	}
}

func GitConfigModelFromApi(ctx context.Context, gitConfig *buddy.PipelineGitConfig) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	gc := &GitConfigModel{}
	gc.loadAPI(ctx, gitConfig)
	r, d := types.ObjectValueFrom(ctx, GitConfigModelAttrs(), &gc)
	diags.Append(d...)
	return r, diags
}

func GitConfigModelToApi(ctx context.Context, gc *types.Object) (*buddy.PipelineGitConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	if gc == nil {
		return nil, diags
	}
	var gitConfig buddy.PipelineGitConfig
	d := gc.As(ctx, &gitConfig, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})
	diags.Append(d...)
	return &gitConfig, diags
}

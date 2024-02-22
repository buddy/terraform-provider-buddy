package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

func GitConfigModelFromApi(ctx context.Context, gitConfig *buddy.PipelineGitConfig) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	var gc GitConfigModel
	if gitConfig != nil {
		gc = GitConfigModel{}
		gc.loadAPI(ctx, gitConfig)
	}
	r, d := types.ObjectValueFrom(ctx, GitConfigModelAttrs(), &gc)
	diags.Append(d...)
	return r, diags
}

func GitConfigModelToApi(ctx context.Context, gc *types.Object) (*buddy.PipelineGitConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	if gc == nil || gc.IsNull() || gc.IsUnknown() {
		return nil, diags
	}
	var gitConfig buddy.PipelineGitConfig
	gcm := GitConfigModel{}
	d := gc.As(ctx, &gcm, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})
	diags.Append(d...)
	gitConfig.Project = gcm.Project.ValueString()
	gitConfig.Branch = gcm.Branch.ValueString()
	gitConfig.Path = gcm.Path.ValueString()
	return &gitConfig, diags
}

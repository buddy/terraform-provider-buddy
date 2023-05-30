package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type remoteParameterModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func RemoteParametersModelToApi(ctx context.Context, s *types.Set) (*[]*buddy.PipelineRemoteParameter, diag.Diagnostics) {
	var rpm []remoteParameterModel
	diags := s.ElementsAs(ctx, &rpm, false)
	remoteParams := make([]*buddy.PipelineRemoteParameter, len(rpm))
	for i, v := range rpm {
		remoteParam := &buddy.PipelineRemoteParameter{}
		if !v.Key.IsNull() && !v.Key.IsUnknown() {
			remoteParam.Key = v.Key.ValueString()
		}
		if !v.Value.IsNull() && !v.Value.IsUnknown() {
			remoteParam.Value = v.Value.ValueString()
		}
		remoteParams[i] = remoteParam
	}
	return &remoteParams, diags
}

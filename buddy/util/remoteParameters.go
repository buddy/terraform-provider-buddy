package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	sourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type remoteParameterModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func remoteParameterModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	}
}

func (r *remoteParameterModel) loadAPI(remoteParam *buddy.PipelineRemoteParameter) {
	r.Key = types.StringValue(remoteParam.Key)
	r.Value = types.StringValue(remoteParam.Value)
}

func SourceRemoteParameterModelAttributes() map[string]sourceschema.Attribute {
	return map[string]sourceschema.Attribute{
		"key": sourceschema.StringAttribute{
			Computed: true,
		},
		"value": sourceschema.StringAttribute{
			Computed: true,
		},
	}
}

func ResourceRemoteParameterModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"key": schema.StringAttribute{
			Required: true,
		},
		"value": schema.StringAttribute{
			Required: true,
		},
	}
}

func RemoteParametersModelFromApi(ctx context.Context, remoteParams *[]*buddy.PipelineRemoteParameter) (basetypes.SetValue, diag.Diagnostics) {
	r := make([]*remoteParameterModel, len(*remoteParams))
	for i, v := range *remoteParams {
		r[i] = &remoteParameterModel{}
		r[i].loadAPI(v)
	}
	return types.SetValueFrom(ctx, types.ObjectType{AttrTypes: remoteParameterModelAttrs()}, &r)
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

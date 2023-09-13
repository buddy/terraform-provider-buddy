package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type triggerConditionModel struct {
	Condition     types.String `tfsdk:"condition"`
	Paths         types.Set    `tfsdk:"paths"`
	VariableKey   types.String `tfsdk:"variable_key"`
	VariableValue types.String `tfsdk:"variable_value"`
	Hours         types.Set    `tfsdk:"hours"`
	Days          types.Set    `tfsdk:"days"`
	ZoneId        types.String `tfsdk:"zone_id"`
	ProjectName   types.String `tfsdk:"project_name"`
	PipelineName  types.String `tfsdk:"pipeline_name"`
	TriggerUser   types.String `tfsdk:"trigger_user"`
	TriggerGroup  types.String `tfsdk:"trigger_group"`
}

func TriggerConditionsModelToApi(ctx context.Context, s *types.Set) (*[]*buddy.PipelineTriggerCondition, diag.Diagnostics) {
	var tcm []triggerConditionModel
	diags := s.ElementsAs(ctx, &tcm, false)
	triggerConditions := make([]*buddy.PipelineTriggerCondition, len(tcm))
	for i, v := range tcm {
		tc := &buddy.PipelineTriggerCondition{}
		if !v.Days.IsNull() && !v.Days.IsUnknown() {
			days, d := Int64SetToApi(ctx, &v.Days)
			diags.Append(d...)
			tc.TriggerDays = *days
		}
		if !v.Hours.IsNull() && !v.Hours.IsUnknown() {
			hours, d := Int64SetToApi(ctx, &v.Hours)
			diags.Append(d...)
			tc.TriggerHours = *hours
		}
		if !v.PipelineName.IsNull() && !v.PipelineName.IsUnknown() {
			tc.TriggerPipelineName = v.PipelineName.ValueString()
		}
		if !v.ProjectName.IsNull() && !v.ProjectName.IsUnknown() {
			tc.TriggerProjectName = v.ProjectName.ValueString()
		}
		if !v.ZoneId.IsNull() && !v.ZoneId.IsUnknown() {
			tc.ZoneId = v.ZoneId.ValueString()
		}
		if !v.VariableValue.IsNull() && !v.VariableValue.IsUnknown() {
			tc.TriggerVariableValue = v.VariableValue.ValueString()
		}
		if !v.VariableKey.IsNull() && !v.VariableKey.IsUnknown() {
			tc.TriggerVariableKey = v.VariableKey.ValueString()
		}
		if !v.Condition.IsNull() && !v.Condition.IsUnknown() {
			tc.TriggerCondition = v.Condition.ValueString()
		}
		if !v.Paths.IsNull() && !v.Paths.IsUnknown() {
			paths, d := StringSetToApi(ctx, &v.Paths)
			diags.Append(d...)
			tc.TriggerConditionPaths = *paths
		}
		if !v.TriggerGroup.IsNull() && !v.TriggerGroup.IsUnknown() {
			tc.TriggerGroup = v.TriggerGroup.ValueString()
		}
		if !v.TriggerUser.IsNull() && !v.TriggerUser.IsUnknown() {
			tc.TriggerUser = v.TriggerUser.ValueString()
		}
		triggerConditions[i] = tc
	}
	return &triggerConditions, diags
}

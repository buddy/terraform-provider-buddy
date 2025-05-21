package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type environmentVariableModel struct {
	Key         types.String `tfsdk:"key"`
	Value       types.String `tfsdk:"value"`
	Description types.String `tfsdk:"description"`
	Settable    types.Bool   `tfsdk:"settable"`
	Encrypted   types.Bool   `tfsdk:"encrypted"`
}

func EnvironmentVariableModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"key": schema.StringAttribute{
			Required: true,
		},
		"value": schema.StringAttribute{
			Required:  true,
			Sensitive: true,
		},
		"description": schema.StringAttribute{
			Optional: true,
		},
		"settable": schema.BoolAttribute{
			Optional: true,
		},
		"encrypted": schema.BoolAttribute{
			Optional: true,
		},
	}
}

func EnvironmentVariableModelToApi(ctx context.Context, s *types.Set) (*[]*buddy.Variable, diag.Diagnostics) {
	var evm []environmentVariableModel
	diags := s.ElementsAs(ctx, &evm, false)
	variables := make([]*buddy.Variable, len(evm))
	for i, v := range evm {
		variable := &buddy.Variable{
			Type:      buddy.VariableTypeVar,
			FilePlace: buddy.VariableSshKeyFilePlaceNone,
		}
		if !v.Key.IsNull() && !v.Key.IsUnknown() {
			variable.Key = v.Key.ValueString()
		}
		if !v.Value.IsNull() && !v.Key.IsUnknown() {
			variable.Value = v.Value.ValueString()
		}
		if !v.Description.IsNull() && !v.Description.IsUnknown() {
			variable.Description = v.Description.ValueString()
		}
		if !v.Settable.IsNull() && !v.Settable.IsUnknown() {
			variable.Settable = v.Settable.ValueBool()
		}
		if !v.Encrypted.IsNull() && !v.Encrypted.IsUnknown() {
			variable.Encrypted = v.Encrypted.ValueBool()
		}
		variables[i] = variable
	}
	return &variables, diags
}

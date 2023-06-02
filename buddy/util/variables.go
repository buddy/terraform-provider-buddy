package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type variableModel struct {
	Key         types.String `tfsdk:"key"`
	Encrypted   types.Bool   `tfsdk:"encrypted"`
	Settable    types.Bool   `tfsdk:"settable"`
	Description types.String `tfsdk:"description"`
	Value       types.String `tfsdk:"value"`
	VariableId  types.Int64  `tfsdk:"variable_id"`
}

func variableModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"key":         types.StringType,
		"encrypted":   types.BoolType,
		"settable":    types.BoolType,
		"description": types.StringType,
		"value":       types.StringType,
		"variable_id": types.Int64Type,
	}
}

func (v *variableModel) loadAPI(variable *buddy.Variable) {
	v.Key = types.StringValue(variable.Key)
	v.Encrypted = types.BoolValue(variable.Encrypted)
	v.Settable = types.BoolValue(variable.Settable)
	v.Description = types.StringValue(variable.Description)
	v.Value = types.StringValue(variable.Value)
	v.VariableId = types.Int64Value(int64(variable.Id))
}

func SourceVariableModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"key": schema.StringAttribute{
			Computed: true,
		},
		"encrypted": schema.BoolAttribute{
			Computed: true,
		},
		"settable": schema.BoolAttribute{
			Computed: true,
		},
		"description": schema.StringAttribute{
			Computed: true,
		},
		"value": schema.StringAttribute{
			Computed:  true,
			Sensitive: true,
		},
		"variable_id": schema.Int64Attribute{
			Computed: true,
		},
	}
}

func VariablesModelFromApi(ctx context.Context, variables *[]*buddy.Variable) (basetypes.SetValue, diag.Diagnostics) {
	r := make([]*variableModel, len(*variables))
	for i, v := range *variables {
		r[i] = &variableModel{}
		r[i].loadAPI(v)
	}
	return types.SetValueFrom(ctx, types.ObjectType{AttrTypes: variableModelAttrs()}, &r)
}

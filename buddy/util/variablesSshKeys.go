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

type variableSshKeyModel struct {
	Key            types.String `tfsdk:"key"`
	Encrypted      types.Bool   `tfsdk:"encrypted"`
	Settable       types.Bool   `tfsdk:"settable"`
	Description    types.String `tfsdk:"description"`
	Value          types.String `tfsdk:"value"`
	VariableId     types.Int64  `tfsdk:"variable_id"`
	FilePlace      types.String `tfsdk:"file_place"`
	FilePath       types.String `tfsdk:"file_path"`
	FileChmod      types.String `tfsdk:"file_chmod"`
	Checksum       types.String `tfsdk:"checksum"`
	KeyFingerprint types.String `tfsdk:"key_fingerprint"`
	PublicValue    types.String `tfsdk:"public_value"`
}

func variableSshKeyModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"key":             types.StringType,
		"encrypted":       types.BoolType,
		"settable":        types.BoolType,
		"description":     types.StringType,
		"value":           types.StringType,
		"variable_id":     types.Int64Type,
		"file_place":      types.StringType,
		"file_path":       types.StringType,
		"file_chmod":      types.StringType,
		"checksum":        types.StringType,
		"key_fingerprint": types.StringType,
		"public_value":    types.StringType,
	}
}

func (v *variableSshKeyModel) loadAPI(variable *buddy.Variable) {
	v.Key = types.StringValue(variable.Key)
	v.Encrypted = types.BoolValue(variable.Encrypted)
	v.Settable = types.BoolValue(variable.Settable)
	v.Description = types.StringValue(variable.Description)
	v.Value = types.StringValue(variable.Value)
	v.VariableId = types.Int64Value(int64(variable.Id))
	v.FilePlace = types.StringValue(variable.FilePlace)
	v.FilePath = types.StringValue(variable.FilePath)
	v.FileChmod = types.StringValue(variable.FileChmod)
	v.Checksum = types.StringValue(variable.Checksum)
	v.KeyFingerprint = types.StringValue(variable.KeyFingerprint)
	v.PublicValue = types.StringValue(variable.PublicValue)
}

func SourceVariableSshKeyModelAttributes() map[string]schema.Attribute {
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
		"file_place": schema.StringAttribute{
			Computed: true,
		},
		"file_path": schema.StringAttribute{
			Computed: true,
		},
		"file_chmod": schema.StringAttribute{
			Computed: true,
		},
		"checksum": schema.StringAttribute{
			Computed: true,
		},
		"key_fingerprint": schema.StringAttribute{
			Computed: true,
		},
		"public_value": schema.StringAttribute{
			Computed: true,
		},
	}
}

func VariablesSshKeysModelFromApi(ctx context.Context, variables *[]*buddy.Variable) (basetypes.SetValue, diag.Diagnostics) {
	r := make([]*variableSshKeyModel, len(*variables))
	for i, v := range *variables {
		r[i] = &variableSshKeyModel{}
		r[i].loadAPI(v)
	}
	return types.SetValueFrom(ctx, types.ObjectType{AttrTypes: variableSshKeyModelAttrs()}, &r)
}

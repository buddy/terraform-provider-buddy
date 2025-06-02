package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type targetAuthModel struct {
	Method     types.String `tfsdk:"method"`
	Username   types.String `tfsdk:"username"`
	Password   types.String `tfsdk:"password"`
	Asset      types.String `tfsdk:"asset"`
	Passphrase types.String `tfsdk:"passphrase"`
	Key        types.String `tfsdk:"key"`
	KeyPath    types.String `tfsdk:"key_path"`
}

func TargetAuthModelToApi(ctx context.Context, s *types.Set) (*buddy.TargetAuth, diag.Diagnostics) {
	var t []targetAuthModel
	diags := s.ElementsAs(ctx, &t, false)
	if len(t) == 0 {
		return nil, diags
	}
	if len(t) != 1 {
		diags.Append(diag.NewErrorDiagnostic("Wrong target auth settings", "There should be only one target auth entry"))
		return nil, diags
	}
	tt := t[0]
	var result buddy.TargetAuth
	result.Method = tt.Method.ValueString()
	if !tt.Username.IsNull() && !tt.Username.IsUnknown() {
		result.Username = tt.Username.ValueString()
	}
	if !tt.Password.IsNull() && !tt.Password.IsUnknown() {
		result.Password = tt.Password.ValueString()
	}
	if !tt.Asset.IsNull() && !tt.Asset.IsUnknown() {
		result.Asset = tt.Asset.ValueString()
	}
	if !tt.Passphrase.IsNull() && !tt.Passphrase.IsUnknown() {
		result.Passphrase = tt.Passphrase.ValueString()
	}
	if !tt.Key.IsNull() && !tt.Key.IsUnknown() {
		result.Key = tt.Key.ValueString()
	}
	if !tt.KeyPath.IsNull() && !tt.KeyPath.IsUnknown() {
		result.KeyPath = tt.KeyPath.ValueString()
	}
	return &result, diags
}

func TargetAuthModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"method": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					buddy.TargetAuthMethodSshKey,
					buddy.TargetAuthMethodPassword,
					buddy.TargetAuthMethodAssetsKey,
					buddy.TargetAuthMethodProxyCredentials,
					buddy.TargetAuthMethodProxyKey,
					buddy.TargetAuthMethodHttp,
				),
			},
		},
		"username": schema.StringAttribute{
			Optional: true,
		},
		"password": schema.StringAttribute{
			Optional:  true,
			Sensitive: true,
		},
		"asset": schema.StringAttribute{
			Optional: true,
		},
		"passphrase": schema.StringAttribute{
			Optional:  true,
			Sensitive: true,
		},
		"key": schema.StringAttribute{
			Optional:  true,
			Sensitive: true,
		},
		"key_path": schema.StringAttribute{
			Optional: true,
		},
	}
}

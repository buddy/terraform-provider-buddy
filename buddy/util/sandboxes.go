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

type sandboxModel struct {
	Name       types.String `tfsdk:"name"`
	Status     types.String `tfsdk:"status"`
	Identifier types.String `tfsdk:"identifier"`
	SandboxId  types.String `tfsdk:"sandbox_id"`
	HtmlUrl    types.String `tfsdk:"html_url"`
}

func (s *sandboxModel) loadAPI(sandbox *buddy.Sandbox) {
	s.Name = types.StringValue(sandbox.Name)
	s.Status = types.StringValue(sandbox.Status)
	s.Identifier = types.StringValue(sandbox.Identifier)
	s.SandboxId = types.StringValue(sandbox.Id)
	s.HtmlUrl = types.StringValue(sandbox.HtmlUrl)
}

func sandboxModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"name":       types.StringType,
		"status":     types.StringType,
		"identifier": types.StringType,
		"sandbox_id": types.StringType,
		"html_url":   types.StringType,
	}
}

func SourceSandboxModelAttributes() map[string]sourceschema.Attribute {
	return map[string]sourceschema.Attribute{
		"name": sourceschema.StringAttribute{
			Computed: true,
		},
		"identifier": sourceschema.StringAttribute{
			Computed: true,
		},
		"sandbox_id": sourceschema.StringAttribute{
			Computed: true,
		},
		"html_url": sourceschema.StringAttribute{
			Computed: true,
		},
		"status": sourceschema.StringAttribute{
			Computed: true,
		},
	}
}

func SandboxesModelFromApi(ctx context.Context, sandboxes *[]*buddy.Sandbox) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	l := make([]*sandboxModel, len(*sandboxes))
	for i, v := range *sandboxes {
		l[i] = &sandboxModel{}
		l[i].loadAPI(v)
	}
	ll, d := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: sandboxModelAttrs()}, &l)
	diags.Append(d...)
	return ll, diags
}

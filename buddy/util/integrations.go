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

type integrationModel struct {
	HtmlUrl       types.String `tfsdk:"html_url"`
	IntegrationId types.String `tfsdk:"integration_id"`
	Name          types.String `tfsdk:"name"`
	Type          types.String `tfsdk:"type"`
}

func integrationModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"html_url":       types.StringType,
		"integration_id": types.StringType,
		"name":           types.StringType,
		"type":           types.StringType,
	}
}

func (i *integrationModel) loadAPI(integration *buddy.Integration) {
	i.HtmlUrl = types.StringValue(integration.HtmlUrl)
	i.IntegrationId = types.StringValue(integration.HashId)
	i.Name = types.StringValue(integration.Name)
	i.Type = types.StringValue(integration.Type)
}

func SourceIntegrationModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"html_url": schema.StringAttribute{
			Computed: true,
		},
		"integration_id": schema.StringAttribute{
			Computed: true,
		},
		"name": schema.StringAttribute{
			Computed: true,
		},
		"type": schema.StringAttribute{
			Computed: true,
		},
	}
}

func IntegrationsModelFromApi(ctx context.Context, integrations *[]*buddy.Integration) (basetypes.SetValue, diag.Diagnostics) {
	r := make([]*integrationModel, len(*integrations))
	for i, v := range *integrations {
		r[i] = &integrationModel{}
		r[i].loadAPI(v)
	}
	return types.SetValueFrom(ctx, types.ObjectType{AttrTypes: integrationModelAttrs()}, &r)
}

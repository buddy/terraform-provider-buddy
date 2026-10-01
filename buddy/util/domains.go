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

type domainModel struct {
	Domain   types.String `tfsdk:"domain"`
	DomainId types.String `tfsdk:"domain_id"`
	Type     types.String `tfsdk:"type"`
	HtmlUrl  types.String `tfsdk:"html_url"`
}

func domainModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"domain":    types.StringType,
		"domain_id": types.StringType,
		"type":      types.StringType,
		"html_url":  types.StringType,
	}
}

func (v *domainModel) loadAPI(domain *buddy.Domain) {
	v.Domain = types.StringValue(domain.Name)
	v.DomainId = types.StringValue(domain.Id)
	v.Type = types.StringValue(domain.Type)
	v.HtmlUrl = types.StringValue(domain.HtmlUrl)
}

func SourceDomainModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"domain": schema.StringAttribute{
			Computed: true,
		},
		"domain_id": schema.StringAttribute{
			Computed: true,
		},
		"type": schema.StringAttribute{
			Computed: true,
		},
		"html_url": schema.StringAttribute{
			Computed: true,
		},
	}
}

func DomainsModelFromApi(ctx context.Context, domains *[]*buddy.Domain) (basetypes.SetValue, diag.Diagnostics) {
	r := make([]*domainModel, len(*domains))
	for i, v := range *domains {
		r[i] = &domainModel{}
		r[i].loadAPI(v)
	}
	return types.SetValueFrom(ctx, types.ObjectType{AttrTypes: domainModelAttrs()}, &r)
}

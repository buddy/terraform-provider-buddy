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

type webhookModel struct {
	TargetUrl types.String `tfsdk:"target_url"`
	WebhookId types.Int64  `tfsdk:"webhook_id"`
	HtmlUrl   types.String `tfsdk:"html_url"`
}

func webhookModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"target_url": types.StringType,
		"webhook_id": types.Int64Type,
		"html_url":   types.StringType,
	}
}

func (v *webhookModel) loadAPI(webhook *buddy.Webhook) {
	v.TargetUrl = types.StringValue(webhook.TargetUrl)
	v.WebhookId = types.Int64Value(int64(webhook.Id))
	v.HtmlUrl = types.StringValue(webhook.HtmlUrl)
}

func SourceWebhookModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"target_url": schema.StringAttribute{
			Computed: true,
		},
		"webhook_id": schema.Int64Attribute{
			Computed: true,
		},
		"html_url": schema.StringAttribute{
			Computed: true,
		},
	}
}

func WebhooksModelFromApi(ctx context.Context, webhooks *[]*buddy.Webhook) (basetypes.SetValue, diag.Diagnostics) {
	r := make([]*webhookModel, len(*webhooks))
	for i, v := range *webhooks {
		r[i] = &webhookModel{}
		r[i].loadAPI(v)
	}
	return types.SetValueFrom(ctx, types.ObjectType{AttrTypes: webhookModelAttrs()}, &r)
}

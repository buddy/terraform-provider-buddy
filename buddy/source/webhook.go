package source

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Webhook() *schema.Resource {
	return &schema.Resource{
		Description: "Get webhook by target URL or webhook ID\n\n" +
			"Token scope required: `WORKSPACE`, `WEBHOOK_INFO`",
		ReadContext: readContextWebhook,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The Terraform resource identifier for this item",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"domain": {
				Description:  "The workspace's URL handle",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: util.ValidateDomain,
			},
			"target_url": {
				Description: "The webhook's target URL",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ExactlyOneOf: []string{
					"webhook_id",
					"target_url",
				},
			},
			"webhook_id": {
				Description: "The webhook's ID",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ExactlyOneOf: []string{
					"webhook_id",
					"target_url",
				},
			},
			"html_url": {
				Description: "The webhook's URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func readContextWebhook(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	var diags diag.Diagnostics
	var webhook *buddy.Webhook
	var err error
	domain := d.Get("domain").(string)
	if webhookId, ok := d.GetOk("webhook_id"); ok {
		webhook, _, err = c.WebhookService.Get(domain, webhookId.(int))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		targetUrl := d.Get("target_url").(string)
		webhooks, _, err := c.WebhookService.GetList(domain)
		if err != nil {
			return diag.FromErr(err)
		}
		for _, w := range webhooks.Webhooks {
			if w.TargetUrl == targetUrl {
				webhook = w
				break
			}
		}
		if webhook == nil {
			return diag.Errorf("Webhook not found")
		}
	}
	err = util.ApiWebhookToResourceData(domain, webhook, d, true)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

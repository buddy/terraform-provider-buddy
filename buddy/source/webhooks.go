package source

import (
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
)

func Webhooks() *schema.Resource {
	return &schema.Resource{
		Description: "List webhooks and optionally filter them by target_url\n\n" +
			"Token scope required: `WORKSPACE`, `WEBHOOK_INFO`",
		ReadContext: readContextWebhooks,
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
			"target_url_regex": {
				Description:  "The webhook's target_url regular expression to match",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"webhooks": {
				Description: "List of webhooks",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_url": {
							Description: "The webhook's target URL",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"webhook_id": {
							Description: "The webhook's ID",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"html_url": {
							Description: "The webhook's URL",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func readContextWebhooks(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	var targetUrlRegex *regexp.Regexp
	domain := d.Get("domain").(string)
	webhooks, _, err := c.WebhookService.GetList(domain)
	if err != nil {
		return diag.FromErr(err)
	}
	var result []interface{}
	if targetUrl, ok := d.GetOk("target_url_regex"); ok {
		targetUrlRegex = regexp.MustCompile(targetUrl.(string))
	}
	for _, w := range webhooks.Webhooks {
		if targetUrlRegex != nil && !targetUrlRegex.MatchString(w.TargetUrl) {
			continue
		}
		result = append(result, util.ApiShortWebhookToMap(w))
	}
	d.SetId(util.UniqueString())
	err = d.Set("webhooks", result)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

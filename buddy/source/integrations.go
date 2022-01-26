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

func Integrations() *schema.Resource {
	return &schema.Resource{
		Description: "List integrations and optionally filter them by name or type\n\n" +
			"Token scope required: `INTEGRATION_INFO`",
		ReadContext: readContextIntegrations,
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
			"name_regex": {
				Description:  "The integration's name regular expression to match",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"type": {
				Description: "The integration's type",
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					api.IntegrationTypeDigitalOcean,
					api.IntegrationTypeAmazon,
					api.IntegrationTypeShopify,
					api.IntegrationTypePushover,
					api.IntegrationTypeRackspace,
					api.IntegrationTypeCloudflare,
					api.IntegrationTypeNewRelic,
					api.IntegrationTypeSentry,
					api.IntegrationTypeRollbar,
					api.IntegrationTypeDatadog,
					api.IntegrationTypeDigitalOceanSpaces,
					api.IntegrationTypeHoneybadger,
					api.IntegrationTypeVultr,
					api.IntegrationTypeSentryEnterprise,
					api.IntegrationTypeLoggly,
					api.IntegrationTypeFirebase,
					api.IntegrationTypeUpcloud,
					api.IntegrationTypeGhostInspector,
					api.IntegrationTypeAzureCloud,
					api.IntegrationTypeDockerHub,
				}, false),
			},
			"integrations": {
				Description: "List of integrations",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"html_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"integration_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func readContextIntegrations(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	var nameRegex *regexp.Regexp
	domain := d.Get("domain").(string)
	typ := d.Get("type").(string)
	integrations, _, err := c.IntegrationService.GetList(domain)
	if err != nil {
		return diag.FromErr(err)
	}
	var result []interface{}
	if name, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(name.(string))
	}
	for _, i := range integrations.Integrations {
		if nameRegex != nil && !nameRegex.MatchString(i.Name) {
			continue
		}
		if typ != "" && typ != i.Type {
			continue
		}
		result = append(result, util.ApiShortIntegrationToMap(i))
	}
	d.SetId(util.UniqueString())
	err = d.Set("integrations", result)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

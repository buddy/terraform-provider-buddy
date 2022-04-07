package source

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Integration() *schema.Resource {
	return &schema.Resource{
		Description: "Get integration by name or integration ID\n\n" +
			"Token scope required: `INTEGRATION_INFO`",
		ReadContext: readContextIntegration,
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
			"name": {
				Description: "The integration's name",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ExactlyOneOf: []string{
					"name",
					"integration_id",
				},
			},
			"type": {
				Description: "The integration's type",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"integration_id": {
				Description: "The integration's ID",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ExactlyOneOf: []string{
					"name",
					"integration_id",
				},
			},
			"html_url": {
				Description: "The integration's URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func readContextIntegration(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	var diags diag.Diagnostics
	var integration *buddy.Integration
	var err error
	domain := d.Get("domain").(string)
	if integrationId, ok := d.GetOk("integration_id"); ok {
		integration, _, err = c.IntegrationService.Get(domain, integrationId.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		name := d.Get("name").(string)
		integrations, _, err := c.IntegrationService.GetList(domain)
		if err != nil {
			return diag.FromErr(err)
		}
		for _, i := range integrations.Integrations {
			if i.Name == name {
				integration = i
				break
			}
		}
		if integration == nil {
			return diag.Errorf("Integration not found")
		}
	}
	err = util.ApiIntegrationToResourceData(domain, integration, d, true)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

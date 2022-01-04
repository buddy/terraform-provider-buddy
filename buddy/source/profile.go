package source

import (
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Profile() *schema.Resource {
	return &schema.Resource{
		Description: "`buddy_profile` data source allows you to fetch details of your account\n\n" +
			"Required scopes for your token: `USER_INFO`",
		ReadContext: readContextProfile,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Compound id of this resource",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Your name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"member_id": {
				Description: "Your user id",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"avatar_url": {
				Description: "Your avatar url",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"html_url": {
				Description: "Url to Buddy myid",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func readContextProfile(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	p, _, err := c.ProfileService.Get()
	if err != nil {
		return diag.FromErr(err)
	}
	err = util.ApiProfileToResourceData(p, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

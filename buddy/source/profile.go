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
		Description: "Get details of a Buddy's user profile\n\n" +
			"Token scope required: `USER_INFO`",
		ReadContext: readContextProfile,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The Terraform resource identifier for this item",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The user's name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"member_id": {
				Description: "The user's ID",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"avatar_url": {
				Description: "The user's avatar URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"html_url": {
				Description: "The user's profile URL",
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

package resource

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Profile() *schema.Resource {
	return &schema.Resource{
		Description: "Manage a user profile\n\n" +
			"Token scope required: `USER_INFO`",
		CreateContext: createContextProfile,
		ReadContext:   readContextProfile,
		UpdateContext: updateContextProfile,
		DeleteContext: deleteContextProfile,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The Terraform resource identifier for this item",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The user's name",
				Type:        schema.TypeString,
				Required:    true,
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
				Description: "The user's URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func createContextProfile(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return updateContextProfile(ctx, d, meta)
}

func readContextProfile(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
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

func updateContextProfile(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	u := buddy.ProfileOps{
		Name: util.InterfaceStringToPointer(d.Get("name")),
	}
	_, _, err := c.ProfileService.Update(&u)
	if err != nil {
		return diag.FromErr(err)
	}
	return readContextProfile(ctx, d, meta)
}

func deleteContextProfile(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

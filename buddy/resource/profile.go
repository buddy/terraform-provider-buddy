package resource

import (
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Profile() *schema.Resource {
	return &schema.Resource{
		Description: "`buddy_profile` allows you to manage your Buddy account.\n\n" +
			"Required scopes for your token: `USER_INFO`",
		CreateContext: createContextProfile,
		ReadContext:   readContextProfile,
		UpdateContext: updateContextProfile,
		DeleteContext: deleteContextProfile,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Compound id of the user",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Name of your user",
				Type:        schema.TypeString,
				Required:    true,
			},
			"member_id": {
				Description: "Id of your user",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"avatar_url": {
				Description: "Avatar url of your user",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"html_url": {
				Description: "Url of your user",
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

func updateContextProfile(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	u := api.ProfileOperationOptions{
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

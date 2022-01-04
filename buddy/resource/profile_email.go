package resource

import (
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ProfileEmail() *schema.Resource {
	return &schema.Resource{
		Description: "`buddy_profile_email` allows you to manage your Buddy account email.\n\n" +
			"Required scopes for your token: `MANAGE_EMAILS`, `USER_EMAIL`",
		CreateContext: createContextProfileEmail,
		ReadContext:   readContextProfileEmail,
		DeleteContext: deleteContextProfileEmail,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Compound id of the email",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email": {
				Description: "Email to add to your account",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"confirmed": {
				Description: "Is email confirmed",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func deleteContextProfileEmail(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	_, err := c.ProfileEmailService.Delete(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func readContextProfileEmail(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	p, _, err := c.ProfileEmailService.GetList()
	if err != nil {
		return diag.FromErr(err)
	}
	found := false
	email := d.Id()
	for _, v := range p.Emails {
		if v.Email == email {
			found = true
			err = util.ApiProfileEmailToResourceData(v, d)
			if err != nil {
				return diag.FromErr(err)
			}
			break
		}
	}
	if !found {
		d.SetId("")
	}
	return diags
}

func createContextProfileEmail(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	p, _, err := c.ProfileEmailService.Create(&api.ProfileEmailOperationOptions{
		Email: util.InterfaceStringToPointer(d.Get("email")),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(p.Email)
	return readContextProfileEmail(ctx, d, meta)
}

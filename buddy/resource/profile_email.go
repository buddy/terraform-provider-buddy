package resource

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ProfileEmail() *schema.Resource {
	return &schema.Resource{
		Description: "Create and manage a user's email\n\n" +
			"Token scopes required: `MANAGE_EMAILS`, `USER_EMAIL`",
		CreateContext: createContextProfileEmail,
		ReadContext:   readContextProfileEmail,
		DeleteContext: deleteContextProfileEmail,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The Terraform resource identifier for this item",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email": {
				Description: "The email to add to the user's profile",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"confirmed": {
				Description: "Is the email confirmed",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func deleteContextProfileEmail(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	var diags diag.Diagnostics
	_, err := c.ProfileEmailService.Delete(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func readContextProfileEmail(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	var diags diag.Diagnostics
	p, resp, err := c.ProfileEmailService.GetList()
	if err != nil {
		if util.IsResourceNotFound(resp, err) {
			d.SetId("")
		}
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
	c := meta.(*buddy.Client)
	p, _, err := c.ProfileEmailService.Create(&buddy.ProfileEmailOps{
		Email: util.InterfaceStringToPointer(d.Get("email")),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(p.Email)
	return readContextProfileEmail(ctx, d, meta)
}

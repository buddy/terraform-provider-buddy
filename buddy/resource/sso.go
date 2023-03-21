package resource

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Sso() *schema.Resource {
	return &schema.Resource{
		Description: "Manage SSO in workspace\n\n" +
			"Workspace administrator rights are required\n\n" +
			"Token scopes required: `WORKSPACE`",
		CreateContext: createContextSso,
		ReadContext:   readContextSso,
		UpdateContext: updateContextSSo,
		DeleteContext: deleteContextSso,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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
				ForceNew:     true,
				ValidateFunc: util.ValidateDomain,
			},
			"sso_url": {
				Description: "The identity provider single sign-on url",
				Type:        schema.TypeString,
				Required:    true,
			},
			"issuer": {
				Description: "The identity provider issuer url",
				Type:        schema.TypeString,
				Required:    true,
			},
			"certificate": {
				Description: "The identity provider certificate",
				Type:        schema.TypeString,
				Required:    true,
			},
			"signature": {
				Description: "The SAML signature algorithm. Allowed: `sha1`, `sha256`, `sha512`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"digest": {
				Description: "The SAML digest algorithm. Allowed: `sha1`, `sha256`, `sha512`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"require_for_all": {
				Description: "Enable mandatory SAML SSO authentication for all workspace members",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"html_url": {
				Description: "The Sso's URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func createContextSso(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	domain := d.Get("domain").(string)
	_, _ = c.SsoService.Enable(domain)
	digest := util.InterfaceStringToPointer(d.Get("digest"))
	ssoUrl := util.InterfaceStringToPointer(d.Get("sso_url"))
	issuer := util.InterfaceStringToPointer(d.Get("issuer"))
	certificate := util.InterfaceStringToPointer(d.Get("certificate"))
	signature := util.InterfaceStringToPointer(d.Get("signature"))
	requireForAll := util.InterfaceBoolToPointer(d.Get("require_for_all"))
	ops := buddy.SsoUpdateOps{
		DigestMethod:            digest,
		SignatureMethod:         signature,
		SsoUrl:                  ssoUrl,
		Issuer:                  issuer,
		Certificate:             certificate,
		RequireSsoForAllMembers: requireForAll,
	}
	_, _, err := c.SsoService.Update(domain, &ops)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(domain)
	return readContextSso(ctx, d, meta)
}

func readContextSso(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	var diags diag.Diagnostics
	sso, resp, err := c.SsoService.Get(d.Id())
	if err != nil {
		if util.IsResourceNotFound(resp, err) {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	err = util.ApiSsoToResourceData(d.Id(), sso, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func updateContextSSo(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return createContextSso(ctx, d, meta)
}

func deleteContextSso(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	var diags diag.Diagnostics
	_, err := c.SsoService.Disable(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

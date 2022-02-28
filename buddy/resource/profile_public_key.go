package resource

import (
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
	"strconv"
)

func ProfilePublicKey() *schema.Resource {
	return &schema.Resource{
		Description: "Create and manage a user's public key\n\n" +
			"Token scope required: `USER_KEY`",
		CreateContext: createContextProfilePublicKey,
		ReadContext:   readContextProfilePublicKey,
		DeleteContext: deleteContextProfilePublicKey,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The public key's ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"content": {
				Description: "The public key's content (starts with ssh-rsa)",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"title": {
				Description: "The public key's title",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"html_url": {
				Description: "The public key's URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func deleteContextProfilePublicKey(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	keyId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = c.PublicKeyService.Delete(keyId)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func readContextProfilePublicKey(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	keyId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	key, resp, err := c.PublicKeyService.Get(keyId)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	err = util.ApiPublicKeyToResourceData(key, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func createContextProfilePublicKey(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	opt := api.PublicKeyOperationOptions{
		Content: util.InterfaceStringToPointer(d.Get("content")),
	}
	if title, ok := d.GetOk("title"); ok {
		opt.Title = util.InterfaceStringToPointer(title)
	}
	key, _, err := c.PublicKeyService.Create(&opt)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(key.Id))
	return readContextProfilePublicKey(ctx, d, meta)
}

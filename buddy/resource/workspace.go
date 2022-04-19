package resource

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Workspace() *schema.Resource {
	return &schema.Resource{
		Description: "Create and manage a workspace\n\n" +
			"Invite-only token is required. Contact support@buddy.works for more details\n\n" +
			"Token scope required: `WORKSPACE`",
		CreateContext: createContextWorkspace,
		ReadContext:   readContextWorkspace,
		UpdateContext: updateContextWorkspace,
		DeleteContext: deleteContextWorkspace,
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
			"name": {
				Description: "The workspace's name",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"encryption_salt": {
				Description: "The workspace's salt to encrypt secrets in YAML & API",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"workspace_id": {
				Description: "The workspace's ID",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"html_url": {
				Description: "The workspace's URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"owner_id": {
				Description: "The workspace's owner ID",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"frozen": {
				Description: "Is the workspace frozen",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"create_date": {
				Description: "The workspace's create date",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func deleteContextWorkspace(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	var diags diag.Diagnostics
	_, err := c.WorkspaceService.Delete(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func readContextWorkspace(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	var diags diag.Diagnostics
	workspace, resp, err := c.WorkspaceService.Get(d.Id())
	if err != nil {
		if util.IsResourceNotFound(resp, err) {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	err = util.ApiWorkspaceToResourceData(workspace, d, false)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func updateContextWorkspace(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	if d.HasChanges("name", "encryption_salt") {
		domain := d.Get("domain").(string)
		opt := buddy.WorkspaceUpdateOps{
			Name: util.InterfaceStringToPointer(d.Get("name")),
		}
		if salt, ok := d.GetOk("encryption_salt"); ok {
			opt.EncryptionSalt = util.InterfaceStringToPointer(salt)
		}
		_, _, err := c.WorkspaceService.Update(domain, &opt)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return readContextWorkspace(ctx, d, meta)
}

func createContextWorkspace(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	opt := buddy.WorkspaceCreateOps{
		Domain: util.InterfaceStringToPointer(d.Get("domain")),
	}
	if salt, ok := d.GetOk("encryption_salt"); ok {
		opt.EncryptionSalt = util.InterfaceStringToPointer(salt)
	}
	if name, ok := d.GetOk("name"); ok {
		opt.Name = util.InterfaceStringToPointer(name)
	}
	workspace, _, err := c.WorkspaceService.Create(&opt)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(workspace.Domain)
	return readContextWorkspace(ctx, d, meta)
}

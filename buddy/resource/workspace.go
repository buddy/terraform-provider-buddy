package resource

import (
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
)

func Workspace() *schema.Resource {
	return &schema.Resource{
		Description: "`buddy_workspace` allows you to create Buddy workspace.\n\n" +
			"You will need special token to manage this resource. Contact support@buddy.works for more info.\n\n" +
			"Required scopes for your token: `WORKSPACE`",
		CreateContext: createContextWorkspace,
		ReadContext:   readContextWorkspace,
		UpdateContext: updateContextWorkspace,
		DeleteContext: deleteContextWorkspace,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Compound id of the workspace",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"domain": {
				Description:  "Domain of the workspace to manage",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: util.ValidateDomain,
			},
			"name": {
				Description: "Name of the workspace",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"encryption_salt": {
				Description: "Salt to encrypt secrets in YAML & API",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"workspace_id": {
				Description: "Id of the workspace",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"html_url": {
				Description: "Url of the workspace",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"owner_id": {
				Description: "Workspace owner id",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"frozen": {
				Description: "Is workspace frozen",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"create_date": {
				Description: "Workspace create date",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func deleteContextWorkspace(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// nie ma usuwania
	var diags diag.Diagnostics
	return diags
}

func readContextWorkspace(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	workspace, resp, err := c.WorkspaceService.Get(d.Id())
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
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
	c := meta.(*api.Client)
	if d.HasChanges("name", "encryption_salt") {
		domain := d.Get("domain").(string)
		opt := api.WorkspaceOperationOptions{}
		if d.HasChange("encryption_salt") {
			opt.EncryptionSalt = util.InterfaceStringToPointer(d.Get("encryption_salt"))
		}
		if d.HasChange("name") {
			opt.Name = util.InterfaceStringToPointer(d.Get("name"))
		}
		_, _, err := c.WorkspaceService.Update(domain, &opt)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return readContextWorkspace(ctx, d, meta)
}

func createContextWorkspace(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	opt := api.WorkspaceOperationOptions{
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

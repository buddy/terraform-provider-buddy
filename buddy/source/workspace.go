package source

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Workspace() *schema.Resource {
	return &schema.Resource{
		Description: "Get workspace by URL handle or name\n\n" +
			"Token scope required: `WORKSPACE`",
		ReadContext: readContextWorkspace,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The Terraform resource identifier for this item",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"domain": {
				Description: "The workspace's URL handle",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ExactlyOneOf: []string{
					"domain",
					"name",
				},
				ValidateFunc: util.ValidateDomain,
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
			"name": {
				Description: "The workspace's name",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ExactlyOneOf: []string{
					"domain",
					"name",
				},
			},
		},
	}
}

func readContextWorkspace(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	var diags diag.Diagnostics
	var workspace *buddy.Workspace
	var err error
	if domain, ok := d.GetOk("domain"); ok {
		workspace, _, err = c.WorkspaceService.Get(domain.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		name := d.Get("name").(string)
		workspaces, _, err := c.WorkspaceService.GetList()
		if err != nil {
			return diag.FromErr(err)
		}
		for _, w := range workspaces.Workspaces {
			if w.Name == name {
				workspace = w
				break
			}
		}
		if workspace == nil {
			return diag.Errorf("Workspace not found")
		}
	}
	err = util.ApiWorkspaceToResourceData(workspace, d, true)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

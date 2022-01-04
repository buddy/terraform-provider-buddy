package source

import (
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Workspace() *schema.Resource {
	return &schema.Resource{
		Description: "`buddy_workspace` data source allows you to find workspace by domain or name\n\n" +
			"Required scopes for your token: `WORKSPACE`",
		ReadContext: readContextWorkspace,
		Schema: map[string]*schema.Schema{
			"domain": {
				Description: "Domain of the workspace",
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
				Description: "Id of the workspace",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"html_url": {
				Description: "Url of the workspace",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Name of the workspace",
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
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	var workspace *api.Workspace
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

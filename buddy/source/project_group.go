package source

import (
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ProjectGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Get project group\n\n" +
			"Token scope required: `WORKSPACE`",
		ReadContext: readContextProjectGroup,
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
				ValidateFunc: util.ValidateDomain,
			},
			"project_name": {
				Description: "The project's name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"group_id": {
				Description: "The group's ID",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"html_url": {
				Description: "The group's URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The group's name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"permission": {
				Description: "The group's permission in the project",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"html_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"permission_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pipeline_access_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repository_access_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sandbox_access_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func readContextProjectGroup(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	domain := d.Get("domain").(string)
	projectName := d.Get("project_name").(string)
	groupId := d.Get("group_id").(int)
	group, _, err := c.ProjectGroupService.GetProjectGroup(domain, projectName, groupId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = util.ApiProjectGroupToResourceData(domain, projectName, group, d, false)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

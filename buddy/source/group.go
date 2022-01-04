package source

import (
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Group() *schema.Resource {
	return &schema.Resource{
		Description: "`buddy_group` data source allows you to find group by name or group_id\n\n" +
			"Required scopes for your token: `WORKSPACE`",
		ReadContext: readContextGroup,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Compound id of the group",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"domain": {
				Description:  "Domain of the workspace",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: util.ValidateDomain,
			},
			"name": {
				Description: "Name of the group",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ExactlyOneOf: []string{
					"name",
					"group_id",
				},
			},
			"group_id": {
				Description: "Id of the group",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ExactlyOneOf: []string{
					"name",
					"group_id",
				},
			},
			"html_url": {
				Description: "Url of the group",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Description of the group",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func readContextGroup(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	var group *api.Group
	var err error
	domain := d.Get("domain").(string)
	if groupId, ok := d.GetOk("group_id"); ok {
		group, _, err = c.GroupService.Get(domain, groupId.(int))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		name := d.Get("name").(string)
		groups, _, err := c.GroupService.GetList(domain)
		if err != nil {
			return diag.FromErr(err)
		}
		for _, g := range groups.Groups {
			if g.Name == name {
				group = g
				break
			}
		}
		if group == nil {
			return diag.Errorf("Group not found")
		}
	}
	err = util.ApiGroupToResourceData(domain, group, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

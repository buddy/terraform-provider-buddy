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
		Description: "Get group by name or group ID\n\n" +
			"Token scope required: `WORKSPACE`",
		ReadContext: readContextGroup,
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
			"name": {
				Description: "The group's name",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ExactlyOneOf: []string{
					"name",
					"group_id",
				},
			},
			"group_id": {
				Description: "The group's ID",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ExactlyOneOf: []string{
					"name",
					"group_id",
				},
			},
			"html_url": {
				Description: "The group's URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "The group's description",
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

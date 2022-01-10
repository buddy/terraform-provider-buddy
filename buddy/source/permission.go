package source

import (
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Permission() *schema.Resource {
	return &schema.Resource{
		Description: "Get permission by name or permission_id\n\n" +
			"Token scope required: `WORKSPACE`",
		ReadContext: readContextPermission,
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
				Description: "The permission's name",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ExactlyOneOf: []string{
					"permission_id",
					"name",
				},
			},
			"pipeline_access_level": {
				Description: "The permission's access level to pipelines",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"repository_access_level": {
				Description: "The permission's access level to repository",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"sandbox_access_level": {
				Description: "The permission's access level to sandboxes",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"permission_id": {
				Description: "The permission's ID",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ExactlyOneOf: []string{
					"permission_id",
					"name",
				},
			},
			"html_url": {
				Description: "The permission's URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "The permission's description",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "The permission's type",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func readContextPermission(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	var permission *api.Permission
	var err error
	domain := d.Get("domain").(string)
	if permissionId, ok := d.GetOk("permission_id"); ok {
		permission, _, err = c.PermissionService.Get(domain, permissionId.(int))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		name := d.Get("name").(string)
		permissions, _, err := c.PermissionService.GetList(domain)
		if err != nil {
			return diag.FromErr(err)
		}
		for _, p := range permissions.PermissionSets {
			if p.Name == name {
				permission = p
				break
			}
		}
		if permission == nil {
			return diag.Errorf("Permission not found")
		}
	}
	err = util.ApiPermissionToResourceData(domain, permission, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

package source

import (
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
)

func Permissions() *schema.Resource {
	return &schema.Resource{
		Description: "List permissions and optionally filter them by name or type\n\n" +
			"Token scope required: `WORKSPACE`",
		ReadContext: readContextPermissions,
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
			"name_regex": {
				Description:  "The permission's name regular expression to match",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"type": {
				Description: "Filter permissions by type (`CUSTOM`, `READ_ONLY`, `DEVELOPER`)",
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					api.PermissionTypeCustom,
					api.PermissionTypeReadOnly,
					api.PermissionTypeDeveloper,
				}, false),
			},
			"permissions": {
				Description: "List of permissions",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The permission's name",
							Type:        schema.TypeString,
							Computed:    true,
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
							Description: "The permission's access level to sandboxes.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"permission_id": {
							Description: "The permission's ID",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"html_url": {
							Description: "The permission's URL",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"type": {
							Description: "The permission's type",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func readContextPermissions(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	var nameRegex *regexp.Regexp
	domain := d.Get("domain").(string)
	typ := d.Get("type").(string)
	permissions, _, err := c.PermissionService.GetList(domain)
	if err != nil {
		return diag.FromErr(err)
	}
	var result []interface{}
	if name, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(name.(string))
	}
	for _, p := range permissions.PermissionSets {
		if nameRegex != nil && !nameRegex.MatchString(p.Name) {
			continue
		}
		if typ != "" && typ != p.Type {
			continue
		}
		result = append(result, util.ApiShortPermissionToMap(p))
	}
	d.SetId(util.UniqueString())
	err = d.Set("permissions", result)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

package source

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
)

func Permissions() *schema.Resource {
	return &schema.Resource{
		Description: "List permissions (roles) and optionally filter them by name or type\n\n" +
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
				Description: "Filter permissions by type (`CUSTOM`, `READ_ONLY`, `DEVELOPER`, `PROJECT_MANAGER`)",
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					buddy.PermissionTypeCustom,
					buddy.PermissionTypeReadOnly,
					buddy.PermissionTypeDeveloper,
					buddy.PermissionTypeProjectManager,
				}, false),
			},
			"permissions": {
				Description: "List of permissions (roles)",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
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
						"project_team_access_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sandbox_access_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"permission_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"html_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func readContextPermissions(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
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

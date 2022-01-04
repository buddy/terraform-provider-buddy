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

func Groups() *schema.Resource {
	return &schema.Resource{
		Description: "`buddy_groups` data source allows you to get list of groups in workspace and filter them by name\n\n" +
			"Required scopes for your token: `WORKSPACE`",
		ReadContext: readContextGroups,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Compound id of the resource",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"domain": {
				Description:  "Domain of the workspace",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: util.ValidateDomain,
			},
			"name_regex": {
				Description:  "Regular expression to match name of the group",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"groups": {
				Description: "List of groups",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"html_url": {
							Description: "Url of the group",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"group_id": {
							Description: "Id of the group",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"name": {
							Description: "Name of the group",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func readContextGroups(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	var nameRegex *regexp.Regexp
	domain := d.Get("domain").(string)
	groups, _, err := c.GroupService.GetList(domain)
	if err != nil {
		return diag.FromErr(err)
	}
	var result []interface{}
	if name, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(name.(string))
	}
	for _, g := range groups.Groups {
		if nameRegex != nil && !nameRegex.MatchString(g.Name) {
			continue
		}
		result = append(result, util.ApiShortGroupToMap(g))
	}
	d.SetId(util.UniqueString())
	err = d.Set("groups", result)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

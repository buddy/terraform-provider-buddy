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

func Members() *schema.Resource {
	return &schema.Resource{
		Description: "`buddy_members` data source allows you to get list of members in workspace and filter them by name\n\n" +
			"Required scopes for your token: `WORKSPACE`",
		ReadContext: readContextMembers,
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
				Description:  "Regular expression to match name of the member",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"members": {
				Description: "List of members",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"html_url": {
							Description: "Url of the member",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"member_id": {
							Description: "Id of the member",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"name": {
							Description: "Name of the member",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"email": {
							Description: "Email of the member",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"avatar_url": {
							Description: "Avatar url of the member",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"admin": {
							Description: "Is member a workspace administrator",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"workspace_owner": {
							Description: "Is member a workspace owner",
							Type:        schema.TypeBool,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func readContextMembers(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	var nameRegex *regexp.Regexp
	domain := d.Get("domain").(string)
	members, _, err := c.MemberService.GetList(domain)
	if err != nil {
		return diag.FromErr(err)
	}
	var result []interface{}
	if name, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(name.(string))
	}
	for _, m := range members.Members {
		if nameRegex != nil && !nameRegex.MatchString(m.Name) {
			continue
		}
		result = append(result, util.ApiShortMemberToMap(m))
	}
	d.SetId(util.UniqueString())
	err = d.Set("members", result)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

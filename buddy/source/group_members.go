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

func GroupMembers() *schema.Resource {
	return &schema.Resource{
		Description: "List members of a group and optionally filter them by name\n\n" +
			"Token scope required: `WORKSPACE`",
		ReadContext: readContextGroupMembers,
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
			"group_id": {
				Description: "The group's ID",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"name_regex": {
				Description:  "The member's name regular expression to match",
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
							Type:     schema.TypeString,
							Computed: true,
						},
						"member_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"avatar_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"admin": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"workspace_owner": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func readContextGroupMembers(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	var nameRegex *regexp.Regexp
	domain := d.Get("domain").(string)
	groupId := d.Get("group_id").(int)
	members, _, err := c.GroupService.GetGroupMembers(domain, groupId)
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

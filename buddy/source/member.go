package source

import (
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Member() *schema.Resource {
	return &schema.Resource{
		Description: "`buddy_member` data source allows you to find member by name, email or member_id\n\n" +
			"Required scopes for your token: `WORKSPACE`",
		ReadContext: readContextMember,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Compound id of the member",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"domain": {
				Description:  "Domain of the workspace",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: util.ValidateDomain,
			},
			"email": {
				Description: "Email of the member",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ExactlyOneOf: []string{
					"email",
					"name",
					"member_id",
				},
			},
			"admin": {
				Description: "Is member a workspace administrator",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"name": {
				Description: "Name of the member",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ExactlyOneOf: []string{
					"email",
					"name",
					"member_id",
				},
			},
			"member_id": {
				Description: "Id of the member",
				Type:        schema.TypeInt,
				Optional:    true,
				ExactlyOneOf: []string{
					"email",
					"name",
					"member_id",
				},
			},
			"html_url": {
				Description: "Url of the member",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"avatar_url": {
				Description: "Avatar url of the member",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"workspace_owner": {
				Description: "Is member a workspace owner",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func readContextMember(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	var member *api.Member
	var err error
	domain := d.Get("domain").(string)
	if memberId, ok := d.GetOk("member_id"); ok {
		member, _, err = c.MemberService.Get(domain, memberId.(int))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		name, nameOk := d.GetOk("name")
		email, emailOk := d.GetOk("email")
		list, _, err := c.MemberService.GetList(domain)
		if err != nil {
			return diag.FromErr(err)
		}
		for _, m := range list.Members {
			if nameOk && m.Name == name.(string) {
				member = m
				break
			}
			if emailOk && m.Email == email.(string) {
				member = m
				break
			}
		}
		if member == nil {
			return diag.Errorf("Member not found %d", len(list.Members))
		}
	}
	err = util.ApiMemberToResourceData(domain, member, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

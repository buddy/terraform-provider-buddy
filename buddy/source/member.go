package source

//
//import (
//	"buddy-terraform/buddy/util"
//	"context"
//	"github.com/buddy/api-go-sdk/buddy"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
//)
//
//func Member() *schema.Resource {
//	return &schema.Resource{
//		Description: "Get member by name, email or member ID\n\n" +
//			"Token scope required: `WORKSPACE`",
//		ReadContext: readContextMember,
//		Schema: map[string]*schema.Schema{
//			"id": {
//				Description: "The Terraform resource identifier for this item",
//				Type:        schema.TypeString,
//				Computed:    true,
//			},
//			"domain": {
//				Description:  "The workspace's URL handle",
//				Type:         schema.TypeString,
//				Required:     true,
//				ValidateFunc: util.ValidateDomain,
//			},
//			"email": {
//				Description: "The member's email",
//				Type:        schema.TypeString,
//				Computed:    true,
//				Optional:    true,
//				ExactlyOneOf: []string{
//					"email",
//					"name",
//					"member_id",
//				},
//			},
//			"admin": {
//				Description: "Is the member a workspace administrator",
//				Type:        schema.TypeBool,
//				Computed:    true,
//			},
//			"name": {
//				Description: "The member's name",
//				Type:        schema.TypeString,
//				Computed:    true,
//				Optional:    true,
//				ExactlyOneOf: []string{
//					"email",
//					"name",
//					"member_id",
//				},
//			},
//			"member_id": {
//				Description: "The member's ID",
//				Type:        schema.TypeInt,
//				Optional:    true,
//				ExactlyOneOf: []string{
//					"email",
//					"name",
//					"member_id",
//				},
//			},
//			"html_url": {
//				Description: "The member's URL",
//				Type:        schema.TypeString,
//				Computed:    true,
//			},
//			"avatar_url": {
//				Description: "The member's avatar URL",
//				Type:        schema.TypeString,
//				Computed:    true,
//			},
//			"workspace_owner": {
//				Description: "Is the member the workspace owner",
//				Type:        schema.TypeBool,
//				Computed:    true,
//			},
//		},
//	}
//}
//
//func readContextMember(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	c := meta.(*buddy.Client)
//	var diags diag.Diagnostics
//	var member *buddy.Member
//	var err error
//	domain := d.Get("domain").(string)
//	if memberId, ok := d.GetOk("member_id"); ok {
//		member, _, err = c.MemberService.Get(domain, memberId.(int))
//		if err != nil {
//			return diag.FromErr(err)
//		}
//	} else {
//		name, nameOk := d.GetOk("name")
//		email, emailOk := d.GetOk("email")
//		list, _, err := c.MemberService.GetListAll(domain)
//		if err != nil {
//			return diag.FromErr(err)
//		}
//		for _, m := range list.Members {
//			if nameOk && m.Name == name.(string) {
//				member = m
//				break
//			}
//			if emailOk && m.Email == email.(string) {
//				member = m
//				break
//			}
//		}
//		if member == nil {
//			return diag.Errorf("Member not found %d", len(list.Members))
//		}
//	}
//	err = util.ApiMemberToResourceData(domain, member, d, true)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	return diags
//}

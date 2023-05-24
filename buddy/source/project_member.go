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
//func ProjectMember() *schema.Resource {
//	return &schema.Resource{
//		Description: "Get project member\n\n" +
//			"Token scope required: `WORKSPACE`",
//		ReadContext: readContextProjectMember,
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
//			"project_name": {
//				Description: "The project's name",
//				Type:        schema.TypeString,
//				Required:    true,
//			},
//			"member_id": {
//				Description: "The member's ID",
//				Type:        schema.TypeInt,
//				Required:    true,
//			},
//			"html_url": {
//				Description: "The member's URL",
//				Type:        schema.TypeString,
//				Computed:    true,
//			},
//			"name": {
//				Description: "The member's name",
//				Type:        schema.TypeString,
//				Computed:    true,
//			},
//			"email": {
//				Description: "The member's email",
//				Type:        schema.TypeString,
//				Computed:    true,
//			},
//			"avatar_url": {
//				Description: "The member's avatar URL",
//				Type:        schema.TypeString,
//				Computed:    true,
//			},
//			"admin": {
//				Description: "Is the member a workspace administrator",
//				Type:        schema.TypeBool,
//				Computed:    true,
//			},
//			"workspace_owner": {
//				Description: "Is the member the workspace owner",
//				Type:        schema.TypeBool,
//				Computed:    true,
//			},
//			"permission": {
//				Description: "The member's permission in the project",
//				Type:        schema.TypeList,
//				Computed:    true,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						"html_url": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"permission_id": {
//							Type:     schema.TypeInt,
//							Computed: true,
//						},
//						"name": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"type": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"pipeline_access_level": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"project_team_access_level": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"repository_access_level": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"sandbox_access_level": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//					},
//				},
//			},
//		},
//	}
//}
//
//func readContextProjectMember(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	c := meta.(*buddy.Client)
//	var diags diag.Diagnostics
//	domain := d.Get("domain").(string)
//	projectName := d.Get("project_name").(string)
//	memberId := d.Get("member_id").(int)
//	member, _, err := c.ProjectMemberService.GetProjectMember(domain, projectName, memberId)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	err = util.ApiProjectMemberToResourceData(domain, projectName, member, d, false)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	return diags
//}

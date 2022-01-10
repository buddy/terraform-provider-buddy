package resource

import (
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
	"strconv"
	"strings"
)

func ProjectMember() *schema.Resource {
	return &schema.Resource{
		Description: "Manage a workspace project member permission\n\n" +
			"Workspace administrator rights are required\n\n" +
			"Token scope required: `WORKSPACE`",
		CreateContext: createContextProjectMember,
		ReadContext:   readContextProjectMember,
		UpdateContext: updateContextProjectMember,
		DeleteContext: deleteContextProjectMember,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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
				ForceNew:     true,
				ValidateFunc: util.ValidateDomain,
			},
			"project_name": {
				Description: "The project's name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"member_id": {
				Description: "The member's ID",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"permission_id": {
				Description: "The permission's ID",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"html_url": {
				Description: "The member's URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The member's name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email": {
				Description: "The member's email",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"avatar_url": {
				Description: "The member's avatar URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"admin": {
				Description: "Is the member a workspace administrator",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"workspace_owner": {
				Description: "Is the member the workspace owner",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"permission": {
				Description: "The member's permission in the project",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"html_url": {
							Description: "The permission's URL",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"permission_id": {
							Description: "The permission's ID",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"name": {
							Description: "The permission's name",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"type": {
							Description: "The permission's type",
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
							Description: "The permission's access level to sandboxes",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func deleteContextProjectMember(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	domain, projectName, mid, err := util.DecomposeTripleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	memberId, err := strconv.Atoi(mid)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = c.ProjectMemberService.DeleteProjectMember(domain, projectName, memberId)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func updateContextProjectMember(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	domain := d.Get("domain").(string)
	projectName := d.Get("project_name").(string)
	memberId := d.Get("member_id").(int)
	_, _, err := c.ProjectMemberService.UpdateProjectMember(domain, projectName, memberId, &api.ProjectMemberOperationOptions{
		PermissionSet: &api.ProjectMemberOperationOptions{
			Id: util.InterfaceIntToPointer(d.Get("permission_id")),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(util.ComposeTripleId(domain, projectName, strconv.Itoa(memberId)))
	return readContextProjectMember(ctx, d, meta)
}

func readContextProjectMember(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	domain, projectName, mid, err := util.DecomposeTripleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	memberId, err := strconv.Atoi(mid)
	if err != nil {
		return diag.FromErr(err)
	}
	m, resp, err := c.ProjectMemberService.GetProjectMember(domain, projectName, memberId)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	err = util.ApiProjectMemberToResourceData(domain, projectName, m, d, true)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func createContextProjectMember(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	domain := d.Get("domain").(string)
	projectName := d.Get("project_name").(string)
	member, _, err := c.ProjectMemberService.CreateProjectMember(domain, projectName, &api.ProjectMemberOperationOptions{
		Id: util.InterfaceIntToPointer(d.Get("member_id")),
		PermissionSet: &api.ProjectMemberOperationOptions{
			Id: util.InterfaceIntToPointer(d.Get("permission_id")),
		},
	})
	if err != nil {
		if strings.Contains(err.Error(), "This user is already assigned to the project") {
			return updateContextProjectMember(ctx, d, meta)
		}
		return diag.FromErr(err)
	}
	d.SetId(util.ComposeTripleId(domain, projectName, strconv.Itoa(member.Id)))
	return readContextProjectMember(ctx, d, meta)
}

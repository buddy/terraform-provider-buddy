package resource

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"strings"
)

func ProjectGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Manage a workspace project group permission\n\n" +
			"Workspace administrator rights are required\n\n" +
			"Token scope required: `WORKSPACE`",
		CreateContext: createContextProjectGroup,
		ReadContext:   readContextProjectGroup,
		UpdateContext: updateContextProjectGroup,
		DeleteContext: deleteContextProjectGroup,
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
			"group_id": {
				Description: "The group's ID",
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
				Description: "The group's URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The group's name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"permission": {
				Description: "The group's permission in the project",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"html_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"permission_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
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
						"sandbox_access_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func deleteContextProjectGroup(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	var diags diag.Diagnostics
	domain, projectName, gid, err := util.DecomposeTripleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	groupId, err := strconv.Atoi(gid)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = c.ProjectGroupService.DeleteProjectGroup(domain, projectName, groupId)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func updateContextProjectGroup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	domain := d.Get("domain").(string)
	projectName := d.Get("project_name").(string)
	groupId := d.Get("group_id").(int)
	_, _, err := c.ProjectGroupService.UpdateProjectGroup(domain, projectName, groupId, &buddy.ProjectGroupOps{
		PermissionSet: &buddy.ProjectGroupOps{
			Id: util.InterfaceIntToPointer(d.Get("permission_id")),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(util.ComposeTripleId(domain, projectName, strconv.Itoa(groupId)))
	return readContextProjectGroup(ctx, d, meta)
}

func readContextProjectGroup(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	var diags diag.Diagnostics
	domain, projectName, gid, err := util.DecomposeTripleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	groupId, err := strconv.Atoi(gid)
	if err != nil {
		return diag.FromErr(err)
	}
	g, resp, err := c.ProjectGroupService.GetProjectGroup(domain, projectName, groupId)
	if err != nil {
		if util.IsResourceNotFound(resp, err) {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	err = util.ApiProjectGroupToResourceData(domain, projectName, g, d, true)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func createContextProjectGroup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	domain := d.Get("domain").(string)
	projectName := d.Get("project_name").(string)
	group, _, err := c.ProjectGroupService.CreateProjectGroup(domain, projectName, &buddy.ProjectGroupOps{
		Id: util.InterfaceIntToPointer(d.Get("group_id")),
		PermissionSet: &buddy.ProjectGroupOps{
			Id: util.InterfaceIntToPointer(d.Get("permission_id")),
		},
	})
	if err != nil {
		if strings.Contains(err.Error(), "This group is already assigned to the project") {
			return updateContextProjectGroup(ctx, d, meta)
		}
		return diag.FromErr(err)
	}
	d.SetId(util.ComposeTripleId(domain, projectName, strconv.Itoa(group.Id)))
	return readContextProjectGroup(ctx, d, meta)
}

package resource

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strconv"
)

func Permission() *schema.Resource {
	return &schema.Resource{
		Description: "Create and manage a workspace permission (role)\n\n" +
			"Workspace administrator rights are required\n\n" +
			"Token scope required: `WORKSPACE`",
		CreateContext: createContextPermission,
		ReadContext:   readContextPermission,
		UpdateContext: updateContextPermission,
		DeleteContext: deleteContextPermission,
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
			"name": {
				Description: "The permission's name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"pipeline_access_level": {
				Description: "The permission's access level to pipelines. Allowed: `DENIED`, `READ_ONLY`, `RUN_ONLY`, `READ_WRITE`",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validation.StringInSlice([]string{
					buddy.PermissionAccessLevelDenied,
					buddy.PermissionAccessLevelReadOnly,
					buddy.PermissionAccessLevelRunOnly,
					buddy.PermissionAccessLevelReadWrite,
				}, false),
			},
			"repository_access_level": {
				Description: "The permission's access level to repository. Allowed: `READ_ONLY`, `READ_WRITE`, `MANAGE`",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validation.StringInSlice([]string{
					buddy.PermissionAccessLevelReadOnly,
					buddy.PermissionAccessLevelReadWrite,
					buddy.PermissionAccessLevelManage,
				}, false),
			},
			"sandbox_access_level": {
				Description: "The permission's access level to sandboxes. Allowed: `DENIED`, `READ_ONLY`, `READ_WRITE`",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validation.StringInSlice([]string{
					buddy.PermissionAccessLevelDenied,
					buddy.PermissionAccessLevelReadOnly,
					buddy.PermissionAccessLevelReadWrite,
				}, false),
			},
			"project_team_access_level": {
				Description: "The permission's access level to team. Allowed: `READ_ONLY`, `MANAGE`",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ValidateFunc: validation.StringInSlice([]string{
					buddy.PermissionAccessLevelReadOnly,
					buddy.PermissionAccessLevelManage,
				}, false),
			},
			"permission_id": {
				Description: "The permission's ID",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"description": {
				Description: "The permission's description",
				Type:        schema.TypeString,
				Optional:    true,
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
	}
}

func deleteContextPermission(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	var diags diag.Diagnostics
	domain, pid, err := util.DecomposeDoubleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	permissionId, err := strconv.Atoi(pid)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = c.PermissionService.Delete(domain, permissionId)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func updateContextPermission(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	domain, pid, err := util.DecomposeDoubleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	permissionId, err := strconv.Atoi(pid)
	if err != nil {
		return diag.FromErr(err)
	}
	opt := buddy.PermissionOps{
		Name:                   util.InterfaceStringToPointer(d.Get("name")),
		PipelineAccessLevel:    util.InterfaceStringToPointer(d.Get("pipeline_access_level")),
		RepositoryAccessLevel:  util.InterfaceStringToPointer(d.Get("repository_access_level")),
		SandboxAccessLevel:     util.InterfaceStringToPointer(d.Get("sandbox_access_level")),
		ProjectTeamAccessLevel: util.InterfaceStringToPointer(d.Get("project_team_access_level")),
		Description:            util.InterfaceStringToPointer(d.Get("description")),
	}
	_, _, err = c.PermissionService.Update(domain, permissionId, &opt)
	if err != nil {
		return diag.FromErr(err)
	}
	return readContextPermission(ctx, d, meta)
}

func readContextPermission(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	var diags diag.Diagnostics
	domain, pid, err := util.DecomposeDoubleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	permissionId, err := strconv.Atoi(pid)
	if err != nil {
		return diag.FromErr(err)
	}
	permission, resp, err := c.PermissionService.Get(domain, permissionId)
	if err != nil {
		if util.IsResourceNotFound(resp, err) {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	err = util.ApiPermissionToResourceData(domain, permission, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func createContextPermission(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	domain := d.Get("domain").(string)
	opt := buddy.PermissionOps{
		Name:                  util.InterfaceStringToPointer(d.Get("name")),
		PipelineAccessLevel:   util.InterfaceStringToPointer(d.Get("pipeline_access_level")),
		RepositoryAccessLevel: util.InterfaceStringToPointer(d.Get("repository_access_level")),
		SandboxAccessLevel:    util.InterfaceStringToPointer(d.Get("sandbox_access_level")),
	}
	if projectTeamAccessLevel, ok := d.GetOk("project_team_access_level"); ok {
		opt.ProjectTeamAccessLevel = util.InterfaceStringToPointer(projectTeamAccessLevel)
	}
	if description, ok := d.GetOk("description"); ok {
		opt.Description = util.InterfaceStringToPointer(description)
	}
	permission, _, err := c.PermissionService.Create(domain, &opt)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(util.ComposeDoubleId(domain, strconv.Itoa(permission.Id)))
	return readContextPermission(ctx, d, meta)
}

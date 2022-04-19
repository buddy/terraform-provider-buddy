package resource

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func Group() *schema.Resource {
	return &schema.Resource{
		Description: "Create and manage a user's group\n\n" +
			"Workspace administrator rights are required\n\n" +
			"Token scope required: `WORKSPACE`",
		CreateContext: createContextGroup,
		ReadContext:   readContextGroup,
		UpdateContext: updateContextGroup,
		DeleteContext: deleteContextGroup,
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
				Description: "The group's name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"auto_assign_to_new_projects": {
				Description: "Defines whether or not to automatically assign group to new projects",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"auto_assign_permission_set_id": {
				Description: "The permission's ID with which the group will be assigned to new projects",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"group_id": {
				Description: "The group's ID",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"html_url": {
				Description: "The group's URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "The group's description",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func createContextGroup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	domain := d.Get("domain").(string)
	assignToProjects := d.Get("auto_assign_to_new_projects").(bool)
	opt := buddy.GroupOps{
		Name: util.InterfaceStringToPointer(d.Get("name")),
	}
	if description, ok := d.GetOk("description"); ok {
		opt.Description = util.InterfaceStringToPointer(description)
	}
	if assignToProjects {
		permissionId := d.Get("auto_assign_permission_set_id").(int)
		opt.AutoAssignToNewProjects = &assignToProjects
		opt.AutoAssignPermissionSetId = &permissionId
	}
	group, _, err := c.GroupService.Create(domain, &opt)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(util.ComposeDoubleId(domain, strconv.Itoa(group.Id)))
	return readContextGroup(ctx, d, meta)
}

func readContextGroup(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	var diags diag.Diagnostics
	domain, gid, err := util.DecomposeDoubleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	groupId, err := strconv.Atoi(gid)
	if err != nil {
		return diag.FromErr(err)
	}
	g, resp, err := c.GroupService.Get(domain, groupId)
	if err != nil {
		if util.IsResourceNotFound(resp, err) {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	err = util.ApiGroupToResourceData(domain, g, d, false)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func updateContextGroup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	domain, gid, err := util.DecomposeDoubleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	groupId, err := strconv.Atoi(gid)
	if err != nil {
		return diag.FromErr(err)
	}
	assignToProjects := d.Get("auto_assign_to_new_projects").(bool)
	u := buddy.GroupOps{
		Name:                    util.InterfaceStringToPointer(d.Get("name")),
		Description:             util.InterfaceStringToPointer(d.Get("description")),
		AutoAssignToNewProjects: &assignToProjects,
	}
	if assignToProjects {
		u.AutoAssignPermissionSetId = util.InterfaceIntToPointer(d.Get("auto_assign_permission_set_id"))
	}
	_, _, err = c.GroupService.Update(domain, groupId, &u)
	if err != nil {
		return diag.FromErr(err)
	}
	return readContextGroup(ctx, d, meta)
}

func deleteContextGroup(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	var diags diag.Diagnostics
	domain, gid, err := util.DecomposeDoubleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	groupId, err := strconv.Atoi(gid)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = c.GroupService.Delete(domain, groupId)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

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

func GroupMember() *schema.Resource {
	return &schema.Resource{
		Description: "Create and manage a workspace group member\n\n" +
			"Workspace administrator rights are required\n\n" +
			"Token scope required: `WORKSPACE`",
		CreateContext: createContextGroupMember,
		ReadContext:   readContextGroupMember,
		UpdateContext: updateContexGroupMember,
		DeleteContext: deleteContextGroupMember,
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
			"group_id": {
				Description: "The group's ID",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"member_id": {
				Description: "The member's ID",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"status": {
				Description: "The member's status. Allowed: `MEMBER`, `MANAGER`",
				Type:        schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					buddy.GroupMemberStatusMember,
					buddy.GroupMemberStatusManager,
				}, false),
				Optional: true,
				Computed: true,
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
		},
	}
}

func deleteContextGroupMember(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	var diags diag.Diagnostics
	domain, gid, mid, err := util.DecomposeTripleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	groupId, err := strconv.Atoi(gid)
	if err != nil {
		return diag.FromErr(err)
	}
	memberId, err := strconv.Atoi(mid)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = c.GroupService.DeleteGroupMember(domain, groupId, memberId)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func readContextGroupMember(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	var diags diag.Diagnostics
	domain, gid, mid, err := util.DecomposeTripleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	groupId, err := strconv.Atoi(gid)
	if err != nil {
		return diag.FromErr(err)
	}
	memberId, err := strconv.Atoi(mid)
	if err != nil {
		return diag.FromErr(err)
	}
	m, resp, err := c.GroupService.GetGroupMember(domain, groupId, memberId)
	if err != nil {
		if util.IsResourceNotFound(resp, err) {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	err = util.ApiGroupMemberToResourceData(domain, groupId, m, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func updateContexGroupMember(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	domain, gid, mid, err := util.DecomposeTripleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	groupId, err := strconv.Atoi(gid)
	if err != nil {
		return diag.FromErr(err)
	}
	memberId, err := strconv.Atoi(mid)
	if err != nil {
		return diag.FromErr(err)
	}
	status := buddy.GroupMemberStatusMember
	if s, ok := d.GetOk("status"); ok {
		status = s.(string)
	}
	_, _, err = c.GroupService.UpdateGroupMember(domain, groupId, memberId, status)
	if err != nil {
		return diag.FromErr(err)
	}
	return readContextGroupMember(ctx, d, meta)
}

func createContextGroupMember(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	domain := d.Get("domain").(string)
	groupId := d.Get("group_id").(int)
	ops := buddy.GroupMemberOps{
		Id: util.InterfaceIntToPointer(d.Get("member_id")),
	}
	if status, ok := d.GetOk("status"); ok {
		ops.Status = util.InterfaceStringToPointer(status)
	}
	member, _, err := c.GroupService.AddGroupMember(domain, groupId, &ops)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(util.ComposeTripleId(domain, strconv.Itoa(groupId), strconv.Itoa(member.Id)))
	return readContextGroupMember(ctx, d, meta)
}

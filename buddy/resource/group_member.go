package resource

import (
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
	"strconv"
)

func GroupMember() *schema.Resource {
	return &schema.Resource{
		Description: "`buddy_group_member` allows you to manage Buddy group membership.\n\n" +
			"You will need admin access in workspace for this resource to work.\n\n" +
			"Required scope for your token: `WORKSPACE`",
		CreateContext: createContextGroupMember,
		ReadContext:   readContextGroupMember,
		DeleteContext: deleteContextGroupMember,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Compound id of the group member",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"domain": {
				Description:  "Workspace domain in which the group membership will be managed",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: util.ValidateDomain,
			},
			"group_id": {
				Description: "Id of the group",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"member_id": {
				Description: "Id of the member",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"html_url": {
				Description: "Url of the member",
				Type:        schema.TypeString,
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
	}
}

func deleteContextGroupMember(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
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
	c := meta.(*api.Client)
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
		if resp.StatusCode == http.StatusNotFound {
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

func createContextGroupMember(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	domain := d.Get("domain").(string)
	groupId := d.Get("group_id").(int)
	member, _, err := c.GroupService.AddGroupMember(domain, groupId, &api.GroupMemberOperationOptions{
		Id: util.InterfaceIntToPointer(d.Get("member_id")),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(util.ComposeTripleId(domain, strconv.Itoa(groupId), strconv.Itoa(member.Id)))
	return readContextGroupMember(ctx, d, meta)
}

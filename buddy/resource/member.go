package resource

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func Member() *schema.Resource {
	return &schema.Resource{
		Description: "Create and manage a workspace member\n\n" +
			"Workspace administrator rights are required\n\n" +
			"Token scope required: `WORKSPACE`",
		CreateContext: createContextMember,
		ReadContext:   readContextMember,
		UpdateContext: updateContextMember,
		DeleteContext: deleteContextMember,
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
			"email": {
				Description:  "The member's email",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: util.ValidateEmail,
			},
			"admin": {
				Description: "Is the member a workspace administrator",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"name": {
				Description: "The member's name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"member_id": {
				Description: "The member's ID",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"html_url": {
				Description: "The member's URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"avatar_url": {
				Description: "The member's avatar URL",
				Type:        schema.TypeString,
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

func deleteContextMember(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	var diags diag.Diagnostics
	domain, mid, err := util.DecomposeDoubleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	memberId, err := strconv.Atoi(mid)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = c.MemberService.Delete(domain, memberId)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func updateContextMember(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("admin") {
		c := meta.(*buddy.Client)
		domain, mid, err := util.DecomposeDoubleId(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		memberId, err := strconv.Atoi(mid)
		if err != nil {
			return diag.FromErr(err)
		}
		_, _, err = c.MemberService.UpdateAdmin(domain, memberId, &buddy.MemberAdminOps{
			Admin: util.InterfaceBoolToPointer(d.Get("admin")),
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return readContextMember(ctx, d, meta)
}

func readContextMember(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	var diags diag.Diagnostics
	domain, mid, err := util.DecomposeDoubleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	memberId, err := strconv.Atoi(mid)
	if err != nil {
		return diag.FromErr(err)
	}
	member, resp, err := c.MemberService.Get(domain, memberId)
	if err != nil {
		if util.IsResourceNotFound(resp, err) {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	err = util.ApiMemberToResourceData(domain, member, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func createContextMember(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	domain := d.Get("domain").(string)
	admin := d.Get("admin").(bool)
	opt := buddy.MemberOps{
		Email: util.InterfaceStringToPointer(d.Get("email")),
	}
	member, _, err := c.MemberService.Create(domain, &opt)
	if err != nil {
		return diag.FromErr(err)
	}
	if admin {
		_, _, err := c.MemberService.UpdateAdmin(domain, member.Id, &buddy.MemberAdminOps{
			Admin: util.InterfaceBoolToPointer(true),
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(util.ComposeDoubleId(domain, strconv.Itoa(member.Id)))
	return readContextMember(ctx, d, meta)
}

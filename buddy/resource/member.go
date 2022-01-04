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

func Member() *schema.Resource {
	return &schema.Resource{
		Description: "`buddy_member` allows you to create and manage Buddy workspace member.\n\n" +
			"You will need admin access in workspace for this resource to work.\n\n" +
			"Required scopes for your token: `WORKSPACE`",
		CreateContext: createContextMember,
		ReadContext:   readContextMember,
		UpdateContext: updateContextMember,
		DeleteContext: deleteContextMember,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Compound id of the member",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"domain": {
				Description:  "Workspace domain in which the member will be created",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: util.ValidateDomain,
			},
			"email": {
				Description:  "Email of the member",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: util.ValidateEmail,
			},
			"admin": {
				Description: "Should member has administrator privileges",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"name": {
				Description: "Name of the member",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"member_id": {
				Description: "Id of the member",
				Type:        schema.TypeInt,
				Computed:    true,
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

func deleteContextMember(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
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
		c := meta.(*api.Client)
		domain, mid, err := util.DecomposeDoubleId(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		memberId, err := strconv.Atoi(mid)
		if err != nil {
			return diag.FromErr(err)
		}
		_, _, err = c.MemberService.UpdateAdmin(domain, memberId, &api.MemberAdminOperationOptions{
			Admin: util.InterfaceBoolToPointer(d.Get("admin")),
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return readContextMember(ctx, d, meta)
}

func readContextMember(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
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
		if resp.StatusCode == http.StatusNotFound {
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
	c := meta.(*api.Client)
	domain := d.Get("domain").(string)
	admin := d.Get("admin").(bool)
	opt := api.MemberOperationOptions{
		Email: util.InterfaceStringToPointer(d.Get("email")),
	}
	member, _, err := c.MemberService.Create(domain, &opt)
	if err != nil {
		return diag.FromErr(err)
	}
	if admin {
		_, _, err := c.MemberService.UpdateAdmin(domain, member.Id, &api.MemberAdminOperationOptions{
			Admin: util.InterfaceBoolToPointer(true),
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(util.ComposeDoubleId(domain, strconv.Itoa(member.Id)))
	return readContextMember(ctx, d, meta)
}

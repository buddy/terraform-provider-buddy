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

func Group() *schema.Resource {
	return &schema.Resource{
		Description: "Create and manage a workspace group\n\n" +
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
				Computed:    true,
			},
		},
	}
}

func createContextGroup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	domain := d.Get("domain").(string)
	opt := api.GroupOperationOptions{
		Name: util.InterfaceStringToPointer(d.Get("name")),
	}
	if description, ok := d.GetOk("description"); ok {
		opt.Description = util.InterfaceStringToPointer(description)
	}
	group, _, err := c.GroupService.Create(domain, &opt)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(util.ComposeDoubleId(domain, strconv.Itoa(group.Id)))
	return readContextGroup(ctx, d, meta)
}

func readContextGroup(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
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
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	err = util.ApiGroupToResourceData(domain, g, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func updateContextGroup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	domain, gid, err := util.DecomposeDoubleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	groupId, err := strconv.Atoi(gid)
	if err != nil {
		return diag.FromErr(err)
	}
	u := api.GroupOperationOptions{
		Name: util.InterfaceStringToPointer(d.Get("name")),
	}
	if d.HasChange("description") {
		u.Description = util.InterfaceStringToPointer(d.Get("description"))
	}
	_, _, err = c.GroupService.Update(domain, groupId, &u)
	if err != nil {
		return diag.FromErr(err)
	}
	return readContextGroup(ctx, d, meta)
}

func deleteContextGroup(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
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

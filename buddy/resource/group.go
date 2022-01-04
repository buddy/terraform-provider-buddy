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
		Description: "`buddy_group` allows you to create and manage Buddy groups.\n\n" +
			"You will need admin access in workspace for this resource to work.\n\n" +
			"Required scope for your token: `WORKSPACE`",
		CreateContext: createContextGroup,
		ReadContext:   readContextGroup,
		UpdateContext: updateContextGroup,
		DeleteContext: deleteContextGroup,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Compound id of the group",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"domain": {
				Description:  "Workspace domain in which the group will be created",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: util.ValidateDomain,
			},
			"name": {
				Description: "Name of the group",
				Type:        schema.TypeString,
				Required:    true,
			},
			"group_id": {
				Description: "Id of the group",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"html_url": {
				Description: "Url of the group",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Description of the group",
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

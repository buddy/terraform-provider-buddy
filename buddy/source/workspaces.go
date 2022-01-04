package source

import (
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
)

func Workspaces() *schema.Resource {
	return &schema.Resource{
		Description: "`buddy_workspaces` data source allows you to get list of workspaces and filter them by name or domain\n\n" +
			"Required scopes for your token: `WORKSPACE`",
		ReadContext: readContextWorkspaces,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Compound id of the resource",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"domain_regex": {
				Description:  "Regular expression to match domain of the workspace",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"name_regex": {
				Description:  "Regular expression to match name of the workspace",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"workspaces": {
				Description: "List of workspaces",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"html_url": {
							Description: "Url of the workspace",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"workspace_id": {
							Description: "Id of the workspace",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"name": {
							Description: "Name of the workspace",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"domain": {
							Description: "Domain of the workspace",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func readContextWorkspaces(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	var nameRegex *regexp.Regexp
	var domainRegex *regexp.Regexp
	workspaces, _, err := c.WorkspaceService.GetList()
	if err != nil {
		return diag.FromErr(err)
	}
	var result []interface{}
	if name, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(name.(string))
	}
	if domain, ok := d.GetOk("domain_regex"); ok {
		domainRegex = regexp.MustCompile(domain.(string))
	}
	for _, w := range workspaces.Workspaces {
		if nameRegex != nil && !nameRegex.MatchString(w.Name) {
			continue
		}
		if domainRegex != nil && !domainRegex.MatchString(w.Domain) {
			continue
		}
		result = append(result, util.ApiShortWorkspaceToMap(w))
	}
	d.SetId(util.UniqueString())
	err = d.Set("workspaces", result)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

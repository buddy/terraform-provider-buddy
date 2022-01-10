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

func Projects() *schema.Resource {
	return &schema.Resource{
		Description: "List projects and optionally filter them by membership, status, name or display_name\n\n" +
			"Required scopes for your token: `WORKSPACE`",
		ReadContext: readContextProjects,
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
				ValidateFunc: util.ValidateDomain,
			},
			"name_regex": {
				Description:  "The project's name regular expression to match",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"display_name_regex": {
				Description:  "The project's display name regular expression to match",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"membership": {
				Description: "In case a user is a workspace administrators, if true lists projects the user belongs to",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"status": {
				Description: "Filter projects by status (`ACTIVE`, `CLOSED`)",
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"ACTIVE",
					"CLOSED",
				}, false),
			},
			"projects": {
				Description: "List of projects",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"html_url": {
							Description: "The project's URL",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "The project's unique name ID",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"display_name": {
							Description: "The project's display name",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"status": {
							Description: "The project's status",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func readContextProjects(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	var nameRegex *regexp.Regexp
	var displayNameRegex *regexp.Regexp
	opt := api.QueryProjectList{}
	if membership, ok := d.GetOk("membership"); ok {
		opt.Membership = membership.(bool)
	}
	if status, ok := d.GetOk("status"); ok {
		opt.Status = status.(string)
	}
	domain := d.Get("domain").(string)
	projects, _, err := c.ProjectService.GetList(domain, &opt)
	if err != nil {
		return diag.FromErr(err)
	}
	var result []interface{}
	if name, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(name.(string))
	}
	if displayName, ok := d.GetOk("display_name_regex"); ok {
		displayNameRegex = regexp.MustCompile(displayName.(string))
	}
	for _, p := range projects.Projects {
		if nameRegex != nil && !nameRegex.MatchString(p.Name) {
			continue
		}
		if displayNameRegex != nil && !displayNameRegex.MatchString(p.DisplayName) {
			continue
		}
		result = append(result, util.ApiShortProjectToMap(p))
	}
	d.SetId(util.UniqueString())
	err = d.Set("projects", result)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

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

func Pipelines() *schema.Resource {
	return &schema.Resource{
		Description: "List pipelines and optionally filter them by name\n\n" +
			"Token scopes required: `WORKSPACE`, `EXECUTION_INFO`",
		ReadContext: readContextPipelines,
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
			"project_name": {
				Description: "The project's name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name_regex": {
				Description:  "The pipeline's name regular expression to match",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"pipelines": {
				Description: "List of pipelines",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pipeline_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"html_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_execution_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_execution_revision": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"refs": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"event": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"refs": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func readContextPipelines(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	var nameRegex *regexp.Regexp
	domain := d.Get("domain").(string)
	projectName := d.Get("project_name").(string)
	pipelines, _, err := c.PipelineService.GetList(domain, projectName)
	if err != nil {
		return diag.FromErr(err)
	}
	var result []interface{}
	if name, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(name.(string))
	}
	for _, p := range pipelines.Pipelines {
		if nameRegex != nil && !nameRegex.MatchString(p.Name) {
			continue
		}
		result = append(result, util.ApiShortPipelineToMap(p))
	}
	d.SetId(util.UniqueString())
	err = d.Set("pipelines", result)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

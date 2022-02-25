package source

import (
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Pipeline() *schema.Resource {
	return &schema.Resource{
		Description: "Get pipeline by name or pipeline ID\n\n" +
			"Token scopes required: `WORKSPACE`, `EXECUTION_INFO`",
		ReadContext: readContextPipeline,
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
			"name": {
				Description: "The pipeline's name",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ExactlyOneOf: []string{
					"pipeline_id",
					"name",
				},
			},
			"priority": {
				Description: "The pipeline's priority",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"pipeline_id": {
				Description: "The pipeline's ID",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ExactlyOneOf: []string{
					"pipeline_id",
					"name",
				},
			},
			"html_url": {
				Description: "The pipeline's URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"on": {
				Description: "The pipeline's trigger mode",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_execution_status": {
				Description: "The pipeline's last run status",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_execution_revision": {
				Description: "The pipeline's last run revision",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"refs": {
				Description: "The pipeline's list of refs",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": {
				Description: "The pipeline's list of tags. Only for `Buddy Enterprise`",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"event": {
				Description: "The pipeline's list of events",
				Type:        schema.TypeList,
				Computed:    true,
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
			"definition_source": {
				Description: "The pipeline's definition source",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"remote_project_name": {
				Description: "The pipeline's remote definition project name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"remote_branch": {
				Description: "The pipeline's remote definition branch name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"remote_path": {
				Description: "The pipeline's remote definition path",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"remote_parameter": {
				Description: "The pipeline's remote definition parameters",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func readContextPipeline(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	var pipeline *api.Pipeline
	var err error
	domain := d.Get("domain").(string)
	projectName := d.Get("project_name").(string)
	if pipelineId, ok := d.GetOk("pipeline_id"); ok {
		pipeline, _, err = c.PipelineService.Get(domain, projectName, pipelineId.(int))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		name := d.Get("name").(string)
		pipelines, _, err := c.PipelineService.GetList(domain, projectName)
		if err != nil {
			return diag.FromErr(err)
		}
		for _, p := range pipelines.Pipelines {
			if p.Name == name {
				pipeline = p
				break
			}
		}
		if pipeline == nil {
			return diag.Errorf("Pipeline not found")
		}
	}
	err = util.ApiPipelineToResourceData(domain, projectName, pipeline, d, true)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

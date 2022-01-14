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

func VariablesSshKeys() *schema.Resource {
	return &schema.Resource{
		Description: "List variables ssh keys and optionally filter them by key, project_name, pipeline_id or action_id\n\n" +
			"Token scope required: `WORKSPACE`, `VARIABLE_INFO`",
		ReadContext: readContextVariablesSshKeys,
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
			"key_regex": {
				Description:  "The variable's key regular expression to match",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"project_name": {
				Description: "The variable's project name",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"pipeline_id": {
				Description: "The variable's pipeline ID",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"action_id": {
				Description: "The variable's action ID",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"variables": {
				Description: "List of variables",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Description: "The variable's name",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"encrypted": {
							Description: "The variable ssh key is always encrypted",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"settable": {
							Description: "The variable ssh key is not changeable",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"description": {
							Description: "The variable's description",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"value": {
							Description: "The variable's encrypted value",
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
						},
						"variable_id": {
							Description: "The variable's ID",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"file_name": {
							Description: "The variable's file name",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"file_place": {
							Description: "Available values: `CONTAINER`, `NONE`",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"file_path": {
							Description: "The variable's file place",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"file_chmod": {
							Description: "The variable's file permission in the container",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"checksum": {
							Description: "The variable's checksum",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"key_fingerprint": {
							Description: "The variable's fingherprint",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"public_value": {
							Description: "The variable's public key",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func readContextVariablesSshKeys(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	var keyRegex *regexp.Regexp
	domain := d.Get("domain").(string)
	opt := api.VariableGetListQuery{}
	if projectName, ok := d.GetOk("project_name"); ok {
		opt.ProjectName = projectName.(string)
	}
	if pipelineId, ok := d.GetOk("pipeline_id"); ok {
		opt.PipelineId = pipelineId.(int)
	}
	if actionId, ok := d.GetOk("action_id"); ok {
		opt.ActionId = actionId.(int)
	}
	variables, _, err := c.VariableService.GetList(domain, &opt)
	if err != nil {
		return diag.FromErr(err)
	}
	if key, ok := d.GetOk("key_regex"); ok {
		keyRegex = regexp.MustCompile(key.(string))
	}
	var result []interface{}
	for _, v := range variables.Variables {
		if v.Type != api.VariableTypeSshKey {
			continue
		}
		if keyRegex != nil && !keyRegex.MatchString(v.Key) {
			continue
		}
		result = append(result, util.ApiShortVariableSshKeyToMap(v))
	}
	d.SetId(util.UniqueString())
	err = d.Set("variables", result)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

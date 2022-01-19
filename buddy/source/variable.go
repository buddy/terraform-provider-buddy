package source

import (
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Variable() *schema.Resource {
	return &schema.Resource{
		Description: "Get variable by key or variable ID\n\n" +
			"Token scope required: `WORKSPACE`, `VARIABLE_INFO`",
		ReadContext: readContextVariable,
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
			"key": {
				Description: "The variable's name",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ExactlyOneOf: []string{
					"variable_id",
					"key",
				},
			},
			"project_name": {
				Description: "The variable's project name",
				Type:        schema.TypeString,
				Optional:    true,
				RequiredWith: []string{
					"key",
				},
			},
			"pipeline_id": {
				Description: "The variable's pipeline ID",
				Type:        schema.TypeInt,
				Optional:    true,
				RequiredWith: []string{
					"key",
				},
			},
			"action_id": {
				Description: "The variable's action ID",
				Type:        schema.TypeInt,
				Optional:    true,
				RequiredWith: []string{
					"key",
				},
			},
			"encrypted": {
				Description: "Is the variable's value encrypted",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"settable": {
				Description: "Is the variable's value changeable",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"description": {
				Description: "The variable's description",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"value": {
				Description: "The variable's value. Encrypted if **encrypted** == true",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			"variable_id": {
				Description: "The variable's ID",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ExactlyOneOf: []string{
					"variable_id",
					"key",
				},
			},
		},
	}
}

func readContextVariable(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	var variable *api.Variable
	var err error
	domain := d.Get("domain").(string)
	if variableId, ok := d.GetOk("variable_id"); ok {
		variable, _, err = c.VariableService.Get(domain, variableId.(int))
		if err != nil {
			return diag.FromErr(err)
		}
		if variable.Type != api.VariableTypeVar {
			return diag.Errorf("Variable not found")
		}
	} else {
		key := d.Get("key").(string)
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
		for _, v := range variables.Variables {
			if v.Type != api.VariableTypeVar {
				continue
			}
			if v.Key == key {
				variable = v
				break
			}
		}
		if variable == nil {
			return diag.Errorf("Variable not found")
		}
	}
	err = util.ApiVariableToResourceData(domain, variable, d, false)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

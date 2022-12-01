package source

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func VariableSshKey() *schema.Resource {
	return &schema.Resource{
		Description: "Get variables of SSH key type by key or variable ID\n\n" +
			"Token scope required: `WORKSPACE`, `VARIABLE_INFO`",
		ReadContext: readContextVariableSshKey,
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
				Description: "Get only from provided project",
				Type:        schema.TypeString,
				Optional:    true,
				RequiredWith: []string{
					"key",
				},
			},
			"pipeline_id": {
				Description: "Get only from provided pipeline",
				Type:        schema.TypeInt,
				Optional:    true,
				RequiredWith: []string{
					"key",
				},
			},
			"action_id": {
				Description: "Get only from provided action",
				Type:        schema.TypeInt,
				Optional:    true,
				RequiredWith: []string{
					"key",
				},
			},
			"value": {
				Description: "The variable's value, always encrypted for buddy_variable_ssh_key",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			"file_place": {
				Description: "Should the variable's be copied to an action's container in **file_path** (`CONTAINER`, `NONE`)",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"file_path": {
				Description: "The variable's path in the action's container",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"file_chmod": {
				Description: "The variable's file permission in an action's container",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"settable": {
				Description: "Is the variable's value changeable, always false for buddy_variable_ssh_key",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"description": {
				Description: "The variable's description",
				Type:        schema.TypeString,
				Computed:    true,
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
			"encrypted": {
				Description: "Is the variable's value encrypted, always true for buddy_variable_ssh_key",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"checksum": {
				Description: "The variable's checksum",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"key_fingerprint": {
				Description: "The variable's fingerprint",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"public_value": {
				Description: "The variable's public key",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func readContextVariableSshKey(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	var diags diag.Diagnostics
	var variable *buddy.Variable
	var err error
	domain := d.Get("domain").(string)
	if variableId, ok := d.GetOk("variable_id"); ok {
		variable, _, err = c.VariableService.Get(domain, variableId.(int))
		if err != nil {
			return diag.FromErr(err)
		}
		if variable.Type != buddy.VariableTypeSshKey {
			return diag.Errorf("Variable not found")
		}
	} else {
		key := d.Get("key").(string)
		opt := buddy.VariableGetListQuery{}
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
			if v.Type != buddy.VariableTypeSshKey {
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
	err = util.ApiVariableSshKeyToResourceData(domain, variable, d, false)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

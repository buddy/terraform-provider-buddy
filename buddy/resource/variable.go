package resource

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func Variable() *schema.Resource {
	return &schema.Resource{
		Description: "Create and manage a variable\n\n" +
			"Workspace administrator rights are required\n\n" +
			"Token scopes required: `WORKSPACE`, `VARIABLE_ADD`, `VARIABLE_MANAGE`, `VARIABLE_INFO`",
		CreateContext: createContextVariable,
		ReadContext:   readContextVariable,
		UpdateContext: updateContextVariable,
		DeleteContext: deleteContextVariable,
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
			"key": {
				Description: "The variable's name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"value": {
				Description: "The variable's value",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"encrypted": {
				Description: "Is the variable's value encrypted",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"project_name": {
				Description: "The variable's project name",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"pipeline_id": {
				Description: "The variable's pipeline ID",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"action_id": {
				Description: "The variable's action ID",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"settable": {
				Description: "Is the variable's value changeable",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"description": {
				Description: "The variable's description",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"value_processed": {
				Description: "The variable's processed value. Encrypted if **encrypted** == true",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			"variable_id": {
				Description: "The variable's ID",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func deleteContextVariable(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	var diags diag.Diagnostics
	domain, vid, err := util.DecomposeDoubleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	variableId, err := strconv.Atoi(vid)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = c.VariableService.Delete(domain, variableId)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func updateContextVariable(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	domain, vid, err := util.DecomposeDoubleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	variableId, err := strconv.Atoi(vid)
	if err != nil {
		return diag.FromErr(err)
	}
	opt := buddy.VariableOps{
		Value: util.InterfaceStringToPointer(d.Get("value")),
	}
	if d.HasChange("encrypted") {
		opt.Encrypted = util.InterfaceBoolToPointer(d.Get("encrypted"))
	}
	if d.HasChange("settable") {
		opt.Settable = util.InterfaceBoolToPointer(d.Get("settable"))
	}
	if d.HasChange("description") {
		opt.Description = util.InterfaceStringToPointer(d.Get("description"))
	}
	_, _, err = c.VariableService.Update(domain, variableId, &opt)
	if err != nil {
		return diag.FromErr(err)
	}
	return readContextVariable(ctx, d, meta)
}

func readContextVariable(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	var diags diag.Diagnostics
	domain, vid, err := util.DecomposeDoubleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	variableId, err := strconv.Atoi(vid)
	if err != nil {
		return diag.FromErr(err)
	}
	variable, resp, err := c.VariableService.Get(domain, variableId)
	if err != nil {
		if util.IsResourceNotFound(resp, err) {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	if variable.Type != buddy.VariableTypeVar {
		return diag.Errorf("Variable not found")
	}
	err = util.ApiVariableToResourceData(domain, variable, d, true)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func createContextVariable(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	domain := d.Get("domain").(string)
	opt := buddy.VariableOps{
		Key:   util.InterfaceStringToPointer(d.Get("key")),
		Value: util.InterfaceStringToPointer(d.Get("value")),
		Type:  util.InterfaceStringToPointer(buddy.VariableTypeVar),
	}
	if settable, ok := d.GetOk("settable"); ok {
		opt.Settable = util.InterfaceBoolToPointer(settable)
	}
	if encrypted, ok := d.GetOk("encrypted"); ok {
		opt.Encrypted = util.InterfaceBoolToPointer(encrypted)
	}
	if description, ok := d.GetOk("description"); ok {
		opt.Description = util.InterfaceStringToPointer(description)
	}
	if projectName, ok := d.GetOk("project_name"); ok {
		opt.Project = &buddy.VariableProject{
			Name: projectName.(string),
		}
	}
	if pipelineId, ok := d.GetOk("pipeline_id"); ok {
		opt.Pipeline = &buddy.VariablePipeline{
			Id: pipelineId.(int),
		}
	}
	if actionId, ok := d.GetOk("action_id"); ok {
		opt.Action = &buddy.VariableAction{
			Id: actionId.(int),
		}
	}
	variable, _, err := c.VariableService.Create(domain, &opt)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(util.ComposeDoubleId(domain, strconv.Itoa(variable.Id)))
	return readContextVariable(ctx, d, meta)
}

package resource

// todo variable ssh key

//
//import (
//	"buddy-terraform/buddy/util"
//	"context"
//	"github.com/buddy/api-go-sdk/buddy"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
//	"strconv"
//)
//
//func VariableSshKey() *schema.Resource {
//	return &schema.Resource{
//		Description: "Create and manage a variable of SSH key type\n\n" +
//			"Workspace administrator rights are required\n\n" +
//			"Token scope required: `WORKSPACE`, `VARIABLE_ADD`, `VARIABLE_MANAGE`, `VARIABLE_INFO`",
//		CreateContext: createContextVariableSshKey,
//		ReadContext:   readContextVariableSshKey,
//		UpdateContext: updateContextVariableSshKey,
//		DeleteContext: deleteContextVariableSshKey,
//		Importer: &schema.ResourceImporter{
//			StateContext: schema.ImportStatePassthroughContext,
//		},
//		Schema: map[string]*schema.Schema{
//			"id": {
//				Description: "The Terraform resource identifier for this item",
//				Type:        schema.TypeString,
//				Computed:    true,
//			},
//			"domain": {
//				Description:  "The workspace's URL handle",
//				Type:         schema.TypeString,
//				Required:     true,
//				ForceNew:     true,
//				ValidateFunc: util.ValidateDomain,
//			},
//			"key": {
//				Description: "The variable's name",
//				Type:        schema.TypeString,
//				Required:    true,
//				ForceNew:    true,
//			},
//			"value": {
//				Description: "The variable's value",
//				Type:        schema.TypeString,
//				Required:    true,
//				Sensitive:   true,
//			},
//			"file_place": {
//				Description: "Should the variable's be copied to an action's container in **file_path** (`CONTAINER`, `NONE`)",
//				Type:        schema.TypeString,
//				Required:    true,
//				ValidateFunc: validation.StringInSlice([]string{
//					buddy.VariableSshKeyFilePlaceContainer,
//					buddy.VariableSshKeyFilePlaceNone,
//				}, false),
//			},
//			"file_path": {
//				Description: "The variable's path in the action's container",
//				Type:        schema.TypeString,
//				Required:    true,
//			},
//			"file_chmod": {
//				Description: "The variable's file permission in an action's container",
//				Type:        schema.TypeString,
//				Required:    true,
//			},
//			"project_name": {
//				Description: "The variable's project name",
//				Type:        schema.TypeString,
//				Optional:    true,
//				ForceNew:    true,
//			},
//			"pipeline_id": {
//				Description: "The variable's pipeline ID",
//				Type:        schema.TypeInt,
//				Optional:    true,
//				ForceNew:    true,
//			},
//			"action_id": {
//				Description: "The variable's action ID",
//				Type:        schema.TypeInt,
//				Optional:    true,
//				ForceNew:    true,
//			},
//			"settable": {
//				Description: "Is the variable's value changeable, always false for buddy_variable_ssh_key",
//				Type:        schema.TypeBool,
//				Computed:    true,
//			},
//			"description": {
//				Description: "The variable's description",
//				Type:        schema.TypeString,
//				Optional:    true,
//			},
//			"variable_id": {
//				Description: "The variable's ID",
//				Type:        schema.TypeInt,
//				Computed:    true,
//			},
//			"value_processed": {
//				Description: "The variable's value, always encrypted for buddy_variable_ssh_key",
//				Type:        schema.TypeString,
//				Computed:    true,
//				Sensitive:   true,
//			},
//			"encrypted": {
//				Description: "Is the variable's value encrypted, always true for buddy_variable_ssh_key",
//				Type:        schema.TypeBool,
//				Computed:    true,
//			},
//			"checksum": {
//				Description: "The variable's checksum",
//				Type:        schema.TypeString,
//				Computed:    true,
//			},
//			"key_fingerprint": {
//				Description: "The variable's fingerprint",
//				Type:        schema.TypeString,
//				Computed:    true,
//			},
//			"public_value": {
//				Description: "The variable's public key",
//				Type:        schema.TypeString,
//				Computed:    true,
//			},
//		},
//	}
//}
//
//func deleteContextVariableSshKey(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	c := meta.(*buddy.Client)
//	var diags diag.Diagnostics
//	domain, vid, err := util.DecomposeDoubleId(d.Id())
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	variableId, err := strconv.Atoi(vid)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	_, err = c.VariableService.Delete(domain, variableId)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	return diags
//}
//
//func updateContextVariableSshKey(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	c := meta.(*buddy.Client)
//	domain, vid, err := util.DecomposeDoubleId(d.Id())
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	variableId, err := strconv.Atoi(vid)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	opt := buddy.VariableOps{
//		Type:      util.InterfaceStringToPointer(buddy.VariableTypeSshKey),
//		Value:     util.InterfaceStringToPointer(d.Get("value")),
//		FilePlace: util.InterfaceStringToPointer(d.Get("file_place")),
//		FilePath:  util.InterfaceStringToPointer(d.Get("file_path")),
//		FileChmod: util.InterfaceStringToPointer(d.Get("file_chmod")),
//		Encrypted: util.InterfaceBoolToPointer(true),
//		Settable:  util.InterfaceBoolToPointer(false),
//	}
//	if d.HasChange("description") {
//		opt.Description = util.InterfaceStringToPointer(d.Get("description"))
//	}
//	_, _, err = c.VariableService.Update(domain, variableId, &opt)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	return readContextVariableSshKey(ctx, d, meta)
//}
//
//func readContextVariableSshKey(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	c := meta.(*buddy.Client)
//	var diags diag.Diagnostics
//	domain, vid, err := util.DecomposeDoubleId(d.Id())
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	variableId, err := strconv.Atoi(vid)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	variable, resp, err := c.VariableService.Get(domain, variableId)
//	if err != nil {
//		if util.IsResourceNotFound(resp, err) {
//			d.SetId("")
//			return diags
//		}
//		return diag.FromErr(err)
//	}
//	if variable.Type != buddy.VariableTypeSshKey {
//		return diag.Errorf("Variable not found")
//	}
//	err = util.ApiVariableSshKeyToResourceData(domain, variable, d, true)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	return diags
//}
//
//func createContextVariableSshKey(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	c := meta.(*buddy.Client)
//	domain := d.Get("domain").(string)
//	opt := buddy.VariableOps{
//		Key:       util.InterfaceStringToPointer(d.Get("key")),
//		Value:     util.InterfaceStringToPointer(d.Get("value")),
//		Type:      util.InterfaceStringToPointer(buddy.VariableTypeSshKey),
//		Encrypted: util.InterfaceBoolToPointer(true),
//		Settable:  util.InterfaceBoolToPointer(false),
//		FilePlace: util.InterfaceStringToPointer(d.Get("file_place")),
//		FilePath:  util.InterfaceStringToPointer(d.Get("file_path")),
//		FileChmod: util.InterfaceStringToPointer(d.Get("file_chmod")),
//	}
//	if description, ok := d.GetOk("description"); ok {
//		opt.Description = util.InterfaceStringToPointer(description)
//	}
//	if projectName, ok := d.GetOk("project_name"); ok {
//		opt.Project = &buddy.VariableProject{
//			Name: projectName.(string),
//		}
//	}
//	if pipelineId, ok := d.GetOk("pipeline_id"); ok {
//		opt.Pipeline = &buddy.VariablePipeline{
//			Id: pipelineId.(int),
//		}
//	}
//	if actionId, ok := d.GetOk("action_id"); ok {
//		opt.Action = &buddy.VariableAction{
//			Id: actionId.(int),
//		}
//	}
//	variable, _, err := c.VariableService.Create(domain, &opt)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	d.SetId(util.ComposeDoubleId(domain, strconv.Itoa(variable.Id)))
//	return readContextVariableSshKey(ctx, d, meta)
//}

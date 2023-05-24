package source

//
//import (
//	"buddy-terraform/buddy/util"
//	"context"
//	"github.com/buddy/api-go-sdk/buddy"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
//	"regexp"
//)
//
//func VariablesSshKeys() *schema.Resource {
//	return &schema.Resource{
//		Description: "List variables of SSH key type and optionally filter them by key, project name, pipeline or action\n\n" +
//			"Token scope required: `WORKSPACE`, `VARIABLE_INFO`",
//		ReadContext: readContextVariablesSshKeys,
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
//				ValidateFunc: util.ValidateDomain,
//			},
//			"key_regex": {
//				Description:  "The variable's key regular expression to match",
//				Type:         schema.TypeString,
//				Optional:     true,
//				ValidateFunc: validation.StringIsValidRegExp,
//			},
//			"project_name": {
//				Description: "Get only from provided project",
//				Type:        schema.TypeString,
//				Optional:    true,
//			},
//			"pipeline_id": {
//				Description: "Get only from provided pipeline",
//				Type:        schema.TypeInt,
//				Optional:    true,
//			},
//			"action_id": {
//				Description: "Get only from provided action",
//				Type:        schema.TypeInt,
//				Optional:    true,
//			},
//			"variables": {
//				Description: "List of variables",
//				Type:        schema.TypeList,
//				Computed:    true,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						"key": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"encrypted": {
//							Type:     schema.TypeBool,
//							Computed: true,
//						},
//						"settable": {
//							Type:     schema.TypeBool,
//							Computed: true,
//						},
//						"description": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"value": {
//							Type:      schema.TypeString,
//							Computed:  true,
//							Sensitive: true,
//						},
//						"variable_id": {
//							Description: "The variable's ID",
//							Type:        schema.TypeInt,
//							Computed:    true,
//						},
//						"file_place": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"file_path": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"file_chmod": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"checksum": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"key_fingerprint": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"public_value": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//					},
//				},
//			},
//		},
//	}
//}
//
//func readContextVariablesSshKeys(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	c := meta.(*buddy.Client)
//	var diags diag.Diagnostics
//	var keyRegex *regexp.Regexp
//	domain := d.Get("domain").(string)
//	opt := buddy.VariableGetListQuery{}
//	if projectName, ok := d.GetOk("project_name"); ok {
//		opt.ProjectName = projectName.(string)
//	}
//	if pipelineId, ok := d.GetOk("pipeline_id"); ok {
//		opt.PipelineId = pipelineId.(int)
//	}
//	if actionId, ok := d.GetOk("action_id"); ok {
//		opt.ActionId = actionId.(int)
//	}
//	variables, _, err := c.VariableService.GetList(domain, &opt)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	if key, ok := d.GetOk("key_regex"); ok {
//		keyRegex = regexp.MustCompile(key.(string))
//	}
//	var result []interface{}
//	for _, v := range variables.Variables {
//		if v.Type != buddy.VariableTypeSshKey {
//			continue
//		}
//		if keyRegex != nil && !keyRegex.MatchString(v.Key) {
//			continue
//		}
//		result = append(result, util.ApiShortVariableSshKeyToMap(v))
//	}
//	d.SetId(util.UniqueString())
//	err = d.Set("variables", result)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	return diags
//}

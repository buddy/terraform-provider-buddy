package resource

// todo webhook

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
//func Webhook() *schema.Resource {
//	return &schema.Resource{
//		Description: "Create and manage a workspace webhook\n\n" +
//			"Workspace administrator rights are required\n\n" +
//			"Token scopes required: `WORKSPACE`, `WEBHOOK_ADD`, `WEBHOOK_MANAGE`, `WEBHOOK_INFO`",
//		CreateContext: createContextWebhook,
//		ReadContext:   readContextWebhook,
//		DeleteContext: deleteContextWebhook,
//		UpdateContext: updateContextWebhook,
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
//			"events": {
//				Description: "The webhook's event's list. Allowed: `PUSH`, `EXECUTION_STARTED`, `EXECUTION_SUCCESSFUL`, `EXECUTION_FAILED`, `EXECUTION_FINISHED`",
//				Type:        schema.TypeSet,
//				Required:    true,
//				MinItems:    1,
//				Elem: &schema.Schema{
//					Type: schema.TypeString,
//					ValidateFunc: validation.StringInSlice([]string{
//						buddy.WebhookEventPush,
//						buddy.WebhookEventExecutionStarted,
//						buddy.WebhookEventExecutionSuccessful,
//						buddy.WebhookEventExecutionFailed,
//						buddy.WebhookEventExecutionFinished,
//					}, false),
//				},
//			},
//			"projects": {
//				Description: "To which projects the webhook should be assigned. If left empty all projects will be used",
//				Type:        schema.TypeSet,
//				Required:    true,
//				MinItems:    0,
//				Elem: &schema.Schema{
//					Type: schema.TypeString,
//				},
//			},
//			"target_url": {
//				Description: "The webhook's target URL",
//				Type:        schema.TypeString,
//				Required:    true,
//			},
//			"secret_key": {
//				Description: "The webhook's secret value sent in the payload",
//				Type:        schema.TypeString,
//				Optional:    true,
//			},
//			"webhook_id": {
//				Description: "The webhook's ID",
//				Type:        schema.TypeInt,
//				Computed:    true,
//			},
//			"html_url": {
//				Description: "The webhook's URL",
//				Type:        schema.TypeString,
//				Computed:    true,
//			},
//		},
//	}
//}
//
//func updateContextWebhook(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	c := meta.(*buddy.Client)
//	domain, wid, err := util.DecomposeDoubleId(d.Id())
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	webhookId, err := strconv.Atoi(wid)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	opt := buddy.WebhookOps{}
//	if d.HasChange("target_url") {
//		opt.TargetUrl = util.InterfaceStringToPointer(d.Get("target_url"))
//	}
//	if d.HasChange("events") {
//		opt.Events = util.InterfaceStringSetToPointer(d.Get("events"))
//	}
//	if d.HasChange("secret_key") {
//		opt.SecretKey = util.InterfaceStringToPointer(d.Get("secret_key"))
//	}
//	if d.HasChange("target_url") {
//		opt.TargetUrl = util.InterfaceStringToPointer(d.Get("target_url"))
//	}
//	_, _, err = c.WebhookService.Update(domain, webhookId, &opt)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	return readContextWebhook(ctx, d, meta)
//}
//
//func deleteContextWebhook(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	c := meta.(*buddy.Client)
//	var diags diag.Diagnostics
//	domain, wid, err := util.DecomposeDoubleId(d.Id())
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	webhookId, err := strconv.Atoi(wid)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	_, err = c.WebhookService.Delete(domain, webhookId)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	return diags
//}
//
//func readContextWebhook(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	c := meta.(*buddy.Client)
//	var diags diag.Diagnostics
//	domain, wid, err := util.DecomposeDoubleId(d.Id())
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	webhookId, err := strconv.Atoi(wid)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	w, resp, err := c.WebhookService.Get(domain, webhookId)
//	if err != nil {
//		if util.IsResourceNotFound(resp, err) {
//			d.SetId("")
//			return diags
//		}
//		return diag.FromErr(err)
//	}
//	err = util.ApiWebhookToResourceData(domain, w, d, false)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	return diags
//}
//
//func createContextWebhook(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	c := meta.(*buddy.Client)
//	domain := d.Get("domain").(string)
//	opt := buddy.WebhookOps{
//		TargetUrl: util.InterfaceStringToPointer(d.Get("target_url")),
//		Events:    util.InterfaceStringSetToPointer(d.Get("events")),
//	}
//	if secretKey, ok := d.GetOk("secret_key"); ok {
//		opt.SecretKey = util.InterfaceStringToPointer(secretKey)
//	}
//	if projects, ok := d.GetOk("projects"); ok {
//		opt.Projects = util.InterfaceStringSetToPointer(projects)
//	}
//	webhook, _, err := c.WebhookService.Create(domain, &opt)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	d.SetId(util.ComposeDoubleId(domain, strconv.Itoa(webhook.Id)))
//	return readContextWebhook(ctx, d, meta)
//}

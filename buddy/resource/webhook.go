package resource

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
)

var (
	_ resource.Resource                = &webhookResource{}
	_ resource.ResourceWithConfigure   = &webhookResource{}
	_ resource.ResourceWithImportState = &webhookResource{}
)

func NewWebhookResource() resource.Resource {
	return &webhookResource{}
}

type webhookResource struct {
	client *buddy.Client
}

type webhookResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Domain    types.String `tfsdk:"domain"`
	TargetUrl types.String `tfsdk:"target_url"`
	SecretKey types.String `tfsdk:"secret_key"`
	Events    types.Set    `tfsdk:"events"`
	Projects  types.Set    `tfsdk:"projects"`
	WebhookId types.Int64  `tfsdk:"webhook_id"`
	HtmlUrl   types.String `tfsdk:"html_url"`
}

func (r *webhookResourceModel) decomposeId() (string, int, error) {
	domain, wid, err := util.DecomposeDoubleId(r.ID.ValueString())
	if err != nil {
		return "", 0, err
	}
	webhookId, err := strconv.Atoi(wid)
	if err != nil {
		return "", 0, err
	}
	return domain, webhookId, nil
}

func (r *webhookResourceModel) loadAPI(ctx context.Context, domain string, webhook *buddy.Webhook) diag.Diagnostics {
	r.ID = types.StringValue(util.ComposeDoubleId(domain, strconv.Itoa(webhook.Id)))
	r.Domain = types.StringValue(domain)
	r.TargetUrl = types.StringValue(webhook.TargetUrl)
	r.SecretKey = types.StringValue(webhook.SecretKey)
	r.HtmlUrl = types.StringValue(webhook.HtmlUrl)
	r.WebhookId = types.Int64Value(int64(webhook.Id))
	events, eventsDiags := types.SetValueFrom(ctx, types.StringType, &webhook.Events)
	r.Events = events
	projects, projectsDiags := types.SetValueFrom(ctx, types.StringType, &webhook.Projects)
	r.Projects = projects
	eventsDiags.Append(projectsDiags...)
	return eventsDiags
}

func (r *webhookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}

func (r *webhookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage a workspace webhook\n\n" +
			"Workspace administrator rights are required\n\n" +
			"Token scopes required: `WORKSPACE`, `WEBHOOK_ADD`, `WEBHOOK_MANAGE`, `WEBHOOK_INFO`",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Terraform resource identifier for this item",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The workspace's URL handle",
				Required:            true,
				Validators:          util.StringValidatorsDomain(),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_url": schema.StringAttribute{
				MarkdownDescription: "The webhook's target URL",
				Required:            true,
			},
			"secret_key": schema.StringAttribute{
				MarkdownDescription: "The webhook's secret value sent in the payload",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"events": schema.SetAttribute{
				MarkdownDescription: "The webhook's event's list. Allowed: `PUSH`, `EXECUTION_STARTED`, `EXECUTION_SUCCESSFUL`, `EXECUTION_FAILED`, `EXECUTION_FINISHED`",
				Required:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(stringvalidator.OneOf(
						buddy.WebhookEventPush,
						buddy.WebhookEventExecutionStarted,
						buddy.WebhookEventExecutionSuccessful,
						buddy.WebhookEventExecutionFailed,
						buddy.WebhookEventExecutionFinished,
					)),
				},
			},
			"projects": schema.SetAttribute{
				MarkdownDescription: "To which projects the webhook should be assigned. If left empty all projects will be used",
				Required:            true,
				ElementType:         types.StringType,
			},
			"webhook_id": schema.Int64Attribute{
				MarkdownDescription: "The webhook's ID",
				Computed:            true,
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The webhook's URL",
				Computed:            true,
			},
		},
	}
}

func (r *webhookResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *webhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *webhookResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	events, d := util.StringSetToApi(ctx, &data.Events)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	projects, d := util.StringSetToApi(ctx, &data.Projects)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	ops := buddy.WebhookOps{
		TargetUrl: data.TargetUrl.ValueStringPointer(),
		Events:    events,
		Projects:  projects,
	}
	if !data.SecretKey.IsNull() && !data.SecretKey.IsUnknown() {
		ops.SecretKey = data.SecretKey.ValueStringPointer()
	}
	webhook, _, err := r.client.WebhookService.Create(domain, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create webhook", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, webhook)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *webhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *webhookResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, webhookId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("webhook", err))
		return
	}
	webhook, httpResp, err := r.client.WebhookService.Get(domain, webhookId)
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get webhook", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, webhook)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *webhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *webhookResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, webhookId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("webhook", err))
		return
	}
	events, d := util.StringSetToApi(ctx, &data.Events)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	projects, d := util.StringSetToApi(ctx, &data.Projects)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	ops := buddy.WebhookOps{
		TargetUrl: data.TargetUrl.ValueStringPointer(),
		Events:    events,
		Projects:  projects,
	}
	if !data.SecretKey.IsNull() && !data.SecretKey.IsUnknown() {
		ops.SecretKey = data.SecretKey.ValueStringPointer()
	}
	webhook, _, err := r.client.WebhookService.Update(domain, webhookId, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update webhook", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, webhook)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *webhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *webhookResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, webhookId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("webhook", err))
		return
	}
	_, err = r.client.WebhookService.Delete(domain, webhookId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete webhook", err))
	}
}

func (r *webhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

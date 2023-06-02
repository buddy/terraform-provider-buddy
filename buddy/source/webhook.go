package source

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
	"strconv"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ datasource.DataSource              = &webhookSource{}
	_ datasource.DataSourceWithConfigure = &webhookSource{}
)

func NewWebhookSource() datasource.DataSource {
	return &webhookSource{}
}

type webhookSource struct {
	client *buddy.Client
}

type webhookSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Domain    types.String `tfsdk:"domain"`
	TargetUrl types.String `tfsdk:"target_url"`
	WebhookId types.Int64  `tfsdk:"webhook_id"`
	HtmlUrl   types.String `tfsdk:"html_url"`
}

func (s *webhookSourceModel) loadAPI(domain string, webhook *buddy.Webhook) {
	s.ID = types.StringValue(util.ComposeDoubleId(domain, strconv.Itoa(webhook.Id)))
	s.Domain = types.StringValue(domain)
	s.TargetUrl = types.StringValue(webhook.TargetUrl)
	s.WebhookId = types.Int64Value(int64(webhook.Id))
	s.HtmlUrl = types.StringValue(webhook.HtmlUrl)
}

func (s *webhookSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}

func (s *webhookSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *webhookSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get webhook by target URL or webhook ID\n\n" +
			"Token scope required: `WORKSPACE`, `WEBHOOK_INFO`",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Terraform resource identifier for this item",
				Computed:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The workspace's URL handle",
				Required:            true,
				Validators:          util.StringValidatorsDomain(),
			},
			"target_url": schema.StringAttribute{
				MarkdownDescription: "The webhook's target URL",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("target_url"),
						path.MatchRoot("webhook_id"),
					}...),
				},
			},
			"webhook_id": schema.Int64Attribute{
				MarkdownDescription: "The webhook's ID",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("target_url"),
						path.MatchRoot("webhook_id"),
					}...),
				},
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The webhook's URL",
				Computed:            true,
			},
		},
	}
}

func (s *webhookSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *webhookSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var webhook *buddy.Webhook
	var err error
	domain := data.Domain.ValueString()
	if !data.WebhookId.IsNull() && !data.WebhookId.IsUnknown() {
		var httpRes *http.Response
		wid := int(data.WebhookId.ValueInt64())
		webhook, httpRes, err = s.client.WebhookService.Get(domain, wid)
		if err != nil {
			if util.IsResourceNotFound(httpRes, err) {
				resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("webhook"))
				return
			}
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get webhook", err))
			return
		}
	} else {
		var webhooks *buddy.Webhooks
		targetUrl := data.TargetUrl.ValueString()
		webhooks, _, err = s.client.WebhookService.GetList(domain)
		if err != nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get webhooks", err))
			return
		}
		for _, w := range webhooks.Webhooks {
			if w.TargetUrl == targetUrl {
				webhook = w
				break
			}
		}
		if webhook == nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("webhook"))
			return
		}
	}
	data.loadAPI(domain, webhook)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

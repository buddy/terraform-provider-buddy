package source

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
)

var (
	_ datasource.DataSource              = &webhooksSource{}
	_ datasource.DataSourceWithConfigure = &webhooksSource{}
)

func NewWebhooksSource() datasource.DataSource {
	return &webhooksSource{}
}

type webhooksSource struct {
	client *buddy.Client
}

type webhooksSourceModel struct {
	ID             types.String `tfsdk:"id"`
	Domain         types.String `tfsdk:"domain"`
	TargetUrlRegex types.String `tfsdk:"target_url_regex"`
	Webhooks       types.Set    `tfsdk:"webhooks"`
}

func (s *webhooksSourceModel) loadAPI(ctx context.Context, domain string, webhooks *[]*buddy.Webhook) diag.Diagnostics {
	s.ID = types.StringValue(util.UniqueString())
	s.Domain = types.StringValue(domain)
	w, d := util.WebhooksModelFromApi(ctx, webhooks)
	s.Webhooks = w
	return d
}

func (s *webhooksSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhooks"
}

func (s *webhooksSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *webhooksSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List webhooks and optionally filter them by target URL\n\n" +
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
			"target_url_regex": schema.StringAttribute{
				MarkdownDescription: "The webhook's target_url regular expression to match",
				Optional:            true,
				Validators: []validator.String{
					util.RegexpValidator(),
				},
			},
			"webhooks": schema.SetNestedAttribute{
				MarkdownDescription: "List of webhooks",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.SourceWebhookModelAttributes(),
				},
			},
		},
	}
}

func (s *webhooksSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *webhooksSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	var targetSiteRegex *regexp.Regexp
	if !data.TargetUrlRegex.IsNull() && !data.TargetUrlRegex.IsUnknown() {
		targetSiteRegex = regexp.MustCompile(data.TargetUrlRegex.ValueString())
	}
	webhooks, _, err := s.client.WebhookService.GetList(domain)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get webhooks", err))
		return
	}
	var result []*buddy.Webhook
	for _, w := range webhooks.Webhooks {
		if targetSiteRegex != nil && !targetSiteRegex.MatchString(w.TargetUrl) {
			continue
		}
		result = append(result, w)
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, &result)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

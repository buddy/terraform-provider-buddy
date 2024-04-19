package source

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ datasource.DataSource              = &integrationSource{}
	_ datasource.DataSourceWithConfigure = &integrationSource{}
)

func NewIntegrationSource() datasource.DataSource {
	return &integrationSource{}
}

type integrationSource struct {
	client *buddy.Client
}

type integrationSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Domain        types.String `tfsdk:"domain"`
	Name          types.String `tfsdk:"name"`
	IntegrationId types.String `tfsdk:"integration_id"`
	Identifier    types.String `tfsdk:"identifier"`
	Type          types.String `tfsdk:"type"`
	HtmlUrl       types.String `tfsdk:"html_url"`
}

func (s *integrationSourceModel) loadAPI(domain string, integration *buddy.Integration) {
	s.ID = types.StringValue(util.ComposeDoubleId(domain, integration.HashId))
	s.Domain = types.StringValue(domain)
	s.Name = types.StringValue(integration.Name)
	s.IntegrationId = types.StringValue(integration.HashId)
	s.Identifier = types.StringValue(integration.Identifier)
	s.Type = types.StringValue(integration.Type)
	s.HtmlUrl = types.StringValue(integration.HtmlUrl)
}

func (s *integrationSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration"
}

func (s *integrationSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *integrationSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get integration by name or integration ID\n\n" +
			"Token scope required: `INTEGRATION_INFO`",
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
			"name": schema.StringAttribute{
				MarkdownDescription: "The integration's name",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("name"),
						path.MatchRoot("integration_id"),
					}...),
				},
			},
			"integration_id": schema.StringAttribute{
				MarkdownDescription: "The integration's ID",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("name"),
						path.MatchRoot("integration_id"),
					}...),
				},
			},
			"identifier": schema.StringAttribute{
				MarkdownDescription: "The integration's identifier",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The integration's type",
				Computed:            true,
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The integration's URL",
				Computed:            true,
			},
		},
	}
}

func (s *integrationSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *integrationSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var integration *buddy.Integration
	domain := data.Domain.ValueString()
	if !data.IntegrationId.IsNull() && !data.IntegrationId.IsUnknown() {
		var httpResp *http.Response
		var err error
		integration, httpResp, err = s.client.IntegrationService.Get(domain, data.IntegrationId.ValueString())
		if err != nil {
			if util.IsResourceNotFound(httpResp, err) {
				resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("integration"))
				return
			}
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get integration", err))
			return
		}
	} else {
		name := data.Name.ValueString()
		integrations, _, err := s.client.IntegrationService.GetList(domain)
		if err != nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get integrations", err))
			return
		}
		for _, i := range integrations.Integrations {
			if i.Name == name {
				integration = i
				break
			}
		}
		if integration == nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("integration"))
			return
		}
	}
	data.loadAPI(domain, integration)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

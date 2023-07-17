package source

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ datasource.DataSource              = &integrationsSource{}
	_ datasource.DataSourceWithConfigure = &integrationsSource{}
)

func NewIntegrationsSource() datasource.DataSource {
	return &integrationsSource{}
}

type integrationsSource struct {
	client *buddy.Client
}

type integrationsSourceModel struct {
	Id           types.String `tfsdk:"id"`
	Domain       types.String `tfsdk:"domain"`
	NameRegex    types.String `tfsdk:"name_regex"`
	Type         types.String `tfsdk:"type"`
	Integrations types.Set    `tfsdk:"integrations"`
}

func (s *integrationsSourceModel) loadAPI(ctx context.Context, domain string, integrations *[]*buddy.Integration) diag.Diagnostics {
	s.Id = types.StringValue(util.UniqueString())
	s.Domain = types.StringValue(domain)
	i, d := util.IntegrationsModelFromApi(ctx, integrations)
	s.Integrations = i
	return d
}

func (s *integrationsSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integrations"
}

func (s *integrationsSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *integrationsSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List integrations and optionally filter them by name or type\n\n" +
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
			"name_regex": schema.StringAttribute{
				MarkdownDescription: "The integration's name regular expression to match",
				Optional:            true,
				Validators: []validator.String{
					util.RegexpValidator(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The integration's type",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						buddy.IntegrationTypeDigitalOcean,
						buddy.IntegrationTypeAmazon,
						buddy.IntegrationTypeShopify,
						buddy.IntegrationTypePushover,
						buddy.IntegrationTypeRackspace,
						buddy.IntegrationTypeCloudflare,
						buddy.IntegrationTypeNewRelic,
						buddy.IntegrationTypeSentry,
						buddy.IntegrationTypeRollbar,
						buddy.IntegrationTypeDatadog,
						buddy.IntegrationTypeDigitalOceanSpaces,
						buddy.IntegrationTypeHoneybadger,
						buddy.IntegrationTypeVultr,
						buddy.IntegrationTypeSentryEnterprise,
						buddy.IntegrationTypeLoggly,
						buddy.IntegrationTypeFirebase,
						buddy.IntegrationTypeUpcloud,
						buddy.IntegrationTypeGhostInspector,
						buddy.IntegrationTypeAzureCloud,
						buddy.IntegrationTypeDockerHub,
						buddy.IntegrationTypeGitHub,
						buddy.IntegrationTypeGitLab,
						buddy.IntegrationTypeStackHawk,
					),
				},
			},
			"integrations": schema.SetNestedAttribute{
				MarkdownDescription: "List of integrations",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.SourceIntegrationModelAttributes(),
				},
			},
		},
	}
}

func (s *integrationsSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *integrationsSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var nameRegexp *regexp.Regexp
	var typ *string
	domain := data.Domain.ValueString()
	if !data.NameRegex.IsNull() && !data.NameRegex.IsUnknown() {
		nameRegexp = regexp.MustCompile(data.NameRegex.ValueString())
	}
	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		typ = data.Type.ValueStringPointer()
	}
	var result []*buddy.Integration
	integrations, _, err := s.client.IntegrationService.GetList(domain)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get integrations", err))
		return
	}
	for _, i := range integrations.Integrations {
		if nameRegexp != nil && !nameRegexp.MatchString(i.Name) {
			continue
		}
		if typ != nil && *typ != i.Type {
			continue
		}
		result = append(result, i)
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, &result)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

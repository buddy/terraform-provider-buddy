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
	_ datasource.DataSource              = &domainsSource{}
	_ datasource.DataSourceWithConfigure = &domainsSource{}
)

func NewDomainsSource() datasource.DataSource {
	return &domainsSource{}
}

type domainsSource struct {
	client *buddy.Client
}

type domainsSourceModel struct {
	ID              types.String `tfsdk:"id"`
	WorkspaceDomain types.String `tfsdk:"workspace_domain"`
	Type            types.String `tfsdk:"type"`
	DomainRegex     types.String `tfsdk:"domain_regex"`
	Domains         types.Set    `tfsdk:"domains"`
}

func (s *domainsSourceModel) loadAPI(ctx context.Context, workspaceDomain string, domains *[]*buddy.Domain) diag.Diagnostics {
	s.ID = types.StringValue(util.UniqueString())
	s.WorkspaceDomain = types.StringValue(workspaceDomain)
	d, diags := util.DomainsModelFromApi(ctx, domains)
	s.Domains = d
	return diags
}

func (s *domainsSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domains"
}

func (s *domainsSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *domainsSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List domains and optionally filter them by type or name\n\n" +
			"Token scope required: `DOMAIN_READ`",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Terraform resource identifier for this item",
				Computed:            true,
			},
			"workspace_domain": schema.StringAttribute{
				MarkdownDescription: "The workspace's URL handle",
				Required:            true,
				Validators:          util.StringValidatorsDomain(),
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "List only domains of the given type. Allowed values: POINTED, PRIVATE, REGISTERED, CLAIMED",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						buddy.DomainTypePointed,
						buddy.DomainTypePrivate,
						buddy.DomainTypeRegistered,
						buddy.DomainTypeClaimed,
					),
				},
			},
			"domain_regex": schema.StringAttribute{
				MarkdownDescription: "The domain's name regular expression to match",
				Optional:            true,
				Validators: []validator.String{
					util.RegexpValidator(),
				},
			},
			"domains": schema.SetNestedAttribute{
				MarkdownDescription: "List of domains",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.SourceDomainModelAttributes(),
				},
			},
		},
	}
}

func (s *domainsSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *domainsSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	workspaceDomain := data.WorkspaceDomain.ValueString()
	var query *buddy.DomainGetListQuery
	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		query = &buddy.DomainGetListQuery{Type: data.Type.ValueString()}
	}
	var domainRegex *regexp.Regexp
	if !data.DomainRegex.IsNull() && !data.DomainRegex.IsUnknown() {
		domainRegex = regexp.MustCompile(data.DomainRegex.ValueString())
	}
	domains, _, err := s.client.DomainService.GetList(workspaceDomain, query)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get domains", err))
		return
	}
	var result []*buddy.Domain
	for _, d := range domains.Domains {
		if domainRegex != nil && !domainRegex.MatchString(d.Name) {
			continue
		}
		result = append(result, d)
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, workspaceDomain, &result)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

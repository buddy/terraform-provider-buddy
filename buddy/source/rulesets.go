package source

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ datasource.DataSource              = &rulesetsSource{}
	_ datasource.DataSourceWithConfigure = &rulesetsSource{}
)

func NewRulesetsSource() datasource.DataSource {
	return &rulesetsSource{}
}

type rulesetsSource struct {
	client *buddy.Client
}

type rulesetsSourceModel struct {
	ID              types.String `tfsdk:"id"`
	Domain          types.String `tfsdk:"domain"`
	ProjectName     types.String `tfsdk:"project_name"`
	IdentifierRegex types.String `tfsdk:"identifier_regex"`
	Rulesets        types.Set    `tfsdk:"rulesets"`
}

func (s *rulesetsSourceModel) loadAPI(ctx context.Context, domain string, projectName string, rulesets []*buddy.Ruleset) diag.Diagnostics {
	s.ID = types.StringValue(util.UniqueString())
	s.Domain = types.StringValue(domain)
	s.ProjectName = types.StringValue(projectName)
	r, d := util.RulesetsModelFromApi(ctx, rulesets)
	s.Rulesets = r
	return d
}

func (s *rulesetsSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rulesets"
}

func (s *rulesetsSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *rulesetsSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List rulesets of a project and optionally filter them by identifier. Only workspace administrators and project managers can read rulesets\n\n" +
			"Token scope required: `RULESET_READ`",
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
			"project_name": schema.StringAttribute{
				MarkdownDescription: "The project's name",
				Required:            true,
			},
			"identifier_regex": schema.StringAttribute{
				MarkdownDescription: "The ruleset's identifier regular expression to match",
				Optional:            true,
				Validators: []validator.String{
					util.RegexpValidator(),
				},
			},
			"rulesets": schema.SetNestedAttribute{
				MarkdownDescription: "List of rulesets",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.SourceRulesetModelAttributes(),
				},
			},
		},
	}
}

func (s *rulesetsSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *rulesetsSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	projectName := data.ProjectName.ValueString()
	var identifierRegex *regexp.Regexp
	if !data.IdentifierRegex.IsNull() && !data.IdentifierRegex.IsUnknown() {
		identifierRegex = regexp.MustCompile(data.IdentifierRegex.ValueString())
	}
	rulesets, _, err := s.client.RulesetService.GetList(domain, projectName)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get rulesets", err))
		return
	}
	result := []*buddy.Ruleset{}
	for _, r := range rulesets.Rulesets {
		if identifierRegex != nil && !identifierRegex.MatchString(r.Identifier) {
			continue
		}
		result = append(result, r)
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, result)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

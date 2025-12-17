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
	_ datasource.DataSource              = &environmentsSource{}
	_ datasource.DataSourceWithConfigure = &environmentsSource{}
)

func NewEnvironmentsSource() datasource.DataSource {
	return &environmentsSource{}
}

type environmentsSource struct {
	client *buddy.Client
}

type environmentsSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Domain       types.String `tfsdk:"domain"`
	ProjectName  types.String `tfsdk:"project_name"`
	NameRegex    types.String `tfsdk:"name_regex"`
	Environments types.Set    `tfsdk:"environments"`
}

func (s *environmentsSourceModel) loadAPI(ctx context.Context, domain string, projectName string, environments *[]*buddy.Environment) diag.Diagnostics {
	s.ID = types.StringValue(util.UniqueString())
	s.Domain = types.StringValue(domain)
	s.ProjectName = types.StringValue(projectName)
	e, d := util.EnvironmentsModelFromApi(ctx, environments)
	s.Environments = e
	return d
}

func (s *environmentsSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environments"
}

func (s *environmentsSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *environmentsSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List environments and optionally filter them by name\n\n" +
			"Token scope required: `WORKSPACE`, `ENVIRONMENT_INFO`",
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
				Optional:            true,
				Computed:            true,
			},
			"name_regex": schema.StringAttribute{
				MarkdownDescription: "The environment's name regular expression to match",
				Optional:            true,
				Validators: []validator.String{
					util.RegexpValidator(),
				},
			},
			"environments": schema.SetNestedAttribute{
				MarkdownDescription: "List of environments",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.SourceEnvironmentModelAttributes(),
				},
			},
		},
	}
}

func (s *environmentsSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *environmentsSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	projectName := data.ProjectName.ValueString()
	var nameRegex *regexp.Regexp
	if !data.NameRegex.IsNull() && !data.NameRegex.IsUnknown() {
		nameRegex = regexp.MustCompile(data.NameRegex.ValueString())
	}
	environments, _, err := s.client.EnvironmentService.GetList(domain, projectName)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get environments", err))
		return
	}
	var result []*buddy.Environment
	for _, e := range environments.Environments {
		if nameRegex != nil && !nameRegex.MatchString(e.Name) {
			continue
		}
		result = append(result, e)
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, &result)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

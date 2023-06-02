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
	_ datasource.DataSource              = &pipelinesSource{}
	_ datasource.DataSourceWithConfigure = &pipelinesSource{}
)

func NewPipelinesSource() datasource.DataSource {
	return &pipelinesSource{}
}

type pipelinesSource struct {
	client *buddy.Client
}

type pipelinesSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Domain      types.String `tfsdk:"domain"`
	ProjectName types.String `tfsdk:"project_name"`
	NameRegex   types.String `tfsdk:"name_regex"`
	Pipelines   types.Set    `tfsdk:"pipelines"`
}

func (s *pipelinesSourceModel) loadAPI(ctx context.Context, domain string, projectName string, pipelines *[]*buddy.Pipeline) diag.Diagnostics {
	s.ID = types.StringValue(util.UniqueString())
	s.Domain = types.StringValue(domain)
	s.ProjectName = types.StringValue(projectName)
	p, d := util.PipelinesModelFromApi(ctx, pipelines)
	s.Pipelines = p
	return d
}

func (s *pipelinesSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pipelines"
}

func (s *pipelinesSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *pipelinesSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List pipelines and optionally filter them by name\n\n" +
			"Token scopes required: `WORKSPACE`, `EXECUTION_INFO`",
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
			"name_regex": schema.StringAttribute{
				MarkdownDescription: "The pipeline's name regular expression to match",
				Optional:            true,
				Validators: []validator.String{
					util.RegexpValidator(),
				},
			},
			"pipelines": schema.SetNestedAttribute{
				MarkdownDescription: "List of pipelines",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.SourcePipelineModelAttributes(),
				},
			},
		},
	}
}

func (s *pipelinesSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *pipelinesSourceModel
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
	pipelines, _, err := s.client.PipelineService.GetListAll(domain, projectName)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get pipelines", err))
		return
	}
	var result []*buddy.Pipeline
	for _, p := range pipelines.Pipelines {
		if nameRegex != nil && !nameRegex.MatchString(p.Name) {
			continue
		}
		result = append(result, p)
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, &result)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

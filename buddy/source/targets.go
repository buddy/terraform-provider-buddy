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
	_ datasource.DataSource              = &targetsSource{}
	_ datasource.DataSourceWithConfigure = &targetsSource{}
)

func NewTargetsSource() datasource.DataSource {
	return &targetsSource{}
}

type targetsSource struct {
	client *buddy.Client
}

type targetsSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Domain        types.String `tfsdk:"domain"`
	ProjectName   types.String `tfsdk:"project_name"`
	PipelineId    types.Int64  `tfsdk:"pipeline_id"`
	ActionId      types.Int64  `tfsdk:"action_id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	NameRegex     types.String `tfsdk:"name_regex"`
	Targets       types.Set    `tfsdk:"targets"`
}

func (s *targetsSourceModel) loadAPI(ctx context.Context, domain string, targets *[]*buddy.Target) diag.Diagnostics {
	s.ID = types.StringValue(util.UniqueString())
	s.Domain = types.StringValue(domain)
	t, d := util.TargetsModelFromApi(ctx, targets)
	s.Targets = t
	return d
}

func (s *targetsSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_targets"
}

func (s *targetsSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List targets\n\n" +
			"Token scope required: `WORKSPACE`, `TARGET_INFO`",
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
			},
			"pipeline_id": schema.Int64Attribute{
				MarkdownDescription: "The pipeline's name",
				Optional:            true,
			},
			"action_id": schema.Int64Attribute{
				MarkdownDescription: "The pipeline action's name",
				Optional:            true,
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "The environment's name",
				Optional:            true,
			},
			"name_regex": schema.StringAttribute{
				MarkdownDescription: "The target's name regular expression to match",
				Optional:            true,
				Validators: []validator.String{
					util.RegexpValidator(),
				},
			},
			"targets": schema.SetNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.SourceTargetModelAttributes(),
				},
			},
		},
	}
}

func (s *targetsSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *targetsSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *targetsSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	query := &buddy.TargetGetListQuery{}
	if !data.ProjectName.IsNull() && !data.ProjectName.IsUnknown() {
		query.ProjectName = data.ProjectName.ValueString()
	}
	if !data.PipelineId.IsNull() && !data.PipelineId.IsUnknown() {
		query.PipelineId = int(data.PipelineId.ValueInt64())
	}
	if !data.ActionId.IsNull() && !data.ActionId.IsUnknown() {
		query.ActionId = int(data.ActionId.ValueInt64())
	}
	if !data.EnvironmentId.IsNull() && !data.EnvironmentId.IsUnknown() {
		query.EnvironmentId = data.EnvironmentId.ValueString()
	}
	var nameRegex *regexp.Regexp
	if !data.NameRegex.IsNull() && !data.NameRegex.IsUnknown() {
		nameRegex = regexp.MustCompile(data.NameRegex.ValueString())
	}
	targets, _, err := s.client.TargetService.GetList(domain, query)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get targets", err))
		return
	}
	var result []*buddy.Target
	for _, t := range targets.Targets {
		if nameRegex != nil && !nameRegex.MatchString(t.Name) {
			continue
		}
		result = append(result, t)
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, &result)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

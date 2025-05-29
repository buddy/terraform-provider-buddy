package source

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	Domain        types.String `tfsdk:"domain"`
	ProjectName   types.String `tfsdk:"project_name"`
	PipelineId    types.Int64  `tfsdk:"pipeline_id"`
	ActionId      types.Int64  `tfsdk:"action_id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	Targets       types.List   `tfsdk:"targets"`
}

func (s *targetsSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_targets"
}

func (s *targetsSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List targets\n\n" +
			"Token scope required: `WORKSPACE`",
		Attributes: map[string]schema.Attribute{
			"domain": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"project_name": schema.StringAttribute{
				Optional: true,
			},
			"pipeline_id": schema.Int64Attribute{
				Optional: true,
			},
			"action_id": schema.Int64Attribute{
				Optional: true,
			},
			"environment_id": schema.StringAttribute{
				Optional: true,
			},
			"targets": schema.ListNestedAttribute{
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

	if !data.ProjectName.IsNull() && data.ProjectName.ValueString() != "" {
		query.ProjectName = data.ProjectName.ValueString()
	}

	if !data.PipelineId.IsNull() {
		query.PipelineId = int(data.PipelineId.ValueInt64())
	}

	if !data.ActionId.IsNull() {
		query.ActionId = int(data.ActionId.ValueInt64())
	}

	if !data.EnvironmentId.IsNull() && data.EnvironmentId.ValueString() != "" {
		query.EnvironmentId = data.EnvironmentId.ValueString()
	}

	targets, _, err := s.client.TargetService.GetList(domain, query)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get targets", err))
		return
	}

	result := make([]*util.TargetModel, 0)
	for _, t := range targets.Targets {
		target := &util.TargetModel{}
		target.LoadAPI(ctx, t)
		result = append(result, target)
	}

	modelTargets, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: util.TargetModelAttrs()}, result)
	resp.Diagnostics.Append(diags...)

	data.Targets = modelTargets
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

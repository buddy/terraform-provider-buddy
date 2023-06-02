package source

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
	"strconv"
)

var (
	_ datasource.DataSource              = &pipelineSource{}
	_ datasource.DataSourceWithConfigure = &pipelineSource{}
)

func NewPipelineSource() datasource.DataSource {
	return &pipelineSource{}
}

type pipelineSource struct {
	client *buddy.Client
}

type pipelineSourceModel struct {
	ID                    types.String `tfsdk:"id"`
	Domain                types.String `tfsdk:"domain"`
	ProjectName           types.String `tfsdk:"project_name"`
	Name                  types.String `tfsdk:"name"`
	PipelineId            types.Int64  `tfsdk:"pipeline_id"`
	Priority              types.String `tfsdk:"priority"`
	HtmlUrl               types.String `tfsdk:"html_url"`
	On                    types.String `tfsdk:"on"`
	LastExecutionStatus   types.String `tfsdk:"last_execution_status"`
	LastExecutionRevision types.String `tfsdk:"last_execution_revision"`
	Disabled              types.Bool   `tfsdk:"disabled"`
	DisablingReason       types.String `tfsdk:"disabling_reason"`
	Refs                  types.Set    `tfsdk:"refs"`
	Event                 types.Set    `tfsdk:"event"`
	Tags                  types.Set    `tfsdk:"tags"`
	DefinitionSource      types.String `tfsdk:"definition_source"`
	RemoteProjectName     types.String `tfsdk:"remote_project_name"`
	RemoteBranch          types.String `tfsdk:"remote_branch"`
	RemotePath            types.String `tfsdk:"remote_path"`
	RemoteParameter       types.Set    `tfsdk:"remote_parameter"`
}

func (s *pipelineSourceModel) loadAPI(ctx context.Context, domain string, projectName string, pipeline *buddy.Pipeline) diag.Diagnostics {
	var diags diag.Diagnostics
	s.ID = types.StringValue(util.ComposeTripleId(domain, projectName, strconv.Itoa(pipeline.Id)))
	s.Domain = types.StringValue(domain)
	s.ProjectName = types.StringValue(projectName)
	s.Name = types.StringValue(pipeline.Name)
	s.PipelineId = types.Int64Value(int64(pipeline.Id))
	s.Priority = types.StringValue(pipeline.Priority)
	s.HtmlUrl = types.StringValue(pipeline.HtmlUrl)
	s.On = types.StringValue(pipeline.On)
	s.LastExecutionRevision = types.StringValue(pipeline.LastExecutionRevision)
	s.LastExecutionStatus = types.StringValue(pipeline.LastExecutionStatus)
	s.Disabled = types.BoolValue(pipeline.Disabled)
	s.DisablingReason = types.StringValue(pipeline.DisabledReason)
	r, d := types.SetValueFrom(ctx, types.StringType, &pipeline.Refs)
	diags.Append(d...)
	s.Refs = r
	e, d := util.EventsModelFromApi(ctx, &pipeline.Events)
	diags.Append(d...)
	s.Event = e
	t, d := types.SetValueFrom(ctx, types.StringType, &pipeline.Tags)
	diags.Append(d...)
	s.Tags = t
	ds := pipeline.DefinitionSource
	if ds == "" {
		ds = buddy.PipelineDefinitionSourceLocal
	}
	s.DefinitionSource = types.StringValue(ds)
	s.RemoteProjectName = types.StringValue(pipeline.RemoteProjectName)
	s.RemoteBranch = types.StringValue(pipeline.RemoteBranch)
	s.RemotePath = types.StringValue(pipeline.RemotePath)
	rp, d := util.RemoteParametersModelFromApi(ctx, &pipeline.RemoteParameters)
	diags.Append(d...)
	s.RemoteParameter = rp
	return diags
}

func (s *pipelineSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pipeline"
}

func (s *pipelineSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *pipelineSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get pipeline by name or pipeline ID\n\n" +
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
			"name": schema.StringAttribute{
				MarkdownDescription: "The pipeline's name",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("pipeline_id"),
						path.MatchRoot("name"),
					}...),
				},
			},
			"pipeline_id": schema.Int64Attribute{
				MarkdownDescription: "The pipeline's ID",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("pipeline_id"),
						path.MatchRoot("name"),
					}...),
				},
			},
			"priority": schema.StringAttribute{
				MarkdownDescription: "The pipeline's priority",
				Computed:            true,
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The pipeline's URL",
				Computed:            true,
			},
			"on": schema.StringAttribute{
				MarkdownDescription: "The pipeline's trigger mode",
				Computed:            true,
			},
			"last_execution_status": schema.StringAttribute{
				MarkdownDescription: "The pipeline's last run status",
				Computed:            true,
			},
			"last_execution_revision": schema.StringAttribute{
				MarkdownDescription: "The pipeline's last run revision",
				Computed:            true,
			},
			"disabled": schema.BoolAttribute{
				MarkdownDescription: "Defines wheter or not the pipeline can be run",
				Computed:            true,
			},
			"disabling_reason": schema.StringAttribute{
				MarkdownDescription: "The pipeline's disabling reason",
				Computed:            true,
			},
			"refs": schema.SetAttribute{
				MarkdownDescription: "The pipeline's list of refs",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "The pipeline's list of tags. Only for `Buddy Enterprise`",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"definition_source": schema.StringAttribute{
				MarkdownDescription: "The pipeline's definition source",
				Computed:            true,
			},
			"remote_project_name": schema.StringAttribute{
				MarkdownDescription: "The pipeline's remote definition project name",
				Computed:            true,
			},
			"remote_branch": schema.StringAttribute{
				MarkdownDescription: "The pipeline's remote definition branch name",
				Computed:            true,
			},
			"remote_path": schema.StringAttribute{
				MarkdownDescription: "The pipeline's remote definition path",
				Computed:            true,
			},
			// singular form for compatibility
			"event": schema.SetNestedAttribute{
				MarkdownDescription: "The pipeline's list of events",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.SourceEventModelAttributes(),
				},
			},
			// singular form for compatibility
			"remote_parameter": schema.SetNestedAttribute{
				MarkdownDescription: "The pipeline's remote definition parameters",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.SourceRemoteParameterModelAttributes(),
				},
			},
		},
	}
}

func (s *pipelineSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *pipelineSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	projectName := data.ProjectName.ValueString()
	var pipeline *buddy.Pipeline
	var err error
	if !data.PipelineId.IsNull() && !data.PipelineId.IsUnknown() {
		var httpRes *http.Response
		pipeline, httpRes, err = s.client.PipelineService.Get(domain, projectName, int(data.PipelineId.ValueInt64()))
		if err != nil {
			if util.IsResourceNotFound(httpRes, err) {
				resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("pipeline"))
				return
			}
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get pipeline", err))
			return
		}
	} else {
		name := data.Name.ValueString()
		var pipelines *buddy.Pipelines
		pipelines, _, err = s.client.PipelineService.GetListAll(domain, projectName)
		if err != nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get pipelines", err))
			return
		}
		for _, p := range pipelines.Pipelines {
			if p.Name == name {
				pipeline = p
				break
			}
		}
		if pipeline == nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("pipeline"))
			return
		}
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, pipeline)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

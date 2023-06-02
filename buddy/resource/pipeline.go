package resource

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
)

var (
	_ resource.Resource                = &pipelineResource{}
	_ resource.ResourceWithConfigure   = &pipelineResource{}
	_ resource.ResourceWithImportState = &pipelineResource{}
)

func NewPipelineResource() resource.Resource {
	return &pipelineResource{}
}

type pipelineResource struct {
	client *buddy.Client
}

type pipelineResourceModel struct {
	ID                        types.String `tfsdk:"id"`
	Domain                    types.String `tfsdk:"domain"`
	ProjectName               types.String `tfsdk:"project_name"`
	HtmlUrl                   types.String `tfsdk:"html_url"`
	PipelineId                types.Int64  `tfsdk:"pipeline_id"`
	Name                      types.String `tfsdk:"name"`
	DefinitionSource          types.String `tfsdk:"definition_source"`
	RemoteProjectName         types.String `tfsdk:"remote_project_name"`
	RemoteBranch              types.String `tfsdk:"remote_branch"`
	RemotePath                types.String `tfsdk:"remote_path"`
	RemoteParameters          types.Set    `tfsdk:"remote_parameter"`
	On                        types.String `tfsdk:"on"`
	Priority                  types.String `tfsdk:"priority"`
	FetchAllRefs              types.Bool   `tfsdk:"fetch_all_refs"`
	AlwaysFromScratch         types.Bool   `tfsdk:"always_from_scratch"`
	Disabled                  types.Bool   `tfsdk:"disabled"`
	DisablingReason           types.String `tfsdk:"disabling_reason"`
	FailOnPrepareEnvWarning   types.Bool   `tfsdk:"fail_on_prepare_env_warning"`
	AutoClearCache            types.Bool   `tfsdk:"auto_clear_cache"`
	NoSkipToMostRecent        types.Bool   `tfsdk:"no_skip_to_most_recent"`
	DoNotCreateCommitStatus   types.Bool   `tfsdk:"do_not_create_commit_status"`
	StartDate                 types.String `tfsdk:"start_date"`
	Delay                     types.Int64  `tfsdk:"delay"`
	CloneDepth                types.Int64  `tfsdk:"clone_depth"`
	Cron                      types.String `tfsdk:"cron"`
	Paused                    types.Bool   `tfsdk:"paused"`
	IgnoreFailOnProjectStatus types.Bool   `tfsdk:"ignore_fail_on_project_status"`
	ExecutionMessageTemplate  types.String `tfsdk:"execution_message_template"`
	Worker                    types.String `tfsdk:"worker"`
	TargetSiteUrl             types.String `tfsdk:"target_site_url"`
	LastExecutionStatus       types.String `tfsdk:"last_execution_status"`
	LastExecutionRevision     types.String `tfsdk:"last_execution_revision"`
	CreateDate                types.String `tfsdk:"create_date"`
	Creator                   types.Set    `tfsdk:"creator"`
	Project                   types.Set    `tfsdk:"project"`
	Refs                      types.Set    `tfsdk:"refs"`
	Tags                      types.Set    `tfsdk:"tags"`
	Events                    types.Set    `tfsdk:"event"`
	TriggerConditions         types.Set    `tfsdk:"trigger_condition"`
	Permissions               types.Set    `tfsdk:"permissions"`
}

func (r *pipelineResourceModel) loadAPI(ctx context.Context, domain string, projectName string, pipeline *buddy.Pipeline) diag.Diagnostics {
	var diags diag.Diagnostics
	r.ID = types.StringValue(util.ComposeTripleId(domain, projectName, strconv.Itoa(pipeline.Id)))
	r.Domain = types.StringValue(domain)
	r.ProjectName = types.StringValue(projectName)
	r.HtmlUrl = types.StringValue(pipeline.HtmlUrl)
	r.Name = types.StringValue(pipeline.Name)
	r.PipelineId = types.Int64Value(int64(pipeline.Id))
	r.On = types.StringValue(pipeline.On)
	refs, d := types.SetValueFrom(ctx, types.StringType, &pipeline.Refs)
	diags.Append(d...)
	r.Refs = refs
	r.Priority = types.StringValue(pipeline.Priority)
	r.LastExecutionStatus = types.StringValue(pipeline.LastExecutionStatus)
	r.LastExecutionRevision = types.StringValue(pipeline.LastExecutionRevision)
	r.CreateDate = types.StringValue(pipeline.CreateDate)
	projectSet := []*buddy.Project{pipeline.Project}
	project, d := util.ProjectsModelFromApi(ctx, &projectSet)
	diags.Append(d...)
	r.Project = project
	creatorSet := []*buddy.Member{pipeline.Creator}
	creator, d := util.MembersModelFromApi(ctx, &creatorSet)
	diags.Append(d...)
	r.Creator = creator
	r.AlwaysFromScratch = types.BoolValue(pipeline.AlwaysFromScratch)
	r.IgnoreFailOnProjectStatus = types.BoolValue(pipeline.IgnoreFailOnProjectStatus)
	r.NoSkipToMostRecent = types.BoolValue(pipeline.NoSkipToMostRecent)
	r.AutoClearCache = types.BoolValue(pipeline.AutoClearCache)
	r.FetchAllRefs = types.BoolValue(pipeline.FetchAllRefs)
	r.FailOnPrepareEnvWarning = types.BoolValue(pipeline.FailOnPrepareEnvWarning)
	r.DoNotCreateCommitStatus = types.BoolValue(pipeline.DoNotCreateCommitStatus)
	r.DefinitionSource = types.StringValue(util.GetPipelineDefinitionSource(pipeline))
	r.RemotePath = types.StringValue(pipeline.RemotePath)
	r.RemoteBranch = types.StringValue(pipeline.RemoteBranch)
	r.RemoteProjectName = types.StringValue(pipeline.RemoteProjectName)
	r.StartDate = types.StringValue(pipeline.StartDate)
	r.Delay = types.Int64Value(int64(pipeline.Delay))
	r.Paused = types.BoolValue(pipeline.Paused)
	r.Cron = types.StringValue(pipeline.Cron)
	r.Disabled = types.BoolValue(pipeline.Disabled)
	r.DisablingReason = types.StringValue(pipeline.DisabledReason)
	r.CloneDepth = types.Int64Value(int64(pipeline.CloneDepth))
	r.ExecutionMessageTemplate = types.StringValue(pipeline.ExecutionMessageTemplate)
	r.TargetSiteUrl = types.StringValue(pipeline.TargetSiteUrl)
	return diags
}

func (r *pipelineResourceModel) decomposeId() (string, string, int, error) {
	domain, projectName, pid, err := util.DecomposeTripleId(r.ID.ValueString())
	if err != nil {
		return "", "", 0, err
	}
	pipelineId, err := strconv.Atoi(pid)
	if err != nil {
		return "", "", 0, err
	}
	return domain, projectName, pipelineId, nil
}

func (r *pipelineResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pipeline"
}

func (r *pipelineResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage a pipeline\n\n" +
			"Token scopes required: `WORKSPACE`, `EXECUTION_MANAGE`, `EXECUTION_INFO`",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Terraform resource identifier for this item",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The workspace's URL handle",
				Required:            true,
				Validators:          util.StringValidatorsDomain(),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_name": schema.StringAttribute{
				MarkdownDescription: "The project's name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The pipeline's URL",
				Computed:            true,
			},
			"pipeline_id": schema.Int64Attribute{
				MarkdownDescription: "The pipeline's ID",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The pipeline's name",
				Required:            true,
			},
			"definition_source": schema.StringAttribute{
				MarkdownDescription: "The pipeline's definition source. Allowed: `LOCAL`, `REMOTE`",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: stringdefault.StaticString(buddy.PipelineDefinitionSourceLocal),
				Validators: []validator.String{
					stringvalidator.OneOf(
						buddy.PipelineDefinitionSourceLocal,
						buddy.PipelineDefinitionSourceRemote,
					),
				},
			},
			"remote_project_name": schema.StringAttribute{
				MarkdownDescription: "The pipeline's remote definition project name. Set it if `definition_source: REMOTE`",
				Optional:            true,
				Computed:            true,
			},
			"remote_branch": schema.StringAttribute{
				MarkdownDescription: "The pipeline's remote definition branch name. Set it if `definition_source: REMOTE`",
				Optional:            true,
				Computed:            true,
			},
			"remote_path": schema.StringAttribute{
				MarkdownDescription: "The pipeline's remote definition path. Set it if `definition_source: REMOTE`",
				Optional:            true,
				Computed:            true,
			},
			"on": schema.StringAttribute{
				MarkdownDescription: "The pipeline's trigger mode. Required when not using remote definition. Allowed: `CLICK`, `EVENT`, `SCHEDULE`",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						buddy.PipelineOnClick,
						buddy.PipelineOnEvent,
						buddy.PipelineOnSchedule,
					),
				},
			},
			"priority": schema.StringAttribute{
				MarkdownDescription: "The pipeline's priority. Allowed: `LOW`, `NORMAL`, `HIGH`",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						buddy.PipelinePriorityHigh,
						buddy.PipelinePriorityNormal,
						buddy.PipelinePriorityLow,
					),
				},
			},
			"fetch_all_refs": schema.BoolAttribute{
				MarkdownDescription: "Defines whether or not fetch all refs from repository",
				Optional:            true,
				Computed:            true,
			},
			"always_from_scratch": schema.BoolAttribute{
				MarkdownDescription: "Defines whether or not to upload everything from scratch on every run",
				Optional:            true,
				Computed:            true,
			},
			"disabled": schema.BoolAttribute{
				MarkdownDescription: "Defines wheter or not the pipeline can be run",
				Optional:            true,
				Computed:            true,
			},
			"disabling_reason": schema.StringAttribute{
				MarkdownDescription: "The pipeline's disabling reason",
				Optional:            true,
				Computed:            true,
			},
			"fail_on_prepare_env_warning": schema.BoolAttribute{
				MarkdownDescription: "Defines either or not run should fail if any warning occurs in prepare environment",
				Optional:            true,
				Computed:            true,
			},
			"auto_clear_cache": schema.BoolAttribute{
				MarkdownDescription: "Defines whether or not to automatically clear cache before running the pipeline",
				Optional:            true,
				Computed:            true,
			},
			"no_skip_to_most_recent": schema.BoolAttribute{
				MarkdownDescription: "Defines whether or not to skip run to the most recent run",
				Optional:            true,
				Computed:            true,
			},
			"do_not_create_commit_status": schema.BoolAttribute{
				MarkdownDescription: "Defines whether or not to omit sending commit statuses to GitHub or GitLab upon execution",
				Optional:            true,
				Computed:            true,
			},
			"start_date": schema.StringAttribute{
				MarkdownDescription: "The pipeline's start date. Required if the pipeline is set to `on: SCHEDULE` and no `cron` is specified. Format: `2016-11-18T12:38:16.000Z`",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("cron"),
					}...),
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("delay"),
					}...),
				},
			},
			"delay": schema.Int64Attribute{
				MarkdownDescription: "The pipeline's runs interval (in minutes). Required if the pipeline is set to `on: SCHEDULE` and no `cron` is specified",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.ConflictsWith(path.Expressions{
						path.MatchRoot("cron"),
					}...),
					int64validator.AlsoRequires(path.Expressions{
						path.MatchRoot("start_date"),
					}...),
				},
			},
			"clone_depth": schema.Int64Attribute{
				MarkdownDescription: "The pipeline's filesystem clone depth. Creates a shallow clone with a history truncated to the specified number of commits",
				Optional:            true,
				Computed:            true,
			},
			"cron": schema.StringAttribute{
				MarkdownDescription: "The pipeline's CRON expression. Required if the pipeline is set to `on: SCHEDULE` and neither `start_date` nor `delay` is specified",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("delay"),
						path.MatchRoot("start_date"),
					}...),
				},
			},
			"paused": schema.BoolAttribute{
				MarkdownDescription: "Is the pipeline's run paused. Restricted to `on: SCHEDULE`",
				Optional:            true,
				Computed:            true,
			},
			"ignore_fail_on_project_status": schema.BoolAttribute{
				MarkdownDescription: "If set to true the status of a given pipeline will be ignored on the projects' dashboard",
				Optional:            true,
				Computed:            true,
			},
			"execution_message_template": schema.StringAttribute{
				MarkdownDescription: "The pipeline's run title. Default: `$BUDDY_EXECUTION_REVISION_SUBJECT`",
				Optional:            true,
				Computed:            true,
			},
			"worker": schema.StringAttribute{
				MarkdownDescription: "The pipeline's worker name. Only for `Buddy Enterprise`",
				Optional:            true,
			},
			"target_site_url": schema.StringAttribute{
				MarkdownDescription: "The pipeline's website target URL",
				Optional:            true,
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
			"create_date": schema.StringAttribute{
				MarkdownDescription: "The pipeline's date of creation",
				Computed:            true,
			},
			// set for compatibility
			"creator": schema.SetNestedAttribute{
				MarkdownDescription: "The pipeline's creator",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.ResourceMemberModelAttributes(),
				},
			},
			// set for compatibility
			"project": schema.SetNestedAttribute{
				MarkdownDescription: "The pipeline's project",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"html_url": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"display_name": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"refs": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The pipeline's list of refs. Set it if `on: CLICK`",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Set{
					setvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("event"),
					}...),
				},
			},
			"tags": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The pipeline's list of tags. Only for `Buddy Enterprise`",
				Optional:            true,
			},
		},
		Blocks: map[string]schema.Block{
			// singular form for compatibility
			"event": schema.SetNestedBlock{
				MarkdownDescription: "The pipeline's list of events. Set it if `on: EVENT`",
				NestedObject: schema.NestedBlockObject{
					Attributes: util.ResourceEventModelAttributes(),
				},
				Validators: []validator.Set{
					setvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("refs"),
					}...),
				},
			},
			"permissions": schema.SetNestedBlock{
				MarkdownDescription: "The pipeline's permissions",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"others": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									buddy.PipelinePermissionDefault,
									buddy.PipelinePermissionDenied,
									buddy.PipelinePermissionReadOnly,
									buddy.PipelinePermissionRunOnly,
									buddy.PipelinePermissionReadWrite,
								),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"user": schema.SetNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: util.PipelinePermissionsAccessModelAttributes(),
							},
						},
						"group": schema.SetNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: util.PipelinePermissionsAccessModelAttributes(),
							},
						},
					},
				},
				Validators: []validator.Set{
					setvalidator.SizeAtMost(1),
				},
			},
			// singular form for compatibility
			"remote_parameter": schema.SetNestedBlock{
				MarkdownDescription: "The pipeline's remote definition parameters. Set it if `definition_source: REMOTE`",
				NestedObject: schema.NestedBlockObject{
					Attributes: util.ResourceRemoteParameterModelAttributes(),
				},
			},
			// singular form for compatibility
			"trigger_condition": schema.SetNestedBlock{
				MarkdownDescription: "The pipeline's list of trigger conditions",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"condition": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									buddy.PipelineTriggerConditionOnChange,
									buddy.PipelineTriggerConditionOnChangeAtPath,
									buddy.PipelineTriggerConditionVarIs,
									buddy.PipelineTriggerConditionVarIsNot,
									buddy.PipelineTriggerConditionVarContains,
									buddy.PipelineTriggerConditionVarNotContains,
									buddy.PipelineTriggerConditionDateTime,
									buddy.PipelineTriggerConditionSuccessPipeline,
								),
							},
						},
						"paths": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
						},
						"variable_key": schema.StringAttribute{
							Optional: true,
						},
						"variable_value": schema.StringAttribute{
							Optional: true,
						},
						"hours": schema.SetAttribute{
							ElementType: types.Int64Type,
							Optional:    true,
						},
						"days": schema.SetAttribute{
							ElementType: types.Int64Type,
							Optional:    true,
						},
						"zone_id": schema.StringAttribute{
							Optional: true,
						},
						"project_name": schema.StringAttribute{
							Optional: true,
						},
						"pipeline_name": schema.StringAttribute{
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func (r *pipelineResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *pipelineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *pipelineResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	projectName := data.ProjectName.ValueString()
	ops := buddy.PipelineOps{
		Name:                    data.Name.ValueStringPointer(),
		FailOnPrepareEnvWarning: data.FailOnPrepareEnvWarning.ValueBoolPointer(),
		FetchAllRefs:            data.FetchAllRefs.ValueBoolPointer(),
	}
	if !data.Permissions.IsNull() && !data.Permissions.IsUnknown() {
		permissions, d := util.PipelinePermissionsModelToApi(ctx, &data.Permissions)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Permissions = permissions
	}
	if !data.Refs.IsNull() && !data.Refs.IsUnknown() {
		refs, d := util.StringSetToApi(ctx, &data.Refs)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Refs = refs
	}
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		tags, d := util.StringSetToApi(ctx, &data.Tags)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Tags = tags
	}
	if !data.Events.IsNull() && !data.Events.IsUnknown() {
		events, d := util.EventsModelToApi(ctx, &data.Events)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Events = events
	}
	if !data.On.IsNull() && !data.On.IsUnknown() {
		ops.On = data.On.ValueStringPointer()
	}
	if !data.AlwaysFromScratch.IsNull() && !data.AlwaysFromScratch.IsUnknown() {
		ops.AlwaysFromScratch = data.AlwaysFromScratch.ValueBoolPointer()
	}
	if !data.Priority.IsNull() && !data.Priority.IsUnknown() {
		ops.Priority = data.Priority.ValueStringPointer()
	}
	if !data.AutoClearCache.IsNull() && !data.AutoClearCache.IsUnknown() {
		ops.AutoClearCache = data.AutoClearCache.ValueBoolPointer()
	}
	if !data.NoSkipToMostRecent.IsNull() && !data.NoSkipToMostRecent.IsUnknown() {
		ops.NoSkipToMostRecent = data.NoSkipToMostRecent.ValueBoolPointer()
	}
	if !data.DoNotCreateCommitStatus.IsNull() && !data.DoNotCreateCommitStatus.IsUnknown() {
		ops.DoNotCreateCommitStatus = data.DoNotCreateCommitStatus.ValueBoolPointer()
	}
	if !data.StartDate.IsNull() && !data.StartDate.IsUnknown() {
		ops.StartDate = data.StartDate.ValueStringPointer()
	}
	if !data.Delay.IsNull() && !data.Delay.IsUnknown() {
		ops.Delay = util.PointerInt(data.Delay.ValueInt64())
	}
	if !data.CloneDepth.IsNull() && !data.CloneDepth.IsUnknown() {
		ops.CloneDepth = util.PointerInt(data.CloneDepth.ValueInt64())
	}
	if !data.Cron.IsNull() && !data.Cron.IsUnknown() {
		ops.Cron = data.Cron.ValueStringPointer()
	}
	if !data.Paused.IsNull() && !data.Paused.IsUnknown() {
		ops.Paused = data.Paused.ValueBoolPointer()
	}
	if !data.IgnoreFailOnProjectStatus.IsNull() && !data.IgnoreFailOnProjectStatus.IsUnknown() {
		ops.IgnoreFailOnProjectStatus = data.IgnoreFailOnProjectStatus.ValueBoolPointer()
	}
	if !data.ExecutionMessageTemplate.IsNull() && !data.ExecutionMessageTemplate.IsUnknown() {
		ops.ExecutionMessageTemplate = data.ExecutionMessageTemplate.ValueStringPointer()
	}
	if !data.Worker.IsNull() && !data.Worker.IsUnknown() {
		ops.Worker = data.Worker.ValueStringPointer()
	}
	if !data.TargetSiteUrl.IsNull() && !data.TargetSiteUrl.IsUnknown() {
		ops.TargetSiteUrl = data.TargetSiteUrl.ValueStringPointer()
	}
	if !data.DefinitionSource.IsNull() && !data.DefinitionSource.IsUnknown() {
		ops.DefinitionSource = data.DefinitionSource.ValueStringPointer()
	}
	if !data.RemotePath.IsNull() && !data.RemotePath.IsUnknown() {
		ops.RemotePath = data.RemotePath.ValueStringPointer()
	}
	if !data.RemoteBranch.IsNull() && !data.RemoteBranch.IsUnknown() {
		ops.RemoteBranch = data.RemoteBranch.ValueStringPointer()
	}
	if !data.RemoteProjectName.IsNull() && !data.RemoteProjectName.IsUnknown() {
		ops.RemoteProjectName = data.RemoteProjectName.ValueStringPointer()
	}
	if !data.RemoteParameters.IsNull() && !data.RemoteParameters.IsUnknown() {
		remoteParams, d := util.RemoteParametersModelToApi(ctx, &data.RemoteParameters)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.RemoteParameters = remoteParams
	}
	if !data.Disabled.IsNull() && !data.Disabled.IsUnknown() {
		ops.Disabled = data.Disabled.ValueBoolPointer()
	}
	if !data.DisablingReason.IsNull() && !data.DisablingReason.IsUnknown() {
		ops.DisabledReason = data.DisablingReason.ValueStringPointer()
	}
	if !data.TriggerConditions.IsNull() && !data.TriggerConditions.IsUnknown() {
		tc, d := util.TriggerConditionsModelToApi(ctx, &data.TriggerConditions)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.TriggerConditions = tc
	}
	pipeline, _, err := r.client.PipelineService.Create(domain, projectName, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create pipeline", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, pipeline)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *pipelineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *pipelineResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain, projectName, pipelineId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("pipeline", err))
		return
	}
	pipeline, httpResp, err := r.client.PipelineService.Get(domain, projectName, pipelineId)
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get pipeline", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, pipeline)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *pipelineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *pipelineResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, projectName, pipelineId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("pipeline", err))
		return
	}
	ops := buddy.PipelineOps{
		Name: data.Name.ValueStringPointer(),
	}
	if !data.Permissions.IsNull() && !data.Permissions.IsUnknown() {
		permissions, d := util.PipelinePermissionsModelToApi(ctx, &data.Permissions)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Permissions = permissions
	}
	if !data.Refs.IsNull() && !data.Refs.IsUnknown() {
		refs, d := util.StringSetToApi(ctx, &data.Refs)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Refs = refs
	}
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		tags, d := util.StringSetToApi(ctx, &data.Tags)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Tags = tags
	}
	if !data.Events.IsNull() && !data.Events.IsUnknown() {
		events, d := util.EventsModelToApi(ctx, &data.Events)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Events = events
	}
	if !data.TriggerConditions.IsNull() && !data.TriggerConditions.IsUnknown() {
		tc, d := util.TriggerConditionsModelToApi(ctx, &data.TriggerConditions)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.TriggerConditions = tc
	}
	if !data.On.IsNull() && !data.On.IsUnknown() {
		ops.On = data.On.ValueStringPointer()
	}
	if !data.AlwaysFromScratch.IsNull() && !data.AlwaysFromScratch.IsUnknown() {
		ops.AlwaysFromScratch = data.AlwaysFromScratch.ValueBoolPointer()
	}
	if !data.Priority.IsNull() && !data.Priority.IsUnknown() {
		ops.Priority = data.Priority.ValueStringPointer()
	}
	if !data.FailOnPrepareEnvWarning.IsNull() && !data.FailOnPrepareEnvWarning.IsUnknown() {
		ops.FailOnPrepareEnvWarning = data.FailOnPrepareEnvWarning.ValueBoolPointer()
	}
	if !data.FetchAllRefs.IsNull() && !data.FetchAllRefs.IsUnknown() {
		ops.FetchAllRefs = data.FetchAllRefs.ValueBoolPointer()
	}
	if !data.AutoClearCache.IsNull() && !data.AutoClearCache.IsUnknown() {
		ops.AutoClearCache = data.AutoClearCache.ValueBoolPointer()
	}
	if !data.NoSkipToMostRecent.IsNull() && !data.NoSkipToMostRecent.IsUnknown() {
		ops.NoSkipToMostRecent = data.NoSkipToMostRecent.ValueBoolPointer()
	}
	if !data.DoNotCreateCommitStatus.IsNull() && !data.DoNotCreateCommitStatus.IsUnknown() {
		ops.DoNotCreateCommitStatus = data.DoNotCreateCommitStatus.ValueBoolPointer()
	}
	if !data.StartDate.IsNull() && !data.StartDate.IsUnknown() {
		ops.StartDate = data.StartDate.ValueStringPointer()
	}
	if !data.Delay.IsNull() && !data.Delay.IsUnknown() {
		ops.Delay = util.PointerInt(data.Delay.ValueInt64())
	}
	if !data.CloneDepth.IsNull() && !data.CloneDepth.IsUnknown() {
		ops.CloneDepth = util.PointerInt(data.CloneDepth.ValueInt64())
	}
	if !data.Cron.IsNull() && !data.Cron.IsUnknown() {
		ops.Cron = data.Cron.ValueStringPointer()
	}
	if !data.Paused.IsNull() && !data.Paused.IsUnknown() {
		ops.Paused = data.Paused.ValueBoolPointer()
	}
	if !data.IgnoreFailOnProjectStatus.IsNull() && !data.IgnoreFailOnProjectStatus.IsUnknown() {
		ops.IgnoreFailOnProjectStatus = data.IgnoreFailOnProjectStatus.ValueBoolPointer()
	}
	if !data.ExecutionMessageTemplate.IsNull() && !data.ExecutionMessageTemplate.IsUnknown() {
		ops.ExecutionMessageTemplate = data.ExecutionMessageTemplate.ValueStringPointer()
	}
	if !data.Worker.IsNull() && !data.Worker.IsUnknown() {
		ops.Worker = data.Worker.ValueStringPointer()
	}
	if !data.TargetSiteUrl.IsNull() && !data.TargetSiteUrl.IsUnknown() {
		ops.TargetSiteUrl = data.TargetSiteUrl.ValueStringPointer()
	}
	if !data.DefinitionSource.IsNull() && !data.DefinitionSource.IsUnknown() {
		ops.DefinitionSource = data.DefinitionSource.ValueStringPointer()
	}
	if !data.RemotePath.IsNull() && !data.RemotePath.IsUnknown() {
		ops.RemotePath = data.RemotePath.ValueStringPointer()
	}
	if !data.RemoteBranch.IsNull() && !data.RemoteBranch.IsUnknown() {
		ops.RemoteBranch = data.RemoteBranch.ValueStringPointer()
	}
	if !data.RemoteProjectName.IsNull() && !data.RemoteProjectName.IsUnknown() {
		ops.RemoteProjectName = data.RemoteProjectName.ValueStringPointer()
	}
	if !data.RemoteParameters.IsNull() && !data.RemoteParameters.IsUnknown() {
		remoteParams, d := util.RemoteParametersModelToApi(ctx, &data.RemoteParameters)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.RemoteParameters = remoteParams
	}
	if !data.Disabled.IsNull() && !data.Disabled.IsUnknown() {
		ops.Disabled = data.Disabled.ValueBoolPointer()
	}
	if !data.DisablingReason.IsNull() && !data.DisablingReason.IsUnknown() {
		ops.DisabledReason = data.DisablingReason.ValueStringPointer()
	}
	pipeline, _, err := r.client.PipelineService.Update(domain, projectName, pipelineId, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update pipeline", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, pipeline)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *pipelineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *pipelineResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, projectName, pipelineId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("pipeline", err))
		return
	}
	_, err = r.client.PipelineService.Delete(domain, projectName, pipelineId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete pipeline", err))
	}
}

func (r *pipelineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

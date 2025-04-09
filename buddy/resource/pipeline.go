package resource

import (
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
	"terraform-provider-buddy/buddy/util"
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
	GitConfigRef              types.String `tfsdk:"git_config_ref"`
	GitConfig                 types.Object `tfsdk:"git_config"`
	DefinitionSource          types.String `tfsdk:"definition_source"`
	RemoteProjectName         types.String `tfsdk:"remote_project_name"`
	RemoteBranch              types.String `tfsdk:"remote_branch"`
	RemotePath                types.String `tfsdk:"remote_path"`
	RemoteParameters          types.Set    `tfsdk:"remote_parameter"`
	Cpu                       types.String `tfsdk:"cpu"`
	Priority                  types.String `tfsdk:"priority"`
	FetchAllRefs              types.Bool   `tfsdk:"fetch_all_refs"`
	AlwaysFromScratch         types.Bool   `tfsdk:"always_from_scratch"`
	ConcurrentPipelineRuns    types.Bool   `tfsdk:"concurrent_pipeline_runs"`
	DescriptionRequired       types.Bool   `tfsdk:"description_required"`
	GitChangesetBase          types.String `tfsdk:"git_changeset_base"`
	FilesystemChangesetBase   types.String `tfsdk:"filesystem_changeset_base"`
	Disabled                  types.Bool   `tfsdk:"disabled"`
	DisablingReason           types.String `tfsdk:"disabling_reason"`
	FailOnPrepareEnvWarning   types.Bool   `tfsdk:"fail_on_prepare_env_warning"`
	AutoClearCache            types.Bool   `tfsdk:"auto_clear_cache"`
	NoSkipToMostRecent        types.Bool   `tfsdk:"no_skip_to_most_recent"`
	DoNotCreateCommitStatus   types.Bool   `tfsdk:"do_not_create_commit_status"`
	CloneDepth                types.Int64  `tfsdk:"clone_depth"`
	Paused                    types.Bool   `tfsdk:"paused"`
	PauseOnRepeatedFailures   types.Int64  `tfsdk:"pause_on_repeated_failures"`
	IgnoreFailOnProjectStatus types.Bool   `tfsdk:"ignore_fail_on_project_status"`
	ExecutionMessageTemplate  types.String `tfsdk:"execution_message_template"`
	Worker                    types.String `tfsdk:"worker"`
	TargetSiteUrl             types.String `tfsdk:"target_site_url"`
	ManageVariablesByYaml     types.Bool   `tfsdk:"manage_variables_by_yaml"`
	ManagePermissionsByYaml   types.Bool   `tfsdk:"manage_permissions_by_yaml"`
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
	r.GitConfigRef = types.StringValue(pipeline.GitConfigRef)
	gitConfig, d := util.GitConfigModelFromApi(ctx, pipeline.GitConfig)
	diags.Append(d...)
	r.GitConfig = gitConfig
	r.PipelineId = types.Int64Value(int64(pipeline.Id))
	r.Cpu = types.StringValue(pipeline.Cpu)
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
	r.DescriptionRequired = types.BoolValue(pipeline.DescriptionRequired)
	r.GitChangesetBase = types.StringValue(pipeline.GitChangesetBase)
	r.FilesystemChangesetBase = types.StringValue(pipeline.FilesystemChangesetBase)
	r.ConcurrentPipelineRuns = types.BoolValue(pipeline.ConcurrentPipelineRuns)
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
	r.Paused = types.BoolValue(pipeline.Paused)
	r.PauseOnRepeatedFailures = types.Int64Value(int64(pipeline.PauseOnRepeatedFailures))
	r.Disabled = types.BoolValue(pipeline.Disabled)
	r.DisablingReason = types.StringValue(pipeline.DisabledReason)
	r.CloneDepth = types.Int64Value(int64(pipeline.CloneDepth))
	r.ExecutionMessageTemplate = types.StringValue(pipeline.ExecutionMessageTemplate)
	r.TargetSiteUrl = types.StringValue(pipeline.TargetSiteUrl)
	r.ManagePermissionsByYaml = types.BoolValue(pipeline.ManagePermissionsByYaml)
	r.ManageVariablesByYaml = types.BoolValue(pipeline.ManageVariablesByYaml)
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
			"git_config_ref": schema.StringAttribute{
				MarkdownDescription: "The pipeline's GIT configuration type. Allowed: `NONE`, `FIXED`, `DYNAMIC`",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(buddy.PipelineGitConfigRefNone),
				Validators: []validator.String{
					stringvalidator.OneOf(
						buddy.PipelineGitConfigRefNone,
						buddy.PipelineGitConfigRefFixed,
						buddy.PipelineGitConfigRefDynamic,
					),
				},
			},
			"git_config": schema.ObjectAttribute{
				MarkdownDescription: "The pipeline's GIT configuration spec for `git_config_ref` = `FIXED`",
				Optional:            true,
				Computed:            true,
				AttributeTypes:      util.GitConfigModelAttrs(),
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
			"cpu": schema.StringAttribute{
				MarkdownDescription: "The pipeline's cpu. Allowed: `X64`, `ARM`",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						buddy.PipelineCpuX64,
						buddy.PipelineCpuArm,
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
			"concurrent_pipeline_runs": schema.BoolAttribute{
				MarkdownDescription: "Defines whether or not pipeline can be run concurrently",
				Optional:            true,
				Computed:            true,
			},
			"description_required": schema.BoolAttribute{
				MarkdownDescription: "Defines whether or not pipeline's execution must be commented",
				Optional:            true,
				Computed:            true,
			},
			"git_changeset_base": schema.StringAttribute{
				MarkdownDescription: "Defines pipeline's GIT changeset",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						buddy.PipelineGitChangeSetBaseLatestRun,
						buddy.PipelineGitChangeSetBaseLatestRunMatchingRef,
						buddy.PipelineGitChangeSetBasePullRequest,
					),
				},
			},
			"filesystem_changeset_base": schema.StringAttribute{
				MarkdownDescription: "Defines pipeline's filesystem changeset",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						buddy.PipelineFilesystemChangeSetBaseDateModified,
						buddy.PipelineFilesystemChangeSetBaseContents,
					),
				},
			},
			"disabled": schema.BoolAttribute{
				MarkdownDescription: "Defines whether or not the pipeline can be run",
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
			"clone_depth": schema.Int64Attribute{
				MarkdownDescription: "The pipeline's filesystem clone depth. Creates a shallow clone with a history truncated to the specified number of commits",
				Optional:            true,
				Computed:            true,
			},
			"paused": schema.BoolAttribute{
				MarkdownDescription: "Is the pipeline's run paused. Restricted schedule",
				Optional:            true,
				Computed:            true,
			},
			"pause_on_repeated_failures": schema.Int64Attribute{
				MarkdownDescription: "The pipeline's max failed executions before it is paused. Restricted to schedule",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.AtMost(100),
					int64validator.AtLeast(1),
				},
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
			"manage_variables_by_yaml": schema.BoolAttribute{
				MarkdownDescription: "If set to true pipeline variables will be managed by yaml",
				Optional:            true,
				Computed:            true,
			},
			"manage_permissions_by_yaml": schema.BoolAttribute{
				MarkdownDescription: "If set to true pipeline permissions will be managed by yaml",
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
				MarkdownDescription: "The pipeline's list of refs for manual mode",
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
				MarkdownDescription: "The pipeline's list of events",
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
									buddy.PipelineTriggerConditionTriggeringUserIsNotInGroup,
									buddy.PipelineTriggerConditionTriggeringUserIsInGroup,
									buddy.PipelineTriggerConditionTriggeringUserIs,
									buddy.PipelineTriggerConditionTriggeringUserIsNot,
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
						"timezone": schema.StringAttribute{
							Optional: true,
						},
						"project_name": schema.StringAttribute{
							Optional: true,
						},
						"pipeline_name": schema.StringAttribute{
							Optional: true,
						},
						"trigger_user": schema.StringAttribute{
							Optional: true,
						},
						"trigger_group": schema.StringAttribute{
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
	if !data.Cpu.IsNull() && !data.Cpu.IsUnknown() {
		ops.Cpu = data.Cpu.ValueStringPointer()
	}
	if !data.DescriptionRequired.IsNull() && !data.DescriptionRequired.IsUnknown() {
		ops.DescriptionRequired = data.DescriptionRequired.ValueBoolPointer()
	}
	if !data.GitChangesetBase.IsNull() && !data.GitChangesetBase.IsUnknown() {
		ops.GitChangesetBase = data.GitChangesetBase.ValueStringPointer()
	}
	if !data.FilesystemChangesetBase.IsNull() && !data.FilesystemChangesetBase.IsUnknown() {
		ops.FilesystemChangesetBase = data.FilesystemChangesetBase.ValueStringPointer()
	}
	if !data.ConcurrentPipelineRuns.IsNull() && !data.ConcurrentPipelineRuns.IsUnknown() {
		ops.ConcurrentPipelineRuns = data.ConcurrentPipelineRuns.ValueBoolPointer()
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
	if !data.CloneDepth.IsNull() && !data.CloneDepth.IsUnknown() {
		ops.CloneDepth = util.PointerInt(data.CloneDepth.ValueInt64())
	}
	if !data.Paused.IsNull() && !data.Paused.IsUnknown() {
		ops.Paused = data.Paused.ValueBoolPointer()
	}
	if !data.PauseOnRepeatedFailures.IsNull() && !data.PauseOnRepeatedFailures.IsUnknown() {
		ops.PauseOnRepeatedFailures = util.PointerInt(data.PauseOnRepeatedFailures.ValueInt64())
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
	if !data.ManageVariablesByYaml.IsNull() && !data.ManageVariablesByYaml.IsUnknown() {
		ops.ManageVariablesByYaml = data.ManageVariablesByYaml.ValueBoolPointer()
	}
	if !data.ManagePermissionsByYaml.IsNull() && !data.ManagePermissionsByYaml.IsUnknown() {
		ops.ManagePermissionsByYaml = data.ManagePermissionsByYaml.ValueBoolPointer()
	}
	if !data.GitConfigRef.IsNull() && !data.GitConfigRef.IsUnknown() {
		ops.GitConfigRef = data.GitConfigRef.ValueStringPointer()
	}
	if !data.GitConfig.IsNull() && !data.GitConfig.IsUnknown() {
		gitConfig, d := util.GitConfigModelToApi(ctx, &data.GitConfig)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.GitConfig = gitConfig
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
	if !data.Cpu.IsNull() && !data.Cpu.IsUnknown() {
		ops.Cpu = data.Cpu.ValueStringPointer()
	}
	if !data.AlwaysFromScratch.IsNull() && !data.AlwaysFromScratch.IsUnknown() {
		ops.AlwaysFromScratch = data.AlwaysFromScratch.ValueBoolPointer()
	}
	if !data.ConcurrentPipelineRuns.IsNull() && !data.ConcurrentPipelineRuns.IsUnknown() {
		ops.ConcurrentPipelineRuns = data.ConcurrentPipelineRuns.ValueBoolPointer()
	}
	if !data.DescriptionRequired.IsNull() && !data.DescriptionRequired.IsUnknown() {
		ops.DescriptionRequired = data.DescriptionRequired.ValueBoolPointer()
	}
	if !data.FilesystemChangesetBase.IsNull() && !data.FilesystemChangesetBase.IsUnknown() {
		ops.FilesystemChangesetBase = data.FilesystemChangesetBase.ValueStringPointer()
	}
	if !data.GitChangesetBase.IsNull() && !data.GitChangesetBase.IsUnknown() {
		ops.GitChangesetBase = data.GitChangesetBase.ValueStringPointer()
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
	if !data.PauseOnRepeatedFailures.IsNull() && !data.PauseOnRepeatedFailures.IsUnknown() {
		ops.PauseOnRepeatedFailures = util.PointerInt(data.PauseOnRepeatedFailures.ValueInt64())
	}
	if !data.CloneDepth.IsNull() && !data.CloneDepth.IsUnknown() {
		ops.CloneDepth = util.PointerInt(data.CloneDepth.ValueInt64())
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
	if !data.ManageVariablesByYaml.IsNull() && !data.ManageVariablesByYaml.IsUnknown() {
		ops.ManageVariablesByYaml = data.ManageVariablesByYaml.ValueBoolPointer()
	}
	if !data.ManagePermissionsByYaml.IsNull() && !data.ManagePermissionsByYaml.IsUnknown() {
		ops.ManagePermissionsByYaml = data.ManagePermissionsByYaml.ValueBoolPointer()
	}
	if !data.GitConfigRef.IsNull() && !data.GitConfigRef.IsUnknown() {
		ops.GitConfigRef = data.GitConfigRef.ValueStringPointer()
	}
	if !data.GitConfig.IsNull() && !data.GitConfig.IsUnknown() {
		gitConfig, d := util.GitConfigModelToApi(ctx, &data.GitConfig)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.GitConfig = gitConfig
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

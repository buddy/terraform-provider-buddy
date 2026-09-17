package resource

import (
	"context"
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ resource.Resource                = &variableResource{}
	_ resource.ResourceWithConfigure   = &variableResource{}
	_ resource.ResourceWithImportState = &variableResource{}
)

func NewVariableResource() resource.Resource {
	return &variableResource{}
}

type variableResource struct {
	client *buddy.Client
}

type variableResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Domain         types.String `tfsdk:"domain"`
	Key            types.String `tfsdk:"key"`
	Value          types.String `tfsdk:"value"`
	Encrypted      types.Bool   `tfsdk:"encrypted"`
	ProjectName    types.String `tfsdk:"project_name"`
	PipelineId     types.Int64  `tfsdk:"pipeline_id"`
	ActionId       types.Int64  `tfsdk:"action_id"`
	EnvironmentId  types.String `tfsdk:"environment_id"`
	SandboxId      types.String `tfsdk:"sandbox_id"`
	Settable       types.Bool   `tfsdk:"settable"`
	Description    types.String `tfsdk:"description"`
	Note           types.String `tfsdk:"note"`
	AgentNote      types.String `tfsdk:"agent_note"`
	ValueProcessed types.String `tfsdk:"value_processed"`
	VariableId     types.Int64  `tfsdk:"variable_id"`

	RunOnlySettable      types.Bool   `tfsdk:"run_only_settable"`
	Disabled             types.Bool   `tfsdk:"disabled"`
	PipelinesAccessLevel types.String `tfsdk:"pipelines_access_level"`
	SandboxesAccessLevel types.String `tfsdk:"sandboxes_access_level"`
	AllowedPipeline      types.Set    `tfsdk:"allowed_pipeline"`
	AllowedSandboxes     types.Set    `tfsdk:"allowed_sandboxes"`
}

func (r *variableResourceModel) decomposeId() (string, int, error) {
	domain, vid, err := util.DecomposeDoubleId(r.ID.ValueString())
	if err != nil {
		return "", 0, err
	}
	variableId, err := strconv.Atoi(vid)
	if err != nil {
		return "", 0, err
	}
	return domain, variableId, nil
}

func (r *variableResourceModel) loadAPI(domain string, variable *buddy.Variable) {
	r.ID = types.StringValue(util.ComposeDoubleId(domain, strconv.Itoa(variable.Id)))
	r.Domain = types.StringValue(domain)
	r.Key = types.StringValue(variable.Key)
	r.Encrypted = types.BoolValue(variable.Encrypted)
	r.Settable = types.BoolValue(variable.Settable)
	r.Note = types.StringValue(variable.Note)
	r.Description = types.StringValue(variable.Note)
	r.AgentNote = types.StringValue(variable.AgentNote)
	r.ValueProcessed = types.StringValue(variable.Value)
	r.VariableId = types.Int64Value(int64(variable.Id))
	r.RunOnlySettable = types.BoolValue(variable.RunOnlySettable)
	r.Disabled = types.BoolValue(variable.Disabled)
	r.PipelinesAccessLevel = types.StringValue(variable.PipelinesAccessLevel)
	r.SandboxesAccessLevel = types.StringValue(variable.SandboxesAccessLevel)
	if variable.Project != nil {
		r.ProjectName = types.StringValue(variable.Project.Name)
	} else {
		r.ProjectName = types.StringNull()
	}
	if variable.Pipeline != nil {
		r.PipelineId = types.Int64Value(int64(variable.Pipeline.Id))
	} else {
		r.PipelineId = types.Int64Null()
	}
	if variable.Action != nil {
		r.ActionId = types.Int64Value(int64(variable.Action.Id))
	} else {
		r.ActionId = types.Int64Null()
	}
	if variable.Environment != nil {
		r.EnvironmentId = types.StringValue(variable.Environment.Id)
	} else {
		r.EnvironmentId = types.StringNull()
	}
	if variable.Sandbox != nil {
		r.SandboxId = types.StringValue(variable.Sandbox.Id)
	} else {
		r.SandboxId = types.StringNull()
	}
}

// applyCommonOps fills in the fields that behave the same on create and update. The scope
// fields are left out on purpose - the API rejects any attempt to change them on update.
func (r *variableResourceModel) applyCommonOps(ctx context.Context, ops *buddy.VariableOps) diag.Diagnostics {
	var diags diag.Diagnostics
	if !r.Encrypted.IsNull() && !r.Encrypted.IsUnknown() {
		ops.Encrypted = r.Encrypted.ValueBoolPointer()
	}
	if !r.Settable.IsNull() && !r.Settable.IsUnknown() {
		ops.Settable = r.Settable.ValueBoolPointer()
	}
	if !r.RunOnlySettable.IsNull() && !r.RunOnlySettable.IsUnknown() {
		ops.RunOnlySettable = r.RunOnlySettable.ValueBoolPointer()
	}
	if !r.Disabled.IsNull() && !r.Disabled.IsUnknown() {
		ops.Disabled = r.Disabled.ValueBoolPointer()
	}
	if !r.Note.IsNull() && !r.Note.IsUnknown() {
		ops.Note = r.Note.ValueStringPointer()
	} else if !r.Description.IsNull() && !r.Description.IsUnknown() {
		ops.Note = r.Description.ValueStringPointer()
	}
	if !r.AgentNote.IsNull() && !r.AgentNote.IsUnknown() {
		ops.AgentNote = r.AgentNote.ValueStringPointer()
	}
	if !r.PipelinesAccessLevel.IsNull() && !r.PipelinesAccessLevel.IsUnknown() {
		ops.PipelinesAccessLevel = r.PipelinesAccessLevel.ValueStringPointer()
	}
	if !r.SandboxesAccessLevel.IsNull() && !r.SandboxesAccessLevel.IsUnknown() {
		ops.SandboxesAccessLevel = r.SandboxesAccessLevel.ValueStringPointer()
	}
	if !r.AllowedPipeline.IsNull() && !r.AllowedPipeline.IsUnknown() {
		pips, d := util.VariableAllowedPipelinesModelToApi(ctx, &r.AllowedPipeline)
		diags.Append(d...)
		ops.AllowedPipelines = pips
	}
	if !r.AllowedSandboxes.IsNull() && !r.AllowedSandboxes.IsUnknown() {
		sbs, d := util.VariableAllowedSandboxesModelToApi(ctx, &r.AllowedSandboxes)
		diags.Append(d...)
		ops.AllowedSandboxes = sbs
	}
	return diags
}

func (r *variableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_variable"
}

func (r *variableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage a variable\n\n" +
			"Workspace administrator rights are required\n\n" +
			"Token scopes required: `VARIABLE_READ`, `VARIABLE_WRITE`, `VARIABLE_MANAGE`",
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
			"key": schema.StringAttribute{
				MarkdownDescription: "The variable's name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "The variable's value",
				Required:            true,
				Sensitive:           true,
			},
			"encrypted": schema.BoolAttribute{
				MarkdownDescription: "Is the variable's value encrypted",
				Optional:            true,
				Computed:            true,
			},
			"project_name": schema.StringAttribute{
				MarkdownDescription: "The variable's project name. Set for project scope",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("pipeline_id"),
						path.MatchRoot("action_id"),
						path.MatchRoot("environment_id"),
						path.MatchRoot("sandbox_id"),
					}...),
				},
			},
			"pipeline_id": schema.Int64Attribute{
				MarkdownDescription: "The variable's pipeline ID. Set for pipeline scope",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.ConflictsWith(path.Expressions{
						path.MatchRoot("project_name"),
						path.MatchRoot("action_id"),
						path.MatchRoot("environment_id"),
						path.MatchRoot("sandbox_id"),
					}...),
				},
			},
			"action_id": schema.Int64Attribute{
				MarkdownDescription: "The variable's action ID. Set for action scope",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.ConflictsWith(path.Expressions{
						path.MatchRoot("project_name"),
						path.MatchRoot("pipeline_id"),
						path.MatchRoot("environment_id"),
						path.MatchRoot("sandbox_id"),
					}...),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "The variable's environmental ID. Set for envrionment scope",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("project_name"),
						path.MatchRoot("pipeline_id"),
						path.MatchRoot("action_id"),
						path.MatchRoot("sandbox_id"),
					}...),
				},
			},
			"sandbox_id": schema.StringAttribute{
				MarkdownDescription: "The variable's sandbox ID. Set for sandbox scope",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("project_name"),
						path.MatchRoot("pipeline_id"),
						path.MatchRoot("action_id"),
						path.MatchRoot("environment_id"),
					}...),
				},
			},
			"settable": schema.BoolAttribute{
				MarkdownDescription: "Is the variable's value changeable",
				Optional:            true,
				Computed:            true,
			},
			"run_only_settable": schema.BoolAttribute{
				MarkdownDescription: "Can the variable's value be changed only by a running pipeline. Requires **settable** == true",
				Optional:            true,
				Computed:            true,
			},
			"disabled": schema.BoolAttribute{
				MarkdownDescription: "Defines whether or not the variable is passed to a pipeline",
				Optional:            true,
				Computed:            true,
			},
			"pipelines_access_level": schema.StringAttribute{
				MarkdownDescription: "The default access level for pipelines. Only for workspace and project scope",
				Validators: []validator.String{
					stringvalidator.OneOf(buddy.VariableAccessLevelUseOnly, buddy.VariableAccessLevelDenied),
				},
				Optional: true,
				Computed: true,
			},
			"sandboxes_access_level": schema.StringAttribute{
				MarkdownDescription: "The default access level for sandboxes. Only for workspace and project scope",
				Validators: []validator.String{
					stringvalidator.OneOf(buddy.VariableAccessLevelUseOnly, buddy.VariableAccessLevelDenied),
				},
				Optional: true,
				Computed: true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The variable's description",
				DeprecationMessage:  "Use note field instead",
				Optional:            true,
				Computed:            true,
				Validators:          util.StringValidatorsAlias("note"),
				PlanModifiers: []planmodifier.String{
					util.AliasPlanModifier("note"),
				},
			},
			"note": schema.StringAttribute{
				MarkdownDescription: "The variable's note",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					util.AliasPlanModifier("description"),
				},
			},
			"agent_note": schema.StringAttribute{
				MarkdownDescription: "The variable's agent note",
				Optional:            true,
				Computed:            true,
			},
			"value_processed": schema.StringAttribute{
				MarkdownDescription: "The variable's processed value. Encrypted if **encrypted** == true",
				Computed:            true,
				Sensitive:           true,
			},
			"variable_id": schema.Int64Attribute{
				MarkdownDescription: "The variable's ID",
				Computed:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"allowed_pipeline": schema.SetNestedBlock{
				MarkdownDescription: "List of exceptions from **pipelines_access_level**. Every rule must carry the opposite access level and only one form - whole pipeline or single **action** - may be used per pipeline. Only for workspace and project scope",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"project": schema.StringAttribute{
							MarkdownDescription: "The pipeline's project name",
							Required:            true,
						},
						"pipeline": schema.StringAttribute{
							MarkdownDescription: "The pipeline's name or identifier",
							Required:            true,
						},
						"action": schema.StringAttribute{
							MarkdownDescription: "The action's identifier. Set to limit the rule to a single action",
							Optional:            true,
						},
						"access_level": schema.StringAttribute{
							MarkdownDescription: "The pipeline's access level",
							Validators: []validator.String{
								stringvalidator.OneOf(buddy.VariableAccessLevelUseOnly, buddy.VariableAccessLevelDenied),
							},
							Required: true,
						},
					},
				},
			},
			"allowed_sandboxes": schema.SetNestedBlock{
				MarkdownDescription: "List of exceptions from **sandboxes_access_level**. Every rule must carry the opposite access level. Only for workspace and project scope",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"project": schema.StringAttribute{
							MarkdownDescription: "The sandbox's project name",
							Required:            true,
						},
						"sandbox": schema.StringAttribute{
							MarkdownDescription: "The sandbox's name or identifier",
							Required:            true,
						},
						"access_level": schema.StringAttribute{
							MarkdownDescription: "The sandbox's access level",
							Validators: []validator.String{
								stringvalidator.OneOf(buddy.VariableAccessLevelUseOnly, buddy.VariableAccessLevelDenied),
							},
							Required: true,
						},
					},
				},
			},
		},
	}
}

func (r *variableResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *variableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *variableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	typ := buddy.VariableTypeVar
	ops := buddy.VariableOps{
		Key:   data.Key.ValueStringPointer(),
		Value: data.Value.ValueStringPointer(),
		Type:  &typ,
	}
	resp.Diagnostics.Append(data.applyCommonOps(ctx, &ops)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !data.ProjectName.IsNull() && !data.ProjectName.IsUnknown() {
		ops.Project = &buddy.VariableProject{
			Name: data.ProjectName.ValueString(),
		}
	}
	if !data.PipelineId.IsNull() && !data.PipelineId.IsUnknown() {
		ops.Pipeline = &buddy.VariablePipeline{
			Id: int(data.PipelineId.ValueInt64()),
		}
	}
	if !data.ActionId.IsNull() && !data.ActionId.IsUnknown() {
		ops.Action = &buddy.VariableAction{
			Id: int(data.ActionId.ValueInt64()),
		}
	}
	if !data.EnvironmentId.IsNull() && !data.EnvironmentId.IsUnknown() {
		ops.Environment = &buddy.VariableEnvironment{
			Id: data.EnvironmentId.ValueString(),
		}
	}
	if !data.SandboxId.IsNull() && !data.SandboxId.IsUnknown() {
		ops.Sandbox = &buddy.VariableSandbox{
			Id: data.SandboxId.ValueString(),
		}
	}
	variable, _, err := r.client.VariableService.Create(domain, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create variable", err))
		return
	}
	data.loadAPI(domain, variable)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *variableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *variableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, variableId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("variable", err))
		return
	}
	variable, httpResp, err := r.client.VariableService.Get(domain, variableId)
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get variable", err))
		return
	}
	if variable.Type != buddy.VariableTypeVar {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get variable", fmt.Errorf("variable not found")))
		return
	}
	data.loadAPI(domain, variable)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *variableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *variableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, variableId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("variable", err))
		return
	}
	ops := buddy.VariableOps{
		Value: data.Value.ValueStringPointer(),
	}
	resp.Diagnostics.Append(data.applyCommonOps(ctx, &ops)...)
	if resp.Diagnostics.HasError() {
		return
	}
	variable, _, err := r.client.VariableService.Update(domain, variableId, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update variable", err))
		return
	}
	data.loadAPI(domain, variable)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *variableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *variableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, variableId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("variable", err))
		return
	}
	_, err = r.client.VariableService.Delete(domain, variableId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete variable", err))
	}
}

func (r *variableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

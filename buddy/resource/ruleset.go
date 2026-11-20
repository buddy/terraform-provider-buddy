package resource

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ resource.Resource                = &rulesetResource{}
	_ resource.ResourceWithConfigure   = &rulesetResource{}
	_ resource.ResourceWithImportState = &rulesetResource{}
)

func NewRulesetResource() resource.Resource {
	return &rulesetResource{}
}

type rulesetResource struct {
	client *buddy.Client
}

type rulesetResourceModel struct {
	ID                        types.String `tfsdk:"id"`
	Domain                    types.String `tfsdk:"domain"`
	ProjectName               types.String `tfsdk:"project_name"`
	HtmlUrl                   types.String `tfsdk:"html_url"`
	RulesetId                 types.String `tfsdk:"ruleset_id"`
	Identifier                types.String `tfsdk:"identifier"`
	Name                      types.String `tfsdk:"name"`
	Disabled                  types.Bool   `tfsdk:"disabled"`
	Refs                      types.Set    `tfsdk:"refs"`
	Note                      types.String `tfsdk:"note"`
	AgentNote                 types.String `tfsdk:"agent_note"`
	BypassAdmins              types.Bool   `tfsdk:"bypass_admins"`
	BypassOthers              types.String `tfsdk:"bypass_others"`
	BypassUser                types.Set    `tfsdk:"bypass_user"`
	BypassGroup               types.Set    `tfsdk:"bypass_group"`
	RequiredPipelines         types.Set    `tfsdk:"required_pipelines"`
	RequiredUnitTestsSuites   types.Set    `tfsdk:"required_unit_tests_suites"`
	RequiredVisualTestsSuites types.Set    `tfsdk:"required_visual_tests_suites"`
	RequireUpToDate           types.Bool   `tfsdk:"require_up_to_date"`
	RequirePullRequest        types.Bool   `tfsdk:"require_pull_request"`
	BlockForcePush            types.Bool   `tfsdk:"block_force_push"`
	RestrictDeletion          types.Bool   `tfsdk:"restrict_deletion"`
	RestrictCreation          types.Bool   `tfsdk:"restrict_creation"`
	RestrictUpdate            types.Bool   `tfsdk:"restrict_update"`
	MissingSources            types.Set    `tfsdk:"missing_sources"`
}

func (r *rulesetResourceModel) decomposeId() (string, string, string, error) {
	return util.DecomposeTripleId(r.ID.ValueString())
}

func (r *rulesetResourceModel) loadAPI(ctx context.Context, domain string, projectName string, ruleset *buddy.Ruleset) diag.Diagnostics {
	var diags diag.Diagnostics
	var d diag.Diagnostics
	r.ID = types.StringValue(util.ComposeTripleId(domain, projectName, ruleset.Id))
	r.Domain = types.StringValue(domain)
	r.ProjectName = types.StringValue(projectName)
	r.HtmlUrl = types.StringValue(ruleset.HtmlUrl)
	r.RulesetId = types.StringValue(ruleset.Id)
	r.Identifier = types.StringValue(ruleset.Identifier)
	r.Name = util.StringFromApiOrNull(ruleset.Name)
	r.Disabled = types.BoolValue(ruleset.Disabled)
	r.Refs, d = types.SetValueFrom(ctx, types.StringType, &ruleset.Refs)
	diags.Append(d...)
	r.Note = util.StringFromApiOrNull(ruleset.Note)
	r.AgentNote = util.StringFromApiOrNull(ruleset.AgentNote)
	bypass := ruleset.Bypass
	if bypass == nil {
		bypass = &buddy.RulesetBypass{}
	}
	r.BypassAdmins = types.BoolPointerValue(bypass.Admins)
	r.BypassOthers = util.StringFromApiOrNull(bypass.Others)
	r.BypassUser, d = util.RulesetBypassEntriesFromApi(ctx, bypass.Users)
	diags.Append(d...)
	r.BypassGroup, d = util.RulesetBypassEntriesFromApi(ctx, bypass.Groups)
	diags.Append(d...)
	rules := ruleset.Rules
	if rules == nil {
		rules = &buddy.RulesetRules{}
	}
	checks := rules.RequiredChecks
	if checks == nil {
		checks = &buddy.RulesetRequiredChecks{}
	}
	r.RequiredPipelines, d = util.StringSetFromApiOrNull(ctx, checks.Pipelines)
	diags.Append(d...)
	r.RequiredUnitTestsSuites, d = util.StringSetFromApiOrNull(ctx, checks.UnitTestsSuites)
	diags.Append(d...)
	r.RequiredVisualTestsSuites, d = util.StringSetFromApiOrNull(ctx, checks.VisualTestsSuites)
	diags.Append(d...)
	r.RequireUpToDate = types.BoolValue(checks.RequireUpToDate)
	r.RequirePullRequest = types.BoolValue(rules.RequirePullRequest)
	r.BlockForcePush = types.BoolValue(rules.BlockForcePush)
	r.RestrictDeletion = types.BoolValue(rules.RestrictDeletion)
	r.RestrictCreation = types.BoolValue(rules.RestrictCreation)
	r.RestrictUpdate = types.BoolValue(rules.RestrictUpdate)
	r.MissingSources, d = util.RulesetSourcesFromApi(ctx, ruleset.MissingSources)
	diags.Append(d...)
	return diags
}

// toApi builds the whole ruleset: the API replaces bypass and rules as a whole, and a value left out of the
// configuration is sent empty so that it is cleared
func (r *rulesetResourceModel) toApi(ctx context.Context) (*buddy.RulesetOps, diag.Diagnostics) {
	var diags diag.Diagnostics
	var d diag.Diagnostics
	empty := ""
	ops := buddy.RulesetOps{
		Identifier: r.Identifier.ValueStringPointer(),
		Name:       &empty,
		Note:       &empty,
		AgentNote:  &empty,
		Disabled:   r.Disabled.ValueBoolPointer(),
		Bypass:     &buddy.RulesetBypass{},
		Rules: &buddy.RulesetRules{
			RequiredChecks:     &buddy.RulesetRequiredChecks{},
			RequirePullRequest: r.RequirePullRequest.ValueBool(),
			BlockForcePush:     r.BlockForcePush.ValueBool(),
			RestrictDeletion:   r.RestrictDeletion.ValueBool(),
			RestrictCreation:   r.RestrictCreation.ValueBool(),
			RestrictUpdate:     r.RestrictUpdate.ValueBool(),
		},
	}
	if !r.Name.IsNull() && !r.Name.IsUnknown() {
		ops.Name = r.Name.ValueStringPointer()
	}
	if !r.Note.IsNull() && !r.Note.IsUnknown() {
		ops.Note = r.Note.ValueStringPointer()
	}
	if !r.AgentNote.IsNull() && !r.AgentNote.IsUnknown() {
		ops.AgentNote = r.AgentNote.ValueStringPointer()
	}
	refs, d := util.StringSetToApiOrNil(ctx, &r.Refs)
	diags.Append(d...)
	ops.Refs = &refs
	if !r.BypassAdmins.IsNull() && !r.BypassAdmins.IsUnknown() {
		ops.Bypass.Admins = r.BypassAdmins.ValueBoolPointer()
	}
	if !r.BypassOthers.IsNull() && !r.BypassOthers.IsUnknown() {
		ops.Bypass.Others = r.BypassOthers.ValueString()
	}
	ops.Bypass.Users, d = util.RulesetBypassEntriesToApi(ctx, &r.BypassUser)
	diags.Append(d...)
	ops.Bypass.Groups, d = util.RulesetBypassEntriesToApi(ctx, &r.BypassGroup)
	diags.Append(d...)
	checks := ops.Rules.RequiredChecks
	checks.Pipelines, d = util.StringSetToApiOrNil(ctx, &r.RequiredPipelines)
	diags.Append(d...)
	checks.UnitTestsSuites, d = util.StringSetToApiOrNil(ctx, &r.RequiredUnitTestsSuites)
	diags.Append(d...)
	checks.VisualTestsSuites, d = util.StringSetToApiOrNil(ctx, &r.RequiredVisualTestsSuites)
	diags.Append(d...)
	checks.RequireUpToDate = r.RequireUpToDate.ValueBool()
	return &ops, diags
}

func (r *rulesetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ruleset"
}

func (r *rulesetResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *rulesetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func rulesFlag(description string) schema.BoolAttribute {
	return schema.BoolAttribute{
		MarkdownDescription: description + ". Default: `false`",
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
	}
}

func requiredChecksSet(description string) schema.SetAttribute {
	return schema.SetAttribute{
		MarkdownDescription: description + ". Each of them must mark commits with its status",
		Optional:            true,
		ElementType:         types.StringType,
		Validators: []validator.Set{
			setvalidator.SizeAtLeast(1),
		},
	}
}

func (r *rulesetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage a ruleset that protects branches and tags of a project. " +
			"Only workspace administrators and project managers can manage rulesets\n\n" +
			"Token scopes required: `RULESET_READ`, `RULESET_MANAGE`",
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
				MarkdownDescription: "The ruleset's URL",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ruleset_id": schema.StringAttribute{
				MarkdownDescription: "The ruleset's ID",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"identifier": schema.StringAttribute{
				MarkdownDescription: "The ruleset's identifier, unique in the project regardless of letter case",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The ruleset's name",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"disabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the ruleset is turned off. Default: `false`",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"refs": schema.SetAttribute{
				MarkdownDescription: "Patterns of the refs the ruleset protects, 1-20. Same syntax as refs in pipelines: `refs/heads/main`, `refs/tags/v*`, " +
					"a short name (matches both a branch and a tag), a wildcard, a regex in parentheses, or `:default` for the default branch",
				Required:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeBetween(1, 20),
				},
			},
			"note": schema.StringAttribute{
				MarkdownDescription: "The ruleset's note",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"agent_note": schema.StringAttribute{
				MarkdownDescription: "The ruleset's note for agents",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"bypass_admins": schema.BoolAttribute{
				MarkdownDescription: "`true` - workspace administrators and project managers bypass the ruleset, `false` - they do not, not set - no entry for them",
				Optional:            true,
			},
			"bypass_others": schema.StringAttribute{
				MarkdownDescription: "Whether the remaining members with access to the project bypass the ruleset (`ALLOWED`, `DENIED`), not set - no entry for them",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(buddy.RulesetBypassAccessAllowed, buddy.RulesetBypassAccessDenied),
				},
			},
			"required_pipelines":           requiredChecksSet("Identifiers of the pipelines that must succeed on the head commit of a pull request before it is merged"),
			"required_unit_tests_suites":   requiredChecksSet("Identifiers of the unit tests suites that must pass on the head commit of a pull request before it is merged"),
			"required_visual_tests_suites": requiredChecksSet("Identifiers of the visual tests suites that must be approved on the head commit of a pull request before it is merged"),
			"require_up_to_date":           rulesFlag("The pull request branch must also contain the current state of the target branch"),
			"require_pull_request":         rulesFlag("Changes reach a protected branch only by merging a pull request. Ignored for tags"),
			"block_force_push":             rulesFlag("Rewriting the history of a protected branch (non-fast-forward push) is not allowed. Ignored for tags"),
			"restrict_deletion":            rulesFlag("Only users who bypass the ruleset may delete a protected ref"),
			"restrict_creation":            rulesFlag("Only users who bypass the ruleset may create a protected ref"),
			"restrict_update":              rulesFlag("Only users who bypass the ruleset may move a protected ref, merging a pull request included"),
			"missing_sources": schema.SetNestedAttribute{
				MarkdownDescription: "Required checks whose source no longer exists or no longer marks commits with its status. Such a check blocks merging until the ruleset is fixed",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.ResourceRulesetSourceModelAttributes(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"bypass_user": schema.SetNestedBlock{
				MarkdownDescription: "A user (agents included) who may or may not bypass the ruleset. The user must have access to the project. " +
					"The user's own entry decides first, then `bypass_admins`, the user's groups, `bypass_others`",
				NestedObject: schema.NestedBlockObject{
					Attributes: util.ResourceRulesetBypassEntryModelAttributes("user"),
				},
			},
			"bypass_group": schema.SetNestedBlock{
				MarkdownDescription: "A group that may or may not bypass the ruleset. The group must be assigned to the project",
				NestedObject: schema.NestedBlockObject{
					Attributes: util.ResourceRulesetBypassEntryModelAttributes("group"),
				},
			},
		},
	}
}

func (r *rulesetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *rulesetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	projectName := data.ProjectName.ValueString()
	ops, d := data.toApi(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	ruleset, _, err := r.client.RulesetService.Create(domain, projectName, ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create ruleset", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, ruleset)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *rulesetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *rulesetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, projectName, rulesetId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("ruleset", err))
		return
	}
	ruleset, httpResp, err := r.client.RulesetService.Get(domain, projectName, rulesetId)
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get ruleset", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, ruleset)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *rulesetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *rulesetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, projectName, rulesetId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("ruleset", err))
		return
	}
	ops, d := data.toApi(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	ruleset, _, err := r.client.RulesetService.Update(domain, projectName, rulesetId, ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update ruleset", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, ruleset)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *rulesetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *rulesetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, projectName, rulesetId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("ruleset", err))
		return
	}
	_, err = r.client.RulesetService.Delete(domain, projectName, rulesetId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete ruleset", err))
	}
}

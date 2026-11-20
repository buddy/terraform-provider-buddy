package source

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ datasource.DataSource              = &rulesetSource{}
	_ datasource.DataSourceWithConfigure = &rulesetSource{}
)

func NewRulesetSource() datasource.DataSource {
	return &rulesetSource{}
}

type rulesetSource struct {
	client *buddy.Client
}

type rulesetSourceModel struct {
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

func (s *rulesetSourceModel) loadAPI(ctx context.Context, domain string, projectName string, ruleset *buddy.Ruleset) diag.Diagnostics {
	var diags diag.Diagnostics
	var d diag.Diagnostics
	s.ID = types.StringValue(util.ComposeTripleId(domain, projectName, ruleset.Id))
	s.Domain = types.StringValue(domain)
	s.ProjectName = types.StringValue(projectName)
	s.HtmlUrl = types.StringValue(ruleset.HtmlUrl)
	s.RulesetId = types.StringValue(ruleset.Id)
	s.Identifier = types.StringValue(ruleset.Identifier)
	s.Name = util.StringFromApiOrNull(ruleset.Name)
	s.Disabled = types.BoolValue(ruleset.Disabled)
	s.Refs, d = types.SetValueFrom(ctx, types.StringType, &ruleset.Refs)
	diags.Append(d...)
	s.Note = util.StringFromApiOrNull(ruleset.Note)
	s.AgentNote = util.StringFromApiOrNull(ruleset.AgentNote)
	bypass := ruleset.Bypass
	if bypass == nil {
		bypass = &buddy.RulesetBypass{}
	}
	s.BypassAdmins = types.BoolPointerValue(bypass.Admins)
	s.BypassOthers = util.StringFromApiOrNull(bypass.Others)
	s.BypassUser, d = util.RulesetBypassEntriesFromApi(ctx, bypass.Users)
	diags.Append(d...)
	s.BypassGroup, d = util.RulesetBypassEntriesFromApi(ctx, bypass.Groups)
	diags.Append(d...)
	rules := ruleset.Rules
	if rules == nil {
		rules = &buddy.RulesetRules{}
	}
	checks := rules.RequiredChecks
	if checks == nil {
		checks = &buddy.RulesetRequiredChecks{}
	}
	s.RequiredPipelines, d = util.StringSetFromApiOrNull(ctx, checks.Pipelines)
	diags.Append(d...)
	s.RequiredUnitTestsSuites, d = util.StringSetFromApiOrNull(ctx, checks.UnitTestsSuites)
	diags.Append(d...)
	s.RequiredVisualTestsSuites, d = util.StringSetFromApiOrNull(ctx, checks.VisualTestsSuites)
	diags.Append(d...)
	s.RequireUpToDate = types.BoolValue(checks.RequireUpToDate)
	s.RequirePullRequest = types.BoolValue(rules.RequirePullRequest)
	s.BlockForcePush = types.BoolValue(rules.BlockForcePush)
	s.RestrictDeletion = types.BoolValue(rules.RestrictDeletion)
	s.RestrictCreation = types.BoolValue(rules.RestrictCreation)
	s.RestrictUpdate = types.BoolValue(rules.RestrictUpdate)
	s.MissingSources, d = util.RulesetSourcesFromApi(ctx, ruleset.MissingSources)
	diags.Append(d...)
	return diags
}

func (s *rulesetSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ruleset"
}

func (s *rulesetSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *rulesetSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	computedStringSet := func(description string) schema.SetAttribute {
		return schema.SetAttribute{MarkdownDescription: description, Computed: true, ElementType: types.StringType}
	}
	computedBool := func(description string) schema.BoolAttribute {
		return schema.BoolAttribute{MarkdownDescription: description, Computed: true}
	}
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get ruleset by identifier or ruleset ID. Only workspace administrators and project managers can read rulesets\n\n" +
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
			"identifier": schema.StringAttribute{
				MarkdownDescription: "The ruleset's identifier (letter case does not matter)",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("identifier"),
						path.MatchRoot("ruleset_id"),
					}...),
				},
			},
			"ruleset_id": schema.StringAttribute{
				MarkdownDescription: "The ruleset's ID",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("identifier"),
						path.MatchRoot("ruleset_id"),
					}...),
				},
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The ruleset's URL",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The ruleset's name, null when not set",
				Computed:            true,
			},
			"disabled": computedBool("Whether the ruleset is turned off"),
			"refs":     computedStringSet("Patterns of the refs the ruleset protects"),
			"note": schema.StringAttribute{
				MarkdownDescription: "The ruleset's note, null when not set",
				Computed:            true,
			},
			"agent_note": schema.StringAttribute{
				MarkdownDescription: "The ruleset's note for agents, null when not set",
				Computed:            true,
			},
			"bypass_admins": computedBool("`true` - workspace administrators and project managers bypass the ruleset, `false` - they do not, null - no entry for them"),
			"bypass_others": schema.StringAttribute{
				MarkdownDescription: "Whether the remaining members with access to the project bypass the ruleset (`ALLOWED`, `DENIED`), null - no entry for them",
				Computed:            true,
			},
			"bypass_user": schema.SetNestedAttribute{
				MarkdownDescription: "Users (agents included) who may or may not bypass the ruleset",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.SourceRulesetBypassEntryModelAttributes(),
				},
			},
			"bypass_group": schema.SetNestedAttribute{
				MarkdownDescription: "Groups that may or may not bypass the ruleset",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.SourceRulesetBypassEntryModelAttributes(),
				},
			},
			"required_pipelines":           computedStringSet("Identifiers of the pipelines that must succeed, null when none"),
			"required_unit_tests_suites":   computedStringSet("Identifiers of the unit tests suites that must pass, null when none"),
			"required_visual_tests_suites": computedStringSet("Identifiers of the visual tests suites that must be approved, null when none"),
			"require_up_to_date":           computedBool("The pull request branch must also contain the current state of the target branch"),
			"require_pull_request":         computedBool("Changes reach a protected branch only by merging a pull request"),
			"block_force_push":             computedBool("Rewriting the history of a protected branch is not allowed"),
			"restrict_deletion":            computedBool("Only users who bypass the ruleset may delete a protected ref"),
			"restrict_creation":            computedBool("Only users who bypass the ruleset may create a protected ref"),
			"restrict_update":              computedBool("Only users who bypass the ruleset may move a protected ref"),
			"missing_sources": schema.SetNestedAttribute{
				MarkdownDescription: "Required checks whose source no longer exists or no longer marks commits with its status",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.SourceRulesetSourceModelAttributes(),
				},
			},
		},
	}
}

func (s *rulesetSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *rulesetSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	projectName := data.ProjectName.ValueString()
	var ruleset *buddy.Ruleset
	if !data.RulesetId.IsNull() && !data.RulesetId.IsUnknown() {
		r, httpResp, err := s.client.RulesetService.Get(domain, projectName, data.RulesetId.ValueString())
		if err != nil {
			if util.IsResourceNotFound(httpResp, err) {
				resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("ruleset"))
				return
			}
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get ruleset", err))
			return
		}
		ruleset = r
	} else {
		rulesets, _, err := s.client.RulesetService.GetList(domain, projectName)
		if err != nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get rulesets", err))
			return
		}
		identifier := data.Identifier.ValueString()
		for _, r := range rulesets.Rulesets {
			if strings.EqualFold(r.Identifier, identifier) {
				ruleset = r
				break
			}
		}
		if ruleset == nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("ruleset"))
			return
		}
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, ruleset)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

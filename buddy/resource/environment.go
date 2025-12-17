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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ resource.Resource                = &environmentResource{}
	_ resource.ResourceWithConfigure   = &environmentResource{}
	_ resource.ResourceWithImportState = &environmentResource{}
)

func NewEnvironmentResource() resource.Resource {
	return &environmentResource{}
}

type environmentResource struct {
	client *buddy.Client
}

func (e *environmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (e *environmentResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	e.client = req.ProviderData.(*buddy.Client)
}

type environmentResourceModel struct {
	ID                     types.String `tfsdk:"id"`
	Domain                 types.String `tfsdk:"domain"`
	ProjectName            types.String `tfsdk:"project_name"`
	HtmlUrl                types.String `tfsdk:"html_url"`
	EnvironmentId          types.String `tfsdk:"environment_id"`
	Identifier             types.String `tfsdk:"identifier"`
	Name                   types.String `tfsdk:"name"`
	Icon                   types.String `tfsdk:"icon"`
	PublicUrl              types.String `tfsdk:"public_url"`
	AllPipelinesAllowed    types.Bool   `tfsdk:"all_pipelines_allowed"`
	AllowedPipeline        types.Set    `tfsdk:"allowed_pipeline"`
	AllEnvironmentsAllowed types.Bool   `tfsdk:"all_environments_allowed"`
	AllowedEnvironment     types.Set    `tfsdk:"allowed_environment"`
	CreateDate             types.String `tfsdk:"create_date"`
	Scope                  types.String `tfsdk:"scope"`
	BaseOnly               types.Bool   `tfsdk:"base_only"`
	BaseEnvironments       types.Set    `tfsdk:"base_environments"`
	Project                types.Set    `tfsdk:"project"`
	Tags                   types.Set    `tfsdk:"tags"`
	Permissions            types.Set    `tfsdk:"permissions"`
}

func (r *environmentResourceModel) loadAPI(ctx context.Context, domain string, environment *buddy.Environment) diag.Diagnostics {
	var diags diag.Diagnostics
	// zeby ograniczyc breaking change zostawmy id jako 3 elementowe z pustym projektem (nie istotnym teraz)
	r.ID = types.StringValue(util.ComposeTripleId(domain, "", environment.Id))
	r.Domain = types.StringValue(domain)
	if environment.Project != nil {
		r.ProjectName = types.StringValue(environment.Project.Name)
	} else {
		r.ProjectName = types.StringNull()
	}
	r.HtmlUrl = types.StringValue(environment.HtmlUrl)
	r.EnvironmentId = types.StringValue(environment.Id)
	r.Identifier = types.StringValue(environment.Identifier)
	r.Name = types.StringValue(environment.Name)
	r.Icon = types.StringValue(environment.Icon)
	r.PublicUrl = types.StringValue(environment.PublicUrl)
	r.AllPipelinesAllowed = types.BoolValue(environment.AllPipelinesAllowed)
	r.AllEnvironmentsAllowed = types.BoolValue(environment.AllEnvironmentsAllowed)
	r.CreateDate = types.StringValue(environment.CreateDate)
	r.Scope = types.StringValue(environment.Scope)
	r.BaseOnly = types.BoolValue(environment.BaseOnly)
	tags, d := types.SetValueFrom(ctx, types.StringType, &environment.Tags)
	diags.Append(d...)
	var envProjects []*buddy.Project
	if environment.Project != nil {
		envProjects = []*buddy.Project{environment.Project}
	} else {
		envProjects = []*buddy.Project{}
	}
	r.Tags = tags
	projects, d := util.ProjectsModelFromApi(ctx, &envProjects)
	diags.Append(d...)
	r.Project = projects
	base, d := types.SetValueFrom(ctx, types.StringType, &environment.BaseEnvironments)
	diags.Append(d...)
	r.BaseEnvironments = base
	return diags
}

func (r *environmentResourceModel) decomposeId() (string, string, string, error) {
	domain, projectName, environmentId, err := util.DecomposeTripleId(r.ID.ValueString())
	if err != nil {
		return "", "", "", err
	}
	return domain, projectName, environmentId, nil
}

func (e *environmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

func (e *environmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage an environment\n\n" +
			"Token scopes required: `WORKSPACE`, `ENVIRONMENT_MANAGE`, `ENVIRONMENT_INFO`",
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
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The environment's URL",
				Computed:            true,
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "The environment's ID",
				Computed:            true,
			},
			"identifier": schema.StringAttribute{
				MarkdownDescription: "The environment's identifier",
				Required:            true,
				Validators:          util.StringValidatorIdentifier(),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The environment's name",
				Required:            true,
			},
			"icon": schema.StringAttribute{
				MarkdownDescription: "The environment's icon",
				Optional:            true,
				Computed:            true,
			},
			"public_url": schema.StringAttribute{
				MarkdownDescription: "The environment's public URL",
				Optional:            true,
				Computed:            true,
			},
			"all_pipelines_allowed": schema.BoolAttribute{
				MarkdownDescription: "Defines whether or not environment can be used in all pipelines",
				Optional:            true,
				Computed:            true,
			},
			"all_environments_allowed": schema.BoolAttribute{
				MarkdownDescription: "Defines whether or not environment can be by inherited by other environements",
				Optional:            true,
				Computed:            true,
			},
			"create_date": schema.StringAttribute{
				MarkdownDescription: "The environment's create date",
				Computed:            true,
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: "The environment's scope",
				Computed:            true,
			},
			"base_only": schema.BoolAttribute{
				MarkdownDescription: "Defines whether or not environment can be only used as base environment",
				Optional:            true,
				Computed:            true,
			},
			"base_environments": schema.SetAttribute{
				MarkdownDescription: "The environment's list of parent environments ID to inherit from",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"project": schema.SetNestedAttribute{
				MarkdownDescription: "The environment's project",
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
			"tags": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The environment's list of tags",
				Optional:            true,
				Computed:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"allowed_pipeline": schema.SetNestedBlock{
				MarkdownDescription: "The environment's allowed pipeline",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"project": schema.StringAttribute{
							MarkdownDescription: "The project's name",
							Required:            true,
						},
						"pipeline": schema.StringAttribute{
							MarkdownDescription: "The pipeline's identifier",
							Required:            true,
						},
					},
				},
			},
			"allowed_environment": schema.SetNestedBlock{
				MarkdownDescription: "The environment's allowed child environment",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"project": schema.StringAttribute{
							MarkdownDescription: "The project's name",
							Optional:            true,
						},
						"environment": schema.StringAttribute{
							MarkdownDescription: "The environment's identifier",
							Required:            true,
						},
					},
				},
			},
			"permissions": schema.SetNestedBlock{
				MarkdownDescription: "The environment's permissions",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"others": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									buddy.EnvironmentPermissionAccessLevelUseOnly,
									buddy.EnvironmentPermissionAccessLevelDefault,
									buddy.EnvironmentPermissionAccessLevelDenied,
									buddy.EnvironmentPermissionAccessLevelManage,
								),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"user": schema.SetNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: util.EnvironmentPermissionsAccessModelAttributes(),
							},
						},
						"group": schema.SetNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: util.EnvironmentPermissionsAccessModelAttributes(),
							},
						},
					},
				},
				Validators: []validator.Set{
					setvalidator.SizeAtMost(1),
				},
			},
		},
	}
}

func (e *environmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *environmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	ops := buddy.EnvironmentOps{
		Name:       data.Name.ValueStringPointer(),
		Identifier: data.Identifier.ValueStringPointer(),
	}
	var scope string
	if !data.ProjectName.IsNull() && !data.ProjectName.IsUnknown() {
		ops.Project = &buddy.ProjectSimple{
			Name: data.ProjectName.ValueString(),
		}
		scope = buddy.EnvironmentScopeProject
	} else {
		scope = buddy.EnvironmentScopeWorkspace
	}
	ops.Scope = &scope
	if !data.Icon.IsNull() && !data.Icon.IsUnknown() {
		ops.Icon = data.Icon.ValueStringPointer()
	}
	if !data.PublicUrl.IsNull() && !data.PublicUrl.IsUnknown() {
		ops.PublicUrl = data.PublicUrl.ValueStringPointer()
	}
	if !data.BaseOnly.IsNull() && !data.BaseOnly.IsUnknown() {
		ops.BaseOnly = data.BaseOnly.ValueBoolPointer()
	}
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		tags, d := util.StringSetToApi(ctx, &data.Tags)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Tags = tags
	}
	if !data.Permissions.IsNull() && !data.Permissions.IsUnknown() {
		permissions, d := util.EnvironmentPermissionsModelToApi(ctx, &data.Permissions)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Permissions = permissions
	}
	if !data.AllowedPipeline.IsNull() && !data.AllowedPipeline.IsUnknown() {
		pipelines, d := util.EnvironmentPipelinesModelToApi(ctx, &data.AllowedPipeline)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.AllowedPipelines = pipelines
		f := false
		ops.AllPipelinesAllowed = &f
	}
	if !data.AllowedEnvironment.IsNull() && !data.AllowedEnvironment.IsUnknown() {
		environments, d := util.EnvironmentEnvironmentsModelToApi(ctx, &data.AllowedEnvironment)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.AllowedEnvironments = environments
		f := false
		ops.AllEnvironmentsAllowed = &f
	}
	if !data.AllPipelinesAllowed.IsNull() && !data.AllPipelinesAllowed.IsUnknown() {
		ops.AllPipelinesAllowed = data.AllPipelinesAllowed.ValueBoolPointer()
	}
	if !data.AllEnvironmentsAllowed.IsNull() && !data.AllEnvironmentsAllowed.IsUnknown() {
		ops.AllEnvironmentsAllowed = data.AllEnvironmentsAllowed.ValueBoolPointer()
	}
	if !data.BaseEnvironments.IsNull() && !data.BaseEnvironments.IsUnknown() {
		base, d := util.StringSetToApi(ctx, &data.BaseEnvironments)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.BaseEnvironments = base
	}
	environment, _, err := e.client.EnvironmentService.Create(domain, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create environment", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, environment)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (e *environmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *environmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, _, environmentId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("environment", err))
		return
	}
	environment, httpResp, err := e.client.EnvironmentService.Get(domain, environmentId)
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get environment", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, environment)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (e *environmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *environmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, _, environmentId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("environment", err))
		return
	}
	ops := buddy.EnvironmentOps{
		Name:       data.Name.ValueStringPointer(),
		Identifier: data.Identifier.ValueStringPointer(),
	}
	var scope string
	if !data.ProjectName.IsNull() && !data.ProjectName.IsUnknown() {
		ops.Project = &buddy.ProjectSimple{
			Name: data.ProjectName.ValueString(),
		}
		scope = buddy.EnvironmentScopeProject
	} else {
		scope = buddy.EnvironmentScopeWorkspace
	}
	ops.Scope = &scope
	if !data.Icon.IsNull() && !data.Icon.IsUnknown() {
		ops.Icon = data.Icon.ValueStringPointer()
	}
	if !data.PublicUrl.IsNull() && !data.PublicUrl.IsUnknown() {
		ops.PublicUrl = data.PublicUrl.ValueStringPointer()
	}
	if !data.BaseOnly.IsNull() && !data.BaseOnly.IsUnknown() {
		ops.BaseOnly = data.BaseOnly.ValueBoolPointer()
	}
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		tags, d := util.StringSetToApi(ctx, &data.Tags)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Tags = tags
	}
	if !data.Permissions.IsNull() && !data.Permissions.IsUnknown() {
		permissions, d := util.EnvironmentPermissionsModelToApi(ctx, &data.Permissions)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Permissions = permissions
	}
	if !data.AllowedPipeline.IsNull() && !data.AllowedPipeline.IsUnknown() {
		pipelines, d := util.EnvironmentPipelinesModelToApi(ctx, &data.AllowedPipeline)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.AllowedPipelines = pipelines
		f := false
		ops.AllPipelinesAllowed = &f
	}
	if !data.AllowedEnvironment.IsNull() && !data.AllowedEnvironment.IsUnknown() {
		environments, d := util.EnvironmentEnvironmentsModelToApi(ctx, &data.AllowedEnvironment)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.AllowedEnvironments = environments
		f := false
		ops.AllEnvironmentsAllowed = &f
	}
	if !data.AllPipelinesAllowed.IsNull() && !data.AllPipelinesAllowed.IsUnknown() {
		ops.AllPipelinesAllowed = data.AllPipelinesAllowed.ValueBoolPointer()
	}
	if !data.AllEnvironmentsAllowed.IsNull() && !data.AllEnvironmentsAllowed.IsUnknown() {
		ops.AllEnvironmentsAllowed = data.AllEnvironmentsAllowed.ValueBoolPointer()
	}
	if !data.BaseEnvironments.IsNull() && !data.BaseEnvironments.IsUnknown() {
		base, d := util.StringSetToApi(ctx, &data.BaseEnvironments)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.BaseEnvironments = base
	}
	environment, _, err := e.client.EnvironmentService.Update(domain, environmentId, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update environment", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, environment)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (e *environmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *environmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, _, environmentId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("environment", err))
		return
	}
	_, err = e.client.EnvironmentService.Delete(domain, environmentId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete environment", err))
	}
}

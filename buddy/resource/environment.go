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
	ID                  types.String `tfsdk:"id"`
	Domain              types.String `tfsdk:"domain"`
	ProjectName         types.String `tfsdk:"project_name"`
	HtmlUrl             types.String `tfsdk:"html_url"`
	EnvironmentId       types.String `tfsdk:"environment_id"`
	Name                types.String `tfsdk:"name"`
	Identifier          types.String `tfsdk:"identifier"`
	Type                types.String `tfsdk:"type"`
	Tags                types.Set    `tfsdk:"tags"`
	PublicUrl           types.String `tfsdk:"public_url"`
	AllPipelinesAllowed types.Bool   `tfsdk:"all_pipelines_allowed"`
	AllowedPipelines    types.Set    `tfsdk:"allowed_pipelines"`
	Project             types.Set    `tfsdk:"project"`
	Variable            types.Set    `tfsdk:"var"`
	Permissions         types.Set    `tfsdk:"permissions"`
}

func (r *environmentResourceModel) loadAPI(ctx context.Context, domain string, projectName string, environment *buddy.Environment) diag.Diagnostics {
	var diags diag.Diagnostics
	r.ID = types.StringValue(util.ComposeTripleId(domain, projectName, environment.Id))
	r.Domain = types.StringValue(domain)
	r.ProjectName = types.StringValue(projectName)
	r.HtmlUrl = types.StringValue(environment.HtmlUrl)
	r.EnvironmentId = types.StringValue(environment.Id)
	r.Name = types.StringValue(environment.Name)
	r.Identifier = types.StringValue(environment.Identifier)
	r.Type = types.StringValue(environment.Type)
	tags, d := types.SetValueFrom(ctx, types.StringType, &environment.Tags)
	diags.Append(d...)
	r.Tags = tags
	r.PublicUrl = types.StringValue(environment.PublicUrl)
	r.AllPipelinesAllowed = types.BoolValue(environment.AllPipelinesAllowed)
	ids := make([]int64, len(environment.AllowedPipelines))
	for i, v := range environment.AllowedPipelines {
		ids[i] = int64(v.Id)
	}
	allowedPipelines, d := types.SetValueFrom(ctx, types.Int64Type, &ids)
	diags.Append(d...)
	r.AllowedPipelines = allowedPipelines
	projects, d := util.ProjectsModelFromApi(ctx, &[]*buddy.Project{environment.Project})
	diags.Append(d...)
	r.Project = projects
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
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The pipeline's URL",
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
			"type": schema.StringAttribute{
				MarkdownDescription: "The environment's typpe. Allowed: `PRODUCTION`, `STAGE`, `DEV`",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{stringvalidator.OneOf(
					buddy.EnvironmentTypeProduction,
					buddy.EnvironmentTypeStage,
					buddy.EnvironmentTypeDev,
				)},
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
			"allowed_pipelines": schema.SetAttribute{
				ElementType:         types.Int64Type,
				MarkdownDescription: "List of pipeline IDs that is allowed to use the environment",
				Optional:            true,
				Computed:            true,
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
			"var": schema.SetNestedBlock{
				MarkdownDescription: "The environment's variables",
				NestedObject: schema.NestedBlockObject{
					Attributes: util.EnvironmentVariableModelAttributes(),
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
	projectName := data.ProjectName.ValueString()
	ops := buddy.EnvironmentOps{
		Name: data.Name.ValueStringPointer(),
		Type: data.Type.ValueStringPointer(),
	}
	if !data.Identifier.IsNull() && !data.Identifier.IsUnknown() {
		ops.Identifier = data.Identifier.ValueStringPointer()
	}
	if !data.PublicUrl.IsNull() && !data.PublicUrl.IsUnknown() {
		ops.PublicUrl = data.PublicUrl.ValueStringPointer()
	}
	if !data.AllPipelinesAllowed.IsNull() && !data.AllPipelinesAllowed.IsUnknown() {
		ops.AllPipelinesAllowed = data.AllPipelinesAllowed.ValueBoolPointer()
	}
	if !data.AllowedPipelines.IsNull() && !data.AllowedPipelines.IsUnknown() {
		ids, d := util.Int64SetToApi(ctx, &data.AllowedPipelines)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		allowedPipelines := make([]*buddy.AllowedPipeline, len(*ids))
		for i, v := range *ids {
			allowedPipelines[i] = &buddy.AllowedPipeline{
				Id: v,
			}
		}
		ops.AllowedPipelines = &allowedPipelines
	}
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		tags, d := util.StringSetToApi(ctx, &data.Tags)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Tags = tags
	}
	if !data.Variable.IsNull() && !data.Variable.IsUnknown() {
		variables, d := util.EnvironmentVariableModelToApi(ctx, &data.Variable)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Variables = variables
	}
	if !data.Permissions.IsNull() && !data.Permissions.IsUnknown() {
		permissions, d := util.EnvironmentPermissionsModelToApi(ctx, &data.Permissions)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Permissions = permissions
	}
	environment, _, err := e.client.EnvironmentService.Create(domain, projectName, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create environment", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, environment)...)
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
	domain, projectName, environmentId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("environment", err))
		return
	}
	environment, httpResp, err := e.client.EnvironmentService.Get(domain, projectName, environmentId)
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get environment", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, environment)...)
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
	domain, projectName, environmentId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("environment", err))
		return
	}
	ops := buddy.EnvironmentOps{
		Name: data.Name.ValueStringPointer(),
		Type: data.Type.ValueStringPointer(),
	}
	if !data.Identifier.IsNull() && !data.Identifier.IsUnknown() {
		ops.Identifier = data.Identifier.ValueStringPointer()
	}
	if !data.PublicUrl.IsNull() && !data.PublicUrl.IsUnknown() {
		ops.PublicUrl = data.PublicUrl.ValueStringPointer()
	}
	if !data.AllPipelinesAllowed.IsNull() && !data.AllPipelinesAllowed.IsUnknown() {
		ops.AllPipelinesAllowed = data.AllPipelinesAllowed.ValueBoolPointer()
	}
	if !data.AllowedPipelines.IsNull() && !data.AllowedPipelines.IsUnknown() {
		ids, d := util.Int64SetToApi(ctx, &data.AllowedPipelines)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		allowedPipelines := make([]*buddy.AllowedPipeline, len(*ids))
		for i, v := range *ids {
			allowedPipelines[i] = &buddy.AllowedPipeline{
				Id: v,
			}
		}
		ops.AllowedPipelines = &allowedPipelines
	}
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		tags, d := util.StringSetToApi(ctx, &data.Tags)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Tags = tags
	}
	if !data.Variable.IsNull() && !data.Variable.IsUnknown() {
		variables, d := util.EnvironmentVariableModelToApi(ctx, &data.Variable)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Variables = variables
	}
	if !data.Permissions.IsNull() && !data.Permissions.IsUnknown() {
		permissions, d := util.EnvironmentPermissionsModelToApi(ctx, &data.Permissions)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Permissions = permissions
	}
	environment, _, err := e.client.EnvironmentService.Update(domain, projectName, environmentId, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update environment", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, environment)...)
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
	domain, projectName, environmentId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("environment", err))
		return
	}
	_, err = e.client.EnvironmentService.Delete(domain, projectName, environmentId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete environment", err))
	}
}

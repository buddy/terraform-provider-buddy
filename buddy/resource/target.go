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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ resource.Resource                = &targetResource{}
	_ resource.ResourceWithConfigure   = &targetResource{}
	_ resource.ResourceWithImportState = &targetResource{}
)

func NewTargetResource() resource.Resource {
	return &targetResource{}
}

type targetResource struct {
	client *buddy.Client
}

type targetResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Domain        types.String `tfsdk:"domain"`
	TargetId      types.String `tfsdk:"target_id"`
	HtmlUrl       types.String `tfsdk:"html_url"`
	Name          types.String `tfsdk:"name"`
	Identifier    types.String `tfsdk:"identifier"`
	Tags          types.Set    `tfsdk:"tags"`
	Type          types.String `tfsdk:"type"`
	Host          types.String `tfsdk:"host"`
	Scope         types.String `tfsdk:"scope"`
	Repository    types.String `tfsdk:"repository"`
	Port          types.String `tfsdk:"port"`
	Path          types.String `tfsdk:"path"`
	Secure        types.Bool   `tfsdk:"secure"`
	Integration   types.String `tfsdk:"integration"`
	Disabled      types.Bool   `tfsdk:"disabled"`
	Auth          types.Set    `tfsdk:"auth"`
	ProjectName   types.String `tfsdk:"project_name"`
	PipelineId    types.Int64  `tfsdk:"pipeline_id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	Proxy         types.Set    `tfsdk:"proxy"`
	Permissions   types.Set    `tfsdk:"permissions"`
}

func (m *targetResourceModel) decomposeId() (string, string, error) {
	domain, environmentId, err := util.DecomposeDoubleId(m.ID.ValueString())
	if err != nil {
		return "", "", err
	}
	return domain, environmentId, nil
}

func (m *targetResourceModel) loadAPI(ctx context.Context, domain string, target *buddy.Target) diag.Diagnostics {
	var diags diag.Diagnostics
	m.ID = types.StringValue(util.ComposeDoubleId(domain, target.Id))
	m.Domain = types.StringValue(domain)
	m.HtmlUrl = types.StringValue(target.HtmlUrl)
	m.TargetId = types.StringValue(target.Id)
	m.Identifier = types.StringValue(target.Identifier)
	tags, d := types.SetValueFrom(ctx, types.StringType, &target.Tags)
	diags.Append(d...)
	m.Tags = tags
	m.Name = types.StringValue(target.Name)
	m.Type = types.StringValue(target.Type)
	m.Host = types.StringValue(target.Host)
	m.Scope = types.StringValue(target.Scope)
	m.Repository = types.StringValue(target.Repository)
	m.Port = types.StringValue(target.Port)
	m.Path = types.StringValue(target.Path)
	m.Secure = types.BoolValue(target.Secure)
	m.Integration = types.StringValue(target.Integration)
	m.Disabled = types.BoolValue(target.Disabled)
	return diags
}

func (m *targetResourceModel) toOps(ctx context.Context) (*buddy.TargetOps, diag.Diagnostics) {
	var diags diag.Diagnostics
	ops := buddy.TargetOps{}
	if !m.Identifier.IsNull() && !m.Identifier.IsUnknown() {
		ops.Identifier = m.Identifier.ValueStringPointer()
	}
	if !m.Name.IsNull() && !m.Name.IsUnknown() {
		ops.Name = m.Name.ValueStringPointer()
	}
	if !m.Tags.IsNull() && !m.Tags.IsUnknown() {
		tags, d := util.StringSetToApi(ctx, &m.Tags)
		diags.Append(d...)
		ops.Tags = tags
	}
	if !m.Type.IsNull() && !m.Type.IsUnknown() {
		ops.Type = m.Type.ValueStringPointer()
	}
	if !m.Host.IsNull() && !m.Host.IsUnknown() {
		ops.Host = m.Host.ValueStringPointer()
	}
	if !m.Scope.IsNull() && !m.Scope.IsUnknown() {
		ops.Scope = m.Scope.ValueStringPointer()
	}
	if !m.Repository.IsNull() && !m.Repository.IsUnknown() {
		ops.Repository = m.Repository.ValueStringPointer()
	}
	if !m.Port.IsNull() && !m.Port.IsUnknown() {
		ops.Port = m.Port.ValueStringPointer()
	}
	if !m.Path.IsNull() && !m.Path.IsUnknown() {
		ops.Path = m.Path.ValueStringPointer()
	}
	if !m.Secure.IsNull() && !m.Secure.IsUnknown() {
		ops.Secure = m.Secure.ValueBoolPointer()
	}
	if !m.Integration.IsNull() && !m.Integration.IsUnknown() {
		ops.Integration = m.Integration.ValueStringPointer()
	}
	if !m.Disabled.IsNull() && !m.Disabled.IsUnknown() {
		ops.Disabled = m.Disabled.ValueBoolPointer()
	}
	if !m.ProjectName.IsNull() && !m.ProjectName.IsUnknown() {
		ops.Project = &buddy.TargetProject{
			Name: m.ProjectName.ValueString(),
		}
	}
	if !m.PipelineId.IsNull() && !m.PipelineId.IsUnknown() {
		ops.Pipeline = &buddy.TargetPipeline{
			Id: int(m.PipelineId.ValueInt64()),
		}
	}
	if !m.EnvironmentId.IsNull() && !m.EnvironmentId.IsUnknown() {
		ops.Environment = &buddy.TargetEnvironment{
			Id: m.EnvironmentId.ValueString(),
		}
	}
	if !m.Permissions.IsNull() && !m.Permissions.IsUnknown() {
		permissions, d := util.TargetPermissionsModelToApi(ctx, &m.Permissions)
		diags.Append(d...)
		ops.Permissions = permissions
	}
	if !m.Auth.IsNull() && !m.Auth.IsUnknown() {
		auth, d := util.TargetAuthModelToApi(ctx, &m.Auth)
		diags.Append(d...)
		ops.Auth = auth
	}
	if !m.Proxy.IsNull() && !m.Proxy.IsUnknown() {
		proxy, d := util.TargetProxyModelToApi(ctx, &m.Proxy)
		diags.Append(d...)
		ops.Proxy = proxy
	}
	return &ops, diags
}

func (r *targetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_target"
}

func (r *targetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage a target\n\n" +
			"Token scope required: `WORKSPACE`, `TARGET_MANAGE`, `TARGET_INFO`",
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"pipeline_id": schema.Int64Attribute{
				MarkdownDescription: "The pipeline's id",
				Optional:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "The environment's id",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The target's URL",
				Computed:            true,
			},
			"target_id": schema.StringAttribute{
				MarkdownDescription: "The targets's ID",
				Computed:            true,
			},
			"identifier": schema.StringAttribute{
				MarkdownDescription: "The target's identifier",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The target's name",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The target's type. Allowed: `FTP`, `SSH`, `MATCH`, `UPCLOUD`, `VULTR`, `DIGITAL_OCEAN`, `GIT`",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						buddy.TargetTypeFtp,
						buddy.TargetTypeSsh,
						buddy.TargetTypeMatch,
						buddy.TargetTypeUpcloud,
						buddy.TargetTypeVultr,
						buddy.TargetTypeDigitalOcean,
						buddy.TargetTypeGit,
					),
				},
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "The target's list of tags",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: "The target's scope. Set for `MATCH`",
				Optional:            true,
				Computed:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "The target's host. Set for `FTP`, `SSH`, `UPCLOUD`, `VULTR`, `DIGITAL_OCEAN`",
				Optional:            true,
				Computed:            true,
			},
			"repository": schema.StringAttribute{
				MarkdownDescription: "The target's repository. Set for `GIT`",
				Optional:            true,
				Computed:            true,
			},
			"port": schema.StringAttribute{
				MarkdownDescription: "The target's port. Set for `FTP`, `SSH`, `UPCLOUD`, `VULTR`, `DIGITAL_OCEAN`",
				Optional:            true,
				Computed:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "The target's path",
				Optional:            true,
				Computed:            true,
			},
			"secure": schema.BoolAttribute{
				MarkdownDescription: "The target's secure setting. Set for `FTP`",
				Optional:            true,
				Computed:            true,
			},
			"integration": schema.StringAttribute{
				MarkdownDescription: "The target's integration. Set for `UPCLOUD`, `VULTR`, `DIGITAL_OCEAN`",
				Optional:            true,
				Computed:            true,
			},
			"disabled": schema.BoolAttribute{
				MarkdownDescription: "Defines whether or not the target can be run",
				Optional:            true,
				Computed:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"permissions": schema.SetNestedBlock{
				MarkdownDescription: "The target's permissions",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"others": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									buddy.TargetPermissionManage,
									buddy.TargetPermissionUseOnly,
								),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"user": schema.SetNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: util.TargetPermissionsAccessModelAttributes(),
							},
						},
						"group": schema.SetNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: util.TargetPermissionsAccessModelAttributes(),
							},
						},
					},
				},
				Validators: []validator.Set{
					setvalidator.SizeAtMost(1),
				},
			},
			"auth": schema.SetNestedBlock{
				MarkdownDescription: "The target's auth. Set for `FTP`, `GIT`, `SSH`, `UPCLOUD`, `VULTR`, `DIGITAL_OCEAN`",
				NestedObject: schema.NestedBlockObject{
					Attributes: util.TargetAuthModelAttributes(),
				},
				Validators: []validator.Set{
					setvalidator.SizeAtMost(1),
				},
			},
			"proxy": schema.SetNestedBlock{
				MarkdownDescription: "The target's proxy. Set for `SSH`",
				NestedObject: schema.NestedBlockObject{
					Attributes: util.TargetProxyModelAttributes(),
					Blocks: map[string]schema.Block{
						"auth": schema.SetNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: util.TargetAuthModelAttributes(),
							},
							Validators: []validator.Set{
								setvalidator.SizeAtMost(1),
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

func (r *targetResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *targetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *targetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	ops, d := data.toOps(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	target, _, err := r.client.TargetService.Create(domain, ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create target", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, target)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *targetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *targetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, targetId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("target", err))
		return
	}
	target, httpResp, err := r.client.TargetService.Get(domain, targetId)
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get target", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, target)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *targetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *targetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, targetId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("target", err))
		return
	}
	ops, d := data.toOps(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	target, _, err := r.client.TargetService.Update(domain, targetId, ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update target", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, target)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *targetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *targetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, targetId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("target", err))
		return
	}
	_, err = r.client.TargetService.Delete(domain, targetId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete target", err))
		return
	}
}

func (r *targetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

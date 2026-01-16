package resource

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ resource.Resource                = &permissionResource{}
	_ resource.ResourceWithConfigure   = &permissionResource{}
	_ resource.ResourceWithImportState = &permissionResource{}
)

func NewPermissionResource() resource.Resource {
	return &permissionResource{}
}

type permissionResource struct {
	client *buddy.Client
}

type permissionResourceModel struct {
	ID                     types.String `tfsdk:"id"`
	Domain                 types.String `tfsdk:"domain"`
	Name                   types.String `tfsdk:"name"`
	PipelineAccessLevel    types.String `tfsdk:"pipeline_access_level"`
	RepositoryAccessLevel  types.String `tfsdk:"repository_access_level"`
	SandboxAccessLevel     types.String `tfsdk:"sandbox_access_level"`
	ProjectTeamAccessLevel types.String `tfsdk:"project_team_access_level"`
	TargetAccessLevel      types.String `tfsdk:"target_access_level"`
	EnvironmentAccessLevel types.String `tfsdk:"environment_access_level"`
	PermissionId           types.Int64  `tfsdk:"permission_id"`
	Description            types.String `tfsdk:"description"`
	HtmlUrl                types.String `tfsdk:"html_url"`
	Type                   types.String `tfsdk:"type"`
}

func (r *permissionResourceModel) decomposeId() (string, int, error) {
	domain, gid, err := util.DecomposeDoubleId(r.ID.ValueString())
	if err != nil {
		return "", 0, err
	}
	permId, err := strconv.Atoi(gid)
	if err != nil {
		return "", 0, err
	}
	return domain, permId, nil
}

func (r *permissionResourceModel) loadAPI(domain string, permission *buddy.Permission) {
	r.ID = types.StringValue(util.ComposeDoubleId(domain, strconv.Itoa(permission.Id)))
	r.Domain = types.StringValue(domain)
	r.Name = types.StringValue(permission.Name)
	r.PipelineAccessLevel = types.StringValue(permission.PipelineAccessLevel)
	r.RepositoryAccessLevel = types.StringValue(permission.RepositoryAccessLevel)
	r.SandboxAccessLevel = types.StringValue(permission.SandboxAccessLevel)
	r.ProjectTeamAccessLevel = types.StringValue(permission.ProjectTeamAccessLevel)
	r.TargetAccessLevel = types.StringValue(permission.TargetAccessLevel)
	r.EnvironmentAccessLevel = types.StringValue(permission.EnvironmentAccessLevel)
	r.PermissionId = types.Int64Value(int64(permission.Id))
	r.HtmlUrl = types.StringValue(permission.HtmlUrl)
	r.Type = types.StringValue(permission.Type)
	r.Description = types.StringValue(permission.Description)
}

func (r *permissionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permission"
}

func (r *permissionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage a workspace permission (role)\n\n" +
			"Workspace administrator rights are required\n\n" +
			"Token scope required: `WORKSPACE`",
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
			"name": schema.StringAttribute{
				MarkdownDescription: "The permission's name",
				Required:            true,
			},
			"pipeline_access_level": schema.StringAttribute{
				MarkdownDescription: "The permission's access level to pipelines. Allowed: `DENIED`, `READ_ONLY`, `RUN_ONLY`, `READ_WRITE`",
				Required:            true,
				Validators: []validator.String{stringvalidator.OneOf(
					buddy.PermissionAccessLevelDenied,
					buddy.PermissionAccessLevelReadOnly,
					buddy.PermissionAccessLevelRunOnly,
					buddy.PermissionAccessLevelReadWrite,
				)},
			},
			"repository_access_level": schema.StringAttribute{
				MarkdownDescription: "The permission's access level to repository. Allowed: `READ_ONLY`, `READ_WRITE`, `MANAGE`",
				Required:            true,
				Validators: []validator.String{stringvalidator.OneOf(
					buddy.PermissionAccessLevelReadOnly,
					buddy.PermissionAccessLevelReadWrite,
					buddy.PermissionAccessLevelManage,
				)},
			},
			"sandbox_access_level": schema.StringAttribute{
				MarkdownDescription: "The permission's access level to sandboxes. Allowed: `DENIED`, `READ_ONLY`, `READ_WRITE`",
				Required:            true,
				Validators: []validator.String{stringvalidator.OneOf(
					buddy.PermissionAccessLevelDenied,
					buddy.PermissionAccessLevelReadOnly,
					buddy.PermissionAccessLevelReadWrite,
				)},
			},
			"project_team_access_level": schema.StringAttribute{
				MarkdownDescription: "The permission's access level to team. Allowed: `READ_ONLY`, `MANAGE`",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{stringvalidator.OneOf(
					buddy.PermissionAccessLevelReadOnly,
					buddy.PermissionAccessLevelManage,
				)},
			},
			"environment_access_level": schema.StringAttribute{
				MarkdownDescription: "The permission's access level to environments. Allowed: `DENIED`, `MANAGE`, `USE_ONLY`",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{stringvalidator.OneOf(
					buddy.PermissionAccessLevelDenied,
					buddy.PermissionAccessLevelManage,
					buddy.PermissionAccessLevelUseOnly,
				)},
			},
			"target_access_level": schema.StringAttribute{
				MarkdownDescription: "The permission's access level to environments. Allowed: `DENIED`, 'READ_ONLY`, `MANAGE`, `USE_ONLY`",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{stringvalidator.OneOf(
					buddy.PermissionAccessLevelDenied,
					buddy.PermissionAccessLevelManage,
					buddy.PermissionAccessLevelUseOnly,
					buddy.PermissionAccessLevelReadOnly,
				)},
			},
			"permission_id": schema.Int64Attribute{
				MarkdownDescription: "The permission's ID",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The permission's description",
				Optional:            true,
				Computed:            true,
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The permission's URL",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The permission's type",
				Computed:            true,
			},
		},
	}
}

func (r *permissionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *permissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *permissionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	ops := buddy.PermissionOps{
		Name:                  data.Name.ValueStringPointer(),
		PipelineAccessLevel:   data.PipelineAccessLevel.ValueStringPointer(),
		RepositoryAccessLevel: data.RepositoryAccessLevel.ValueStringPointer(),
		SandboxAccessLevel:    data.SandboxAccessLevel.ValueStringPointer(),
	}
	if !data.ProjectTeamAccessLevel.IsNull() && !data.ProjectTeamAccessLevel.IsUnknown() {
		ops.ProjectTeamAccessLevel = data.ProjectTeamAccessLevel.ValueStringPointer()
	}
	if !data.TargetAccessLevel.IsNull() && !data.TargetAccessLevel.IsUnknown() {
		ops.TargetAccessLevel = data.TargetAccessLevel.ValueStringPointer()
	}
	if !data.EnvironmentAccessLevel.IsNull() && !data.EnvironmentAccessLevel.IsUnknown() {
		ops.EnvironmentAccessLevel = data.EnvironmentAccessLevel.ValueStringPointer()
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		ops.Description = data.Description.ValueStringPointer()
	}
	permission, _, err := r.client.PermissionService.Create(domain, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create permission", err))
		return
	}
	data.loadAPI(domain, permission)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *permissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *permissionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, permissionId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("permission", err))
		return
	}
	permission, httpResp, err := r.client.PermissionService.Get(domain, permissionId)
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get permission", err))
		return
	}
	data.loadAPI(domain, permission)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *permissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *permissionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, permissionId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("permission", err))
		return
	}
	ops := buddy.PermissionOps{}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		ops.Name = data.Name.ValueStringPointer()
	}
	if !data.PipelineAccessLevel.IsNull() && !data.PipelineAccessLevel.IsUnknown() {
		ops.PipelineAccessLevel = data.PipelineAccessLevel.ValueStringPointer()
	}
	if !data.RepositoryAccessLevel.IsNull() && !data.RepositoryAccessLevel.IsUnknown() {
		ops.RepositoryAccessLevel = data.RepositoryAccessLevel.ValueStringPointer()
	}
	if !data.SandboxAccessLevel.IsNull() && !data.SandboxAccessLevel.IsUnknown() {
		ops.SandboxAccessLevel = data.SandboxAccessLevel.ValueStringPointer()
	}
	if !data.ProjectTeamAccessLevel.IsNull() && !data.ProjectTeamAccessLevel.IsUnknown() {
		ops.ProjectTeamAccessLevel = data.ProjectTeamAccessLevel.ValueStringPointer()
	}
	if !data.TargetAccessLevel.IsNull() && !data.TargetAccessLevel.IsUnknown() {
		ops.TargetAccessLevel = data.TargetAccessLevel.ValueStringPointer()
	}
	if !data.EnvironmentAccessLevel.IsNull() && !data.EnvironmentAccessLevel.IsUnknown() {
		ops.EnvironmentAccessLevel = data.EnvironmentAccessLevel.ValueStringPointer()
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		ops.Description = data.Description.ValueStringPointer()
	}
	permission, _, err := r.client.PermissionService.Update(domain, permissionId, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update permission", err))
		return
	}
	data.loadAPI(domain, permission)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *permissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *permissionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, permissionId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("permission", err))
		return
	}
	_, err = r.client.PermissionService.Delete(domain, permissionId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete permission", err))
	}
}

func (r *permissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

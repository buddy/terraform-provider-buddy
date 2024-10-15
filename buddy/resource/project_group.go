package resource

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ resource.Resource                = &projectGroupResource{}
	_ resource.ResourceWithConfigure   = &projectGroupResource{}
	_ resource.ResourceWithImportState = &projectGroupResource{}
)

func NewProjectGroupResource() resource.Resource {
	return &projectGroupResource{}
}

type projectGroupResource struct {
	client *buddy.Client
}

type projectGroupResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Domain       types.String `tfsdk:"domain"`
	ProjectName  types.String `tfsdk:"project_name"`
	GroupId      types.Int64  `tfsdk:"group_id"`
	PermissionId types.Int64  `tfsdk:"permission_id"`
	HtmlUrl      types.String `tfsdk:"html_url"`
	Name         types.String `tfsdk:"name"`
	Permission   types.Set    `tfsdk:"permission"`
}

func (r *projectGroupResourceModel) loadAPI(ctx context.Context, domain string, projectName string, projectGroup *buddy.ProjectGroup) diag.Diagnostics {
	r.ID = types.StringValue(util.ComposeTripleId(domain, projectName, strconv.Itoa(projectGroup.Id)))
	r.Domain = types.StringValue(domain)
	r.ProjectName = types.StringValue(projectName)
	r.GroupId = types.Int64Value(int64(projectGroup.Id))
	r.PermissionId = types.Int64Value(int64(projectGroup.PermissionSet.Id))
	r.HtmlUrl = types.StringValue(projectGroup.HtmlUrl)
	r.Name = types.StringValue(projectGroup.Name)
	permissionSet := []*buddy.Permission{projectGroup.PermissionSet}
	permission, diags := util.PermissionsModelFromApi(ctx, &permissionSet)
	r.Permission = permission
	return diags
}

func (r *projectGroupResourceModel) decomposeId() (string, string, int, error) {
	domain, projectName, gid, err := util.DecomposeTripleId(r.ID.ValueString())
	if err != nil {
		return "", "", 0, err
	}
	groupId, err := strconv.Atoi(gid)
	if err != nil {
		return "", "", 0, err
	}
	return domain, projectName, groupId, nil
}

func (r *projectGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_group"
}

func (r *projectGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage a workspace project group permission\n\n" +
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
			"project_name": schema.StringAttribute{
				MarkdownDescription: "The project's name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"group_id": schema.Int64Attribute{
				MarkdownDescription: "The group's ID",
				Required:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"permission_id": schema.Int64Attribute{
				MarkdownDescription: "The permission's ID",
				Required:            true,
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The group's URL",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The group's name",
				Computed:            true,
			},
			// for compatibility, it's a set
			"permission": schema.SetNestedAttribute{
				MarkdownDescription: "The group's permission in the project",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.PermissionModelAttributes(),
				},
			},
		},
	}
}

func (r *projectGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *projectGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *projectGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	projectName := data.ProjectName.ValueString()
	projectGroup, _, err := r.client.ProjectGroupService.CreateProjectGroup(domain, projectName, &buddy.ProjectGroupOps{
		Id: util.PointerInt(data.GroupId.ValueInt64()),
		PermissionSet: &buddy.ProjectGroupOps{
			Id: util.PointerInt(data.PermissionId.ValueInt64()),
		},
	})
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create project group", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, projectGroup)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *projectGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, projectName, groupId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("project group", err))
		return
	}
	projectGroup, httpResp, err := r.client.ProjectGroupService.GetProjectGroup(domain, projectName, groupId)
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get project group", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, projectGroup)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *projectGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, projectName, groupId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("project group", err))
		return
	}
	projectGroup, _, err := r.client.ProjectGroupService.UpdateProjectGroup(domain, projectName, groupId, &buddy.ProjectGroupOps{
		PermissionSet: &buddy.ProjectGroupOps{
			Id: util.PointerInt(data.PermissionId.ValueInt64()),
		},
	})
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update project group", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, projectGroup)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *projectGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, projectName, groupId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("project group", err))
		return
	}
	_, err = r.client.ProjectGroupService.DeleteProjectGroup(domain, projectName, groupId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete project group", err))
	}
}

func (r *projectGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

package resource

import (
	"buddy-terraform/buddy/util"
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
)

var (
	_ resource.Resource                = &projectMemberResource{}
	_ resource.ResourceWithConfigure   = &projectMemberResource{}
	_ resource.ResourceWithImportState = &projectMemberResource{}
)

func NewProjectMemberResource() resource.Resource {
	return &projectMemberResource{}
}

type projectMemberResource struct {
	client *buddy.Client
}

type projectMemberResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Domain         types.String `tfsdk:"domain"`
	ProjectName    types.String `tfsdk:"project_name"`
	MemberId       types.Int64  `tfsdk:"member_id"`
	PermissionId   types.Int64  `tfsdk:"permission_id"`
	HtmlUrl        types.String `tfsdk:"html_url"`
	Name           types.String `tfsdk:"name"`
	Email          types.String `tfsdk:"email"`
	AvatarUrl      types.String `tfsdk:"avatar_url"`
	Admin          types.Bool   `tfsdk:"admin"`
	WorkspaceOwner types.Bool   `tfsdk:"workspace_owner"`
	Permission     types.Object `tfsdk:"permission"`
}

func (r *projectMemberResourceModel) loadAPI(ctx context.Context, domain string, projectName string, projectMember *buddy.ProjectMember) diag.Diagnostics {
	r.ID = types.StringValue(util.ComposeTripleId(domain, projectName, strconv.Itoa(projectMember.Id)))
	r.Domain = types.StringValue(domain)
	r.ProjectName = types.StringValue(projectName)
	r.MemberId = types.Int64Value(int64(projectMember.Id))
	r.PermissionId = types.Int64Value(int64(projectMember.PermissionSet.Id))
	r.HtmlUrl = types.StringValue(projectMember.HtmlUrl)
	r.Name = types.StringValue(projectMember.Name)
	r.Email = types.StringValue(projectMember.Email)
	r.AvatarUrl = types.StringValue(projectMember.AvatarUrl)
	r.Admin = types.BoolValue(projectMember.Admin)
	r.WorkspaceOwner = types.BoolValue(projectMember.WorkspaceOwner)
	permission, diags := util.PermissionTypeValueFrom(ctx, projectMember.PermissionSet)
	r.Permission = permission
	return diags
}

func (r *projectMemberResourceModel) decomposeId() (string, string, int, error) {
	domain, projectName, mid, err := util.DecomposeTripleId(r.ID.ValueString())
	if err != nil {
		return "", "", 0, err
	}
	memberId, err := strconv.Atoi(mid)
	if err != nil {
		return "", "", 0, err
	}
	return domain, projectName, memberId, nil
}

func (r *projectMemberResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_member"
}

func (r *projectMemberResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage a member's permission (role) in a project\n\n" +
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
			"member_id": schema.Int64Attribute{
				MarkdownDescription: "The member's ID",
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
				MarkdownDescription: "The member's URL",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The member's name",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The member's email",
				Computed:            true,
			},
			"avatar_url": schema.StringAttribute{
				MarkdownDescription: "The member's avatar URL",
				Computed:            true,
			},
			"admin": schema.BoolAttribute{
				MarkdownDescription: "Is the member a workspace administrator",
				Computed:            true,
			},
			"workspace_owner": schema.BoolAttribute{
				MarkdownDescription: "Is the member the workspace owner",
				Computed:            true,
			},
			"permission": schema.SingleNestedAttribute{
				MarkdownDescription: "The member's permission in the project",
				Computed:            true,
				Attributes:          util.PermissionTypeComputedAttributes(),
			},
		},
	}
}

func (r *projectMemberResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *projectMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *projectMemberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	projectName := data.ProjectName.ValueString()
	projectMember, _, err := r.client.ProjectMemberService.CreateProjectMember(domain, projectName, &buddy.ProjectMemberOps{
		Id: util.PointerInt(data.MemberId.ValueInt64()),
		PermissionSet: &buddy.ProjectMemberOps{
			Id: util.PointerInt(data.PermissionId.ValueInt64()),
		},
	})
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create project member", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, projectMember)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *projectMemberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, projectName, memberId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("project member", err))
		return
	}
	projectMember, httpResp, err := r.client.ProjectMemberService.GetProjectMember(domain, projectName, memberId)
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get project member", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, projectMember)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *projectMemberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, projectName, memberId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("project member", err))
		return
	}
	projectMember, _, err := r.client.ProjectMemberService.UpdateProjectMember(domain, projectName, memberId, &buddy.ProjectMemberOps{
		PermissionSet: &buddy.ProjectMemberOps{
			Id: util.PointerInt(data.PermissionId.ValueInt64()),
		},
	})
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update project member", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, projectMember)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *projectMemberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, projectName, memberId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("project member", err))
		return
	}
	_, err = r.client.ProjectMemberService.DeleteProjectMember(domain, projectName, memberId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete project member", err))
	}
}

func (r *projectMemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

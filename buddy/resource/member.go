package resource

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ resource.Resource                = &memberResource{}
	_ resource.ResourceWithConfigure   = &memberResource{}
	_ resource.ResourceWithImportState = &memberResource{}
)

func NewMemberResource() resource.Resource {
	return &memberResource{}
}

type memberResource struct {
	client *buddy.Client
}

type memberResourceModel struct {
	ID                        types.String `tfsdk:"id"`
	Domain                    types.String `tfsdk:"domain"`
	Email                     types.String `tfsdk:"email"`
	Admin                     types.Bool   `tfsdk:"admin"`
	AutoAssignToNewProjects   types.Bool   `tfsdk:"auto_assign_to_new_projects"`
	AutoAssignPermissionSetId types.Int64  `tfsdk:"auto_assign_permission_set_id"`
	Name                      types.String `tfsdk:"name"`
	MemberId                  types.Int64  `tfsdk:"member_id"`
	HtmlUrl                   types.String `tfsdk:"html_url"`
	AvatarUrl                 types.String `tfsdk:"avatar_url"`
	WorkspaceOwner            types.Bool   `tfsdk:"workspace_owner"`
}

func (r *memberResourceModel) decomposeId() (string, int, error) {
	domain, mid, err := util.DecomposeDoubleId(r.ID.ValueString())
	if err != nil {
		return "", 0, err
	}
	memberId, err := strconv.Atoi(mid)
	if err != nil {
		return "", 0, err
	}
	return domain, memberId, nil
}

func (r *memberResourceModel) loadAPI(domain string, member *buddy.Member) {
	r.ID = types.StringValue(util.ComposeDoubleId(domain, strconv.Itoa(member.Id)))
	r.Domain = types.StringValue(domain)
	r.Email = types.StringValue(member.Email)
	r.Admin = types.BoolValue(member.Admin)
	r.AutoAssignToNewProjects = types.BoolValue(member.AutoAssignToNewProjects)
	r.Name = types.StringValue(member.Name)
	r.MemberId = types.Int64Value(int64(member.Id))
	r.HtmlUrl = types.StringValue(member.HtmlUrl)
	r.AvatarUrl = types.StringValue(member.AvatarUrl)
	r.WorkspaceOwner = types.BoolValue(member.WorkspaceOwner)
	// auto_assign_permission_set_id we are leaving this prop value as set by client
}

func (r *memberResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_member"
}

func (r *memberResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage a workspace member\n\n" +
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
			"email": schema.StringAttribute{
				MarkdownDescription: "The member's email",
				Required:            true,
				Validators:          util.StringValidatorsEmail(),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"admin": schema.BoolAttribute{
				MarkdownDescription: "Is the member a workspace administrator",
				Optional:            true,
				Computed:            true,
			},
			"auto_assign_to_new_projects": schema.BoolAttribute{
				MarkdownDescription: "Defines whether or not to automatically assign member to new projects",
				Optional:            true,
				Computed:            true,
			},
			"auto_assign_permission_set_id": schema.Int64Attribute{
				MarkdownDescription: "The permission's ID with which the member will be assigned to new projects",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The member's name",
				Computed:            true,
			},
			"member_id": schema.Int64Attribute{
				MarkdownDescription: "The member's ID",
				Computed:            true,
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The member's URL",
				Computed:            true,
			},
			"avatar_url": schema.StringAttribute{
				MarkdownDescription: "The member's avatar URL",
				Computed:            true,
			},
			"workspace_owner": schema.BoolAttribute{
				MarkdownDescription: "Is the member the workspace owner",
				Computed:            true,
			},
		},
	}
}

func (r *memberResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *memberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *memberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	ops := buddy.MemberCreateOps{
		Email: data.Email.ValueStringPointer(),
	}
	member, _, err := r.client.MemberService.Create(domain, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create member", err))
		return
	}
	updateAdmin := !data.Admin.IsNull() && !data.Admin.IsUnknown()
	updateAutoAssign := !data.AutoAssignToNewProjects.IsNull() && !data.AutoAssignToNewProjects.IsUnknown()
	if updateAdmin || updateAutoAssign {
		updateOps := buddy.MemberUpdateOps{}
		if updateAdmin {
			updateOps.Admin = data.Admin.ValueBoolPointer()
		}
		if updateAutoAssign {
			updateOps.AutoAssignToNewProjects = data.AutoAssignToNewProjects.ValueBoolPointer()
			if !data.AutoAssignPermissionSetId.IsNull() && !data.AutoAssignPermissionSetId.IsUnknown() {
				updateOps.AutoAssignPermissionSetId = util.PointerInt(data.AutoAssignPermissionSetId.ValueInt64())
			}
		}
		member, _, err = r.client.MemberService.Update(domain, member.Id, &updateOps)
		if err != nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiError("create Member", err))
			return
		}
	}
	data.loadAPI(domain, member)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *memberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *memberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, memberId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("member", err))
		return
	}
	member, httpResp, err := r.client.MemberService.Get(domain, memberId)
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get member", err))
		return
	}
	data.loadAPI(domain, member)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *memberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *memberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, memberId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("member", err))
		return
	}
	ops := buddy.MemberUpdateOps{}
	if !data.Admin.IsNull() && !data.Admin.IsUnknown() {
		ops.Admin = data.Admin.ValueBoolPointer()
	}
	if !data.AutoAssignToNewProjects.IsNull() && !data.AutoAssignToNewProjects.IsUnknown() {
		ops.AutoAssignToNewProjects = data.AutoAssignToNewProjects.ValueBoolPointer()
		if !data.AutoAssignPermissionSetId.IsNull() && !data.AutoAssignPermissionSetId.IsUnknown() {
			ops.AutoAssignPermissionSetId = util.PointerInt(data.AutoAssignPermissionSetId.ValueInt64())
		}
	}
	member, _, err := r.client.MemberService.Update(domain, memberId, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update member", err))
		return
	}
	data.loadAPI(domain, member)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *memberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *memberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, memberId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("member", err))
		return
	}
	_, err = r.client.MemberService.Delete(domain, memberId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete member", err))
	}
}

func (r *memberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

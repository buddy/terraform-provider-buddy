package resource

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	_ resource.Resource                = &groupMemberResource{}
	_ resource.ResourceWithConfigure   = &groupMemberResource{}
	_ resource.ResourceWithImportState = &groupMemberResource{}
)

func NewGroupMemberResource() resource.Resource {
	return &groupMemberResource{}
}

type groupMemberResource struct {
	client *buddy.Client
}

type groupMemberResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Domain         types.String `tfsdk:"domain"`
	GroupId        types.Int64  `tfsdk:"group_id"`
	MemberId       types.Int64  `tfsdk:"member_id"`
	Status         types.String `tfsdk:"status"`
	HtmlUrl        types.String `tfsdk:"html_url"`
	Name           types.String `tfsdk:"name"`
	Email          types.String `tfsdk:"email"`
	AvatarUrl      types.String `tfsdk:"avatar_url"`
	Admin          types.Bool   `tfsdk:"admin"`
	WorkspaceOwner types.Bool   `tfsdk:"workspace_owner"`
}

func (r *groupMemberResourceModel) decomposeId() (string, int, int, error) {
	domain, gid, mid, err := util.DecomposeTripleId(r.ID.ValueString())
	if err != nil {
		return "", 0, 0, err
	}
	groupId, err := strconv.Atoi(gid)
	if err != nil {
		return "", 0, 0, err
	}
	memberId, err := strconv.Atoi(mid)
	if err != nil {
		return "", 0, 0, err
	}
	return domain, groupId, memberId, nil
}

func (r *groupMemberResourceModel) loadAPI(domain string, groupId int, member *buddy.Member) {
	r.ID = types.StringValue(util.ComposeTripleId(domain, strconv.Itoa(groupId), strconv.Itoa(member.Id)))
	r.Domain = types.StringValue(domain)
	r.GroupId = types.Int64Value(int64(groupId))
	r.MemberId = types.Int64Value(int64(member.Id))
	r.Status = types.StringValue(member.Status)
	r.HtmlUrl = types.StringValue(member.HtmlUrl)
	r.Name = types.StringValue(member.Name)
	r.Email = types.StringValue(member.Email)
	r.AvatarUrl = types.StringValue(member.AvatarUrl)
	r.Admin = types.BoolValue(member.Admin)
	r.WorkspaceOwner = types.BoolValue(member.WorkspaceOwner)
}

func (r *groupMemberResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_member"
}

func (r *groupMemberResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage a workspace group member\n\n" +
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
			"group_id": schema.Int64Attribute{
				MarkdownDescription: "The group's ID",
				Required:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"member_id": schema.Int64Attribute{
				MarkdownDescription: "The member's ID",
				Required:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The member's status. Allowed: `MEMBER`, `MANAGER`",
				Validators: []validator.String{stringvalidator.OneOf(
					buddy.GroupMemberStatusMember,
					buddy.GroupMemberStatusManager,
				)},
				Optional: true,
				Computed: true,
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
		},
	}
}

func (r *groupMemberResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *groupMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *groupMemberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	groupId := int(data.GroupId.ValueInt64())
	ops := buddy.GroupMemberOps{
		Id: util.PointerInt(data.MemberId.ValueInt64()),
	}
	if !data.Status.IsNull() && !data.Status.IsUnknown() {
		ops.Status = data.Status.ValueStringPointer()
	}
	member, _, err := r.client.GroupService.AddGroupMember(domain, groupId, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create group member", err))
		return
	}
	data.loadAPI(domain, groupId, member)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *groupMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *groupMemberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, groupId, memberId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("group_member", err))
		return
	}
	member, httpResp, err := r.client.GroupService.GetGroupMember(domain, groupId, memberId)
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get group member", err))
		return
	}
	data.loadAPI(domain, groupId, member)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *groupMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *groupMemberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, groupId, memberId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("group_member", err))
		return
	}
	status := buddy.GroupMemberStatusMember
	if !data.Status.IsNull() && !data.Status.IsUnknown() {
		status = data.Status.ValueString()
	}
	member, _, err := r.client.GroupService.UpdateGroupMember(domain, groupId, memberId, status)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update group member", err))
		return
	}
	data.loadAPI(domain, groupId, member)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *groupMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *groupMemberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, groupId, memberId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("group_member", err))
		return
	}
	_, err = r.client.GroupService.DeleteGroupMember(domain, groupId, memberId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete group_member", err))
	}
}

func (r *groupMemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

package resource

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
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
	_ resource.Resource                = &groupResource{}
	_ resource.ResourceWithConfigure   = &groupResource{}
	_ resource.ResourceWithImportState = &groupResource{}
)

func NewGroupResource() resource.Resource {
	return &groupResource{}
}

type groupResource struct {
	client *buddy.Client
}

type groupResourceModel struct {
	ID                        types.String `tfsdk:"id"`
	Domain                    types.String `tfsdk:"domain"`
	Name                      types.String `tfsdk:"name"`
	AutoAssignToNewProjects   types.Bool   `tfsdk:"auto_assign_to_new_projects"`
	AutoAssignPermissionSetId types.Int64  `tfsdk:"auto_assign_permission_set_id"`
	GroupId                   types.Int64  `tfsdk:"group_id"`
	HtmlUrl                   types.String `tfsdk:"html_url"`
	Description               types.String `tfsdk:"description"`
}

func (r *groupResourceModel) decomposeId() (string, int, error) {
	domain, gid, err := util.DecomposeDoubleId(r.ID.ValueString())
	if err != nil {
		return "", 0, err
	}
	groupId, err := strconv.Atoi(gid)
	if err != nil {
		return "", 0, err
	}
	return domain, groupId, nil
}

func (r *groupResourceModel) loadAPI(domain string, group *buddy.Group) {
	r.ID = types.StringValue(util.ComposeDoubleId(domain, strconv.Itoa(group.Id)))
	r.Domain = types.StringValue(domain)
	r.Name = types.StringValue(group.Name)
	r.GroupId = types.Int64Value(int64(group.Id))
	r.HtmlUrl = types.StringValue(group.HtmlUrl)
	r.Description = types.StringValue(group.Description)
	r.AutoAssignToNewProjects = types.BoolValue(group.AutoAssignToNewProjects)
	// auto_assign_permission_set_id we are leaving this prop value as set by client
}

func (r *groupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *groupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage a user's group\n\n" +
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
				MarkdownDescription: "The group's name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The group's description",
				Optional:            true,
				Computed:            true,
			},
			"auto_assign_to_new_projects": schema.BoolAttribute{
				MarkdownDescription: "Defines whether or not to automatically assign group to new projects",
				Optional:            true,
				Computed:            true,
			},
			"auto_assign_permission_set_id": schema.Int64Attribute{
				MarkdownDescription: "The permission's ID with which the group will be assigned to new projects",
				Optional:            true,
			},
			"group_id": schema.Int64Attribute{
				MarkdownDescription: "The group's ID",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The group's URL",
				Computed:            true,
			},
		},
	}
}

func (r *groupResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *groupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	ops := buddy.GroupOps{
		Name: data.Name.ValueStringPointer(),
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		ops.Description = data.Description.ValueStringPointer()
	}
	if !data.AutoAssignToNewProjects.IsNull() && !data.AutoAssignToNewProjects.IsUnknown() {
		ops.AutoAssignToNewProjects = data.AutoAssignToNewProjects.ValueBoolPointer()
		if !data.AutoAssignPermissionSetId.IsNull() && !data.AutoAssignPermissionSetId.IsUnknown() {
			ops.AutoAssignPermissionSetId = util.PointerInt(data.AutoAssignPermissionSetId.ValueInt64())
		}
	}
	group, _, err := r.client.GroupService.Create(domain, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create group", err))
		return
	}
	data.loadAPI(domain, group)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *groupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, groupId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("group", err))
		return
	}
	group, httpResp, err := r.client.GroupService.Get(domain, groupId)
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get group", err))
		return
	}
	data.loadAPI(domain, group)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *groupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, groupId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("group", err))
		return
	}
	ops := buddy.GroupOps{
		Name: data.Name.ValueStringPointer(),
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		ops.Description = data.Description.ValueStringPointer()
	}
	if !data.AutoAssignToNewProjects.IsNull() && !data.AutoAssignToNewProjects.IsUnknown() {
		ops.AutoAssignToNewProjects = data.AutoAssignToNewProjects.ValueBoolPointer()
		if !data.AutoAssignPermissionSetId.IsNull() && !data.AutoAssignPermissionSetId.IsUnknown() {
			ops.AutoAssignPermissionSetId = util.PointerInt(data.AutoAssignPermissionSetId.ValueInt64())
		}
	}
	group, _, err := r.client.GroupService.Update(domain, groupId, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update group", err))
		return
	}
	data.loadAPI(domain, group)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *groupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, groupId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("group", err))
		return
	}
	_, err = r.client.GroupService.Delete(domain, groupId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete group", err))
	}
}

func (r *groupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

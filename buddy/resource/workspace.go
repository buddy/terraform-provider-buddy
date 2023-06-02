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
	"terraform-provider-buddy/buddy/util"
	"time"
)

var (
	_ resource.Resource                = &workspaceResource{}
	_ resource.ResourceWithConfigure   = &workspaceResource{}
	_ resource.ResourceWithImportState = &workspaceResource{}
)

func NewWorkspaceResource() resource.Resource {
	return &workspaceResource{}
}

type workspaceResource struct {
	client *buddy.Client
}

type workspaceResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Domain         types.String `tfsdk:"domain"`
	Name           types.String `tfsdk:"name"`
	EncryptionSalt types.String `tfsdk:"encryption_salt"`
	WorkspaceId    types.Int64  `tfsdk:"workspace_id"`
	HtmlUrl        types.String `tfsdk:"html_url"`
	OwnerId        types.Int64  `tfsdk:"owner_id"`
	Frozen         types.Bool   `tfsdk:"frozen"`
	CreateDate     types.String `tfsdk:"create_date"`
}

func (r *workspaceResourceModel) loadAPI(workspace *buddy.Workspace) {
	r.ID = types.StringValue(workspace.Domain)
	r.Domain = types.StringValue(workspace.Domain)
	r.Name = types.StringValue(workspace.Name)
	r.WorkspaceId = types.Int64Value(int64(workspace.Id))
	r.HtmlUrl = types.StringValue(workspace.HtmlUrl)
	r.OwnerId = types.Int64Value(int64(workspace.OwnerId))
	r.Frozen = types.BoolValue(workspace.Frozen)
	cd, err := time.Parse(time.RFC3339, workspace.CreateDate)
	if err == nil {
		// fix seconds
		r.CreateDate = types.StringValue(cd.Format(time.RFC3339))
	} else {
		r.CreateDate = types.StringValue(workspace.CreateDate)
	}
	// EncryptionSalt is not returned from api
}

func (r *workspaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace"
}

func (r *workspaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage a workspace\n\n" +
			"Invite-only token is required. Contact support@buddy.works for more details\n\n" +
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
				MarkdownDescription: "The workspace's name",
				Optional:            true,
				Computed:            true,
			},
			"encryption_salt": schema.StringAttribute{
				MarkdownDescription: "The workspace's salt to encrypt secrets in YAML & API",
				Optional:            true,
			},
			"workspace_id": schema.Int64Attribute{
				MarkdownDescription: "The workspace's ID",
				Computed:            true,
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The workspace's URL",
				Computed:            true,
			},
			"owner_id": schema.Int64Attribute{
				MarkdownDescription: "The workspace's owner ID",
				Computed:            true,
			},
			"frozen": schema.BoolAttribute{
				MarkdownDescription: "Is the workspace frozen",
				Computed:            true,
			},
			"create_date": schema.StringAttribute{
				MarkdownDescription: "The workspace's create date",
				Computed:            true,
			},
		},
	}
}

func (r *workspaceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *workspaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *workspaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *workspaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ops := buddy.WorkspaceCreateOps{
		Domain: data.Domain.ValueStringPointer(),
	}
	if !data.EncryptionSalt.IsNull() && !data.EncryptionSalt.IsUnknown() {
		ops.EncryptionSalt = data.EncryptionSalt.ValueStringPointer()
	}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		ops.Name = data.Name.ValueStringPointer()
	}
	workspace, _, err := r.client.WorkspaceService.Create(&ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create workspace", err))
		return
	}
	data.loadAPI(workspace)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *workspaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *workspaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	workspace, httpResp, err := r.client.WorkspaceService.Get(data.ID.ValueString())
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get workspace", err))
		return
	}
	data.loadAPI(workspace)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *workspaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *workspaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ops := buddy.WorkspaceUpdateOps{}
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		ops.Name = data.Name.ValueStringPointer()
	}
	if !data.EncryptionSalt.IsNull() && !data.EncryptionSalt.IsUnknown() {
		ops.EncryptionSalt = data.EncryptionSalt.ValueStringPointer()
	}
	workspace, _, err := r.client.WorkspaceService.Update(data.ID.ValueString(), &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update workspace", err))
		return
	}
	data.loadAPI(workspace)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *workspaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *workspaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, err := r.client.WorkspaceService.Delete(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete workspace", err))
	}
}

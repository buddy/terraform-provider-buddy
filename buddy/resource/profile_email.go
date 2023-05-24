package resource

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &profileEmailResource{}
	_ resource.ResourceWithConfigure   = &profileEmailResource{}
	_ resource.ResourceWithImportState = &profileEmailResource{}
)

func NewProfileEmailResource() resource.Resource {
	return &profileEmailResource{}
}

type profileEmailResource struct {
	client *buddy.Client
}

type profileEmailResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Email     types.String `tfsdk:"email"`
	Confirmed types.Bool   `tfsdk:"confirmed"`
}

func (r *profileEmailResourceModel) loadAPI(pe *buddy.ProfileEmail) {
	r.ID = types.StringValue(pe.Email)
	r.Email = types.StringValue(pe.Email)
	r.Confirmed = types.BoolValue(pe.Confirmed)
}

func (r *profileEmailResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_profile_email"
}

func (r *profileEmailResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage a user's email\n\n" +
			"Token scopes required: `MANAGE_EMAILS`, `USER_EMAIL`",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Terraform resource identifier for this item",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The email to add to the user's profile",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"confirmed": schema.BoolAttribute{
				MarkdownDescription: "Is the email confirmed",
				Computed:            true,
			},
		},
	}
}

func (r *profileEmailResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *profileEmailResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *profileEmailResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	p, _, err := r.client.ProfileEmailService.Create(&buddy.ProfileEmailOps{
		Email: data.Email.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create profile email", err))
		return
	}
	data.loadAPI(p)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *profileEmailResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// do nothing
}

func (r *profileEmailResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *profileEmailResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	p, httpResp, err := r.client.ProfileEmailService.GetList()
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get profile email", err))
		return
	}
	found := false
	email := data.ID.ValueString()
	for _, v := range p.Emails {
		if v.Email == email {
			found = true
			data.loadAPI(v)
			break
		}
	}
	if !found {
		resp.State.RemoveResource(ctx)
	} else {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
}

func (r *profileEmailResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *profileEmailResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, err := r.client.ProfileEmailService.Delete(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete profile email", err))
		return
	}
}

func (r *profileEmailResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

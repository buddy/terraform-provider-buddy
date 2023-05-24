package resource

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
)

var (
	_ resource.Resource                = &profileResource{}
	_ resource.ResourceWithConfigure   = &profileResource{}
	_ resource.ResourceWithImportState = &profileResource{}
)

func NewProfileResource() resource.Resource {
	return &profileResource{}
}

type profileResource struct {
	client *buddy.Client
}

type profileResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	MemberId  types.Int64  `tfsdk:"member_id"`
	AvatarUrl types.String `tfsdk:"avatar_url"`
	HtmlUrl   types.String `tfsdk:"html_url"`
}

func (r *profileResourceModel) loadAPI(profile *buddy.Profile) {
	r.ID = types.StringValue(strconv.Itoa(profile.Id))
	r.Name = types.StringValue(profile.Name)
	r.MemberId = types.Int64Value(int64(profile.Id))
	r.AvatarUrl = types.StringValue(profile.AvatarUrl)
	r.HtmlUrl = types.StringValue(profile.HtmlUrl)
}

func (r *profileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_profile"
}

func (r *profileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage a user profile\n\n" +
			"Token scope required: `USER_INFO`",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Terraform resource identifier for this item",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The user's name",
				Required:            true,
			},
			"member_id": schema.Int64Attribute{
				MarkdownDescription: "The user's ID",
				Computed:            true,
			},
			"avatar_url": schema.StringAttribute{
				MarkdownDescription: "The user's avatar URL",
				Computed:            true,
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The user's URL",
				Computed:            true,
			},
		},
	}
}

func (r *profileResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *profileResource) update(ctx context.Context, diagnostics *diag.Diagnostics, plan *tfsdk.Plan, state *tfsdk.State) {
	var data *profileResourceModel
	diagnostics.Append(plan.Get(ctx, &data)...)
	if diagnostics.HasError() {
		return
	}
	ops := buddy.ProfileOps{
		Name: data.Name.ValueStringPointer(),
	}
	profile, _, err := r.client.ProfileService.Update(&ops)
	if err != nil {
		diagnostics.Append(util.NewDiagnosticApiError("update profile", err))
		return
	}
	data.loadAPI(profile)
	diagnostics.Append(state.Set(ctx, &data)...)
}

func (r *profileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.update(ctx, &resp.Diagnostics, &req.Plan, &resp.State)
}

func (r *profileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *profileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	profile, _, err := r.client.ProfileService.Get()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get profile", err))
		return
	}
	data.loadAPI(profile)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *profileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.update(ctx, &resp.Diagnostics, &req.Plan, &resp.State)
}

func (r *profileResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// do nothing
}

func (r *profileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

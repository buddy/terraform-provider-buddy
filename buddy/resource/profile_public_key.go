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
	_ resource.Resource                = &profilePublicKeyResource{}
	_ resource.ResourceWithConfigure   = &profilePublicKeyResource{}
	_ resource.ResourceWithImportState = &profilePublicKeyResource{}
)

func NewProfilePublicKeyResource() resource.Resource {
	return &profilePublicKeyResource{}
}

type profilePublicKeyResource struct {
	client *buddy.Client
}

type profilePublicKeyResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Content types.String `tfsdk:"content"`
	Title   types.String `tfsdk:"title"`
	HtmlUrl types.String `tfsdk:"html_url"`
}

func (r *profilePublicKeyResourceModel) loadAPI(key *buddy.PublicKey) {
	r.ID = types.StringValue(strconv.Itoa(key.Id))
	r.Content = types.StringValue(key.Content)
	r.Title = types.StringValue(key.Title)
	r.HtmlUrl = types.StringValue(key.HtmlUrl)
}

func (r *profilePublicKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_profile_public_key"
}

func (r *profilePublicKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage a user's public key\n\n" +
			"Token scope required: `USER_KEY`",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Terraform resource identifier for this item",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "The public key's content (starts with ssh-rsa)",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "The public key's title",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The public key's URL",
				Computed:            true,
			},
		},
	}
}

func (r *profilePublicKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *profilePublicKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *profilePublicKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ops := buddy.PublicKeyOps{
		Content: data.Content.ValueStringPointer(),
	}
	if !data.Title.IsNull() && !data.Title.IsUnknown() {
		ops.Title = data.Title.ValueStringPointer()
	}
	p, _, err := r.client.PublicKeyService.Create(&ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create profile public key", err))
		return
	}
	data.loadAPI(p)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *profilePublicKeyResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// do nothing
}

func (r *profilePublicKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *profilePublicKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	keyId, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("profile public key", err))
		return
	}
	p, httpResp, err := r.client.PublicKeyService.Get(keyId)
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get profile public key", err))
		return
	}
	data.loadAPI(p)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *profilePublicKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *profilePublicKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	keyId, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("profile public key", err))
		return
	}
	_, err = r.client.PublicKeyService.Delete(keyId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete profile public key", err))
		return
	}
}

func (r *profilePublicKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

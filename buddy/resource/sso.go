package resource

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ resource.Resource                = &ssoResource{}
	_ resource.ResourceWithConfigure   = &ssoResource{}
	_ resource.ResourceWithImportState = &ssoResource{}
)

func NewSsoResoruce() resource.Resource {
	return &ssoResource{}
}

type ssoResource struct {
	client *buddy.Client
}

type ssoResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Domain        types.String `tfsdk:"domain"`
	Type          types.String `tfsdk:"type"`
	SsoUrl        types.String `tfsdk:"sso_url"`
	Issuer        types.String `tfsdk:"issuer"`
	ClientId      types.String `tfsdk:"client_id"`
	ClientSecret  types.String `tfsdk:"client_secret"`
	Certificate   types.String `tfsdk:"certificate"`
	Signature     types.String `tfsdk:"signature"`
	Digest        types.String `tfsdk:"digest"`
	RequireForAll types.Bool   `tfsdk:"require_for_all"`
	HtmlUrl       types.String `tfsdk:"html_url"`
}

func (r *ssoResourceModel) loadAPI(domain string, sso *buddy.Sso) {
	r.ID = types.StringValue(domain)
	r.Domain = types.StringValue(domain)
	r.Type = types.StringValue(sso.Type)
	r.SsoUrl = types.StringValue(sso.SsoUrl)
	r.Issuer = types.StringValue(sso.Issuer)
	r.Certificate = types.StringValue(sso.Certificate)
	r.Signature = types.StringValue(sso.SignatureMethod)
	r.Digest = types.StringValue(sso.DigestMethod)
	r.RequireForAll = types.BoolValue(sso.RequireSsoForAllMembers)
	r.HtmlUrl = types.StringValue(sso.HtmlUrl)
}

func (r *ssoResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sso"
}

func (r *ssoResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage SSO in workspace\n\n" +
			"Workspace administrator rights are required\n\n" +
			"Token scopes required: `WORKSPACE`",
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
			"type": schema.StringAttribute{
				MarkdownDescription: "The SSO type. Allowed: `SAML`, `OIDC`. Default: `SAML`",
				Optional:            true,
				Computed:            true,
			},
			"sso_url": schema.StringAttribute{
				MarkdownDescription: "The identity provider single sign-on url",
				Optional:            true,
				Computed:            true,
			},
			"issuer": schema.StringAttribute{
				MarkdownDescription: "The identity provider issuer url",
				Required:            true,
			},
			"certificate": schema.StringAttribute{
				MarkdownDescription: "The identity provider certificate",
				Sensitive:           true,
				Optional:            true,
				Computed:            true,
			},
			"signature": schema.StringAttribute{
				MarkdownDescription: "The SAML signature algorithm. Allowed: `sha1`, `sha256`, `sha512`",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("sha1", "sha256", "sha512"),
				},
			},
			"digest": schema.StringAttribute{
				MarkdownDescription: "The SAML digest algorithm. Allowed: `sha1`, `sha256`, `sha512`",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("sha1", "sha256", "sha512"),
				},
			},
			"require_for_all": schema.BoolAttribute{
				MarkdownDescription: "Enable mandatory SAML SSO authentication for all workspace members",
				Optional:            true,
				Computed:            true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "The OIDC application's Client ID",
				Optional:            true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "The OIDC application's Client Secret",
				Sensitive:           true,
				Optional:            true,
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The Sso's URL",
				Computed:            true,
			},
		},
	}
}

func (r *ssoResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *ssoResource) updateSso(ctx context.Context, diag *diag.Diagnostics, plan *tfsdk.Plan, state *tfsdk.State) {
	var data *ssoResourceModel
	diag.Append(plan.Get(ctx, &data)...)
	if diag.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	_, _ = r.client.SsoService.Enable(domain)
	ops := buddy.SsoUpdateOps{}
	typ := buddy.SsoTypeSaml
	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		typ = data.Type.ValueString()
	}
	ops.Type = &typ
	ops.Issuer = data.Issuer.ValueStringPointer()
	if typ == buddy.SsoTypeSaml {
		ops.SsoUrl = data.SsoUrl.ValueStringPointer()
		ops.Certificate = data.Certificate.ValueStringPointer()
		ops.SignatureMethod = data.Signature.ValueStringPointer()
		ops.DigestMethod = data.Digest.ValueStringPointer()
	} else {
		ops.ClientId = data.ClientId.ValueStringPointer()
		ops.ClientSecret = data.ClientSecret.ValueStringPointer()
	}
	if !data.RequireForAll.IsNull() && !data.RequireForAll.IsUnknown() {
		ops.RequireSsoForAllMembers = data.RequireForAll.ValueBoolPointer()
	}
	sso, _, err := r.client.SsoService.Update(domain, &ops)
	if err != nil {
		diag.Append(util.NewDiagnosticApiError("update sso", err))
		return
	}
	data.loadAPI(domain, sso)
	diag.Append(state.Set(ctx, &data)...)
}

func (r *ssoResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.updateSso(ctx, &resp.Diagnostics, &req.Plan, &resp.State)
}

func (r *ssoResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.updateSso(ctx, &resp.Diagnostics, &req.Plan, &resp.State)
}

func (r *ssoResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ssoResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.ID.ValueString()
	sso, httpResp, err := r.client.SsoService.Get(domain)
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get sso", err))
		return
	}
	data.loadAPI(domain, sso)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ssoResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ssoResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, err := r.client.SsoService.Disable(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete sso", err))
	}
}

func (r *ssoResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

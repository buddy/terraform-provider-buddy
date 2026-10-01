package resource

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ resource.Resource                   = &domainResource{}
	_ resource.ResourceWithConfigure      = &domainResource{}
	_ resource.ResourceWithImportState    = &domainResource{}
	_ resource.ResourceWithValidateConfig = &domainResource{}
)

func NewDomainResource() resource.Resource {
	return &domainResource{}
}

type domainResource struct {
	client *buddy.Client
}

type domainResourceModel struct {
	ID              types.String `tfsdk:"id"`
	WorkspaceDomain types.String `tfsdk:"workspace_domain"`
	Domain          types.String `tfsdk:"domain"`
	Type            types.String `tfsdk:"type"`
	DomainId        types.String `tfsdk:"domain_id"`
	HtmlUrl         types.String `tfsdk:"html_url"`
	AutoRenew       types.Bool   `tfsdk:"auto_renew"`
	OnOwnerBehalf   types.Bool   `tfsdk:"on_owner_behalf"`
}

func (r *domainResourceModel) decomposeId() (string, string, error) {
	workspaceDomain, domainId, err := util.DecomposeDoubleId(r.ID.ValueString())
	if err != nil {
		return "", "", err
	}
	return workspaceDomain, domainId, nil
}

func (r *domainResourceModel) loadAPI(workspaceDomain string, domain *buddy.Domain) {
	r.ID = types.StringValue(util.ComposeDoubleId(workspaceDomain, domain.Id))
	r.WorkspaceDomain = types.StringValue(workspaceDomain)
	r.Domain = types.StringValue(domain.Name)
	r.DomainId = types.StringValue(domain.Id)
	r.Type = types.StringValue(domain.Type)
	r.HtmlUrl = types.StringValue(domain.HtmlUrl)
	r.AutoRenew = types.BoolValue(domain.AutoRenew)
	// on_owner_behalf is not returned by the API, it stays as configured
}

func (r *domainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

func (r *domainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a domain\n\n" +
			"Invite-only token is required. Contact support@buddy.works for more details\n\n" +
			"Destroying the resource deletes the domain with all its records from the workspace\n\n" +
			"Token scopes required: `DOMAIN_READ`, `DOMAIN_MANAGE`",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Terraform resource identifier for this item",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace_domain": schema.StringAttribute{
				MarkdownDescription: "The workspace's URL handle",
				Required:            true,
				Validators:          util.StringValidatorsDomain(),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The domain's name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The domain's type. Allowed values: POINTED (default), PRIVATE, REGISTERED, CLAIMED. PRIVATE requires a plan with private zones",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(buddy.DomainTypePointed),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						buddy.DomainTypePointed,
						buddy.DomainTypePrivate,
						buddy.DomainTypeRegistered,
						buddy.DomainTypeClaimed,
					),
				},
			},
			"auto_renew": schema.BoolAttribute{
				MarkdownDescription: "Renew the domain automatically, allowed only for REGISTERED type",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplace(),
				},
			},
			"on_owner_behalf": schema.BoolAttribute{
				MarkdownDescription: "Register or claim the domain on the workspace owner's behalf, allowed only for REGISTERED and CLAIMED types. Not returned by the API, used only on create",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"domain_id": schema.StringAttribute{
				MarkdownDescription: "The domain's id",
				Computed:            true,
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The domain's URL",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *domainResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data *domainResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() || data.Type.IsUnknown() {
		return
	}
	typ := buddy.DomainTypePointed
	if !data.Type.IsNull() {
		typ = data.Type.ValueString()
	}
	if !data.AutoRenew.IsNull() && typ != buddy.DomainTypeRegistered {
		resp.Diagnostics.AddAttributeError(path.Root("auto_renew"), "Invalid attribute combination", "auto_renew is allowed only for type REGISTERED")
	}
	if !data.OnOwnerBehalf.IsNull() && typ != buddy.DomainTypeRegistered && typ != buddy.DomainTypeClaimed {
		resp.Diagnostics.AddAttributeError(path.Root("on_owner_behalf"), "Invalid attribute combination", "on_owner_behalf is allowed only for types REGISTERED and CLAIMED")
	}
}

func (r *domainResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *domainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *domainResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	workspaceDomain := data.WorkspaceDomain.ValueString()
	domain := data.Domain.ValueString()
	typ := buddy.DomainTypePointed
	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		typ = data.Type.ValueString()
	}
	ops := buddy.DomainCreateOps{
		Name: &domain,
		Type: &typ,
	}
	if !data.AutoRenew.IsNull() && !data.AutoRenew.IsUnknown() {
		ops.AutoRenew = data.AutoRenew.ValueBoolPointer()
	}
	if !data.OnOwnerBehalf.IsNull() && !data.OnOwnerBehalf.IsUnknown() {
		ops.OnOwnerBehalf = data.OnOwnerBehalf.ValueBoolPointer()
	}
	d, _, err := r.client.DomainService.Create(workspaceDomain, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create domain", err))
		return
	}
	data.loadAPI(workspaceDomain, d)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *domainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *domainResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	workspaceDomain, domainId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("domain", err))
		return
	}
	d, httpResp, err := r.client.DomainService.Get(workspaceDomain, domainId)
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get domain", err))
		return
	}
	data.loadAPI(workspaceDomain, d)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *domainResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// do nothing
}

func (r *domainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *domainResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	workspaceDomain, domainId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("domain", err))
		return
	}
	httpResp, err := r.client.DomainService.Delete(workspaceDomain, domainId)
	if err != nil && !util.IsResourceNotFound(httpResp, err) {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete domain", err))
	}
}

func (r *domainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

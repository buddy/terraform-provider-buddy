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
)

var (
	_ resource.Resource                = &domainResource{}
	_ resource.ResourceWithConfigure   = &domainResource{}
	_ resource.ResourceWithImportState = &domainResource{}
)

func NewDomaindResource() resource.Resource {
	return &domainResource{}
}

type domainResource struct {
	client *buddy.Client
}

type domainResourceModel struct {
	ID              types.String `tfsdk:"id"`
	WorkspaceDomain types.String `tfsdk:"workspace_domain"`
	Domain          types.String `tfsdk:"domain"`
}

func (r *domainResourceModel) loadAPI(workspaceDomain string, domain string) {
	r.ID = types.StringValue(util.ComposeDoubleId(workspaceDomain, domain))
	r.WorkspaceDomain = types.StringValue(workspaceDomain)
	r.Domain = types.StringValue(domain)
}

func (r *domainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

func (r *domainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a domain\n\n" +
			"Invite-only token is required. Contact support@buddy.works for more details\n\n" +
			"Token scope required: `ZONE_MANAGE`",
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
		},
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
	ops := buddy.DomainCreateOps{
		Name: &domain,
	}
	d, _, err := r.client.DomainService.Create(workspaceDomain, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create domain", err))
		return
	}
	data.loadAPI(workspaceDomain, d.Name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *domainResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	// do nothing
}

func (r *domainResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// do nothing
}

func (r *domainResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// do nothing
}

func (r *domainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

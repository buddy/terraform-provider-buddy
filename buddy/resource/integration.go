package resource

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ resource.Resource                = &integrationResource{}
	_ resource.ResourceWithConfigure   = &integrationResource{}
	_ resource.ResourceWithImportState = &integrationResource{}
)

func NewIntegrationResource() resource.Resource {
	return &integrationResource{}
}

type integrationResource struct {
	client *buddy.Client
}

type integrationResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Domain          types.String `tfsdk:"domain"`
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	Scope           types.String `tfsdk:"scope"`
	ProjectName     types.String `tfsdk:"project_name"`
	GroupId         types.Int64  `tfsdk:"group_id"`
	Username        types.String `tfsdk:"username"`
	Shop            types.String `tfsdk:"shop"`
	Token           types.String `tfsdk:"token"`
	PartnerToken    types.String `tfsdk:"partner_token"`
	AccessKey       types.String `tfsdk:"access_key"`
	SecretKey       types.String `tfsdk:"secret_key"`
	Audience        types.String `tfsdk:"audience"`
	AuthType        types.String `tfsdk:"auth_type"`
	AppId           types.String `tfsdk:"app_id"`
	TenantId        types.String `tfsdk:"tenant_id"`
	Password        types.String `tfsdk:"password"`
	ApiKey          types.String `tfsdk:"api_key"`
	Email           types.String `tfsdk:"email"`
	RoleAssumptions types.List   `tfsdk:"role_assumption"`
	GoogleConfig    types.String `tfsdk:"google_config"`
	GoogleProject   types.String `tfsdk:"google_project"`
	IntegrationId   types.String `tfsdk:"integration_id"`
	HtmlUrl         types.String `tfsdk:"html_url"`
}

func (r *integrationResourceModel) decomposeId() (string, string, error) {
	domain, iid, err := util.DecomposeDoubleId(r.ID.ValueString())
	if err != nil {
		return "", "", err
	}
	return domain, iid, nil
}

func (r *integrationResourceModel) loadAPI(domain string, integration *buddy.Integration) {
	r.ID = types.StringValue(util.ComposeDoubleId(domain, integration.HashId))
	r.Domain = types.StringValue(domain)
	r.Name = types.StringValue(integration.Name)
	r.Type = types.StringValue(integration.Type)
	r.AuthType = types.StringValue(integration.AuthType)
	r.Scope = types.StringValue(integration.Scope)
	r.ProjectName = types.StringValue(integration.ProjectName)
	r.GroupId = types.Int64Value(int64(integration.GroupId))
	r.HtmlUrl = types.StringValue(integration.HtmlUrl)
	r.IntegrationId = types.StringValue(integration.HashId)
	// rest of the attributes are not returned by api
}

func (r *integrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration"
}

func (r *integrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage an integration\n\n" +
			"Token scopes required: `INTEGRATION_ADD`, `INTEGRATION_MANAGE`, `INTEGRATION_INFO`",
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
				MarkdownDescription: "The integration's name",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The integration's type. Allowed: `DIGITAL_OCEAN`, `AMAZON`, `SHOPIFY`, `PUSHOVER`, " +
					"`RACKSPACE`, `CLOUDFLARE`, `NEW_RELIC`, `SENTRY`, `ROLLBAR`, `DATADOG`, `DO_SPACES`, `HONEYBADGER`, " +
					"`VULTR`, `SENTRY_ENTERPRISE`, `LOGGLY`, `FIREBASE`, `UPCLOUD`, `GHOST_INSPECTOR`, `AZURE_CLOUD`, " +
					"`DOCKER_HUB`, `GOOGLE_SERVICE_ACCOUNT`, `GIT_HUB`, `GIT_LAB`, `STACK_HAWK`",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{stringvalidator.OneOf(
					buddy.IntegrationTypeDigitalOcean,
					buddy.IntegrationTypeAmazon,
					buddy.IntegrationTypeShopify,
					buddy.IntegrationTypePushover,
					buddy.IntegrationTypeRackspace,
					buddy.IntegrationTypeCloudflare,
					buddy.IntegrationTypeNewRelic,
					buddy.IntegrationTypeSentry,
					buddy.IntegrationTypeRollbar,
					buddy.IntegrationTypeDatadog,
					buddy.IntegrationTypeDigitalOceanSpaces,
					buddy.IntegrationTypeHoneybadger,
					buddy.IntegrationTypeVultr,
					buddy.IntegrationTypeSentryEnterprise,
					buddy.IntegrationTypeLoggly,
					buddy.IntegrationTypeFirebase,
					buddy.IntegrationTypeUpcloud,
					buddy.IntegrationTypeGhostInspector,
					buddy.IntegrationTypeAzureCloud,
					buddy.IntegrationTypeDockerHub,
					buddy.IntegrationTypeGoogleServiceAccount,
					buddy.IntegrationTypeGitHub,
					buddy.IntegrationTypeGitLab,
					buddy.IntegrationTypeStackHawk,
				)},
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: "The integration's scope. Allowed:\n\n" +
					"`PRIVATE` - only creator of the integration can use it\n\n" +
					"`WORKSPACE` - all workspace members can use the integration\n\n" +
					"`ADMIN` - only workspace administrators can use the integration\n\n" +
					"`GROUP` - only group members can use the integration\n\n" +
					"`PROJECT` - only project members can use the integration\n\n" +
					"`PRIVATE_IN_PROJECT` - only creator of the integration in specified project can use it\n\n" +
					"`ADMIN_IN_PROJECT` - only workspace administrators in specified project can use the integration\n\n" +
					"`GROUP_IN_PROJECT` - only group members in specified project can use the integration",
				Required: true,
				Validators: []validator.String{stringvalidator.OneOf(
					buddy.IntegrationScopePrivate,
					buddy.IntegrationScopeWorkspace,
					buddy.IntegrationScopeAdmin,
					buddy.IntegrationScopeGroup,
					buddy.IntegrationScopeProject,
					buddy.IntegrationScopePrivateInProject,
					buddy.IntegrationScopeAdminInProject,
					buddy.IntegrationScopeGroupInProject,
				)},
			},
			"project_name": schema.StringAttribute{
				MarkdownDescription: "The project's name. Provide along with scopes: `PROJECT`, `ADMIN_IN_PROJECT`, `GROUP_IN_PROJECT`, `PRIVATE_IN_PROJECT`",
				Optional:            true,
				Computed:            true,
			},
			"group_id": schema.Int64Attribute{
				MarkdownDescription: "The group's ID. Provide along with scopes: `GROUP`, `GROUP_IN_PROJECT`",
				Optional:            true,
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The integration's username. Provide for: `UPCLOUD`, `RACKSPACE`, `DOCKER_HUB`",
				Optional:            true,
			},
			"shop": schema.StringAttribute{
				MarkdownDescription: "The integration's shop. Provide for: `SHOPIFY`",
				Optional:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "The integration's token. Provide for: `DIGITAL_OCEAN`, `SHOPIFY`, `RACKSPACE`, `CLOUDFLARE`, " +
					"`NEW_RELIC`, `SENTRY`, `ROLLBAR`, `DATADOG`, `HONEYBADGER`, `VULTR`, `SENTRY_ENTERPRISE`, " +
					"`LOGGLY`, `FIREBASE`, `GHOST_INSPECTOR`, `PUSHOVER`, `GIT_LAB`, `GIT_HUB`",
				Optional:  true,
				Sensitive: true,
			},
			"partner_token": schema.StringAttribute{
				MarkdownDescription: "The integration's partner token. Provide for: `SHOPIFY`",
				Optional:            true,
				Sensitive:           true,
			},
			"access_key": schema.StringAttribute{
				MarkdownDescription: "The integration's access key. Provide for: `DO_SPACES`, `AMAZON`, `PUSHOVER`",
				Optional:            true,
				Sensitive:           true,
			},
			"secret_key": schema.StringAttribute{
				MarkdownDescription: "The integration's secret key. Provide for: `DO_SPACES`, `AMAZON`",
				Optional:            true,
				Sensitive:           true,
			},
			"audience": schema.StringAttribute{
				MarkdownDescription: "The integration's audience. Provide for OIDC with: `AMAZON`, `AZURE_CLOUD`, `GOOGLE_SERVICE_ACCOUNT`",
				Optional:            true,
			},
			"google_config": schema.StringAttribute{
				MarkdownDescription: "The integration's google config. Provide for `GOOGLE_SERVICE_ACCOUNT` OIDC",
				Optional:            true,
			},
			"google_project": schema.StringAttribute{
				MarkdownDescription: "The integration's google project. Provide for `GOOGLE_SERVICE_ACCOUNT` OIDC",
				Optional:            true,
			},
			"auth_type": schema.StringAttribute{
				MarkdownDescription: "The integration's auth type. Provide for: `AMAZON`, `AZURE_CLOUD`, `GOOGLE_SERVICE_ACCOUNT`. Allowed: `DEFAULT, TRUSTED, OIDC`",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{stringvalidator.OneOf(
					buddy.IntegrationAuthTypeDefault,
					buddy.IntegrationAuthTypeTrusted,
					buddy.IntegrationAuthTypeOidc,
				)},
			},
			"app_id": schema.StringAttribute{
				MarkdownDescription: "The integration's application's ID. Provide for: `AZURE_CLOUD`",
				Optional:            true,
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The integration's tenant's ID. Provide for: `AZURE_CLOUD`",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The integration's password. Provide for: `AZURE_CLOUD`, `UPCLOUD`, `DOCKER_HUB`",
				Optional:            true,
				Sensitive:           true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The integration's API key. Provide for: `CLOUDFLARE`, `GOOGLE_SERVICE_ACCOUNT`, `STACK_HAWK`",
				Optional:            true,
				Sensitive:           true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The integration's email. Provide for: `CLOUDFLARE`",
				Optional:            true,
				Sensitive:           true,
			},
			"integration_id": schema.StringAttribute{
				MarkdownDescription: "The integration's ID",
				Computed:            true,
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The integration's URL",
				Computed:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"role_assumption": schema.ListNestedBlock{
				MarkdownDescription: "The integration's AWS role to assume. Provide for: `AMAZON`",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"arn": schema.StringAttribute{
							MarkdownDescription: "The integration's AWS role ARN to assume",
							Required:            true,
						},
						"external_id": schema.StringAttribute{
							MarkdownDescription: "The integration's AWS external ID to send when assuming AWS role",
							Optional:            true,
						},
						"duration": schema.Int64Attribute{
							MarkdownDescription: "The integration's AWS session duration in seconds",
							Optional:            true,
						},
					},
				},
			},
		},
	}
}

func (r *integrationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *integrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *integrationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	ops := buddy.IntegrationOps{
		Name:  data.Name.ValueStringPointer(),
		Type:  data.Type.ValueStringPointer(),
		Scope: data.Scope.ValueStringPointer(),
	}
	if !data.ProjectName.IsNull() && !data.ProjectName.IsUnknown() {
		ops.ProjectName = data.ProjectName.ValueStringPointer()
	}
	if !data.GroupId.IsNull() && !data.GroupId.IsUnknown() {
		ops.GroupId = util.PointerInt(data.GroupId.ValueInt64())
	}
	if !data.Username.IsNull() && !data.Username.IsUnknown() {
		ops.Username = data.Username.ValueStringPointer()
	}
	if !data.Shop.IsNull() && !data.Shop.IsUnknown() {
		ops.Shop = data.Shop.ValueStringPointer()
	}
	if !data.Token.IsNull() && !data.Token.IsUnknown() {
		ops.Token = data.Token.ValueStringPointer()
	}
	var authType string
	if !data.AuthType.IsNull() && !data.AuthType.IsUnknown() {
		authType = data.AuthType.ValueString()
	}
	if !data.PartnerToken.IsNull() && !data.PartnerToken.IsUnknown() {
		ops.PartnerToken = data.PartnerToken.ValueStringPointer()
		authType = buddy.IntegrationAuthTypeTokenAppExtension
	}
	if !data.AccessKey.IsNull() && !data.AccessKey.IsUnknown() {
		ops.AccessKey = data.AccessKey.ValueStringPointer()
	}
	if !data.SecretKey.IsNull() && !data.SecretKey.IsUnknown() {
		ops.SecretKey = data.SecretKey.ValueStringPointer()
	}
	if !data.Audience.IsNull() && !data.Audience.IsUnknown() {
		ops.Audience = data.Audience.ValueStringPointer()
	}
	if !data.AppId.IsNull() && !data.AppId.IsUnknown() {
		ops.AppId = data.AppId.ValueStringPointer()
	}
	if !data.TenantId.IsNull() && !data.TenantId.IsUnknown() {
		ops.TenantId = data.TenantId.ValueStringPointer()
	}
	if !data.Password.IsNull() && !data.Password.IsUnknown() {
		ops.Password = data.Password.ValueStringPointer()
	}
	if !data.ApiKey.IsNull() && !data.ApiKey.IsUnknown() {
		ops.ApiKey = data.ApiKey.ValueStringPointer()
	}
	if !data.Email.IsNull() && !data.Email.IsUnknown() {
		ops.Email = data.Email.ValueStringPointer()
	}
	if !data.GoogleProject.IsNull() && !data.GoogleProject.IsUnknown() {
		ops.GoogleProject = data.GoogleProject.ValueStringPointer()
	}
	if !data.GoogleConfig.IsNull() && !data.GoogleConfig.IsUnknown() {
		ops.Config = data.GoogleConfig.ValueStringPointer()
	}
	if !data.RoleAssumptions.IsNull() && !data.RoleAssumptions.IsUnknown() {
		roles, diags := util.RoleAssumptionsModelToAPI(ctx, &data.RoleAssumptions)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.RoleAssumptions = roles
	}
	if authType != "" {
		ops.AuthType = &authType
	}
	integration, _, err := r.client.IntegrationService.Create(domain, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create integration", err))
		return
	}
	data.loadAPI(domain, integration)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *integrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *integrationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, integrationId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("integration", err))
		return
	}
	integration, httpResp, err := r.client.IntegrationService.Get(domain, integrationId)
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get integration", err))
		return
	}
	data.loadAPI(domain, integration)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *integrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *integrationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, integrationId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("integration", err))
		return
	}
	ops := buddy.IntegrationOps{
		Name:  data.Name.ValueStringPointer(),
		Type:  data.Type.ValueStringPointer(),
		Scope: data.Scope.ValueStringPointer(),
	}
	if !data.ProjectName.IsNull() && !data.ProjectName.IsUnknown() {
		ops.ProjectName = data.ProjectName.ValueStringPointer()
	}
	if !data.GroupId.IsNull() && !data.GroupId.IsUnknown() {
		ops.GroupId = util.PointerInt(data.GroupId.ValueInt64())
	}
	if !data.Username.IsNull() && !data.Username.IsUnknown() {
		ops.Username = data.Username.ValueStringPointer()
	}
	if !data.Shop.IsNull() && !data.Shop.IsUnknown() {
		ops.Shop = data.Shop.ValueStringPointer()
	}
	if !data.Token.IsNull() && !data.Token.IsUnknown() {
		ops.Token = data.Token.ValueStringPointer()
	}
	var authType string
	if !data.AuthType.IsNull() && !data.AuthType.IsUnknown() {
		authType = data.AuthType.ValueString()
	}
	if !data.PartnerToken.IsNull() && !data.PartnerToken.IsUnknown() {
		ops.PartnerToken = data.PartnerToken.ValueStringPointer()
		authType = buddy.IntegrationAuthTypeTokenAppExtension
	}
	if !data.AccessKey.IsNull() && !data.AccessKey.IsUnknown() {
		ops.AccessKey = data.AccessKey.ValueStringPointer()
	}
	if !data.SecretKey.IsNull() && !data.SecretKey.IsUnknown() {
		ops.SecretKey = data.SecretKey.ValueStringPointer()
	}
	if !data.Audience.IsNull() && !data.Audience.IsUnknown() {
		ops.Audience = data.Audience.ValueStringPointer()
	}
	if !data.AppId.IsNull() && !data.AppId.IsUnknown() {
		ops.AppId = data.AppId.ValueStringPointer()
	}
	if !data.TenantId.IsNull() && !data.TenantId.IsUnknown() {
		ops.TenantId = data.TenantId.ValueStringPointer()
	}
	if !data.Password.IsNull() && !data.Password.IsUnknown() {
		ops.Password = data.Password.ValueStringPointer()
	}
	if !data.ApiKey.IsNull() && !data.ApiKey.IsUnknown() {
		ops.ApiKey = data.ApiKey.ValueStringPointer()
	}
	if !data.Email.IsNull() && !data.Email.IsUnknown() {
		ops.Email = data.Email.ValueStringPointer()
	}
	if !data.GoogleProject.IsNull() && !data.GoogleProject.IsUnknown() {
		ops.GoogleProject = data.GoogleProject.ValueStringPointer()
	}
	if !data.GoogleConfig.IsNull() && !data.GoogleConfig.IsUnknown() {
		ops.Config = data.GoogleConfig.ValueStringPointer()
	}
	if !data.RoleAssumptions.IsNull() && !data.RoleAssumptions.IsUnknown() {
		roles, diags := util.RoleAssumptionsModelToAPI(ctx, &data.RoleAssumptions)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.RoleAssumptions = roles
	}
	if authType != "" {
		ops.AuthType = &authType
	}
	integration, _, err := r.client.IntegrationService.Update(domain, integrationId, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update integration", err))
		return
	}
	data.loadAPI(domain, integration)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *integrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *integrationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, integrationId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("integration", err))
		return
	}
	_, err = r.client.IntegrationService.Delete(domain, integrationId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete integration", err))
	}
}

func (r *integrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

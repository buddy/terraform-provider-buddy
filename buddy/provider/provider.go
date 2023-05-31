package provider

import (
	buddyresource "buddy-terraform/buddy/resource"
	buddysource "buddy-terraform/buddy/source"
	"context"
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"os"
)

var _ provider.Provider = &BuddyProvider{}

type BuddyProvider struct {
	version string
}

type BuddyProviderModel struct {
	Token    types.String `tfsdk:"token"`
	BaseUrl  types.String `tfsdk:"base_url"`
	Insecure types.Bool   `tfsdk:"insecure"`
}

func (p *BuddyProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "buddy"
	resp.Version = p.version
}

func (p *BuddyProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				MarkdownDescription: "The OAuth2 token or Personal Access Token. Can be specified with the `BUDDY_TOKEN` environment variable.",
				Sensitive:           true,
				Optional:            true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "The Buddy API base url. You may need to set this to your Buddy On-Premises API endpoint. Can be specified with the `BUDDY_BASE_URL` environment variable. Default: `https://api.buddy.works`",
				Optional:            true,
			},
			"insecure": schema.BoolAttribute{
				MarkdownDescription: "Disable SSL verification of API calls. You may need to set this to `true` if you are using Buddy On-Premises without signed certificate. Can be specified with the `BUDDY_INSECURE` environmental variable",
				Optional:            true,
			},
		},
	}
}

func (p *BuddyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config BuddyProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Buddy Token",
			"The provider cannot create the Buddy API client as there is unknown configuration value for the Buddy Token",
		)
	}
	if config.BaseUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("base_url"),
			"Unknown Buddy Base URL for the API endpoint",
			"The provider cannot create the Buddy API client as there is unknown configuration value for the Buddy Base URL",
		)
	}
	if config.BaseUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("insecure"),
			"Unknown Buddy insecure value for the API endpoint",
			"The provider cannot create the Buddy API client as there is unknown configuration value for the Buddy insecure attribute",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	token := os.Getenv("BUDDY_TOKEN")
	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}
	baseUrl := os.Getenv("BUDDY_BASE_URL")
	if !config.BaseUrl.IsNull() {
		baseUrl = config.BaseUrl.ValueString()
	}
	insecure := os.Getenv("BUDDY_INSECURE") == "true"
	if !config.Insecure.IsNull() {
		insecure = config.Insecure.ValueBool()
	}

	client, err := buddy.NewClient(token, baseUrl, insecure)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Buddy Client from provider configuration", fmt.Sprintf("The provider failed to create a new Buddy Client from the giver configuration: %s", err.Error()))
		return
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *BuddyProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		buddyresource.NewWorkspaceResource,
		buddyresource.NewGroupResource,
		buddyresource.NewGroupMemberResource,
		buddyresource.NewPermissionResource,
		buddyresource.NewMemberResource,
		buddyresource.NewProfileResource,
		buddyresource.NewProfileEmailResource,
		buddyresource.NewProfilePublicKeyResource,
		buddyresource.NewIntegrationResource,
		buddyresource.NewProjectResource,
		buddyresource.NewProjectGroupResource,
		buddyresource.NewProjectMemberResource,
		buddyresource.NewSsoResoruce,
		buddyresource.NewVariableResource,
		buddyresource.NewVariableSshResource,
		buddyresource.NewWebhookResource,
		buddyresource.NewPipelineResource,
	}
}

func (p *BuddyProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		buddysource.NewGroupSource,
		buddysource.NewGroupMembersSource,
		buddysource.NewGroupsSource,
		buddysource.NewIntegrationSource,
		buddysource.NewIntegrationsSource,
		buddysource.NewMemberSource,
		buddysource.NewMembersSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &BuddyProvider{
			version: version,
		}
	}
}

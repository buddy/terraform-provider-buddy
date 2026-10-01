package provider

import (
	"context"
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"os"
	"strconv"
	"strings"
	buddyresource "terraform-provider-buddy/buddy/resource"
	buddysource "terraform-provider-buddy/buddy/source"
	"time"
)

var _ provider.Provider = &BuddyProvider{}

type BuddyProvider struct {
	version string
}

type BuddyProviderModel struct {
	Token    types.String `tfsdk:"token"`
	BaseUrl  types.String `tfsdk:"base_url"`
	Region   types.String `tfsdk:"region"`
	Insecure types.Bool   `tfsdk:"insecure"`
	Timeout  types.Int64  `tfsdk:"timeout"`
}

var regionBaseUrls = map[string]string{
	"us": "https://api.buddy.works",
	"as": "https://api.asia.buddy.works",
	"eu": "https://api.eu.buddy.works",
}

func regionBaseUrl(region string) (string, bool) {
	baseUrl, ok := regionBaseUrls[strings.ToLower(strings.TrimSpace(region))]
	return baseUrl, ok
}

// resolveBaseUrl picks the API base url out of the base_url & region settings.
// An explicitly given base url always wins, region is only consulted when there is none.
// An empty result means neither was set and the client falls back to its own default.
func resolveBaseUrl(configBaseUrl types.String, configRegion types.String) (string, error) {
	baseUrl := os.Getenv("BUDDY_BASE_URL")
	if !configBaseUrl.IsNull() {
		baseUrl = configBaseUrl.ValueString()
	}
	if baseUrl != "" {
		return baseUrl, nil
	}
	region := os.Getenv("BUDDY_REGION")
	if !configRegion.IsNull() {
		region = configRegion.ValueString()
	}
	if region == "" {
		return "", nil
	}
	regionUrl, ok := regionBaseUrl(region)
	if !ok {
		return "", fmt.Errorf("`%s` is not a known Buddy region, must be one of: us, as, eu", region)
	}
	return regionUrl, nil
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
				MarkdownDescription: "The Buddy API base url. You may need to set this to your Buddy On-Premises API endpoint. Can be specified with the `BUDDY_BASE_URL` environment variable. Takes precedence over `region`. Default: `https://api.buddy.works`",
				Optional:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The Buddy cloud region to connect to. Sets the API base url to the endpoint of the given region: `us` - `https://api.buddy.works`, `as` - `https://api.asia.buddy.works`, `eu` - `https://api.eu.buddy.works`. Case insensitive. Ignored when `base_url` is set. Can be specified with the `BUDDY_REGION` environment variable",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("us", "as", "eu"),
				},
			},
			"insecure": schema.BoolAttribute{
				MarkdownDescription: "Disable SSL verification of API calls. You may need to set this to `true` if you are using Buddy On-Premises without signed certificate. Can be specified with the `BUDDY_INSECURE` environmental variable",
				Optional:            true,
			},
			"timeout": schema.Int64Attribute{
				MarkdownDescription: "The Buddy API client timeout in seconds. Can be specified with the `BUDDY_TIMEOUT` environmental variable. Default: 30s",
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
	if config.Region.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("region"),
			"Unknown Buddy region for the API endpoint",
			"The provider cannot create the Buddy API client as there is unknown configuration value for the Buddy region",
		)
	}
	if config.BaseUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("insecure"),
			"Unknown Buddy insecure value for the API endpoint",
			"The provider cannot create the Buddy API client as there is unknown configuration value for the Buddy insecure attribute",
		)
	}
	if config.Timeout.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("timeout"),
			"Unknown Buddy timeout value for the API endpoint",
			"The provider cannot create the Buddy API client as there is unknown configuration value for the Buddy timeout attribute",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	token := os.Getenv("BUDDY_TOKEN")
	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}
	baseUrl, err := resolveBaseUrl(config.BaseUrl, config.Region)
	if err != nil {
		resp.Diagnostics.AddError("Wrong value for the Buddy region", fmt.Sprintf("The provider cannot create the Buddy API client as %s", err.Error()))
		return
	}
	insecure := os.Getenv("BUDDY_INSECURE") == "true"
	if !config.Insecure.IsNull() {
		insecure = config.Insecure.ValueBool()
	}
	timeout := 30
	t := os.Getenv("BUDDY_TIMEOUT")
	if t != "" {
		var err error
		timeout, err = strconv.Atoi(t)
		if err != nil {
			resp.Diagnostics.AddError("Wrong value in BUDDY_TIMEOUT env variable", "The provider cannot create the Buddy API client as there is wrong value for the BUDDY_TIMEOUT env variable")
			return
		}
	}
	if !config.Timeout.IsNull() {
		timeout = int(config.Timeout.ValueInt64())
	}

	client, err := buddy.NewClientWithTimeout(token, baseUrl, insecure, time.Duration(timeout)*time.Second)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Buddy Client from provider configuration", fmt.Sprintf("The provider failed to create a new Buddy Client from the giver configuration: %s", err.Error()))
		return
	}
	client.SetLogger(newApiLogger(ctx))
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
		buddyresource.NewDomainRecordResource,
		buddyresource.NewDomainResource,
		buddyresource.NewProjectResource,
		buddyresource.NewProjectGroupResource,
		buddyresource.NewProjectMemberResource,
		buddyresource.NewSsoResource,
		buddyresource.NewVariableResource,
		buddyresource.NewVariableSshResource,
		buddyresource.NewWebhookResource,
		buddyresource.NewPipelineResource,
		buddyresource.NewSandboxResource,
		buddyresource.NewSandboxStatusResource,
		buddyresource.NewEnvironmentResource,
		buddyresource.NewTargetResource,
	}
}

func (p *BuddyProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		buddysource.NewDomainsSource,
		buddysource.NewEnvironmentSource,
		buddysource.NewGroupSource,
		buddysource.NewGroupMembersSource,
		buddysource.NewGroupsSource,
		buddysource.NewIntegrationSource,
		buddysource.NewIntegrationsSource,
		buddysource.NewMemberSource,
		buddysource.NewMembersSource,
		buddysource.NewPermissionSource,
		buddysource.NewPermissionsSource,
		buddysource.NewProfileSource,
		buddysource.NewProjectSource,
		buddysource.NewProjectGroupSource,
		buddysource.NewProjectGroupsSource,
		buddysource.NewProjectMemberSource,
		buddysource.NewProjectMembersSource,
		buddysource.NewProjectsSource,
		buddysource.NewVariableSource,
		buddysource.NewVariableSshKeySource,
		buddysource.NewVariablesSource,
		buddysource.NewVariablesSshKeysSource,
		buddysource.NewWebhookSource,
		buddysource.NewWebhooksSource,
		buddysource.NewWorkspaceSource,
		buddysource.NewWorkspacesSource,
		buddysource.NewPipelineSource,
		buddysource.NewPipelinesSource,
		buddysource.NewSandboxesSource,
		buddysource.NewSandboxSource,
		buddysource.NewEnvironmentsSource,
		buddysource.NewTargetSource,
		buddysource.NewTargetsSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &BuddyProvider{
			version: version,
		}
	}
}

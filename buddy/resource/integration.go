package resource

import (
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"net/http"
)

func Integration() *schema.Resource {
	return &schema.Resource{
		Description: "Create and manage a workspace integration\n\n" +
			"Token scopes required: `INTEGRATION_ADD`, `INTEGRATION_MANAGE`, `INTEGRATION_INFO`",
		CreateContext: createContextIntegration,
		ReadContext:   readContextIntegration,
		UpdateContext: updateContextIntegration,
		DeleteContext: deleteContextIntegration,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The Terraform resource identifier for this item",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"domain": {
				Description:  "The workspace's URL handle",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: util.ValidateDomain,
			},
			"name": {
				Description: "The integration's name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"type": {
				Description: "The integration's type. Allowed: `DIGITAL_OCEAN`, `AMAZON`, `SHOPIFY`, `PUSHOVER`, " +
					"`RACKSPACE`, `CLOUDFLARE`, `NEW_RELIC`, `SENTRY`, `ROLLBAR`, `DATADOG`, `DO_SPACES`, `HONEYBADGER`, " +
					"`VULTR`, `SENTRY_ENTERPRISE`, `LOGGLY`, `FIREBASE`, `UPCLOUD`, `GHOST_INSPECTOR`, `AZURE_CLOUD`, " +
					"`DOCKER_HUB`, `GOOGLE_SERVICE_ACCOUNT`",
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					api.IntegrationTypeDigitalOcean,
					api.IntegrationTypeAmazon,
					api.IntegrationTypeShopify,
					api.IntegrationTypePushover,
					api.IntegrationTypeRackspace,
					api.IntegrationTypeCloudflare,
					api.IntegrationTypeNewRelic,
					api.IntegrationTypeSentry,
					api.IntegrationTypeRollbar,
					api.IntegrationTypeDatadog,
					api.IntegrationTypeDigitalOceanSpaces,
					api.IntegrationTypeHoneybadger,
					api.IntegrationTypeVultr,
					api.IntegrationTypeSentryEnterprise,
					api.IntegrationTypeLoggly,
					api.IntegrationTypeFirebase,
					api.IntegrationTypeUpcloud,
					api.IntegrationTypeGhostInspector,
					api.IntegrationTypeAzureCloud,
					api.IntegrationTypeDockerHub,
					api.IntegrationTypeGoogleServiceAccount,
				}, false),
			},
			"scope": {
				Description: "The integration's scope. Allowed:\n\n" +
					"`PRIVATE` - only creator of the integration can use it\n\n" +
					"`WORKSPACE` - all workspace members can use the integration\n\n" +
					"`ADMIN` - only workspace administrators can use the integration\n\n" +
					"`GROUP` - only group members can use the integration\n\n" +
					"`PROJECT` - only project members can use the integration\n\n" +
					"`ADMIN_IN_PROJECT` - only workspace administrators in specified project can use the integration\n\n" +
					"`GROUP_IN_PROJECT` - only group members in specified project can use the integration",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					api.IntegrationScopePrivate,
					api.IntegrationScopeWorkspace,
					api.IntegrationScopeAdmin,
					api.IntegrationScopeGroup,
					api.IntegrationScopeProject,
					api.IntegrationScopeAdminInProject,
					api.IntegrationScopeGroupInProject,
				}, false),
			},
			"project_name": {
				Description: "The project's name. Provide along with scopes: `PROJECT`, `ADMIN_IN_PROJECT`, `GROUP_IN_PROJECT`",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"group_id": {
				Description: "The group's ID. Provide along with scopes: `GROUP`, `GROUP_IN_PROJECT`",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"username": {
				Description: "The integration's username. Provide for integrations: `UPCLOUD`, `RACKSPACE`, `DOCKER_HUB`",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"shop": {
				Description: "The integration's shop. Provide for integration: `SHOPIFY`",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"token": {
				Description: "The integration's token. Provide for integration: `DIGITAL_OCEAN`, `SHOPIFY`, `RACKSPACE`, `CLOUDFLARE`, " +
					"`NEW_RELIC`, `SENTRY`, `ROLLBAR`, `DATADOG`, `HONEYBADGER`, `VULTR`, `SENTRY_ENTERPRISE`, " +
					"`LOGGLY`, `FIREBASE`, `GHOST_INSPECTOR`, `PUSHOVER`",
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"access_key": {
				Description: "The integration's access key. Provide for integration: `DO_SPACES`, `AMAZON`, `PUSHOVER`",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"secret_key": {
				Description: "The integration's secret key. Provide for integration: `DO_SPACES`, `AMAZON`",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"app_id": {
				Description: "The integration's application's ID. Provide for integration: `AZURE_CLOUD`",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"tenant_id": {
				Description: "The integration's tenant's ID. Provide for integration: `AZURE_CLOUD`",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"password": {
				Description: "The integration's password. Provide for integration: `AZURE_CLOUD`, `UPCLOUD`, `DOCKER_HUB`",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"api_key": {
				Description: "The integration's API key. Provide for integration: `CLOUDFLARE`, `GOOGLE_SERVICE_ACCOUNT`",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"email": {
				Description: "The integration's email. Provide for integration: `CLOUDFLARE`",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"role_assumption": {
				Description: "The integration's AWS role to assume. Provide for integration: `AMAZON`",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"external_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"duration": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"integration_id": {
				Description: "The integration's ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"html_url": {
				Description: "The integration's URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func deleteContextIntegration(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	domain, hashId, err := util.DecomposeDoubleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = c.IntegrationService.Delete(domain, hashId)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func updateContextIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	domain, hashId, err := util.DecomposeDoubleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	opt := api.IntegrationOperationOptions{
		Name:  util.InterfaceStringToPointer(d.Get("name")),
		Type:  util.InterfaceStringToPointer(d.Get("type")),
		Scope: util.InterfaceStringToPointer(d.Get("scope")),
	}
	if projectName, ok := d.GetOk("project_name"); ok {
		opt.ProjectName = util.InterfaceStringToPointer(projectName)
	}
	if groupId, ok := d.GetOk("group_id"); ok {
		opt.GroupId = util.InterfaceIntToPointer(groupId)
	}
	if username, ok := d.GetOk("username"); ok {
		opt.Username = util.InterfaceStringToPointer(username)
	}
	if shop, ok := d.GetOk("shop"); ok {
		opt.Shop = util.InterfaceStringToPointer(shop)
	}
	if token, ok := d.GetOk("token"); ok {
		opt.Token = util.InterfaceStringToPointer(token)
	}
	if accessKey, ok := d.GetOk("access_key"); ok {
		opt.AccessKey = util.InterfaceStringToPointer(accessKey)
	}
	if secretKey, ok := d.GetOk("secret_key"); ok {
		opt.SecretKey = util.InterfaceStringToPointer(secretKey)
	}
	if appId, ok := d.GetOk("app_id"); ok {
		opt.AppId = util.InterfaceStringToPointer(appId)
	}
	if tenantId, ok := d.GetOk("tenant_id"); ok {
		opt.TenantId = util.InterfaceStringToPointer(tenantId)
	}
	if password, ok := d.GetOk("password"); ok {
		opt.Password = util.InterfaceStringToPointer(password)
	}
	if apiKey, ok := d.GetOk("api_key"); ok {
		opt.ApiKey = util.InterfaceStringToPointer(apiKey)
	}
	if email, ok := d.GetOk("email"); ok {
		opt.Email = util.InterfaceStringToPointer(email)
	}
	if roleAssumptions, ok := d.GetOk("role_assumption"); ok {
		opt.RoleAssumptions = util.MapRoleAssumptionsToApi(roleAssumptions)
	}
	_, _, err = c.IntegrationService.Update(domain, hashId, &opt)
	if err != nil {
		return diag.FromErr(err)
	}
	return readContextIntegration(ctx, d, meta)
}

func readContextIntegration(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	domain, hashId, err := util.DecomposeDoubleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	integration, resp, err := c.IntegrationService.Get(domain, hashId)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	err = util.ApiIntegrationToResourceData(domain, integration, d, false)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func createContextIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	domain := d.Get("domain").(string)
	opt := api.IntegrationOperationOptions{
		Name:  util.InterfaceStringToPointer(d.Get("name")),
		Type:  util.InterfaceStringToPointer(d.Get("type")),
		Scope: util.InterfaceStringToPointer(d.Get("scope")),
	}
	if projectName, ok := d.GetOk("project_name"); ok {
		opt.ProjectName = util.InterfaceStringToPointer(projectName)
	}
	if groupId, ok := d.GetOk("group_id"); ok {
		opt.GroupId = util.InterfaceIntToPointer(groupId)
	}
	if username, ok := d.GetOk("username"); ok {
		opt.Username = util.InterfaceStringToPointer(username)
	}
	if shop, ok := d.GetOk("shop"); ok {
		opt.Shop = util.InterfaceStringToPointer(shop)
	}
	if token, ok := d.GetOk("token"); ok {
		opt.Token = util.InterfaceStringToPointer(token)
	}
	if accessKey, ok := d.GetOk("access_key"); ok {
		opt.AccessKey = util.InterfaceStringToPointer(accessKey)
	}
	if secretKey, ok := d.GetOk("secret_key"); ok {
		opt.SecretKey = util.InterfaceStringToPointer(secretKey)
	}
	if appId, ok := d.GetOk("app_id"); ok {
		opt.AppId = util.InterfaceStringToPointer(appId)
	}
	if tenantId, ok := d.GetOk("tenant_id"); ok {
		opt.TenantId = util.InterfaceStringToPointer(tenantId)
	}
	if password, ok := d.GetOk("password"); ok {
		opt.Password = util.InterfaceStringToPointer(password)
	}
	if apiKey, ok := d.GetOk("api_key"); ok {
		opt.ApiKey = util.InterfaceStringToPointer(apiKey)
	}
	if email, ok := d.GetOk("email"); ok {
		opt.Email = util.InterfaceStringToPointer(email)
	}
	if roleAssumptions, ok := d.GetOk("role_assumption"); ok {
		opt.RoleAssumptions = util.MapRoleAssumptionsToApi(roleAssumptions)
	}
	integration, _, err := c.IntegrationService.Create(domain, &opt)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(util.ComposeDoubleId(domain, integration.HashId))
	return readContextIntegration(ctx, d, meta)
}

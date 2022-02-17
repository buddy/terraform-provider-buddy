package api

import (
	"net/http"
)

const (
	IntegrationTypeDigitalOcean         = "DIGITAL_OCEAN"
	IntegrationTypeAmazon               = "AMAZON"
	IntegrationTypeShopify              = "SHOPIFY"
	IntegrationTypePushover             = "PUSHOVER"
	IntegrationTypeRackspace            = "RACKSPACE"
	IntegrationTypeCloudflare           = "CLOUDFLARE"
	IntegrationTypeNewRelic             = "NEW_RELIC"
	IntegrationTypeSentry               = "SENTRY"
	IntegrationTypeRollbar              = "ROLLBAR"
	IntegrationTypeDatadog              = "DATADOG"
	IntegrationTypeDigitalOceanSpaces   = "DO_SPACES"
	IntegrationTypeHoneybadger          = "HONEYBADGER"
	IntegrationTypeVultr                = "VULTR"
	IntegrationTypeSentryEnterprise     = "SENTRY_ENTERPRISE"
	IntegrationTypeLoggly               = "LOGGLY"
	IntegrationTypeFirebase             = "FIREBASE"
	IntegrationTypeUpcloud              = "UPCLOUD"
	IntegrationTypeGhostInspector       = "GHOST_INSPECTOR"
	IntegrationTypeAzureCloud           = "AZURE_CLOUD"
	IntegrationTypeDockerHub            = "DOCKER_HUB"
	IntegrationTypeGoogleServiceAccount = "GOOGLE_SERVICE_ACCOUNT"

	IntegrationScopePrivate        = "PRIVATE"
	IntegrationScopeWorkspace      = "WORKSPACE"
	IntegrationScopeAdmin          = "ADMIN"
	IntegrationScopeGroup          = "GROUP"
	IntegrationScopeProject        = "PROJECT"
	IntegrationScopeAdminInProject = "ADMIN_IN_PROJECT"
	IntegrationScopeGroupInProject = "GROUP_IN_PROJECT"
)

type Integration struct {
	HtmlUrl     string `json:"html_url"`
	HashId      string `json:"hash_id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Scope       string `json:"scope"`
	ProjectName string `json:"project_name"`
	GroupId     int    `json:"group_id"`
}

type Integrations struct {
	Url          string         `json:"url"`
	HtmlUrl      string         `json:"html_url"`
	Integrations []*Integration `json:"integrations"`
}

type IntegrationOperationOptions struct {
	Name            *string            `json:"name"`
	Type            *string            `json:"type"`
	Scope           *string            `json:"scope"`
	ProjectName     *string            `json:"project_name,omitempty"`
	GroupId         *int               `json:"group_id,omitempty"`
	Username        *string            `json:"username,omitempty"`
	Shop            *string            `json:"shop,omitempty"`
	Token           *string            `json:"token,omitempty"`
	AccessKey       *string            `json:"access_key,omitempty"`
	SecretKey       *string            `json:"secret_key,omitempty"`
	AppId           *string            `json:"app_id,omitempty"`
	TenantId        *string            `json:"tenant_id,omitempty"`
	Password        *string            `json:"password,omitempty"`
	ApiKey          *string            `json:"api_key,omitempty"`
	Email           *string            `json:"email,omitempty"`
	RoleAssumptions *[]*RoleAssumption `json:"role_assumptions,omitempty"`
}

type RoleAssumption struct {
	Arn        string `json:"arn"`
	ExternalId string `json:"external_id,omitempty"`
	Duration   int    `json:"duration,omitempty"`
}

type IntegrationService struct {
	client *Client
}

func (s *IntegrationService) Create(domain string, opt *IntegrationOperationOptions) (*Integration, *http.Response, error) {
	var i *Integration
	resp, err := s.client.Create(s.client.NewUrlPath("/workspaces/%s/integrations", domain), &opt, &i)
	return i, resp, err
}

func (s *IntegrationService) Delete(domain string, hashId string) (*http.Response, error) {
	return s.client.Delete(s.client.NewUrlPath("/workspaces/%s/integrations/%s", domain, hashId))
}

func (s *IntegrationService) Update(domain string, hashId string, opt *IntegrationOperationOptions) (*Integration, *http.Response, error) {
	var i *Integration
	resp, err := s.client.Update(s.client.NewUrlPath("/workspaces/%s/integrations/%s", domain, hashId), &opt, &i)
	return i, resp, err
}

func (s *IntegrationService) Get(domain string, hashId string) (*Integration, *http.Response, error) {
	var i *Integration
	resp, err := s.client.Get(s.client.NewUrlPath("/workspaces/%s/integrations/%s", domain, hashId), &i, nil)
	return i, resp, err
}

func (s *IntegrationService) GetList(domain string) (*Integrations, *http.Response, error) {
	var l *Integrations
	resp, err := s.client.Get(s.client.NewUrlPath("/workspaces/%s/integrations", domain), &l, nil)
	return l, resp, err
}

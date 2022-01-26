package provider

import (
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/resource"
	"buddy-terraform/buddy/source"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"strconv"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BUDDY_TOKEN", nil),
				Description: descriptions["token"],
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BUDDY_BASE_URL", ""),
				Description: descriptions["base_url"],
			},
			"insecure": {
				Type:     schema.TypeBool,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					v := os.Getenv("BUDDY_INSECURE")
					if v != "" {
						return strconv.ParseBool(v)
					}
					return false, nil
				},
				Description: descriptions["insecure"],
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"buddy_profile":            resource.Profile(),
			"buddy_profile_email":      resource.ProfileEmail(),
			"buddy_profile_public_key": resource.ProfilePublicKey(),
			"buddy_workspace":          resource.Workspace(),
			"buddy_group":              resource.Group(),
			"buddy_group_member":       resource.GroupMember(),
			"buddy_member":             resource.Member(),
			"buddy_permission":         resource.Permission(),
			"buddy_project":            resource.Project(),
			"buddy_project_member":     resource.ProjectMember(),
			"buddy_project_group":      resource.ProjectGroup(),
			"buddy_webhook":            resource.Webhook(),
			"buddy_variable":           resource.Variable(),
			"buddy_variable_ssh_key":   resource.VariableSshKey(),
			"buddy_integration":        resource.Integration(),
			"buddy_pipeline":           resource.Pipeline(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"buddy_workspaces":         source.Workspaces(),
			"buddy_profile":            source.Profile(),
			"buddy_group":              source.Group(),
			"buddy_groups":             source.Groups(),
			"buddy_group_members":      source.GroupMembers(),
			"buddy_members":            source.Members(),
			"buddy_member":             source.Member(),
			"buddy_workspace":          source.Workspace(),
			"buddy_permission":         source.Permission(),
			"buddy_permissions":        source.Permissions(),
			"buddy_project":            source.Project(),
			"buddy_projects":           source.Projects(),
			"buddy_project_member":     source.ProjectMember(),
			"buddy_project_members":    source.ProjectMembers(),
			"buddy_project_group":      source.ProjectGroup(),
			"buddy_project_groups":     source.ProjectGroups(),
			"buddy_variable":           source.Variable(),
			"buddy_variable_ssh_key":   source.VariableSshKey(),
			"buddy_variables":          source.Variables(),
			"buddy_variables_ssh_keys": source.VariablesSshKeys(),
			"buddy_webhook":            source.Webhook(),
			"buddy_webhooks":           source.Webhooks(),
			"buddy_integration":        source.Integration(),
			"buddy_integrations":       source.Integrations(),
			"buddy_pipeline":           source.Pipeline(),
			"buddy_pipelines":          source.Pipelines(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"token":    "The OAuth2 token or Personal Access Token. Can be specified with the `BUDDY_TOKEN` environment variable.",
		"base_url": "The Buddy API base url. You may need to set this to your Buddy On-Premises API endpoint. Can be specified with the `BUDDY_BASE_URL` environment variable. Default: `https://api.buddy.works`",
		"insecure": "Disable SSL verification of API calls. You may need to set this to `true` if you are using Buddy On-Premises without signed certificate. Can be specified with the `BUDDY_INSECURE` environmental variable",
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("token").(string)
	baseUrl := d.Get("base_url").(string)
	insecure := d.Get("insecure").(bool)
	var diags diag.Diagnostics
	c, err := api.NewClient(token, baseUrl, insecure)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return c, diags
}

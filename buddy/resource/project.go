package resource

import (
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
)

func Project() *schema.Resource {
	return &schema.Resource{
		Description: "Create and manage a workspace project\n\n" +
			"Workspace administrator rights are required\n\n" +
			"Token scopes required: `WORKSPACE`, `PROJECT_DELETE`",
		CreateContext: createContextProject,
		ReadContext:   readContextProject,
		UpdateContext: updateContextProject,
		DeleteContext: deleteContextProject,
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
			"display_name": {
				Description: "The project's display name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"integration_id": {
				Description: "The project's integration ID. Needed when cloning from a GitHub, GitLab or BitBucket",
				Type:        schema.TypeString,
				Optional:    true,
				ConflictsWith: []string{
					"custom_repo_url",
					"custom_repo_user",
					"custom_repo_pass",
				},
				RequiredWith: []string{
					"external_project_id",
				},
			},
			"external_project_id": {
				Description: "The project's external project ID. Needed when cloning from GitHub, GitLab or BitBucket",
				Type:        schema.TypeString,
				Optional:    true,
				ConflictsWith: []string{
					"custom_repo_url",
					"custom_repo_user",
					"custom_repo_pass",
				},
				RequiredWith: []string{
					"integration_id",
				},
			},
			"git_lab_project_id": {
				Description: "The project's GitLab project ID. Needed when cloning from a GitLab",
				Type:        schema.TypeString,
				Optional:    true,
				ConflictsWith: []string{
					"custom_repo_url",
					"custom_repo_user",
					"custom_repo_pass",
				},
				RequiredWith: []string{
					"integration_id",
					"external_project_id",
				},
			},
			"custom_repo_url": {
				Description: "The project's custom repository URL. Needed when cloning from a custom repository",
				Type:        schema.TypeString,
				Optional:    true,
				ConflictsWith: []string{
					"integration_id",
					"external_project_id",
					"git_lab_project_id",
				},
			},
			"custom_repo_user": {
				Description: "The project's custom repository user. Needed when cloning from a custom repository",
				Type:        schema.TypeString,
				Optional:    true,
				ConflictsWith: []string{
					"integration_id",
					"external_project_id",
					"git_lab_project_id",
				},
				RequiredWith: []string{
					"custom_repo_url",
					"custom_repo_pass",
				},
			},
			"custom_repo_pass": {
				Description: "The project's custom repository password. Needed when cloning from a custom repository",
				Type:        schema.TypeString,
				Optional:    true,
				ConflictsWith: []string{
					"integration_id",
					"external_project_id",
					"git_lab_project_id",
				},
				RequiredWith: []string{
					"custom_repo_url",
					"custom_repo_user",
				},
			},
			"name": {
				Description: "The project's unique name ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"html_url": {
				Description: "The project's URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"status": {
				Description: "The project's status. Possible values: `CLOSED`, `ACTIVE`",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"create_date": {
				Description: "The project's date of creation",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_by": {
				Description: "The project's creator",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"html_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"member_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"admin": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"workspace_owner": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"avatar_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"http_repository": {
				Description: "The project's Git HTTP endpoint",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ssh_repository": {
				Description: "The project's Git SSH endpoint",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"default_branch": {
				Description: "The project's Git default branch",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ssh_public_key": {
				Description: "The project's public key",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"key_fingerprint": {
				Description: "The project's key fingerprint",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func deleteContextProject(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	domain, name, err := util.DecomposeDoubleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = c.ProjectService.Delete(domain, name)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func updateContextProject(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	domain, name, err := util.DecomposeDoubleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	u := api.ProjectUpdateOptions{
		DisplayName: util.InterfaceStringToPointer(d.Get("display_name")),
	}
	p, _, err := c.ProjectService.Update(domain, name, &u)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(util.ComposeDoubleId(domain, p.Name))
	return readContextProject(ctx, d, meta)
}

func readContextProject(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	var diags diag.Diagnostics
	domain, name, err := util.DecomposeDoubleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	p, resp, err := c.ProjectService.Get(domain, name)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	err = util.ApiProjectToResourceData(domain, p, d, false)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func createContextProject(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*api.Client)
	domain := d.Get("domain").(string)
	opt := api.ProjectCreateOptions{
		DisplayName: util.InterfaceStringToPointer(d.Get("display_name")),
	}
	if integrationId, ok := d.GetOk("integration_id"); ok {
		opt.Integration = &api.ProjectIntegration{
			HashId: integrationId.(string),
		}
	}
	if externalProjectId, ok := d.GetOk("external_project_id"); ok {
		opt.ExternalProjectId = util.InterfaceStringToPointer(externalProjectId)
	}
	if gitLabProjectId, ok := d.GetOk("git_lab_project_id"); ok {
		opt.GitLabProjectId = util.InterfaceStringToPointer(gitLabProjectId)
	}
	if customRepoUrl, ok := d.GetOk("custom_repo_url"); ok {
		opt.CustomRepoUrl = util.InterfaceStringToPointer(customRepoUrl)
	}
	if customRepoUser, ok := d.GetOk("custom_repo_user"); ok {
		opt.CustomRepoUser = util.InterfaceStringToPointer(customRepoUser)
	}
	if customRepoPass, ok := d.GetOk("custom_repo_pass"); ok {
		opt.CustomRepoPass = util.InterfaceStringToPointer(customRepoPass)
	}
	project, _, err := c.ProjectService.Create(domain, &opt)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(util.ComposeDoubleId(domain, project.Name))
	return readContextProject(ctx, d, meta)
}

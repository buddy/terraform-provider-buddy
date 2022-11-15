package resource

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
					"custom_repo_ssh_key_id",
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
					"custom_repo_ssh_key_id",
					"custom_repo_pass",
				},
				RequiredWith: []string{
					"integration_id",
				},
			},
			"update_default_branch_from_external": {
				Description: "Defines whether or not update default branch from external repository (GitHub, GitLab, BitBucket)",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"fetch_submodules": {
				Description: "Defines wheter or not fetch submodules in repository",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				RequiredWith: []string{
					"fetch_submodules_env_key",
				},
			},
			"fetch_submodules_env_key": {
				Description: "The project's environmental key name for fetching submodules",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				RequiredWith: []string{
					"fetch_submodules",
				},
			},
			"access": {
				Description: "The project's access. Possible values: `PRIVATE`, `PUBLIC`",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"allow_pull_requests": {
				Description: "Defines whether or not pull requests are enabled (GitHub)",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"git_lab_project_id": {
				Description: "The project's GitLab project ID. Needed when cloning from a GitLab",
				Type:        schema.TypeString,
				Optional:    true,
				ConflictsWith: []string{
					"custom_repo_url",
					"custom_repo_user",
					"custom_repo_ssh_key_id",
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
			"custom_repo_ssh_key_id": {
				Description: "The project's custom repository SSH key ID. Needed when cloning from a custom repository",
				Type:        schema.TypeInt,
				Optional:    true,
				ConflictsWith: []string{
					"integration_id",
					"external_project_id",
					"git_lab_project_id",
				},
				RequiredWith: []string{
					"custom_repo_url",
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
		},
	}
}

func deleteContextProject(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
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
	c := meta.(*buddy.Client)
	domain, name, err := util.DecomposeDoubleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	u := buddy.ProjectUpdateOps{
		DisplayName: util.InterfaceStringToPointer(d.Get("display_name")),
	}
	if d.HasChange("update_default_branch_from_external") {
		u.UpdateDefaultBranchFromExternal = util.InterfaceBoolToPointer(d.Get("update_default_branch_from_external"))
	}
	if d.HasChange("allow_pull_requests") {
		u.AllowPullRequests = util.InterfaceBoolToPointer(d.Get("allow_pull_requests"))
	}
	if d.HasChange("access") {
		u.Access = util.InterfaceStringToPointer(d.Get("access"))
	}
	if d.HasChanges("fetch_submodules", "fetch_submodules_env_key") {
		u.FetchSubmodules = util.InterfaceBoolToPointer(d.Get("fetch_submodules"))
		u.FetchSubmodulesEnvKey = util.InterfaceStringToPointer(d.Get("fetch_submodules_env_key"))
	}
	p, _, err := c.ProjectService.Update(domain, name, &u)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(util.ComposeDoubleId(domain, p.Name))
	return readContextProject(ctx, d, meta)
}

func readContextProject(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*buddy.Client)
	var diags diag.Diagnostics
	domain, name, err := util.DecomposeDoubleId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	p, resp, err := c.ProjectService.Get(domain, name)
	if err != nil {
		if util.IsResourceNotFound(resp, err) {
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
	c := meta.(*buddy.Client)
	domain := d.Get("domain").(string)
	opt := buddy.ProjectCreateOps{
		DisplayName: util.InterfaceStringToPointer(d.Get("display_name")),
	}
	if integrationId, ok := d.GetOk("integration_id"); ok {
		opt.Integration = &buddy.ProjectIntegration{
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
	if customRepoSshKeyId, ok := d.GetOk("custom_repo_ssh_key_id"); ok {
		opt.CustomRepoSshKeyId = util.InterfaceIntToPointer(customRepoSshKeyId)
	}
	if access, ok := d.GetOk("access"); ok {
		opt.Access = util.InterfaceStringToPointer(access)
	}
	if fetchSubmodulesEnv, ok := d.GetOk("fetch_submodules_env_key"); ok {
		opt.FetchSubmodulesEnvKey = util.InterfaceStringToPointer(fetchSubmodulesEnv)
	}
	updateDefaultBranch := d.Get("update_default_branch_from_external")
	if util.IsBoolPointerSet(updateDefaultBranch) {
		opt.UpdateDefaultBranchFromExternal = util.InterfaceBoolToPointer(updateDefaultBranch)
	}
	allowPullRequests := d.Get("allow_pull_requests")
	if util.IsBoolPointerSet(allowPullRequests) {
		opt.AllowPullRequests = util.InterfaceBoolToPointer(allowPullRequests)
	}
	fetchSubmodules := d.Get("fetch_submodules")
	if util.IsBoolPointerSet(fetchSubmodules) {
		opt.FetchSubmodules = util.InterfaceBoolToPointer(fetchSubmodules)
	}
	project, _, err := c.ProjectService.Create(domain, &opt)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(util.ComposeDoubleId(domain, project.Name))
	return readContextProject(ctx, d, meta)
}

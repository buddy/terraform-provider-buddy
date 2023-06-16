package resource

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	_ resource.Resource                = &projectResource{}
	_ resource.ResourceWithConfigure   = &projectResource{}
	_ resource.ResourceWithImportState = &projectResource{}
)

func NewProjectResource() resource.Resource {
	return &projectResource{}
}

type projectResource struct {
	client *buddy.Client
}

type projectResourceModel struct {
	ID                              types.String `tfsdk:"id"`
	Domain                          types.String `tfsdk:"domain"`
	DisplayName                     types.String `tfsdk:"display_name"`
	WithoutRepository               types.Bool   `tfsdk:"without_repository"`
	IntegrationId                   types.String `tfsdk:"integration_id"`
	ExternalProjectId               types.String `tfsdk:"external_project_id"`
	UpdateDefaultBranchFromExternal types.Bool   `tfsdk:"update_default_branch_from_external"`
	FetchSubmodules                 types.Bool   `tfsdk:"fetch_submodules"`
	FetchSubmodulesEnvKey           types.String `tfsdk:"fetch_submodules_env_key"`
	Access                          types.String `tfsdk:"access"`
	AllowPullRequests               types.Bool   `tfsdk:"allow_pull_requests"`
	GitLabProjectId                 types.String `tfsdk:"git_lab_project_id"`
	CustomRepoUrl                   types.String `tfsdk:"custom_repo_url"`
	CustomRepoSshKeyId              types.Int64  `tfsdk:"custom_repo_ssh_key_id"`
	CustomRepoUser                  types.String `tfsdk:"custom_repo_user"`
	CustomRepoPass                  types.String `tfsdk:"custom_repo_pass"`
	Name                            types.String `tfsdk:"name"`
	HtmlUrl                         types.String `tfsdk:"html_url"`
	Status                          types.String `tfsdk:"status"`
	CreateDate                      types.String `tfsdk:"create_date"`
	CreatedBy                       types.Set    `tfsdk:"created_by"`
	HttpRepository                  types.String `tfsdk:"http_repository"`
	SshRepository                   types.String `tfsdk:"ssh_repository"`
	DefaultBranch                   types.String `tfsdk:"default_branch"`
}

func (r *projectResourceModel) decomposeId() (string, string, error) {
	domain, name, err := util.DecomposeDoubleId(r.ID.ValueString())
	if err != nil {
		return "", "", err
	}
	return domain, name, nil
}

func (r *projectResourceModel) loadAPI(ctx context.Context, domain string, project *buddy.Project) diag.Diagnostics {
	r.ID = types.StringValue(util.ComposeDoubleId(domain, project.Name))
	r.Domain = types.StringValue(domain)
	r.DisplayName = types.StringValue(project.DisplayName)
	r.UpdateDefaultBranchFromExternal = types.BoolValue(project.UpdateDefaultBranchFromExternal)
	r.FetchSubmodules = types.BoolValue(project.FetchSubmodules)
	r.FetchSubmodulesEnvKey = types.StringValue(project.FetchSubmodulesEnvKey)
	r.Access = types.StringValue(project.Access)
	r.AllowPullRequests = types.BoolValue(project.AllowPullRequests)
	r.Name = types.StringValue(project.Name)
	r.HtmlUrl = types.StringValue(project.HtmlUrl)
	r.Status = types.StringValue(project.Status)
	r.CreateDate = types.StringValue(project.CreateDate)
	r.HttpRepository = types.StringValue(project.HttpRepository)
	r.SshRepository = types.StringValue(project.SshRepository)
	r.DefaultBranch = types.StringValue(project.DefaultBranch)
	creatorSet := []*buddy.Member{project.CreatedBy}
	createdBy, diags := util.MembersModelFromApi(ctx, &creatorSet)
	r.CreatedBy = createdBy
	return diags
}

func (r *projectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *projectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage a workspace project\n\n" +
			"Workspace administrator rights are required\n\n" +
			"Token scopes required: `WORKSPACE`, `PROJECT_DELETE`",
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
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The project's display name",
				Required:            true,
			},
			"without_repository": schema.BoolAttribute{
				MarkdownDescription: "Defines wheter or not create GIT repository",
				Optional:            true,
			},
			"integration_id": schema.StringAttribute{
				MarkdownDescription: "The project's integration ID. Needed when cloning from a GitHub, GitLab or BitBucket",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("custom_repo_url"),
						path.MatchRoot("custom_repo_user"),
						path.MatchRoot("custom_repo_ssh_key_id"),
						path.MatchRoot("custom_repo_pass"),
					}...),
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("external_project_id"),
					}...),
				},
			},
			"external_project_id": schema.StringAttribute{
				MarkdownDescription: "The project's external project ID. Needed when cloning from GitHub, GitLab or BitBucket",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("custom_repo_url"),
						path.MatchRoot("custom_repo_user"),
						path.MatchRoot("custom_repo_ssh_key_id"),
						path.MatchRoot("custom_repo_pass"),
					}...),
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("integration_id"),
					}...),
				},
			},
			"update_default_branch_from_external": schema.BoolAttribute{
				MarkdownDescription: "Defines whether or not update default branch from external repository (GitHub, GitLab, BitBucket)",
				Optional:            true,
				Computed:            true,
			},
			"fetch_submodules": schema.BoolAttribute{
				MarkdownDescription: "Defines wheter or not fetch submodules in repository",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Bool{
					boolvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("fetch_submodules_env_key"),
					}...),
				},
			},
			"fetch_submodules_env_key": schema.StringAttribute{
				MarkdownDescription: "The project's environmental key name for fetching submodules",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("fetch_submodules"),
					}...),
				},
			},
			"access": schema.StringAttribute{
				MarkdownDescription: "The project's access. Possible values: `PRIVATE`, `PUBLIC`",
				Optional:            true,
				Computed:            true,
			},
			"allow_pull_requests": schema.BoolAttribute{
				MarkdownDescription: "Defines whether or not pull requests are enabled (GitHub)",
				Optional:            true,
				Computed:            true,
			},
			"git_lab_project_id": schema.StringAttribute{
				MarkdownDescription: "The project's GitLab project ID. Needed when cloning from a GitLab",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("custom_repo_url"),
						path.MatchRoot("custom_repo_user"),
						path.MatchRoot("custom_repo_ssh_key_id"),
						path.MatchRoot("custom_repo_pass"),
					}...),
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("integration_id"),
						path.MatchRoot("external_project_id"),
					}...),
				},
			},
			"custom_repo_url": schema.StringAttribute{
				MarkdownDescription: "The project's custom repository URL. Needed when cloning from a custom repository",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("integration_id"),
						path.MatchRoot("external_project_id"),
						path.MatchRoot("git_lab_project_id"),
					}...),
				},
			},
			"custom_repo_ssh_key_id": schema.Int64Attribute{
				MarkdownDescription: "The project's custom repository SSH key ID. Needed when cloning from a custom repository",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.ConflictsWith(path.Expressions{
						path.MatchRoot("integration_id"),
						path.MatchRoot("external_project_id"),
						path.MatchRoot("git_lab_project_id"),
					}...),
					int64validator.AlsoRequires(path.Expressions{
						path.MatchRoot("custom_repo_url"),
					}...),
				},
			},
			"custom_repo_user": schema.StringAttribute{
				MarkdownDescription: "The project's custom repository user. Needed when cloning from a custom repository",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("integration_id"),
						path.MatchRoot("external_project_id"),
						path.MatchRoot("git_lab_project_id"),
					}...),
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("custom_repo_url"),
						path.MatchRoot("custom_repo_pass"),
					}...),
				},
			},
			"custom_repo_pass": schema.StringAttribute{
				MarkdownDescription: "The project's custom repository password. Needed when cloning from a custom repository",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("integration_id"),
						path.MatchRoot("external_project_id"),
						path.MatchRoot("git_lab_project_id"),
					}...),
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("custom_repo_url"),
						path.MatchRoot("custom_repo_user"),
					}...),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The project's unique name ID",
				Computed:            true,
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The project's URL",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The project's status. Possible values: `CLOSED`, `ACTIVE`",
				Computed:            true,
			},
			"create_date": schema.StringAttribute{
				MarkdownDescription: "The project's date of creation",
				Computed:            true,
			},
			"http_repository": schema.StringAttribute{
				MarkdownDescription: "The project's Git HTTP endpoint",
				Computed:            true,
			},
			"ssh_repository": schema.StringAttribute{
				MarkdownDescription: "The project's Git SSH endpoint",
				Computed:            true,
			},
			"default_branch": schema.StringAttribute{
				MarkdownDescription: "The project's Git default branch",
				Computed:            true,
			},
			// for compatibility it's a set
			"created_by": schema.SetNestedAttribute{
				MarkdownDescription: "The project's creator",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.ResourceMemberModelAttributes(),
				},
			},
		},
	}
}

func (r *projectResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *projectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *projectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	ops := buddy.ProjectCreateOps{
		DisplayName: data.DisplayName.ValueStringPointer(),
	}
	if !data.IntegrationId.IsNull() && !data.IntegrationId.IsUnknown() {
		ops.Integration = &buddy.ProjectIntegration{
			HashId: data.IntegrationId.ValueString(),
		}
	}
	if !data.ExternalProjectId.IsNull() && !data.ExternalProjectId.IsUnknown() {
		ops.ExternalProjectId = data.ExternalProjectId.ValueStringPointer()
	}
	if !data.GitLabProjectId.IsNull() && !data.GitLabProjectId.IsUnknown() {
		ops.GitLabProjectId = data.GitLabProjectId.ValueStringPointer()
	}
	if !data.CustomRepoUrl.IsNull() && !data.CustomRepoUrl.IsUnknown() {
		ops.CustomRepoUrl = data.CustomRepoUrl.ValueStringPointer()
	}
	if !data.CustomRepoUser.IsNull() && !data.CustomRepoUser.IsUnknown() {
		ops.CustomRepoUser = data.CustomRepoUser.ValueStringPointer()
	}
	if !data.CustomRepoPass.IsNull() && !data.CustomRepoPass.IsUnknown() {
		ops.CustomRepoPass = data.CustomRepoPass.ValueStringPointer()
	}
	if !data.CustomRepoSshKeyId.IsNull() && !data.CustomRepoSshKeyId.IsUnknown() {
		ops.CustomRepoSshKeyId = util.PointerInt(data.CustomRepoSshKeyId.ValueInt64())
	}
	if !data.Access.IsNull() && !data.Access.IsUnknown() {
		ops.Access = data.Access.ValueStringPointer()
	}
	if !data.FetchSubmodulesEnvKey.IsNull() && !data.FetchSubmodulesEnvKey.IsUnknown() {
		ops.FetchSubmodulesEnvKey = data.FetchSubmodulesEnvKey.ValueStringPointer()
	}
	if !data.WithoutRepository.IsNull() && !data.WithoutRepository.IsUnknown() {
		ops.WithoutRepository = data.WithoutRepository.ValueBoolPointer()
	}
	if !data.UpdateDefaultBranchFromExternal.IsNull() && !data.UpdateDefaultBranchFromExternal.IsUnknown() {
		ops.UpdateDefaultBranchFromExternal = data.UpdateDefaultBranchFromExternal.ValueBoolPointer()
	}
	if !data.AllowPullRequests.IsNull() && !data.AllowPullRequests.IsUnknown() {
		ops.AllowPullRequests = data.AllowPullRequests.ValueBoolPointer()
	}
	if !data.FetchSubmodules.IsNull() && !data.FetchSubmodules.IsUnknown() {
		ops.FetchSubmodules = data.FetchSubmodules.ValueBoolPointer()
	}
	project, _, err := r.client.ProjectService.Create(domain, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create project", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, project)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *projectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, projectName, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("project", err))
		return
	}
	project, httpResp, err := r.client.ProjectService.Get(domain, projectName)
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get project", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, project)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *projectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, projectName, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("project", err))
		return
	}
	ops := buddy.ProjectUpdateOps{
		DisplayName: data.DisplayName.ValueStringPointer(),
	}
	if !data.UpdateDefaultBranchFromExternal.IsNull() && !data.UpdateDefaultBranchFromExternal.IsUnknown() {
		ops.UpdateDefaultBranchFromExternal = data.UpdateDefaultBranchFromExternal.ValueBoolPointer()
	}
	if !data.AllowPullRequests.IsNull() && !data.AllowPullRequests.IsUnknown() {
		ops.AllowPullRequests = data.AllowPullRequests.ValueBoolPointer()
	}
	if !data.Access.IsNull() && !data.Access.IsUnknown() {
		ops.Access = data.Access.ValueStringPointer()
	}
	if !data.FetchSubmodules.IsNull() && !data.FetchSubmodules.IsUnknown() {
		ops.FetchSubmodules = data.FetchSubmodules.ValueBoolPointer()
	}
	if !data.FetchSubmodulesEnvKey.IsNull() && !data.FetchSubmodulesEnvKey.IsUnknown() {
		ops.FetchSubmodulesEnvKey = data.FetchSubmodulesEnvKey.ValueStringPointer()
	}
	project, _, err := r.client.ProjectService.Update(domain, projectName, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update project", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, project)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *projectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, projectName, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("project", err))
		return
	}
	_, err = r.client.ProjectService.Delete(domain, projectName)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete project", err))
	}
}

func (r *projectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

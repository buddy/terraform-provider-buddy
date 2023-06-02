package source

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
	"strconv"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ datasource.DataSource              = &permissionSource{}
	_ datasource.DataSourceWithConfigure = &permissionSource{}
)

func NewPermissionSource() datasource.DataSource {
	return &permissionSource{}
}

type permissionSource struct {
	client *buddy.Client
}

type permissionSourceModel struct {
	ID                     types.String `tfsdk:"id"`
	Domain                 types.String `tfsdk:"domain"`
	Name                   types.String `tfsdk:"name"`
	PermissionId           types.Int64  `tfsdk:"permission_id"`
	PipelineAccessLevel    types.String `tfsdk:"pipeline_access_level"`
	RepositoryAccessLevel  types.String `tfsdk:"repository_access_level"`
	ProjectTeamAccessLevel types.String `tfsdk:"project_team_access_level"`
	SandboxAccessLevel     types.String `tfsdk:"sandbox_access_level"`
	HtmlUrl                types.String `tfsdk:"html_url"`
	Description            types.String `tfsdk:"description"`
	Type                   types.String `tfsdk:"type"`
}

func (s *permissionSourceModel) loadAPI(domain string, permission *buddy.Permission) {
	s.ID = types.StringValue(util.ComposeDoubleId(domain, strconv.Itoa(permission.Id)))
	s.Domain = types.StringValue(domain)
	s.Name = types.StringValue(permission.Name)
	s.PermissionId = types.Int64Value(int64(permission.Id))
	s.PipelineAccessLevel = types.StringValue(permission.PipelineAccessLevel)
	s.RepositoryAccessLevel = types.StringValue(permission.RepositoryAccessLevel)
	s.ProjectTeamAccessLevel = types.StringValue(permission.ProjectTeamAccessLevel)
	s.SandboxAccessLevel = types.StringValue(permission.SandboxAccessLevel)
	s.HtmlUrl = types.StringValue(permission.HtmlUrl)
	s.Description = types.StringValue(permission.Description)
	s.Type = types.StringValue(permission.Type)
}

func (s *permissionSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permission"
}

func (s *permissionSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *permissionSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get permission (role) by name or permission ID\n\n" +
			"Token scope required: `WORKSPACE`",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Terraform resource identifier for this item",
				Computed:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The workspace's URL handle",
				Required:            true,
				Validators:          util.StringValidatorsDomain(),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The permission's name",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("permission_id"),
						path.MatchRoot("name"),
					}...),
				},
			},
			"permission_id": schema.Int64Attribute{
				MarkdownDescription: "The permission's ID",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("permission_id"),
						path.MatchRoot("name"),
					}...),
				},
			},
			"pipeline_access_level": schema.StringAttribute{
				MarkdownDescription: "The permission's access level to pipelines",
				Computed:            true,
			},
			"repository_access_level": schema.StringAttribute{
				MarkdownDescription: "The permission's access level to repository",
				Computed:            true,
			},
			"project_team_access_level": schema.StringAttribute{
				MarkdownDescription: "The permission's access level to team",
				Computed:            true,
			},
			"sandbox_access_level": schema.StringAttribute{
				MarkdownDescription: "The permission's access level to sandboxes",
				Computed:            true,
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The permission's URL",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The permission's description",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The permission's type",
				Computed:            true,
			},
		},
	}
}

func (s *permissionSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *permissionSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var permission *buddy.Permission
	var err error
	domain := data.Domain.ValueString()
	if !data.PermissionId.IsNull() && !data.PermissionId.IsUnknown() {
		var httpRes *http.Response
		permId := int(data.PermissionId.ValueInt64())
		permission, httpRes, err = s.client.PermissionService.Get(domain, permId)
		if err != nil {
			if util.IsResourceNotFound(httpRes, err) {
				resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("permission"))
				return
			}
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get permission", err))
			return
		}
	} else {
		name := data.Name.ValueString()
		var permissions *buddy.Permissions
		permissions, _, err = s.client.PermissionService.GetList(domain)
		if err != nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get permissions", err))
			return
		}
		for _, p := range permissions.PermissionSets {
			if p.Name == name {
				permission = p
				break
			}
		}
		if permission == nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("permission"))
			return
		}
	}
	data.loadAPI(domain, permission)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

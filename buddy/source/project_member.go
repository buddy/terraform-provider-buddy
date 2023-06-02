package source

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ datasource.DataSource              = &projectMemberSource{}
	_ datasource.DataSourceWithConfigure = &projectMemberSource{}
)

func NewProjectMemberSource() datasource.DataSource {
	return &projectMemberSource{}
}

type projectMemberSource struct {
	client *buddy.Client
}

type projectMemberSourceModel struct {
	ID             types.String `tfsdk:"id"`
	Domain         types.String `tfsdk:"domain"`
	ProjectName    types.String `tfsdk:"project_name"`
	MemberId       types.Int64  `tfsdk:"member_id"`
	HtmlUrl        types.String `tfsdk:"html_url"`
	Name           types.String `tfsdk:"name"`
	Email          types.String `tfsdk:"email"`
	AvatarUrl      types.String `tfsdk:"avatar_url"`
	Admin          types.Bool   `tfsdk:"admin"`
	WorkspaceOwner types.Bool   `tfsdk:"workspace_owner"`
	Permission     types.Set    `tfsdk:"permission"`
}

func (s *projectMemberSourceModel) loadAPI(ctx context.Context, domain string, projectName string, projectMember *buddy.ProjectMember) diag.Diagnostics {
	s.ID = types.StringValue(util.ComposeTripleId(domain, projectName, strconv.Itoa(projectMember.Id)))
	s.Domain = types.StringValue(domain)
	s.ProjectName = types.StringValue(projectName)
	s.MemberId = types.Int64Value(int64(projectMember.Id))
	s.HtmlUrl = types.StringValue(projectMember.HtmlUrl)
	s.Name = types.StringValue(projectMember.Name)
	s.Email = types.StringValue(projectMember.Email)
	s.AvatarUrl = types.StringValue(projectMember.AvatarUrl)
	s.Admin = types.BoolValue(projectMember.Admin)
	s.WorkspaceOwner = types.BoolValue(projectMember.WorkspaceOwner)
	perms := []*buddy.Permission{projectMember.PermissionSet}
	p, d := util.PermissionsModelFromApi(ctx, &perms)
	s.Permission = p
	return d
}

func (s *projectMemberSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_member"
}

func (s *projectMemberSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *projectMemberSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get project member\n\n" +
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
			"project_name": schema.StringAttribute{
				MarkdownDescription: "The project's name",
				Required:            true,
			},
			"member_id": schema.Int64Attribute{
				MarkdownDescription: "The member's ID",
				Required:            true,
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The member's URL",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The member's name",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The member's email",
				Computed:            true,
			},
			"avatar_url": schema.StringAttribute{
				MarkdownDescription: "The member's avatar URL",
				Computed:            true,
			},
			"admin": schema.BoolAttribute{
				MarkdownDescription: "Is the member a workspace administrator",
				Computed:            true,
			},
			"workspace_owner": schema.BoolAttribute{
				MarkdownDescription: "Is the member the workspace owner",
				Computed:            true,
			},
			"permission": schema.SetNestedAttribute{
				MarkdownDescription: "The member's permission in the project",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.SourcePermissionModelAttributes(),
				},
			},
		},
	}
}

func (s *projectMemberSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *projectMemberSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	projectName := data.ProjectName.ValueString()
	memberId := int(data.MemberId.ValueInt64())
	projectMember, httpRes, err := s.client.ProjectMemberService.GetProjectMember(domain, projectName, memberId)
	if err != nil {
		if util.IsResourceNotFound(httpRes, err) {
			resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("project member"))
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get project member", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, projectMember)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

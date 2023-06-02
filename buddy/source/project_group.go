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
	_ datasource.DataSource              = &projectGroupSource{}
	_ datasource.DataSourceWithConfigure = &projectGroupSource{}
)

func NewProjectGroupSource() datasource.DataSource {
	return &projectGroupSource{}
}

type projectGroupSource struct {
	client *buddy.Client
}

type projectGroupSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Domain      types.String `tfsdk:"domain"`
	ProjectName types.String `tfsdk:"project_name"`
	GroupId     types.Int64  `tfsdk:"group_id"`
	HtmlUrl     types.String `tfsdk:"html_url"`
	Name        types.String `tfsdk:"name"`
	Permission  types.Set    `tfsdk:"permission"`
}

func (s *projectGroupSourceModel) loadAPI(ctx context.Context, domain string, projectName string, projectGroup *buddy.ProjectGroup) diag.Diagnostics {
	s.ID = types.StringValue(util.ComposeTripleId(domain, projectName, strconv.Itoa(projectGroup.Id)))
	s.Domain = types.StringValue(domain)
	s.ProjectName = types.StringValue(projectName)
	s.GroupId = types.Int64Value(int64(projectGroup.Id))
	s.HtmlUrl = types.StringValue(projectGroup.HtmlUrl)
	s.Name = types.StringValue(projectGroup.Name)
	perms := []*buddy.Permission{projectGroup.PermissionSet}
	p, d := util.PermissionsModelFromApi(ctx, &perms)
	s.Permission = p
	return d
}

func (s *projectGroupSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_group"
}

func (s *projectGroupSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *projectGroupSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get project group\n\n" +
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
			"group_id": schema.Int64Attribute{
				MarkdownDescription: "The group's ID",
				Required:            true,
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The group's URL",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The group's name",
				Computed:            true,
			},
			"permission": schema.SetNestedAttribute{
				MarkdownDescription: "The group's permission in the project",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.SourcePermissionModelAttributes(),
				},
			},
		},
	}
}

func (s *projectGroupSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *projectGroupSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	projectName := data.ProjectName.ValueString()
	groupId := int(data.GroupId.ValueInt64())
	projectGroup, httpRes, err := s.client.ProjectGroupService.GetProjectGroup(domain, projectName, groupId)
	if err != nil {
		if util.IsResourceNotFound(httpRes, err) {
			resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("project group"))
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get project group", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, projectGroup)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

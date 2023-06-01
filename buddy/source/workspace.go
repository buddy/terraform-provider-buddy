package source

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
)

var (
	_ datasource.DataSource              = &workspaceSource{}
	_ datasource.DataSourceWithConfigure = &workspaceSource{}
)

func NewWorkspaceSource() datasource.DataSource {
	return &workspaceSource{}
}

type workspaceSource struct {
	client *buddy.Client
}

type workspaceSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Domain      types.String `tfsdk:"domain"`
	Name        types.String `tfsdk:"name"`
	WorkspaceId types.Int64  `tfsdk:"workspace_id"`
	HtmlUrl     types.String `tfsdk:"html_url"`
}

func (s *workspaceSourceModel) loadAPI(workspace *buddy.Workspace) {
	s.ID = types.StringValue(workspace.Domain)
	s.Domain = types.StringValue(workspace.Domain)
	s.Name = types.StringValue(workspace.Name)
	s.WorkspaceId = types.Int64Value(int64(workspace.Id))
	s.HtmlUrl = types.StringValue(workspace.HtmlUrl)
}

func (s *workspaceSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace"
}

func (s *workspaceSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *workspaceSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get workspace by URL handle or name\n\n" +
			"Token scope required: `WORKSPACE`",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Terraform resource identifier for this item",
				Computed:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The workspace's URL handle",
				Optional:            true,
				Computed:            true,
				Validators: append([]validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("domain"),
						path.MatchRoot("name"),
					}...),
				}, util.StringValidatorsDomain()...),
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The workspace's name",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("domain"),
						path.MatchRoot("name"),
					}...),
				},
			},
			"workspace_id": schema.Int64Attribute{
				MarkdownDescription: "The workspace's ID",
				Computed:            true,
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The workspace's URL",
				Computed:            true,
			},
		},
	}
}

func (s *workspaceSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *workspaceSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var workspace *buddy.Workspace
	var err error
	if !data.Domain.IsNull() && !data.Domain.IsUnknown() {
		var httpRes *http.Response
		workspace, httpRes, err = s.client.WorkspaceService.Get(data.Domain.ValueString())
		if err != nil {
			if util.IsResourceNotFound(httpRes, err) {
				resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("workspace"))
				return
			}
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get workspace", err))
			return
		}
	} else {
		name := data.Name.ValueString()
		var workspaces *buddy.Workspaces
		workspaces, _, err = s.client.WorkspaceService.GetList()
		if err != nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get workspaces", err))
			return
		}
		for _, w := range workspaces.Workspaces {
			if w.Name == name {
				workspace = w
				break
			}
		}
		if workspace == nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("workspace"))
			return
		}
	}
	data.loadAPI(workspace)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

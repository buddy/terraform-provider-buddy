package source

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ datasource.DataSource              = &projectSource{}
	_ datasource.DataSourceWithConfigure = &projectSource{}
)

func NewProjectSource() datasource.DataSource {
	return &projectSource{}
}

type projectSource struct {
	client *buddy.Client
}

type projectSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Domain      types.String `tfsdk:"domain"`
	DisplayName types.String `tfsdk:"display_name"`
	Name        types.String `tfsdk:"name"`
	HtmlUrl     types.String `tfsdk:"html_url"`
	Status      types.String `tfsdk:"status"`
}

func (s *projectSourceModel) loadAPI(domain string, project *buddy.Project) {
	s.ID = types.StringValue(util.ComposeDoubleId(domain, project.Name))
	s.Domain = types.StringValue(domain)
	s.DisplayName = types.StringValue(project.DisplayName)
	s.Name = types.StringValue(project.Name)
	s.HtmlUrl = types.StringValue(project.HtmlUrl)
	s.Status = types.StringValue(project.Status)
}

func (s *projectSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (s *projectSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *projectSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get project by name or display_name\n\n" +
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
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The project's display name",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("display_name"),
						path.MatchRoot("name"),
					}...),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The project's unique name ID",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("display_name"),
						path.MatchRoot("name"),
					}...),
				},
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The project's URL",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The project's status",
				Computed:            true,
			},
		},
	}
}

func (s *projectSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *projectSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var project *buddy.Project
	var err error
	domain := data.Domain.ValueString()
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		var httpResp *http.Response
		project, httpResp, err = s.client.ProjectService.Get(domain, data.Name.ValueString())
		if err != nil {
			if util.IsResourceNotFound(httpResp, err) {
				resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("project"))
				return
			}
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get project", err))
			return
		}
	} else {
		displayName := data.DisplayName.ValueString()
		var projects *buddy.Projects
		projects, _, err = s.client.ProjectService.GetListAll(domain, nil)
		if err != nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get projects", err))
			return
		}
		for _, p := range projects.Projects {
			if p.DisplayName == displayName {
				project = p
				break
			}
		}
		if project == nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("project"))
			return
		}
	}
	data.loadAPI(domain, project)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

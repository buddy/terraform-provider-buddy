package source

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
)

var (
	_ datasource.DataSource              = &projectsSource{}
	_ datasource.DataSourceWithConfigure = &projectsSource{}
)

func NewProjectsSource() datasource.DataSource {
	return &projectsSource{}
}

type projectsSource struct {
	client *buddy.Client
}

type projectsSourceModel struct {
	ID               types.String `tfsdk:"id"`
	Domain           types.String `tfsdk:"domain"`
	NameRegex        types.String `tfsdk:"name_regex"`
	DisplayNameRegex types.String `tfsdk:"display_name_regex"`
	Membership       types.Bool   `tfsdk:"membership"`
	Status           types.String `tfsdk:"status"`
	Projects         types.Set    `tfsdk:"projects"`
}

func (s *projectsSourceModel) loadAPI(ctx context.Context, domain string, projects *[]*buddy.Project) diag.Diagnostics {
	s.ID = types.StringValue(util.UniqueString())
	s.Domain = types.StringValue(domain)
	p, d := util.ProjectsModelFromApi(ctx, projects)
	s.Projects = p
	return d
}

func (s *projectsSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

func (s *projectsSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *projectsSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List projects and optionally filter them by membership, status, name or display name\n\n" +
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
			"name_regex": schema.StringAttribute{
				MarkdownDescription: "The project's name regular expression to match",
				Optional:            true,
				Validators: []validator.String{
					util.RegexpValidator(),
				},
			},
			"display_name_regex": schema.StringAttribute{
				MarkdownDescription: "The project's display name regular expression to match",
				Optional:            true,
				Validators: []validator.String{
					util.RegexpValidator(),
				},
			},
			"membership": schema.BoolAttribute{
				MarkdownDescription: "For workspace administrators all workspace projects are returned, set true to lists projects the user actually belongs to",
				Optional:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Filter projects by status (`ACTIVE`, `CLOSED`)",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						buddy.ProjectStatusActive,
						buddy.ProjectStatusClosed,
					),
				},
			},
			"projects": schema.SetNestedAttribute{
				MarkdownDescription: "List of projects",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.SourceProjectsModelAttributes(),
				},
			},
		},
	}
}

func (s *projectsSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *projectsSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	var nameRegex *regexp.Regexp
	var displayNameRegex *regexp.Regexp
	if !data.NameRegex.IsNull() && !data.NameRegex.IsUnknown() {
		nameRegex = regexp.MustCompile(data.NameRegex.ValueString())
	}
	if !data.DisplayNameRegex.IsNull() && !data.DisplayNameRegex.IsUnknown() {
		displayNameRegex = regexp.MustCompile(data.DisplayNameRegex.ValueString())
	}
	ops := buddy.ProjectListQuery{}
	if !data.Membership.IsNull() && !data.Membership.IsUnknown() {
		ops.Membership = data.Membership.ValueBool()
	}
	if !data.Status.IsNull() && !data.Status.IsUnknown() {
		ops.Status = data.Status.ValueString()
	}
	projects, _, err := s.client.ProjectService.GetListAll(domain, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get projects", err))
		return
	}
	var result []*buddy.Project
	for _, p := range projects.Projects {
		if nameRegex != nil && !nameRegex.MatchString(p.Name) {
			continue
		}
		if displayNameRegex != nil && !displayNameRegex.MatchString(p.DisplayName) {
			continue
		}
		result = append(result, p)
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, &result)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

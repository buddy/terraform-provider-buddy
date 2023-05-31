package source

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
)

var (
	_ datasource.DataSource              = &projectGroupsSource{}
	_ datasource.DataSourceWithConfigure = &projectGroupsSource{}
)

func NewProjectGroupsSource() datasource.DataSource {
	return &projectGroupsSource{}
}

type projectGroupsSource struct {
	client *buddy.Client
}

type projectGroupsSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Domain      types.String `tfsdk:"domain"`
	ProjectName types.String `tfsdk:"project_name"`
	NameRegex   types.String `tfsdk:"name_regex"`
	Groups      types.Set    `tfsdk:"groups"`
}

func (s *projectGroupsSourceModel) loadAPI(ctx context.Context, domain string, projectName string, groups *[]*buddy.Group) diag.Diagnostics {
	s.ID = types.StringValue(domain)
	s.Domain = types.StringValue(domain)
	s.ProjectName = types.StringValue(projectName)
	g, d := util.GroupsModelFromApi(ctx, groups)
	s.Groups = g
	return d
}

func (s *projectGroupsSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_groups"
}

func (s *projectGroupsSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *projectGroupsSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List project groups and optionally filter them by name\n\n" +
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
			"name_regex": schema.StringAttribute{
				MarkdownDescription: "The group's name regular expression to match",
				Optional:            true,
				Validators: []validator.String{
					util.RegexpValidator(),
				},
			},
			"groups": schema.SetNestedAttribute{
				MarkdownDescription: "List of groups",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.SourceGroupModelAttributes(),
				},
			},
		},
	}
}

func (s *projectGroupsSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *projectGroupsSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	projectName := data.ProjectName.ValueString()
	groups, _, err := s.client.ProjectGroupService.GetProjectGroups(domain, projectName)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get project groups", err))
		return
	}
	var nameRegex *regexp.Regexp
	if !data.NameRegex.IsNull() && !data.NameRegex.IsUnknown() {
		nameRegex = regexp.MustCompile(data.NameRegex.ValueString())
	}
	var result []*buddy.Group
	for _, g := range groups.Groups {
		if nameRegex != nil && !nameRegex.MatchString(g.Name) {
			continue
		}
		result = append(result, g)
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, &result)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

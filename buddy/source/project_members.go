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
	_ datasource.DataSource              = &projectMembersSource{}
	_ datasource.DataSourceWithConfigure = &projectMembersSource{}
)

func NewProjectMembersSource() datasource.DataSource {
	return &projectMembersSource{}
}

type projectMembersSource struct {
	client *buddy.Client
}

type projectMembersSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Domain      types.String `tfsdk:"domain"`
	ProjectName types.String `tfsdk:"project_name"`
	NameRegex   types.String `tfsdk:"name_regex"`
	Members     types.Set    `tfsdk:"members"`
}

func (s *projectMembersSourceModel) loadAPI(ctx context.Context, domain string, projectName string, members *[]*buddy.Member) diag.Diagnostics {
	s.ID = types.StringValue(util.UniqueString())
	s.Domain = types.StringValue(domain)
	s.ProjectName = types.StringValue(projectName)
	m, d := util.MembersModelFromApi(ctx, members)
	s.Members = m
	return d
}

func (s *projectMembersSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_members"
}

func (s *projectMembersSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *projectMembersSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List project members and optionally filter them by name\n\n" +
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
				MarkdownDescription: "The project member's name regular expression to match",
				Optional:            true,
				Validators: []validator.String{
					util.RegexpValidator(),
				},
			},
			"members": schema.SetNestedAttribute{
				MarkdownDescription: "List of members",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.SourceMemberModelAttributes(),
				},
			},
		},
	}
}

func (s *projectMembersSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *projectMembersSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	projectName := data.ProjectName.ValueString()
	members, _, err := s.client.ProjectMemberService.GetProjectMembersAll(domain, projectName)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get project members", err))
		return
	}
	var nameRegex *regexp.Regexp
	if !data.NameRegex.IsNull() && !data.NameRegex.IsUnknown() {
		nameRegex = regexp.MustCompile(data.NameRegex.ValueString())
	}
	var result []*buddy.Member
	for _, m := range members.Members {
		if nameRegex != nil && !nameRegex.MatchString(m.Name) {
			continue
		}
		result = append(result, m)
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, &result)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

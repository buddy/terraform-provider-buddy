package source

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ datasource.DataSource              = &sandboxesSource{}
	_ datasource.DataSourceWithConfigure = &sandboxesSource{}
)

type sandboxesSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Domain      types.String `tfsdk:"domain"`
	ProjectName types.String `tfsdk:"project_name"`
	NameRegex   types.String `tfsdk:"name_regex"`
	Sandboxes   types.Set    `tfsdk:"sandboxes"`
}

func (s *sandboxesSourceModel) loadAPI(ctx context.Context, domain string, projectName string, sandboxes *[]*buddy.Sandbox) diag.Diagnostics {
	s.ID = types.StringValue(util.UniqueString())
	s.Domain = types.StringValue(domain)
	s.ProjectName = types.StringValue(projectName)
	ss, d := util.SandboxesModelFromApi(ctx, sandboxes)
	s.Sandboxes = ss
	return d
}

type sandboxesSource struct {
	client *buddy.Client
}

func NewSandboxesSource() datasource.DataSource {
	return &sandboxesSource{}
}

func (s *sandboxesSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sandboxes"
}

func (s *sandboxesSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *sandboxesSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List sandboxes and optionally filter them by name\n\n" +
			"Token scopes required: `WORKSPACE`, `SANDBOX_INFO`",
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
				MarkdownDescription: "The sandbox's name regular expression to match",
				Optional:            true,
				Validators: []validator.String{
					util.RegexpValidator(),
				},
			},
			"sandboxes": schema.SetNestedAttribute{
				MarkdownDescription: "List of sandboxes",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.SourceSandboxModelAttributes(),
				},
			},
		},
	}
}

func (s *sandboxesSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *sandboxesSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	projectName := data.ProjectName.ValueString()
	var nameRegex *regexp.Regexp
	if !data.NameRegex.IsNull() && !data.NameRegex.IsUnknown() {
		nameRegex = regexp.MustCompile(data.NameRegex.ValueString())
	}
	sandboxes, _, err := s.client.SandboxService.GetList(domain, buddy.Query{
		ProjectName: &projectName,
	})
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get sandboxes", err))
		return
	}
	var result []*buddy.Sandbox
	for _, s := range sandboxes.Sandboxes {
		if nameRegex != nil && !nameRegex.MatchString(s.Name) {
			continue
		}
		result = append(result, s)
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, &result)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

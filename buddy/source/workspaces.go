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
	_ datasource.DataSource              = &workspacesSource{}
	_ datasource.DataSourceWithConfigure = &workspacesSource{}
)

func NewWorkspacesSource() datasource.DataSource {
	return &workspacesSource{}
}

type workspacesSource struct {
	client *buddy.Client
}

type workspacesSourceModel struct {
	ID          types.String `tfsdk:"id"`
	DomainRegex types.String `tfsdk:"domain_regex"`
	NameRegex   types.String `tfsdk:"name_regex"`
	Workspaces  types.Set    `tfsdk:"workspaces"`
}

func (s *workspacesSourceModel) loadAPI(ctx context.Context, workspaces *[]*buddy.Workspace) diag.Diagnostics {
	s.ID = types.StringValue(util.UniqueString())
	w, d := util.WorkspacesModelFromApi(ctx, workspaces)
	s.Workspaces = w
	return d
}

func (s *workspacesSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspaces"
}

func (s *workspacesSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *workspacesSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List workspaces and optionally filter them by name or URL handle\n\n" +
			"Token scope required: `WORKSPACE`",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Terraform resource identifier for this item",
				Computed:            true,
			},
			"domain_regex": schema.StringAttribute{
				MarkdownDescription: "The workspace URL handle regular expression to match",
				Optional:            true,
				Validators: []validator.String{
					util.RegexpValidator(),
				},
			},
			"name_regex": schema.StringAttribute{
				MarkdownDescription: "The workspace name regular expression to match",
				Optional:            true,
				Validators: []validator.String{
					util.RegexpValidator(),
				},
			},
			"workspaces": schema.SetNestedAttribute{
				MarkdownDescription: "List of workspaces",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.SourceWorkspaceModelAttributes(),
				},
			},
		},
	}
}

func (s *workspacesSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *workspacesSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var nameRegex *regexp.Regexp
	var domainRegex *regexp.Regexp
	workspaces, _, err := s.client.WorkspaceService.GetList()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get workspaces", err))
		return
	}
	var result []*buddy.Workspace
	if !data.NameRegex.IsNull() && !data.NameRegex.IsUnknown() {
		nameRegex = regexp.MustCompile(data.NameRegex.ValueString())
	}
	if !data.DomainRegex.IsNull() && !data.DomainRegex.IsUnknown() {
		domainRegex = regexp.MustCompile(data.DomainRegex.ValueString())
	}
	for _, w := range workspaces.Workspaces {
		if nameRegex != nil && !nameRegex.MatchString(w.Name) {
			continue
		}
		if domainRegex != nil && !domainRegex.MatchString(w.Domain) {
			continue
		}
		result = append(result, w)
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, &result)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

package source

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ datasource.DataSource              = &environmentSource{}
	_ datasource.DataSourceWithConfigure = &environmentSource{}
)

func NewEnvironmentSource() datasource.DataSource {
	return &environmentSource{}
}

type environmentSource struct {
	client *buddy.Client
}

type environmentSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Domain        types.String `tfsdk:"domain"`
	ProjectName   types.String `tfsdk:"project_name"`
	HtmlUrl       types.String `tfsdk:"html_url"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	Name          types.String `tfsdk:"name"`
	Identifier    types.String `tfsdk:"identifier"`
	Type          types.String `tfsdk:"type"`
	Tags          types.Set    `tfsdk:"tags"`
	PublicUrl     types.String `tfsdk:"public_url"`
}

func (e *environmentSourceModel) loadAPI(ctx context.Context, domain string, projectName string, environment *buddy.Environment) diag.Diagnostics {
	e.ID = types.StringValue(util.ComposeTripleId(domain, projectName, environment.Id))
	e.Domain = types.StringValue(domain)
	e.ProjectName = types.StringValue(projectName)
	e.HtmlUrl = types.StringValue(environment.HtmlUrl)
	e.EnvironmentId = types.StringValue(environment.Id)
	e.Name = types.StringValue(environment.Name)
	e.Identifier = types.StringValue(environment.Identifier)
	e.Type = types.StringValue(environment.Type)
	e.PublicUrl = types.StringValue(environment.PublicUrl)
	t, d := types.SetValueFrom(ctx, types.StringType, &environment.Tags)
	e.Tags = t
	return d
}

func (s *environmentSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

func (s *environmentSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *environmentSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get environment by name or environment ID\n\n" +
			"Token scope required: `WORKSPACE`, `ENVIRONMENT_INFO`",
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
			"name": schema.StringAttribute{
				MarkdownDescription: "The environment's name",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("name"),
						path.MatchRoot("environment_id"),
					}...),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "The environment's ID",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("name"),
						path.MatchRoot("environment_id"),
					}...),
				},
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The environment's URL",
				Computed:            true,
			},
			"identifier": schema.StringAttribute{
				MarkdownDescription: "The environment's identifier",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The environment's type",
				Computed:            true,
			},
			"public_url": schema.StringAttribute{
				MarkdownDescription: "The environment's public URL",
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The environment's list of tags",
				Computed:            true,
			},
		},
	}
}

func (s *environmentSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *environmentSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var environment *buddy.Environment
	var err error
	domain := data.Domain.ValueString()
	projectName := data.ProjectName.ValueString()
	if !data.EnvironmentId.IsNull() && !data.EnvironmentId.IsUnknown() {
		var httpResp *http.Response
		environmentId := data.EnvironmentId.ValueString()
		environment, httpResp, err = s.client.EnvironmentService.Get(domain, projectName, environmentId)
		if err != nil {
			if util.IsResourceNotFound(httpResp, err) {
				resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("environment"))
				return
			}
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get environment", err))
			return
		}
	} else {
		name := data.Name.ValueString()
		var environments *buddy.Environments
		environments, _, err = s.client.EnvironmentService.GetList(domain, projectName)
		if err != nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get environments", err))
			return
		}
		for _, e := range environments.Environments {
			if e.Name == name {
				environment = e
				break
			}
		}
		if environment == nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("environment"))
			return
		}
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, environment)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

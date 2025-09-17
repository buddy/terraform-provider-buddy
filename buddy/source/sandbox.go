package source

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ datasource.DataSource              = &sandboxSource{}
	_ datasource.DataSourceWithConfigure = &sandboxSource{}
)

func NewSandboxSource() datasource.DataSource {
	return &sandboxSource{}
}

type sandboxSource struct {
	client *buddy.Client
}

type sandboxSourceModel struct {
	ID         types.String `tfsdk:"id"`
	Domain     types.String `tfsdk:"domain"`
	SandboxId  types.String `tfsdk:"sandbox_id"`
	Name       types.String `tfsdk:"name"`
	Status     types.String `tfsdk:"status"`
	Identifier types.String `tfsdk:"identifier"`
	HtmlUrl    types.String `tfsdk:"html_url"`
}

func (s *sandboxSourceModel) loadAPI(domain string, sandbox *buddy.Sandbox) {
	s.ID = types.StringValue(util.ComposeDoubleId(domain, sandbox.Id))
	s.Domain = types.StringValue(domain)
	s.SandboxId = types.StringValue(sandbox.Id)
	s.Name = types.StringValue(sandbox.Name)
	s.Status = types.StringValue(sandbox.Status)
	s.Identifier = types.StringValue(sandbox.Identifier)
	s.HtmlUrl = types.StringValue(sandbox.HtmlUrl)
}

func (s *sandboxSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sandbox"
}

func (s *sandboxSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *sandboxSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get sandbox by sandbox ID\n\n" +
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
			"sandbox_id": schema.StringAttribute{
				MarkdownDescription: "The sandbox's ID",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The sandbox's name",
				Computed:            true,
			},
			"identifier": schema.StringAttribute{
				MarkdownDescription: "The sandbox's identifier",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The sandbox's status",
				Computed:            true,
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The sandbox's URL",
				Computed:            true,
			},
		},
	}
}

func (s *sandboxSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *sandboxSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	sandboxId := data.SandboxId.ValueString()
	sandbox, httpRes, err := s.client.SandboxService.Get(domain, sandboxId)
	if err != nil {
		if util.IsResourceNotFound(httpRes, err) {
			resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("sandbox"))
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get sandbox", err))
		return
	}
	data.loadAPI(domain, sandbox)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

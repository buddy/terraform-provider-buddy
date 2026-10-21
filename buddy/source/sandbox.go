package source

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	ID                   types.String `tfsdk:"id"`
	Domain               types.String `tfsdk:"domain"`
	SandboxId            types.String `tfsdk:"sandbox_id"`
	Name                 types.String `tfsdk:"name"`
	Note                 types.String `tfsdk:"note"`
	AgentNote            types.String `tfsdk:"agent_note"`
	Status               types.String `tfsdk:"status"`
	Identifier           types.String `tfsdk:"identifier"`
	HtmlUrl              types.String `tfsdk:"html_url"`
	Scope                types.String `tfsdk:"scope"`
	SshHost              types.String `tfsdk:"ssh_host"`
	SshPort              types.Int32  `tfsdk:"ssh_port"`
	CreateDate           types.String `tfsdk:"create_date"`
	Fetches              types.List   `tfsdk:"fetches"`
	PipelinesAccessLevel types.String `tfsdk:"pipelines_access_level"`
	AllowedPipelines     types.Set    `tfsdk:"allowed_pipelines"`
}

func (s *sandboxSourceModel) loadAPI(ctx context.Context, domain string, sandbox *buddy.Sandbox) diag.Diagnostics {
	s.ID = types.StringValue(util.ComposeDoubleId(domain, sandbox.Id))
	s.Domain = types.StringValue(domain)
	s.SandboxId = types.StringValue(sandbox.Id)
	s.Name = types.StringValue(sandbox.Name)
	s.Note = types.StringValue(sandbox.Note)
	s.AgentNote = types.StringValue(sandbox.AgentNote)
	s.Status = types.StringValue(sandbox.Status)
	s.Identifier = types.StringValue(sandbox.Identifier)
	s.HtmlUrl = types.StringValue(sandbox.HtmlUrl)
	s.Scope = types.StringValue(sandbox.Scope)
	s.SshHost = types.StringValue(sandbox.SshHost)
	s.SshPort = types.Int32Value(int32(sandbox.SshPort))
	s.CreateDate = types.StringValue(sandbox.CreateDate)
	s.PipelinesAccessLevel = types.StringValue(sandbox.PipelinesAccessLevel)
	fetches, diags := util.SandboxFetchesFromApi(ctx, &sandbox.Fetch)
	s.Fetches = fetches
	pipelines, d := util.SandboxPipelinesFromApi(ctx, &sandbox.AllowedPipelines)
	diags.Append(d...)
	s.AllowedPipelines = pipelines
	return diags
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
			"Token scope required: `SANDBOX_READ`",
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
			"note": schema.StringAttribute{
				MarkdownDescription: "The sandbox's note",
				Computed:            true,
			},
			"agent_note": schema.StringAttribute{
				MarkdownDescription: "The sandbox's agent note",
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
			"scope": schema.StringAttribute{
				MarkdownDescription: "The sandbox's scope: `PROJECT`, `ENVIRONMENT` or `WORKSPACE`",
				Computed:            true,
			},
			"ssh_host": schema.StringAttribute{
				MarkdownDescription: "The sandbox's SSH host",
				Computed:            true,
			},
			"ssh_port": schema.Int32Attribute{
				MarkdownDescription: "The sandbox's SSH port",
				Computed:            true,
			},
			"create_date": schema.StringAttribute{
				MarkdownDescription: "The sandbox's creation date",
				Computed:            true,
			},
			"fetches": schema.ListNestedAttribute{
				MarkdownDescription: "The sources fetched into the sandbox",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.SourceSandboxFetchModelAttributes(),
				},
			},
			"pipelines_access_level": schema.StringAttribute{
				MarkdownDescription: "Indicates if all pipelines are allowed to use this sandbox: `USE_ONLY` or `DENIED`, empty for sandboxes in an environment",
				Computed:            true,
			},
			"allowed_pipelines": schema.SetNestedAttribute{
				MarkdownDescription: "List of specific pipelines allowed to use this sandbox",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.SourceSandboxPipelineModelAttributes(),
				},
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
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, sandbox)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

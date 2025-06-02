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
	_ datasource.DataSource              = &targetSource{}
	_ datasource.DataSourceWithConfigure = &targetSource{}
)

func NewTargetSource() datasource.DataSource {
	return &targetSource{}
}

type targetSource struct {
	client *buddy.Client
}

type targetSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Domain      types.String `tfsdk:"domain"`
	TargetId    types.String `tfsdk:"target_id"`
	HtmlUrl     types.String `tfsdk:"html_url"`
	Name        types.String `tfsdk:"name"`
	Identifier  types.String `tfsdk:"identifier"`
	Tags        types.Set    `tfsdk:"tags"`
	Type        types.String `tfsdk:"type"`
	Host        types.String `tfsdk:"host"`
	Scope       types.String `tfsdk:"scope"`
	Repository  types.String `tfsdk:"repository"`
	Port        types.String `tfsdk:"port"`
	Integration types.String `tfsdk:"integration"`
	Path        types.String `tfsdk:"path"`
	Secure      types.Bool   `tfsdk:"secure"`
	Disabled    types.Bool   `tfsdk:"disabled"`
}

func (m *targetSourceModel) loadAPI(ctx context.Context, domain string, target *buddy.Target) diag.Diagnostics {
	m.ID = types.StringValue(util.ComposeDoubleId(domain, target.Id))
	m.Domain = types.StringValue(domain)
	m.HtmlUrl = types.StringValue(target.HtmlUrl)
	m.TargetId = types.StringValue(target.Id)
	m.Identifier = types.StringValue(target.Identifier)
	tags, d := types.SetValueFrom(ctx, types.StringType, &target.Tags)
	m.Tags = tags
	m.Name = types.StringValue(target.Name)
	m.Type = types.StringValue(target.Type)
	m.Host = types.StringValue(target.Host)
	m.Scope = types.StringValue(target.Scope)
	m.Repository = types.StringValue(target.Repository)
	m.Port = types.StringValue(target.Port)
	m.Path = types.StringValue(target.Path)
	m.Secure = types.BoolValue(target.Secure)
	m.Disabled = types.BoolValue(target.Disabled)
	m.Integration = types.StringValue(target.Integration)
	return d
}

func (s *targetSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_target"
}

func (s *targetSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get target by name or target ID\n\n" +
			"Token scope required: `WORKSPACE`, TARGET_INFO",
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
			"target_id": schema.StringAttribute{
				MarkdownDescription: "The target's ID",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The target's name",
				Computed:            true,
			},
			"identifier": schema.StringAttribute{
				MarkdownDescription: "The target's identifier",
				Computed:            true,
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The target's URL",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The target's type",
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "The target's list of tags",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: "The target's scope",
				Computed:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "The target's host",
				Computed:            true,
			},
			"repository": schema.StringAttribute{
				MarkdownDescription: "The target's repository",
				Computed:            true,
			},
			"port": schema.StringAttribute{
				MarkdownDescription: "The target's port",
				Computed:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "The target's path",
				Computed:            true,
			},
			"secure": schema.BoolAttribute{
				MarkdownDescription: "The target's secure setting",
				Computed:            true,
			},
			"integration": schema.StringAttribute{
				MarkdownDescription: "The target's integration",
				Computed:            true,
			},
			"disabled": schema.BoolAttribute{
				MarkdownDescription: "Defines whether or not the target can be run",
				Computed:            true,
			},
		},
	}
}

func (s *targetSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *targetSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *targetSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	targetId := data.TargetId.ValueString()
	target, _, err := s.client.TargetService.Get(domain, targetId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get target by id", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, target)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

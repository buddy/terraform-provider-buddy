package source

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
	ID           types.String `tfsdk:"id"`
	Domain       types.String `tfsdk:"domain"`
	TargetId     types.String `tfsdk:"target_id"`
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	Host         types.String `tfsdk:"host"`
	Port         types.String `tfsdk:"port"`
	Path         types.String `tfsdk:"path"`
	Secure       types.Bool   `tfsdk:"secure"`
	Scope        types.String `tfsdk:"scope"`
	Repository   types.String `tfsdk:"repository"`
	Integration  types.String `tfsdk:"integration"`
	Tags         types.Set    `tfsdk:"tags"`
	Disabled     types.Bool   `tfsdk:"disabled"`
	AuthMethod   types.String `tfsdk:"auth_method"`
	AuthUsername types.String `tfsdk:"auth_username"`
	HtmlUrl      types.String `tfsdk:"html_url"`
}

func (s *targetSourceModel) loadAPI(ctx context.Context, domain string, target *buddy.Target) diag.Diagnostics {
	var diags diag.Diagnostics

	s.ID = types.StringValue(util.ComposeDoubleId(domain, target.Id))
	s.Domain = types.StringValue(domain)
	s.TargetId = types.StringValue(target.Id)
	s.Name = types.StringValue(target.Name)
	s.Type = types.StringValue(target.Type)
	s.HtmlUrl = types.StringValue(target.HtmlUrl)
	s.Scope = types.StringValue(target.Scope)
	s.Disabled = types.BoolValue(target.Disabled)
	s.Secure = types.BoolValue(target.Secure)

	if target.Host != "" {
		s.Host = types.StringValue(target.Host)
	} else {
		s.Host = types.StringNull()
	}

	if target.Port != "" {
		s.Port = types.StringValue(target.Port)
	} else {
		s.Port = types.StringNull()
	}

	if target.Path != "" {
		s.Path = types.StringValue(target.Path)
	} else {
		s.Path = types.StringNull()
	}

	if target.Repository != "" {
		s.Repository = types.StringValue(target.Repository)
	} else {
		s.Repository = types.StringNull()
	}

	if target.Integration != "" {
		s.Integration = types.StringValue(target.Integration)
	} else {
		s.Integration = types.StringNull()
	}

	tags, d := types.SetValueFrom(ctx, types.StringType, &target.Tags)
	diags.Append(d...)
	s.Tags = tags

	// Handle auth fields if auth is present
	if target.Auth != nil {
		s.AuthMethod = types.StringValue(target.Auth.Method)

		if target.Auth.Username != "" {
			s.AuthUsername = types.StringValue(target.Auth.Username)
		} else {
			s.AuthUsername = types.StringNull()
		}
	} else {
		s.AuthMethod = types.StringNull()
		s.AuthUsername = types.StringNull()
	}

	return diags
}

func (s *targetSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_target"
}

func (s *targetSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *targetSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get target by target ID or name\n\n" +
			"Token scope required: `WORKSPACE`",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Terraform resource identifier for this item",
				Computed:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The workspace's URL handle",
				Required:            true,
			},
			"target_id": schema.StringAttribute{
				MarkdownDescription: "The target's ID",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The target's name",
				Optional:            true,
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The target's type",
				Computed:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "The target's hostname or IP address",
				Computed:            true,
			},
			"port": schema.StringAttribute{
				MarkdownDescription: "The target's port",
				Computed:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "The remote path on the target",
				Computed:            true,
			},
			"secure": schema.BoolAttribute{
				MarkdownDescription: "Whether to use secure connection",
				Computed:            true,
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: "The target's scope",
				Computed:            true,
			},
			"repository": schema.StringAttribute{
				MarkdownDescription: "The repository for registry targets",
				Computed:            true,
			},
			"integration": schema.StringAttribute{
				MarkdownDescription: "The integration ID used for cloud targets",
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "The target's tags",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"disabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the target is disabled",
				Computed:            true,
			},
			"auth_method": schema.StringAttribute{
				MarkdownDescription: "The authentication method",
				Computed:            true,
			},
			"auth_username": schema.StringAttribute{
				MarkdownDescription: "The authentication username",
				Computed:            true,
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The target's URL",
				Computed:            true,
			},
		},
	}
}

func (s *targetSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("target_id"),
			path.MatchRoot("name"),
		),
		datasourcevalidator.AtLeastOneOf(
			path.MatchRoot("target_id"),
			path.MatchRoot("name"),
		),
	}
}

func (s *targetSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data targetSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := data.Domain.ValueString()
	var target *buddy.Target
	var err error

	if !data.TargetId.IsNull() {
		target, _, err = s.client.TargetService.Get(domain, data.TargetId.ValueString())
	} else {
		targets, _, err := s.client.TargetService.GetList(domain, &buddy.TargetGetListQuery{})
		if err == nil {
			if found := util.FilterTargetByName(targets.Targets, data.Name.ValueString()); found != nil {
				target = found
			} else {
				resp.Diagnostics.AddError("Target not found", "Target with name '"+data.Name.ValueString()+"' not found")
				return
			}
		}
	}

	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get target", err))
		return
	}

	resp.Diagnostics.Append(data.loadAPI(ctx, domain, target)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

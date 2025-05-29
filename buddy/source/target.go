package source

import (
	"context"
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	Domain     types.String `tfsdk:"domain"`
	TargetId   types.String `tfsdk:"target_id"`
	Name       types.String `tfsdk:"name"`
	Identifier types.String `tfsdk:"identifier"`
	Tags       types.List   `tfsdk:"tags"`
	Type       types.String `tfsdk:"type"`
	HtmlUrl    types.String `tfsdk:"html_url"`
	Host       types.String `tfsdk:"host"`
	Scope      types.String `tfsdk:"scope"`
	Repository types.String `tfsdk:"repository"`
	Port       types.String `tfsdk:"port"`
	Path       types.String `tfsdk:"path"`
	Secure     types.Bool   `tfsdk:"secure"`
	Disabled   types.Bool   `tfsdk:"disabled"`
}

func (s *targetSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_target"
}

func (s *targetSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get target by target ID or Identifier\n\n" +
			"Token scope required: `WORKSPACE`",
		Attributes: map[string]schema.Attribute{
			"domain": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"target_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("target_id"),
						path.MatchRoot("identifier"),
					}...),
				},
			},
			"identifier": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"tags": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"type": schema.StringAttribute{
				Computed: true,
			},
			"html_url": schema.StringAttribute{
				Computed: true,
			},
			"host": schema.StringAttribute{
				Computed: true,
			},
			"scope": schema.StringAttribute{
				Computed: true,
			},
			"repository": schema.StringAttribute{
				Computed: true,
			},
			"port": schema.StringAttribute{
				Computed: true,
			},
			"path": schema.StringAttribute{
				Computed: true,
			},
			"secure": schema.BoolAttribute{
				Computed: true,
			},
			"disabled": schema.BoolAttribute{
				Computed: true,
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

	var target *buddy.Target
	var err error

	// Determine whether to search by ID or identifier
	if !data.TargetId.IsNull() && data.TargetId.ValueString() != "" {
		// Get by ID
		target, _, err = s.client.TargetService.Get(domain, data.TargetId.ValueString())
		if err != nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get target by id", err))
			return
		}
	} else if !data.Identifier.IsNull() && data.Identifier.ValueString() != "" {
		// Get all targets and find by identifier
		targets, _, err := s.client.TargetService.GetList(domain, nil)
		if err != nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get targets", err))
			return
		}

		identifier := data.Identifier.ValueString()
		for _, t := range targets.Targets {
			if t.Identifier == identifier {
				target = t
				break
			}
		}

		if target == nil {
			resp.Diagnostics.AddError(
				"target not found",
				"target with identifier "+identifier+" not found",
			)
			return
		}
	}

	// Load data from API response
	if err := data.loadAPI(ctx, domain, target); err != nil {
		resp.Diagnostics.AddError("failed to load target data", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m *targetSourceModel) loadAPI(ctx context.Context, domain string, target *buddy.Target) error {
	m.Domain = types.StringValue(domain)
	m.TargetId = types.StringValue(target.Id)
	m.Identifier = types.StringValue(target.Identifier)
	m.Name = types.StringValue(target.Name)
	m.Type = types.StringValue(target.Type)
	m.HtmlUrl = types.StringValue(target.HtmlUrl)
	m.Scope = types.StringValue(target.Scope)

	if target.Tags != nil && len(target.Tags) > 0 {
		tags, diags := types.ListValueFrom(ctx, types.StringType, target.Tags)
		if diags.HasError() {
			return fmt.Errorf("failed to convert tags: %s", diags[0].Summary())
		}
		m.Tags = tags
	} else {
		m.Tags = types.ListNull(types.StringType)
	}

	m.Host = types.StringValue(target.Host)
	m.Repository = types.StringValue(target.Repository)
	m.Port = types.StringValue(target.Port)
	m.Path = types.StringValue(target.Path)
	m.Secure = types.BoolValue(target.Secure)
	m.Disabled = types.BoolValue(target.Disabled)

	return nil
}

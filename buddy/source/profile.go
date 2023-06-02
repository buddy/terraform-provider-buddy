package source

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ datasource.DataSource              = &profileSource{}
	_ datasource.DataSourceWithConfigure = &profileSource{}
)

func NewProfileSource() datasource.DataSource {
	return &profileSource{}
}

type profileSource struct {
	client *buddy.Client
}

type profileSourceModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	MemberId  types.Int64  `tfsdk:"member_id"`
	AvatarUrl types.String `tfsdk:"avatar_url"`
	HtmlUrl   types.String `tfsdk:"html_url"`
}

func (s *profileSourceModel) loadAPI(profile *buddy.Profile) {
	s.Id = types.StringValue(strconv.Itoa(profile.Id))
	s.Name = types.StringValue(profile.Name)
	s.MemberId = types.Int64Value(int64(profile.Id))
	s.AvatarUrl = types.StringValue(profile.AvatarUrl)
	s.HtmlUrl = types.StringValue(profile.HtmlUrl)
}

func (s *profileSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_profile"
}

func (s *profileSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *profileSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get details of a Buddy's user profile\n\n" +
			"Token scope required: `USER_INFO`",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Terraform resource identifier for this item",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The user's name",
				Computed:            true,
			},
			"member_id": schema.Int64Attribute{
				MarkdownDescription: "The user's ID",
				Computed:            true,
			},
			"avatar_url": schema.StringAttribute{
				MarkdownDescription: "The user's avatar URL",
				Computed:            true,
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The user's profile URL",
				Computed:            true,
			},
		},
	}
}

func (s *profileSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *profileSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	p, _, err := s.client.ProfileService.Get()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get profile", err))
		return
	}
	data.loadAPI(p)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

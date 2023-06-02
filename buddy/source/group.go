package source

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
	"strconv"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ datasource.DataSource              = &groupSource{}
	_ datasource.DataSourceWithConfigure = &groupSource{}
)

func NewGroupSource() datasource.DataSource {
	return &groupSource{}
}

type groupSource struct {
	client *buddy.Client
}

type groupSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Domain      types.String `tfsdk:"domain"`
	Name        types.String `tfsdk:"name"`
	GroupId     types.Int64  `tfsdk:"group_id"`
	HtmlUrl     types.String `tfsdk:"html_url"`
	Description types.String `tfsdk:"description"`
}

func (s *groupSourceModel) loadAPI(domain string, group *buddy.Group) {
	s.ID = types.StringValue(util.ComposeDoubleId(domain, strconv.Itoa(group.Id)))
	s.Domain = types.StringValue(domain)
	s.Name = types.StringValue(group.Name)
	s.GroupId = types.Int64Value(int64(group.Id))
	s.HtmlUrl = types.StringValue(group.HtmlUrl)
	s.Description = types.StringValue(group.Description)
}

func (s *groupSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (s *groupSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *groupSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get group by name or group ID\n\n" +
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
			"name": schema.StringAttribute{
				MarkdownDescription: "The group's name",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("name"),
						path.MatchRoot("group_id"),
					}...),
				},
			},
			"group_id": schema.Int64Attribute{
				MarkdownDescription: "The group's ID",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("name"),
						path.MatchRoot("group_id"),
					}...),
				},
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The group's URL",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The group's description",
				Computed:            true,
			},
		},
	}
}

func (s *groupSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *groupSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var group *buddy.Group
	var err error
	domain := data.Domain.ValueString()
	if !data.GroupId.IsNull() && !data.GroupId.IsUnknown() {
		var httpResp *http.Response
		groupId := int(data.GroupId.ValueInt64())
		group, httpResp, err = s.client.GroupService.Get(domain, groupId)
		if err != nil {
			if util.IsResourceNotFound(httpResp, err) {
				resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("group"))
				return
			}
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get group", err))
			return
		}
	} else {
		name := data.Name.ValueString()
		var groups *buddy.Groups
		groups, _, err = s.client.GroupService.GetList(domain)
		if err != nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get groups", err))
			return
		}
		for _, g := range groups.Groups {
			if g.Name == name {
				group = g
				break
			}
		}
		if group == nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("group"))
			return
		}
	}
	data.loadAPI(domain, group)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

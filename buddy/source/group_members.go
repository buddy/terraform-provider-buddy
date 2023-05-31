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
	"strconv"
)

var (
	_ datasource.DataSource              = &groupMembersSource{}
	_ datasource.DataSourceWithConfigure = &groupMembersSource{}
)

func NewGroupMembersSource() datasource.DataSource {
	return &groupMembersSource{}
}

type groupMembersSource struct {
	client *buddy.Client
}

type groupMembersSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Domain    types.String `tfsdk:"domain"`
	GroupId   types.Int64  `tfsdk:"group_id"`
	NameRegex types.String `tfsdk:"name_regex"`
	Members   types.Set    `tfsdk:"members"`
}

func (s *groupMembersSourceModel) loadAPI(ctx context.Context, domain string, groupId int, members *[]*buddy.Member) diag.Diagnostics {
	s.ID = types.StringValue(util.ComposeDoubleId(domain, strconv.Itoa(groupId)))
	s.Domain = types.StringValue(domain)
	s.GroupId = types.Int64Value(int64(groupId))
	list, d := util.MembersModelFromApi(ctx, members)
	s.Members = list
	return d
}

func (s *groupMembersSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_members"
}

func (s *groupMembersSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *groupMembersSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List members of a group and optionally filter them by name\n\n" +
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
			"group_id": schema.Int64Attribute{
				MarkdownDescription: "The group's ID",
				Required:            true,
			},
			"name_regex": schema.StringAttribute{
				MarkdownDescription: "The member's name regular expression to match",
				Optional:            true,
				Validators: []validator.String{
					util.RegexpValidator(),
				},
			},
			"members": schema.SetNestedAttribute{
				MarkdownDescription: "List of members",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.SourceMemberModelAttributes(),
				},
			},
		},
	}
}

func (s *groupMembersSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *groupMembersSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	groupId := int(data.GroupId.ValueInt64())
	members, _, err := s.client.GroupService.GetGroupMembers(domain, groupId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get group members", err))
		return
	}
	var nameRegex *regexp.Regexp
	if !data.NameRegex.IsNull() && !data.NameRegex.IsUnknown() {
		nameRegex = regexp.MustCompile(data.NameRegex.ValueString())
	}
	var result []*buddy.Member
	for _, m := range members.Members {
		if nameRegex != nil && !nameRegex.MatchString(m.Name) {
			continue
		}
		result = append(result, m)
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, groupId, &result)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

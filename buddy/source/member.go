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
	_ datasource.DataSource              = &memberSource{}
	_ datasource.DataSourceWithConfigure = &memberSource{}
)

func NewMemberSource() datasource.DataSource {
	return &memberSource{}
}

type memberSource struct {
	client *buddy.Client
}

type memberSourceModel struct {
	ID             types.String `tfsdk:"id"`
	Domain         types.String `tfsdk:"domain"`
	Email          types.String `tfsdk:"email"`
	Name           types.String `tfsdk:"name"`
	MemberId       types.Int64  `tfsdk:"member_id"`
	Admin          types.Bool   `tfsdk:"admin"`
	HtmlUrl        types.String `tfsdk:"html_url"`
	AvatarUrl      types.String `tfsdk:"avatar_url"`
	WorkspaceOwner types.Bool   `tfsdk:"workspace_owner"`
}

func (s *memberSourceModel) loadAPI(domain string, member *buddy.Member) {
	s.ID = types.StringValue(util.ComposeDoubleId(domain, strconv.Itoa(member.Id)))
	s.Domain = types.StringValue(domain)
	s.Email = types.StringValue(member.Email)
	s.Name = types.StringValue(member.Name)
	s.MemberId = types.Int64Value(int64(member.Id))
	s.Admin = types.BoolValue(member.Admin)
	s.HtmlUrl = types.StringValue(member.HtmlUrl)
	s.AvatarUrl = types.StringValue(member.AvatarUrl)
	s.WorkspaceOwner = types.BoolValue(member.WorkspaceOwner)
}

func (s *memberSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_member"
}

func (s *memberSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *memberSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get member by name, email or member ID\n\n" +
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
			"email": schema.StringAttribute{
				MarkdownDescription: "The member's email",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("email"),
						path.MatchRoot("name"),
						path.MatchRoot("member_id"),
					}...),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The member's name",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("email"),
						path.MatchRoot("name"),
						path.MatchRoot("member_id"),
					}...),
				},
			},
			"member_id": schema.Int64Attribute{
				MarkdownDescription: "The member's ID",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("email"),
						path.MatchRoot("name"),
						path.MatchRoot("member_id"),
					}...),
				},
			},
			"admin": schema.BoolAttribute{
				MarkdownDescription: "Is the member a workspace administrator",
				Computed:            true,
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The member's URL",
				Computed:            true,
			},
			"avatar_url": schema.StringAttribute{
				MarkdownDescription: "The member's avatar URL",
				Computed:            true,
			},
			"workspace_owner": schema.BoolAttribute{
				MarkdownDescription: "Is the member the workspace owner",
				Computed:            true,
			},
		},
	}
}

func (s *memberSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *memberSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	var member *buddy.Member
	var err error
	if !data.MemberId.IsNull() && !data.MemberId.IsUnknown() {
		var httpResp *http.Response
		memberId := int(data.MemberId.ValueInt64())
		member, httpResp, err = s.client.MemberService.Get(domain, memberId)
		if err != nil {
			if util.IsResourceNotFound(httpResp, err) {
				resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("member"))
				return
			}
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get member", err))
			return
		}
	} else {
		var name *string
		var email *string
		var members *buddy.Members
		if !data.Name.IsNull() && !data.Name.IsUnknown() {
			name = data.Name.ValueStringPointer()
		}
		if !data.Email.IsNull() && !data.Email.IsUnknown() {
			email = data.Email.ValueStringPointer()
		}
		members, _, err = s.client.MemberService.GetListAll(domain)
		if err != nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get membmers", err))
			return
		}
		for _, m := range members.Members {
			if name != nil && *name == m.Name {
				member = m
				break
			}
			if email != nil && *email == m.Email {
				member = m
				break
			}
		}
		if member == nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("member"))
			return
		}
	}
	data.loadAPI(domain, member)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

package source

import (
	"buddy-terraform/buddy/util"
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
)

var (
	_ datasource.DataSource              = &permissionsSource{}
	_ datasource.DataSourceWithConfigure = &permissionsSource{}
)

func NewPermissionsSource() datasource.DataSource {
	return &permissionsSource{}
}

type permissionsSource struct {
	client *buddy.Client
}

type permissionsSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Domain      types.String `tfsdk:"domain"`
	NameRegex   types.String `tfsdk:"name_regex"`
	Type        types.String `tfsdk:"type"`
	Permissions types.Set    `tfsdk:"permissions"`
}

func (s *permissionsSourceModel) loadAPI(ctx context.Context, domain string, permissions *[]*buddy.Permission) diag.Diagnostics {
	s.ID = types.StringValue(util.UniqueString())
	s.Domain = types.StringValue(domain)
	p, d := util.PermissionsModelFromApi(ctx, permissions)
	s.Permissions = p
	return d
}

func (s *permissionsSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permissions"
}

func (s *permissionsSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *permissionsSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List permissions (roles) and optionally filter them by name or type\n\n" +
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
			"name_regex": schema.StringAttribute{
				MarkdownDescription: "The permission's name regular expression to match",
				Optional:            true,
				Validators: []validator.String{
					util.RegexpValidator(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Filter permissions by type (`CUSTOM`, `READ_ONLY`, `DEVELOPER`, `PROJECT_MANAGER`)",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						buddy.PermissionTypeCustom,
						buddy.PermissionTypeReadOnly,
						buddy.PermissionTypeDeveloper,
						buddy.PermissionTypeProjectManager,
					),
				},
			},
			"permissions": schema.SetNestedAttribute{
				MarkdownDescription: "List of permissions (roles)",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.SourcePermissionModelAttributes(),
				},
			},
		},
	}
}

func (s *permissionsSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *permissionsSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	var nameRegex *regexp.Regexp
	var typ *string
	if !data.NameRegex.IsNull() && !data.NameRegex.IsUnknown() {
		nameRegex = regexp.MustCompile(data.NameRegex.ValueString())
	}
	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		typ = data.Type.ValueStringPointer()
	}
	permissions, _, err := s.client.PermissionService.GetList(domain)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get permissions", err))
		return
	}
	var result []*buddy.Permission
	for _, p := range permissions.PermissionSets {
		if nameRegex != nil && !nameRegex.MatchString(p.Name) {
			continue
		}
		if typ != nil && *typ != p.Type {
			continue
		}
		result = append(result, p)
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, &result)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

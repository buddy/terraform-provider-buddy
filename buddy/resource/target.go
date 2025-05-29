package resource

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"regexp"
	"strings"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ resource.Resource                = &targetResource{}
	_ resource.ResourceWithConfigure   = &targetResource{}
	_ resource.ResourceWithImportState = &targetResource{}
)

func NewTargetResource() resource.Resource {
	return &targetResource{}
}

type targetResource struct {
	client *buddy.Client
}

type targetResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Domain        types.String `tfsdk:"domain"`
	TargetId      types.String `tfsdk:"target_id"`
	Name          types.String `tfsdk:"name"`
	Identifier    types.String `tfsdk:"identifier"`
	Tags          types.List   `tfsdk:"tags"`
	Type          types.String `tfsdk:"type"`
	Host          types.String `tfsdk:"host"`
	Scope         types.String `tfsdk:"scope"`
	Repository    types.String `tfsdk:"repository"`
	Port          types.String `tfsdk:"port"`
	Path          types.String `tfsdk:"path"`
	Secure        types.Bool   `tfsdk:"secure"`
	Integration   types.String `tfsdk:"integration"`
	Disabled      types.Bool   `tfsdk:"disabled"`
	Auth          types.Object `tfsdk:"auth"`
	ProjectName   types.String `tfsdk:"project_name"`
	PipelineId    types.Int64  `tfsdk:"pipeline_id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	Proxy         types.Object `tfsdk:"proxy"`
	Permissions   types.Object `tfsdk:"permissions"`
}

type targetAuthModel struct {
	Method     types.String `tfsdk:"method"`
	Username   types.String `tfsdk:"username"`
	Password   types.String `tfsdk:"password"`
	Asset      types.String `tfsdk:"asset"`
	Passphrase types.String `tfsdk:"passphrase"`
	Key        types.String `tfsdk:"key"`
}

type targetProxyModel struct {
	Name types.String `tfsdk:"name"`
	Host types.String `tfsdk:"host"`
	Port types.String `tfsdk:"port"`
	Auth types.Object `tfsdk:"auth"`
}

type targetPermissionsModel struct {
	Others types.String `tfsdk:"others"`
	Users  types.List   `tfsdk:"users"`
	Groups types.List   `tfsdk:"groups"`
}

type targetResourcePermissionModel struct {
	Id          types.Int64  `tfsdk:"id"`
	AccessLevel types.String `tfsdk:"access_level"`
}

func (r *targetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_target"
}

func (r *targetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage a target\n\n" +
			"Token scope required: `WORKSPACE`",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"domain": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"identifier": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9-]+$`), "must be a valid slug"),
				},
			},
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						buddy.TargetTypeFtp,
						buddy.TargetTypeSsh,
						buddy.TargetTypeMatch,
						buddy.TargetTypeUpcloud,
						buddy.TargetTypeVultr,
						buddy.TargetTypeDigitalOcean,
						buddy.TargetTypeGit,
					),
				},
			},
			"host": schema.StringAttribute{
				Optional: true,
			},
			"scope": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(buddy.TargetScopeWorkspace),
				Validators: []validator.String{
					stringvalidator.OneOf(
						buddy.TargetScopeWorkspace,
						buddy.TargetScopeProject,
						buddy.TargetScopeEnvironment,
						buddy.TargetScopePipeline,
						buddy.TargetScopeAction,
						buddy.TargetScopeAny,
					),
				},
			},
			"repository": schema.StringAttribute{
				Optional: true,
			},
			"port": schema.StringAttribute{
				Optional: true,
			},
			"path": schema.StringAttribute{
				Optional: true,
			},
			"secure": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"integration": schema.StringAttribute{
				Optional: true,
			},
			"disabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"auth": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"method": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf(
								buddy.TargetAuthMethodPassword,
								buddy.TargetAuthMethodSshKey,
								buddy.TargetAuthMethodAssetsKey,
								buddy.TargetAuthMethodProxyCredentials,
								buddy.TargetAuthMethodProxyKey,
								buddy.TargetAuthMethodHttp,
							),
						},
					},
					"username": schema.StringAttribute{
						Optional: true,
					},
					"password": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
					"asset": schema.StringAttribute{
						Optional: true,
					},
					"passphrase": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
					"key": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
				},
			},
			"project_name": schema.StringAttribute{
				Optional: true,
			},
			"pipeline_id": schema.Int64Attribute{
				Optional: true,
			},
			"environment_id": schema.StringAttribute{
				Optional: true,
			},
			"proxy": schema.SingleNestedAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required: true,
					},
					"host": schema.StringAttribute{
						Required: true,
					},
					"port": schema.StringAttribute{
						Required: true,
					},
					"auth": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"method": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf(
										buddy.TargetAuthMethodPassword,
										buddy.TargetAuthMethodSshKey,
									),
								},
							},
							"username": schema.StringAttribute{
								Optional: true,
							},
							"password": schema.StringAttribute{
								Optional:  true,
								Sensitive: true,
							},
							"asset": schema.StringAttribute{
								Optional: true,
							},
							"passphrase": schema.StringAttribute{
								Optional:  true,
								Sensitive: true,
							},
							"key": schema.StringAttribute{
								Optional:  true,
								Sensitive: true,
							},
						},
					},
				},
			},
			"permissions": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"others": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(buddy.TargetPermissionUseOnly),
						Validators: []validator.String{
							stringvalidator.OneOf(
								buddy.TargetPermissionManage,
								buddy.TargetPermissionUseOnly,
							),
						},
					},
					"users": schema.ListNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.Int64Attribute{
									Required: true,
								},
								"access_level": schema.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf(
											buddy.TargetPermissionManage,
											buddy.TargetPermissionUseOnly,
										),
									},
								},
							},
						},
					},
					"groups": schema.ListNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.Int64Attribute{
									Required: true,
								},
								"access_level": schema.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf(
											buddy.TargetPermissionManage,
											buddy.TargetPermissionUseOnly,
										),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *targetResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *targetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *targetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := data.Domain.ValueString()
	ops, d := data.toOps(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	target, _, err := r.client.TargetService.Create(domain, ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create target", err))
		return
	}

	resp.Diagnostics.Append(data.loadAPI(ctx, domain, target)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *targetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *targetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := data.Domain.ValueString()
	targetId := data.TargetId.ValueString()

	target, httpResp, err := r.client.TargetService.Get(domain, targetId)
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get target", err))
		return
	}

	resp.Diagnostics.Append(data.loadAPI(ctx, domain, target)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *targetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *targetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := data.Domain.ValueString()
	targetId := data.TargetId.ValueString()

	ops, d := data.toOps(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	target, _, err := r.client.TargetService.Update(domain, targetId, ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update target", err))
		return
	}

	resp.Diagnostics.Append(data.loadAPI(ctx, domain, target)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *targetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *targetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := data.Domain.ValueString()
	targetId := data.TargetId.ValueString()

	_, err := r.client.TargetService.Delete(domain, targetId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete target", err))
		return
	}
}

func (r *targetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 {
		resp.Diagnostics.AddAttributeError(
			path.Root("id"),
			"Invalid Import ID",
			"Import ID must be in the format: domain:target_id",
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("target_id"), parts[1])...)
}

func targetAuthModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"method":     types.StringType,
		"username":   types.StringType,
		"password":   types.StringType,
		"asset":      types.StringType,
		"passphrase": types.StringType,
		"key":        types.StringType,
	}
}

func targetProxyModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
		"host": types.StringType,
		"port": types.StringType,
		"auth": types.ObjectType{AttrTypes: targetAuthModelAttrs()},
	}
}

func targetResourcePermissionModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"id":           types.Int64Type,
		"access_level": types.StringType,
	}
}

func targetPermissionsModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"others": types.StringType,
		"users":  types.ListType{ElemType: types.ObjectType{AttrTypes: targetResourcePermissionModelAttrs()}},
		"groups": types.ListType{ElemType: types.ObjectType{AttrTypes: targetResourcePermissionModelAttrs()}},
	}
}

func (m *targetResourceModel) toOps(ctx context.Context) (*buddy.TargetOps, diag.Diagnostics) {
	var diags diag.Diagnostics
	ops := &buddy.TargetOps{}

	if !m.Name.IsNull() && !m.Name.IsUnknown() {
		name := m.Name.ValueString()
		ops.Name = &name
	}

	if !m.Identifier.IsNull() && !m.Identifier.IsUnknown() {
		identifier := m.Identifier.ValueString()
		ops.Identifier = &identifier
	}

	if !m.Tags.IsNull() && !m.Tags.IsUnknown() {
		var tags []string
		diags.Append(m.Tags.ElementsAs(ctx, &tags, false)...)
		ops.Tags = &tags
	}

	if !m.Type.IsNull() && !m.Type.IsUnknown() {
		typ := m.Type.ValueString()
		ops.Type = &typ
	}

	if !m.Host.IsNull() && !m.Host.IsUnknown() {
		host := m.Host.ValueString()
		ops.Host = &host
	}

	if !m.Scope.IsNull() && !m.Scope.IsUnknown() {
		scope := m.Scope.ValueString()
		ops.Scope = &scope
	}

	if !m.Repository.IsNull() && !m.Repository.IsUnknown() {
		repository := m.Repository.ValueString()
		ops.Repository = &repository
	}

	if !m.Port.IsNull() && !m.Port.IsUnknown() {
		port := m.Port.ValueString()
		ops.Port = &port
	}

	if !m.Path.IsNull() && !m.Path.IsUnknown() {
		path := m.Path.ValueString()
		ops.Path = &path
	}

	if !m.Secure.IsNull() && !m.Secure.IsUnknown() {
		secure := m.Secure.ValueBool()
		ops.Secure = &secure
	}

	if !m.Integration.IsNull() && !m.Integration.IsUnknown() {
		integration := m.Integration.ValueString()
		ops.Integration = &integration
	}

	if !m.Disabled.IsNull() && !m.Disabled.IsUnknown() {
		disabled := m.Disabled.ValueBool()
		ops.Disabled = &disabled
	}

	if !m.Auth.IsNull() && !m.Auth.IsUnknown() {
		var auth targetAuthModel
		diags.Append(m.Auth.As(ctx, &auth, basetypes.ObjectAsOptions{})...)
		ops.Auth = &buddy.TargetAuth{
			Method:     auth.Method.ValueString(),
			Username:   auth.Username.ValueString(),
			Password:   auth.Password.ValueString(),
			Asset:      auth.Asset.ValueString(),
			Passphrase: auth.Passphrase.ValueString(),
			Key:        auth.Key.ValueString(),
		}
	}

	if !m.ProjectName.IsNull() && !m.ProjectName.IsUnknown() {
		ops.Project = &buddy.TargetProject{
			Name: m.ProjectName.ValueString(),
		}
	}

	if !m.PipelineId.IsNull() && !m.PipelineId.IsUnknown() {
		ops.Pipeline = &buddy.TargetPipeline{
			Id: int(m.PipelineId.ValueInt64()),
		}
	}

	if !m.EnvironmentId.IsNull() && !m.EnvironmentId.IsUnknown() {
		ops.Environment = &buddy.TargetEnvironment{
			Id: m.EnvironmentId.ValueString(),
		}
	}

	if !m.Proxy.IsNull() && !m.Proxy.IsUnknown() {
		var proxy targetProxyModel
		diags.Append(m.Proxy.As(ctx, &proxy, basetypes.ObjectAsOptions{})...)

		proxyOps := &buddy.TargetProxy{
			Name: proxy.Name.ValueString(),
			Host: proxy.Host.ValueString(),
			Port: proxy.Port.ValueString(),
		}

		if !proxy.Auth.IsNull() && !proxy.Auth.IsUnknown() {
			var proxyAuth targetAuthModel
			diags.Append(proxy.Auth.As(ctx, &proxyAuth, basetypes.ObjectAsOptions{})...)
			proxyOps.Auth = &buddy.TargetAuth{
				Method:     proxyAuth.Method.ValueString(),
				Username:   proxyAuth.Username.ValueString(),
				Password:   proxyAuth.Password.ValueString(),
				Asset:      proxyAuth.Asset.ValueString(),
				Passphrase: proxyAuth.Passphrase.ValueString(),
				Key:        proxyAuth.Key.ValueString(),
			}
		}

		ops.Proxy = proxyOps
	}

	if !m.Permissions.IsNull() && !m.Permissions.IsUnknown() {
		var perms targetPermissionsModel
		diags.Append(m.Permissions.As(ctx, &perms, basetypes.ObjectAsOptions{})...)

		permsOps := &buddy.TargetPermissions{
			Others: perms.Others.ValueString(),
		}

		if !perms.Users.IsNull() && !perms.Users.IsUnknown() {
			var users []targetResourcePermissionModel
			diags.Append(perms.Users.ElementsAs(ctx, &users, false)...)

			var userPerms []*buddy.TargetResourcePermission
			for _, u := range users {
				userPerms = append(userPerms, &buddy.TargetResourcePermission{
					Id:          int(u.Id.ValueInt64()),
					AccessLevel: u.AccessLevel.ValueString(),
				})
			}
			permsOps.Users = userPerms
		}

		if !perms.Groups.IsNull() && !perms.Groups.IsUnknown() {
			var groups []targetResourcePermissionModel
			diags.Append(perms.Groups.ElementsAs(ctx, &groups, false)...)

			var groupPerms []*buddy.TargetResourcePermission
			for _, g := range groups {
				groupPerms = append(groupPerms, &buddy.TargetResourcePermission{
					Id:          int(g.Id.ValueInt64()),
					AccessLevel: g.AccessLevel.ValueString(),
				})
			}
			permsOps.Groups = groupPerms
		}

		ops.Permissions = permsOps
	}

	return ops, diags
}

func (m *targetResourceModel) loadAPI(ctx context.Context, domain string, target *buddy.Target) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(util.ComposeDoubleId(domain, target.Id))
	m.Domain = types.StringValue(domain)
	m.TargetId = types.StringValue(target.Id)
	m.Name = types.StringValue(target.Name)
	m.Identifier = types.StringValue(target.Identifier)
	m.Type = types.StringValue(target.Type)
	m.Scope = types.StringValue(target.Scope)

	if target.Tags != nil && len(target.Tags) > 0 {
		tags, d := types.ListValueFrom(ctx, types.StringType, target.Tags)
		diags.Append(d...)
		m.Tags = tags
	} else {
		m.Tags = types.ListNull(types.StringType)
	}

	if target.Host != "" {
		m.Host = types.StringValue(target.Host)
	}

	if target.Repository != "" {
		m.Repository = types.StringValue(target.Repository)
	}

	if target.Port != "" {
		m.Port = types.StringValue(target.Port)
	}

	if target.Path != "" {
		m.Path = types.StringValue(target.Path)
	}

	m.Secure = types.BoolValue(target.Secure)

	if target.Integration != "" {
		m.Integration = types.StringValue(target.Integration)
	}

	m.Disabled = types.BoolValue(target.Disabled)

	if target.Auth != nil {
		authModel := targetAuthModel{
			Method:   types.StringValue(target.Auth.Method),
			Username: types.StringValue(target.Auth.Username),
		}

		// Only set password if it's already in state (passwords are write-only)
		if !m.Auth.IsNull() {
			var currentAuth targetAuthModel
			diags.Append(m.Auth.As(ctx, &currentAuth, basetypes.ObjectAsOptions{})...)
			authModel.Password = currentAuth.Password
			authModel.Passphrase = currentAuth.Passphrase
			authModel.Key = currentAuth.Key
		}

		if target.Auth.Asset != "" {
			authModel.Asset = types.StringValue(target.Auth.Asset)
		}

		authObj, d := types.ObjectValueFrom(ctx, targetAuthModelAttrs(), authModel)
		diags.Append(d...)
		m.Auth = authObj
	}

	if target.Proxy != nil {
		proxyModel := targetProxyModel{
			Name: types.StringValue(target.Proxy.Name),
			Host: types.StringValue(target.Proxy.Host),
			Port: types.StringValue(target.Proxy.Port),
		}

		if target.Proxy.Auth != nil {
			proxyAuthModel := targetAuthModel{
				Method:   types.StringValue(target.Proxy.Auth.Method),
				Username: types.StringValue(target.Proxy.Auth.Username),
			}

			// Only set password if it's already in state (passwords are write-only)
			if !m.Proxy.IsNull() {
				var currentProxy targetProxyModel
				diags.Append(m.Proxy.As(ctx, &currentProxy, basetypes.ObjectAsOptions{})...)
				if !currentProxy.Auth.IsNull() {
					var currentProxyAuth targetAuthModel
					diags.Append(currentProxy.Auth.As(ctx, &currentProxyAuth, basetypes.ObjectAsOptions{})...)
					proxyAuthModel.Password = currentProxyAuth.Password
					proxyAuthModel.Passphrase = currentProxyAuth.Passphrase
					proxyAuthModel.Key = currentProxyAuth.Key
				}
			}

			if target.Proxy.Auth.Asset != "" {
				proxyAuthModel.Asset = types.StringValue(target.Proxy.Auth.Asset)
			}

			proxyAuthObj, d := types.ObjectValueFrom(ctx, targetAuthModelAttrs(), proxyAuthModel)
			diags.Append(d...)
			proxyModel.Auth = proxyAuthObj
		}

		proxyObj, d := types.ObjectValueFrom(ctx, targetProxyModelAttrs(), proxyModel)
		diags.Append(d...)
		m.Proxy = proxyObj
	}

	if target.Permissions != nil {
		permsModel := targetPermissionsModel{
			Others: types.StringValue(target.Permissions.Others),
		}

		if target.Permissions.Users != nil && len(target.Permissions.Users) > 0 {
			var users []targetResourcePermissionModel
			for _, u := range target.Permissions.Users {
				users = append(users, targetResourcePermissionModel{
					Id:          types.Int64Value(int64(u.Id)),
					AccessLevel: types.StringValue(u.AccessLevel),
				})
			}
			usersList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: targetResourcePermissionModelAttrs()}, users)
			diags.Append(d...)
			permsModel.Users = usersList
		} else {
			permsModel.Users = types.ListNull(types.ObjectType{AttrTypes: targetResourcePermissionModelAttrs()})
		}

		if target.Permissions.Groups != nil && len(target.Permissions.Groups) > 0 {
			var groups []targetResourcePermissionModel
			for _, g := range target.Permissions.Groups {
				groups = append(groups, targetResourcePermissionModel{
					Id:          types.Int64Value(int64(g.Id)),
					AccessLevel: types.StringValue(g.AccessLevel),
				})
			}
			groupsList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: targetResourcePermissionModelAttrs()}, groups)
			diags.Append(d...)
			permsModel.Groups = groupsList
		} else {
			permsModel.Groups = types.ListNull(types.ObjectType{AttrTypes: targetResourcePermissionModelAttrs()})
		}

		permsObj, d := types.ObjectValueFrom(ctx, targetPermissionsModelAttrs(), permsModel)
		diags.Append(d...)
		m.Permissions = permsObj
	}

	return diags
}

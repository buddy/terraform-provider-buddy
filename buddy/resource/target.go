package resource

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	ID             types.String `tfsdk:"id"`
	Domain         types.String `tfsdk:"domain"`
	TargetId       types.String `tfsdk:"target_id"`
	Name           types.String `tfsdk:"name"`
	Type           types.String `tfsdk:"type"`
	Host           types.String `tfsdk:"host"`
	Port           types.String `tfsdk:"port"`
	Path           types.String `tfsdk:"path"`
	Secure         types.Bool   `tfsdk:"secure"`
	Scope          types.String `tfsdk:"scope"`
	Repository     types.String `tfsdk:"repository"`
	Integration    types.String `tfsdk:"integration"`
	Tags           types.Set    `tfsdk:"tags"`
	Disabled       types.Bool   `tfsdk:"disabled"`
	AuthMethod     types.String `tfsdk:"auth_method"`
	AuthUsername   types.String `tfsdk:"auth_username"`
	AuthPassword   types.String `tfsdk:"auth_password"`
	AuthKey        types.String `tfsdk:"auth_key"`
	AuthPassphrase types.String `tfsdk:"auth_passphrase"`
	HtmlUrl        types.String `tfsdk:"html_url"`
}

func (r *targetResourceModel) loadAPI(ctx context.Context, domain string, target *buddy.Target) diag.Diagnostics {
	var diags diag.Diagnostics

	r.ID = types.StringValue(util.ComposeDoubleId(domain, target.Id))
	r.Domain = types.StringValue(domain)
	r.TargetId = types.StringValue(target.Id)
	r.Name = types.StringValue(target.Name)
	r.Type = types.StringValue(target.Type)
	r.HtmlUrl = types.StringValue(target.HtmlUrl)
	r.Scope = types.StringValue(target.Scope)
	r.Disabled = types.BoolValue(target.Disabled)
	r.Secure = types.BoolValue(target.Secure)

	if target.Host != "" {
		r.Host = types.StringValue(target.Host)
	} else {
		r.Host = types.StringNull()
	}

	if target.Port != "" {
		r.Port = types.StringValue(target.Port)
	} else {
		r.Port = types.StringNull()
	}

	if target.Path != "" {
		r.Path = types.StringValue(target.Path)
	} else {
		r.Path = types.StringNull()
	}

	if target.Repository != "" {
		r.Repository = types.StringValue(target.Repository)
	} else {
		r.Repository = types.StringNull()
	}

	if target.Integration != "" {
		r.Integration = types.StringValue(target.Integration)
	} else {
		r.Integration = types.StringNull()
	}

	tags, d := types.SetValueFrom(ctx, types.StringType, &target.Tags)
	diags.Append(d...)
	r.Tags = tags

	// Handle auth fields if auth is present
	if target.Auth != nil {
		r.AuthMethod = types.StringValue(target.Auth.Method)

		if target.Auth.Username != "" {
			r.AuthUsername = types.StringValue(target.Auth.Username)
		} else {
			r.AuthUsername = types.StringNull()
		}

		// Password is not returned by API, keep existing value

		if target.Auth.Key != "" {
			r.AuthKey = types.StringValue(target.Auth.Key)
		} else {
			r.AuthKey = types.StringNull()
		}

		// Passphrase is not returned by API, keep existing value
	} else {
		r.AuthMethod = types.StringNull()
		r.AuthUsername = types.StringNull()
		r.AuthKey = types.StringNull()
	}

	return diags
}

func (r *targetResourceModel) decomposeId() (string, string, error) {
	domain, tid, err := util.DecomposeDoubleId(r.ID.ValueString())
	if err != nil {
		return "", "", err
	}
	return domain, tid, nil
}

func (r *targetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_target"
}

func (r *targetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage a deployment target\n\n" +
			"Workspace administrator rights are required\n\n" +
			"Token scope required: `WORKSPACE`",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Terraform resource identifier for this item",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The workspace's URL handle",
				Required:            true,
				Validators:          util.StringValidatorsDomain(),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_id": schema.StringAttribute{
				MarkdownDescription: "The target's ID",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The target's name",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The target's type",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"SSH",
						"SHELL",
						"FTP",
						"FTPS",
						"SFTP",
						"AMAZON_S3",
						"GOOGLE_CLOUD_STORAGE",
						"AZURE_STORAGE",
						"DOCKER_REGISTRY",
						"KUBERNETES",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "The target's hostname or IP address",
				Optional:            true,
			},
			"port": schema.StringAttribute{
				MarkdownDescription: "The target's port",
				Optional:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "The remote path on the target",
				Optional:            true,
			},
			"secure": schema.BoolAttribute{
				MarkdownDescription: "Whether to use secure connection",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: "The target's scope",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("PRIVATE"),
				Validators: []validator.String{
					stringvalidator.OneOf("PRIVATE", "PUBLIC", "WORKSPACE"),
				},
			},
			"repository": schema.StringAttribute{
				MarkdownDescription: "The repository for registry targets",
				Optional:            true,
			},
			"integration": schema.StringAttribute{
				MarkdownDescription: "The integration ID to use for cloud targets",
				Optional:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "The target's tags",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"disabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the target is disabled",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"auth_method": schema.StringAttribute{
				MarkdownDescription: "The authentication method",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("PASS", "KEY", "BOTH"),
				},
			},
			"auth_username": schema.StringAttribute{
				MarkdownDescription: "The authentication username",
				Optional:            true,
			},
			"auth_password": schema.StringAttribute{
				MarkdownDescription: "The authentication password",
				Optional:            true,
				Sensitive:           true,
			},
			"auth_key": schema.StringAttribute{
				MarkdownDescription: "The SSH key content for authentication",
				Optional:            true,
				Sensitive:           true,
			},
			"auth_passphrase": schema.StringAttribute{
				MarkdownDescription: "The SSH key passphrase",
				Optional:            true,
				Sensitive:           true,
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The target's URL",
				Computed:            true,
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

	ops := buddy.TargetOps{
		Name: data.Name.ValueStringPointer(),
		Type: data.Type.ValueStringPointer(),
	}

	if !data.Host.IsNull() {
		ops.Host = data.Host.ValueStringPointer()
	}

	if !data.Port.IsNull() {
		ops.Port = data.Port.ValueStringPointer()
	}

	if !data.Path.IsNull() {
		ops.Path = data.Path.ValueStringPointer()
	}

	if !data.Secure.IsNull() {
		secure := data.Secure.ValueBool()
		ops.Secure = &secure
	}

	if !data.Scope.IsNull() {
		ops.Scope = data.Scope.ValueStringPointer()
	}

	if !data.Repository.IsNull() {
		ops.Repository = data.Repository.ValueStringPointer()
	}

	if !data.Integration.IsNull() {
		ops.Integration = data.Integration.ValueStringPointer()
	}

	if !data.Disabled.IsNull() {
		disabled := data.Disabled.ValueBool()
		ops.Disabled = &disabled
	}

	if !data.Tags.IsNull() {
		var tags []string
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Tags = &tags
	}

	// Handle auth if any auth fields are set
	if !data.AuthMethod.IsNull() || !data.AuthUsername.IsNull() || !data.AuthPassword.IsNull() || !data.AuthKey.IsNull() || !data.AuthPassphrase.IsNull() {
		auth := &buddy.TargetAuth{}

		if !data.AuthMethod.IsNull() {
			auth.Method = data.AuthMethod.ValueString()
		}

		if !data.AuthUsername.IsNull() {
			auth.Username = data.AuthUsername.ValueString()
		}

		if !data.AuthPassword.IsNull() {
			auth.Password = data.AuthPassword.ValueString()
		}

		if !data.AuthKey.IsNull() {
			auth.Key = data.AuthKey.ValueString()
		}

		if !data.AuthPassphrase.IsNull() {
			auth.Passphrase = data.AuthPassphrase.ValueString()
		}

		ops.Auth = auth
	}

	target, _, err := r.client.TargetService.Create(domain, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create target", err))
		return
	}

	resp.Diagnostics.Append(data.loadAPI(ctx, domain, target)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve password and passphrase as they are not returned by API
	if ops.Auth != nil {
		if ops.Auth.Password != "" {
			data.AuthPassword = types.StringValue(ops.Auth.Password)
		}
		if ops.Auth.Passphrase != "" {
			data.AuthPassphrase = types.StringValue(ops.Auth.Passphrase)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *targetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *targetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain, targetId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("target", err))
		return
	}

	target, _, err := r.client.TargetService.Get(domain, targetId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get target", err))
		return
	}

	// Save current password and passphrase as they are not returned by API
	currentPassword := data.AuthPassword
	currentPassphrase := data.AuthPassphrase

	resp.Diagnostics.Append(data.loadAPI(ctx, domain, target)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Restore password and passphrase
	data.AuthPassword = currentPassword
	data.AuthPassphrase = currentPassphrase

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *targetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *targetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain, targetId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("target", err))
		return
	}

	ops := buddy.TargetOps{}

	if !data.Name.IsNull() {
		ops.Name = data.Name.ValueStringPointer()
	}

	if !data.Host.IsNull() {
		ops.Host = data.Host.ValueStringPointer()
	}

	if !data.Port.IsNull() {
		ops.Port = data.Port.ValueStringPointer()
	}

	if !data.Path.IsNull() {
		ops.Path = data.Path.ValueStringPointer()
	}

	if !data.Secure.IsNull() {
		secure := data.Secure.ValueBool()
		ops.Secure = &secure
	}

	if !data.Scope.IsNull() {
		ops.Scope = data.Scope.ValueStringPointer()
	}

	if !data.Repository.IsNull() {
		ops.Repository = data.Repository.ValueStringPointer()
	}

	if !data.Integration.IsNull() {
		ops.Integration = data.Integration.ValueStringPointer()
	}

	if !data.Disabled.IsNull() {
		disabled := data.Disabled.ValueBool()
		ops.Disabled = &disabled
	}

	if !data.Tags.IsNull() {
		var tags []string
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Tags = &tags
	}

	// Handle auth if any auth fields are set
	if !data.AuthMethod.IsNull() || !data.AuthUsername.IsNull() || !data.AuthPassword.IsNull() || !data.AuthKey.IsNull() || !data.AuthPassphrase.IsNull() {
		auth := &buddy.TargetAuth{}

		if !data.AuthMethod.IsNull() {
			auth.Method = data.AuthMethod.ValueString()
		}

		if !data.AuthUsername.IsNull() {
			auth.Username = data.AuthUsername.ValueString()
		}

		if !data.AuthPassword.IsNull() {
			auth.Password = data.AuthPassword.ValueString()
		}

		if !data.AuthKey.IsNull() {
			auth.Key = data.AuthKey.ValueString()
		}

		if !data.AuthPassphrase.IsNull() {
			auth.Passphrase = data.AuthPassphrase.ValueString()
		}

		ops.Auth = auth
	}

	target, _, err := r.client.TargetService.Update(domain, targetId, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update target", err))
		return
	}

	resp.Diagnostics.Append(data.loadAPI(ctx, domain, target)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve password and passphrase as they are not returned by API
	if ops.Auth != nil {
		if ops.Auth.Password != "" {
			data.AuthPassword = types.StringValue(ops.Auth.Password)
		}
		if ops.Auth.Passphrase != "" {
			data.AuthPassphrase = types.StringValue(ops.Auth.Passphrase)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *targetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *targetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain, targetId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("target", err))
		return
	}

	_, err = r.client.TargetService.Delete(domain, targetId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete target", err))
		return
	}
}

func (r *targetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

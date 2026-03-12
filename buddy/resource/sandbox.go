package resource

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ resource.Resource                = &sandboxResource{}
	_ resource.ResourceWithConfigure   = &sandboxResource{}
	_ resource.ResourceWithImportState = &sandboxResource{}
)

type sandboxResourceModel struct {
	ID                       types.String `tfsdk:"id"`
	Identifier               types.String `tfsdk:"identifier"`
	Domain                   types.String `tfsdk:"domain"`
	ProjectName              types.String `tfsdk:"project_name"`
	HtmlUrl                  types.String `tfsdk:"html_url"`
	SandboxId                types.String `tfsdk:"sandbox_id"`
	Name                     types.String `tfsdk:"name"`
	Status                   types.String `tfsdk:"status"`
	SetupStatus              types.String `tfsdk:"setup_status"`
	BootLogs                 types.String `tfsdk:"boot_logs"`
	InstallCommands          types.String `tfsdk:"install_commands"`
	AppDir                   types.String `tfsdk:"app_dir"`
	Os                       types.String `tfsdk:"os"`
	Resources                types.String `tfsdk:"resources"`
	Timeout                  types.Int32  `tfsdk:"timeout"`
	Tags                     types.Set    `tfsdk:"tags"`
	AppCommands              types.Set    `tfsdk:"app_commands"`
	Apps                     types.Set    `tfsdk:"apps"`
	Endpoints                types.Map    `tfsdk:"endpoints"`
	WaitForRunning           types.Bool   `tfsdk:"wait_for_running"`
	WaitForRunningTimeout    types.Int32  `tfsdk:"wait_for_running_timeout"`
	WaitForConfigured        types.Bool   `tfsdk:"wait_for_configured"`
	WaitForConfiguredTimeout types.Int32  `tfsdk:"wait_for_configured_timeout"`
	WaitForApps              types.Bool   `tfsdk:"wait_for_apps"`
	WaitForAppsTimeout       types.Int32  `tfsdk:"wait_for_apps_timeout"`
	Permissions              types.Set    `tfsdk:"permissions"`
}

func (r *sandboxResourceModel) decomposeId() (string, string, error) {
	domain, sandboxId, err := util.DecomposeDoubleId(r.ID.ValueString())
	if err != nil {
		return "", "", err
	}
	return domain, sandboxId, nil
}

func (r *sandboxResourceModel) loadAPI(ctx context.Context, domain string, sandbox *buddy.Sandbox, waitForRunning bool, waitForRunningTimeout int32, waitForConfigured bool, waitForConfiguredTimeout int32, waitForApps bool, waitForAppsTimeout int32) diag.Diagnostics {
	var diags diag.Diagnostics
	r.ID = types.StringValue(util.ComposeDoubleId(domain, sandbox.Id))
	r.Domain = types.StringValue(domain)
	r.ProjectName = types.StringValue(sandbox.Project.Name)
	r.Identifier = types.StringValue(sandbox.Identifier)
	r.SandboxId = types.StringValue(sandbox.Id)
	r.HtmlUrl = types.StringValue(sandbox.HtmlUrl)
	r.Name = types.StringValue(sandbox.Name)
	r.Status = types.StringValue(sandbox.Status)
	r.SetupStatus = types.StringValue(sandbox.SetupStatus)
	r.BootLogs = types.StringValue(strings.Join(sandbox.BootLogs, "\n"))
	r.InstallCommands = types.StringValue(sandbox.FirstBootCommands)
	r.AppDir = types.StringValue(sandbox.AppDir)
	r.Os = types.StringValue(sandbox.Os)
	r.Resources = types.StringValue(sandbox.Resources)
	tags, d := types.SetValueFrom(ctx, types.StringType, &sandbox.Tags)
	diags.Append(d...)
	r.Tags = tags
	endpoints, d := util.SandboxEndpointsFromApi(ctx, &sandbox.Endpoints)
	diags.Append(d...)
	r.Endpoints = endpoints
	apps, d := util.SandboxAppsFromApi(ctx, &sandbox.Apps)
	diags.Append(d...)
	r.Apps = apps
	r.Timeout = types.Int32Value(int32(sandbox.Timeout))
	r.WaitForRunning = types.BoolValue(waitForRunning)
	r.WaitForRunningTimeout = types.Int32Value(waitForRunningTimeout)
	r.WaitForConfigured = types.BoolValue(waitForConfigured)
	r.WaitForConfiguredTimeout = types.Int32Value(waitForConfiguredTimeout)
	r.WaitForApps = types.BoolValue(waitForApps)
	r.WaitForAppsTimeout = types.Int32Value(waitForAppsTimeout)
	return diags
}

func NewSandboxResource() resource.Resource {
	return &sandboxResource{}
}

type sandboxResource struct {
	client *buddy.Client
}

func (r *sandboxResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sandbox"
}

func (r *sandboxResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage a sandbox\n\n" +
			"Token scopes required: `WORKSPACE`, `SANDBOX_MANAGE`, `SANDBOX_INFO`",
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
			"project_name": schema.StringAttribute{
				MarkdownDescription: "The project's name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The sandbox's URL",
				Computed:            true,
			},
			"sandbox_id": schema.StringAttribute{
				MarkdownDescription: "The sandbox's ID",
				Computed:            true,
			},
			"identifier": schema.StringAttribute{
				MarkdownDescription: "The sandbox's identifier",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The sandbox's name",
				Required:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The sandbox's status",
				Computed:            true,
			},
			"setup_status": schema.StringAttribute{
				MarkdownDescription: "The sandbox's setup status",
				Computed:            true,
			},
			"install_commands": schema.StringAttribute{
				MarkdownDescription: "The sandbox's install commands",
				Optional:            true,
				Computed:            true,
			},
			"app_dir": schema.StringAttribute{
				MarkdownDescription: "The sandbox's app dir",
				Optional:            true,
				Computed:            true,
			},
			"app_commands": schema.SetAttribute{
				MarkdownDescription: "The sandbox's app commands",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"timeout": schema.Int32Attribute{
				MarkdownDescription: "The sandbox's start timeout",
				Optional:            true,
				Computed:            true,
			},
			"boot_logs": schema.StringAttribute{
				MarkdownDescription: "The sandbox's boot logs",
				Computed:            true,
			},
			"apps": schema.SetNestedAttribute{
				MarkdownDescription: "The sandbox's apps",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"command": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"os": schema.StringAttribute{
				MarkdownDescription: "The sandbox's operating system",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(buddy.SandboxOsUbuntu2404),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						buddy.SandboxOsUbuntu2204,
						buddy.SandboxOsUbuntu2404,
					),
				},
			},
			"resources": schema.StringAttribute{
				MarkdownDescription: "The sandbox's resources (cpu, ram)",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						buddy.SandboxResource1X2,
						buddy.SandboxResource2X4,
						buddy.SandboxResource3X6,
						buddy.SandboxResource4X8,
						buddy.SandboxResource5X10,
						buddy.SandboxResource6X12,
						buddy.SandboxResource7X14,
						buddy.SandboxResource8X16,
						buddy.SandboxResource9X18,
						buddy.SandboxResource10X20,
						buddy.SandboxResource11X22,
						buddy.SandboxResource12X24,
					),
				},
			},
			"tags": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The sandbox's list of tags",
				Optional:            true,
				Computed:            true,
			},
			"endpoints": schema.MapNestedAttribute{
				MarkdownDescription: "The sandbox's map of endpoints",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: util.ResourceSandboxEndpointModelAttributes(),
				},
			},
			"wait_for_running": schema.BoolAttribute{
				MarkdownDescription: "Wait until sandbox is running",
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				Computed:            true,
			},
			"wait_for_running_timeout": schema.Int32Attribute{
				MarkdownDescription: "Seconds to wait until sandbox is running",
				Optional:            true,
				Default:             int32default.StaticInt32(120),
				Computed:            true,
			},
			"wait_for_configured": schema.BoolAttribute{
				MarkdownDescription: "Wait until sandbox ran setup commands",
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				Computed:            true,
			},
			"wait_for_configured_timeout": schema.Int32Attribute{
				MarkdownDescription: "Seconds to wait until sandbox ran setup commands",
				Optional:            true,
				Default:             int32default.StaticInt32(120),
				Computed:            true,
			},
			"wait_for_apps": schema.BoolAttribute{
				MarkdownDescription: "Wait until sandbox running apps commands",
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				Computed:            true,
			},
			"wait_for_apps_timeout": schema.Int32Attribute{
				MarkdownDescription: "Seconds to wait until sandbox ran apps commands",
				Optional:            true,
				Default:             int32default.StaticInt32(120),
				Computed:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"permissions": schema.SetNestedBlock{
				MarkdownDescription: "The sandbox's permissions",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"others": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									buddy.SandboxPermissionDenied,
									buddy.SandboxPermissionReadOnly,
									buddy.SandboxPermissionDefault,
									buddy.SandboxPermissionManage,
								),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"user": schema.SetNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: util.SandboxPermissionsAccessModelAttributes(),
							},
						},
						"group": schema.SetNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: util.SandboxPermissionsAccessModelAttributes(),
							},
						},
					},
				},
				Validators: []validator.Set{
					setvalidator.SizeAtMost(1),
				},
			},
		},
	}
}

func (r *sandboxResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *sandboxResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *sandboxResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	projectName := data.ProjectName.ValueString()
	waitForRunning := data.WaitForRunning.ValueBool()
	waitForRunningTimeout := data.WaitForRunningTimeout.ValueInt32()
	waitForConfigured := data.WaitForConfigured.ValueBool()
	waitForConfiguredTimeout := data.WaitForConfiguredTimeout.ValueInt32()
	waitForApps := data.WaitForApps.ValueBool()
	waitForAppsTimeout := data.WaitForAppsTimeout.ValueInt32()
	ops := buddy.SandboxOps{
		Name: data.Name.ValueStringPointer(),
		Os:   data.Os.ValueStringPointer(),
	}
	if !data.Identifier.IsUnknown() && !data.Identifier.IsNull() {
		ops.Identifier = data.Identifier.ValueStringPointer()
	}
	if !data.InstallCommands.IsUnknown() && !data.InstallCommands.IsNull() {
		ops.FirstBootCommands = data.InstallCommands.ValueStringPointer()
	}
	if !data.Permissions.IsUnknown() && !data.Permissions.IsNull() {
		permissions, d := util.SandboxPermissionsModelToApi(ctx, &data.Permissions)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Permissions = permissions
	}
	if !data.AppDir.IsUnknown() && !data.AppDir.IsNull() {
		ops.AppDir = data.AppDir.ValueStringPointer()
	}
	if !data.AppCommands.IsUnknown() && !data.AppCommands.IsNull() {
		commands, d := util.StringSetToApi(ctx, &data.AppCommands)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Apps = commands
	}
	if !data.Timeout.IsUnknown() && !data.Timeout.IsNull() {
		timeout := int(data.Timeout.ValueInt32())
		ops.Timeout = &timeout
	}
	if !data.Resources.IsUnknown() && !data.Resources.IsNull() {
		ops.Resources = data.Resources.ValueStringPointer()
	}
	if !data.Tags.IsUnknown() && !data.Tags.IsNull() {
		tags, d := util.StringSetToApi(ctx, &data.Tags)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Tags = tags
	}
	if !data.Endpoints.IsUnknown() && !data.Endpoints.IsNull() {
		endpoints, d := util.SandboxEndpointsToApi(ctx, &data.Endpoints)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Endpoints = endpoints
	}
	sandbox, _, err := r.client.SandboxService.Create(domain, projectName, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create sandbox", err))
		return
	}
	if waitForRunning {
		sb, d := r.waitForRunning(domain, sandbox.Id, true, waitForRunningTimeout)
		resp.Diagnostics.Append(d...)
		if sb == nil || resp.Diagnostics.HasError() {
			return
		}
		sandbox = sb
	}
	if waitForConfigured {
		sb, d := r.waitForConfigured(domain, sandbox.Id, waitForConfiguredTimeout)
		resp.Diagnostics.Append(d...)
		if sb == nil || resp.Diagnostics.HasError() {
			return
		}
		sandbox = sb
	}
	if waitForApps {
		sb, d := r.waitForApps(domain, sandbox.Id, waitForAppsTimeout)
		resp.Diagnostics.Append(d...)
		if sb == nil || resp.Diagnostics.HasError() {
			return
		}
		sandbox = sb
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, sandbox, waitForRunning, waitForRunningTimeout, waitForConfigured, waitForConfiguredTimeout, waitForApps, waitForAppsTimeout)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sandboxResource) waitForApps(domain string, sandboxId string, timeout int32) (*buddy.Sandbox, diag.Diagnostics) {
	var diags diag.Diagnostics
	sandbox, err := r.client.SandboxService.WaitForAppStatuses(domain, sandboxId, int(timeout), []string{
		buddy.SandboxAppStatusRunning,
		buddy.SandboxAppStatusEnded,
	})
	if err != nil || sandbox == nil {
		diags.Append(util.NewDiagnosticSandboxTimeout("timeout waiting for app to start"))
		return nil, diags
	}
	return sandbox, diags
}

func (r *sandboxResource) waitForConfigured(domain string, sandboxId string, timeout int32) (*buddy.Sandbox, diag.Diagnostics) {
	var diags diag.Diagnostics
	sandbox, err := r.client.SandboxService.WaitForSetupStatuses(domain, sandboxId, int(timeout), []string{
		buddy.SandboxSetupStatusSuccess,
		buddy.SandboxSetupStatusFailed,
	})
	if err != nil || sandbox == nil {
		diags.Append(util.NewDiagnosticSandboxTimeout("timeout waiting for sandbox to run setup commands"))
		return nil, diags
	}
	if sandbox.SetupStatus == buddy.SandboxSetupStatusFailed {
		diags.Append(util.NewDiagnosticSandboxTimeout("setup commands failed to run"))
		return nil, diags
	}
	return sandbox, diags
}

func (r *sandboxResource) waitForRunning(domain string, sandboxId string, start bool, timeout int32) (*buddy.Sandbox, diag.Diagnostics) {
	var diags diag.Diagnostics
	sandbox, err := r.client.SandboxService.WaitForStatuses(domain, sandboxId, int(timeout), []string{
		buddy.SandboxStatusFailed,
		buddy.SandboxStatusRunning,
		buddy.SandboxStatusStopped,
	})
	if err != nil || sandbox == nil {
		diags.Append(util.NewDiagnosticSandboxTimeout("timeout waiting for sandbox to start"))
		return nil, diags
	}
	if sandbox.Status == buddy.SandboxStatusFailed {
		diags.Append(util.NewDiagnosticSandboxTimeout("sandbox failed to start"))
		return nil, diags
	}
	if !start || sandbox.Status == buddy.SandboxStatusRunning {
		return sandbox, diags
	}
	_, _, err = r.client.SandboxService.Start(domain, sandboxId)
	if err != nil {
		diags.Append(util.NewDiagnosticSandboxTimeout("sandbox failed to start"))
		return nil, diags
	}
	sandbox, err = r.client.SandboxService.WaitForStatuses(domain, sandboxId, 120, []string{
		buddy.SandboxStatusRunning,
	})
	if err != nil || sandbox == nil {
		diags.Append(util.NewDiagnosticSandboxTimeout("timeout waiting for sandbox to start"))
		return nil, diags
	}
	return sandbox, diags
}

func (r *sandboxResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *sandboxResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, sandboxId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("sandbox", err))
		return
	}
	sandbox, httpRes, err := r.client.SandboxService.Get(domain, sandboxId)
	if err != nil {
		if util.IsResourceNotFound(httpRes, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get sandbox", err))
		return
	}
	waitForRunning := data.WaitForRunning.ValueBool()
	waitForRunningTimeout := data.WaitForRunningTimeout.ValueInt32()
	waitForConfigured := data.WaitForConfigured.ValueBool()
	waitForConfiguredTimeout := data.WaitForConfiguredTimeout.ValueInt32()
	waitForApps := data.WaitForApps.ValueBool()
	waitForAppsTimeout := data.WaitForAppsTimeout.ValueInt32()
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, sandbox, waitForRunning, waitForRunningTimeout, waitForConfigured, waitForConfiguredTimeout, waitForApps, waitForAppsTimeout)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sandboxResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *sandboxResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, sandboxId, err := data.decomposeId()
	waitForRunning := data.WaitForRunning.ValueBool()
	waitForRunningTimeout := data.WaitForRunningTimeout.ValueInt32()
	waitForConfigured := data.WaitForConfigured.ValueBool()
	waitForConfiguredTimeout := data.WaitForConfiguredTimeout.ValueInt32()
	waitForApps := data.WaitForApps.ValueBool()
	waitForAppsTimeout := data.WaitForAppsTimeout.ValueInt32()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("sandbox", err))
		return
	}
	ops := buddy.SandboxOps{
		Name: data.Name.ValueStringPointer(),
	}
	if !data.Identifier.IsUnknown() && !data.Identifier.IsNull() {
		ops.Identifier = data.Identifier.ValueStringPointer()
	}
	if !data.Permissions.IsUnknown() && !data.Permissions.IsNull() {
		permissions, d := util.SandboxPermissionsModelToApi(ctx, &data.Permissions)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Permissions = permissions
	}
	if !data.InstallCommands.IsUnknown() && !data.InstallCommands.IsNull() {
		ops.FirstBootCommands = data.InstallCommands.ValueStringPointer()
	}
	if !data.AppDir.IsUnknown() && !data.AppDir.IsNull() {
		ops.AppDir = data.AppDir.ValueStringPointer()
	}
	if !data.Timeout.IsUnknown() && !data.Timeout.IsNull() {
		timeout := int(data.Timeout.ValueInt32())
		ops.Timeout = &timeout
	}
	if !data.AppCommands.IsUnknown() && !data.AppCommands.IsNull() {
		commands, d := util.StringSetToApi(ctx, &data.AppCommands)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Apps = commands
	}
	if !data.Os.IsUnknown() && !data.Os.IsNull() {
		ops.Os = data.Os.ValueStringPointer()
	}
	if !data.Resources.IsUnknown() && !data.Resources.IsNull() {
		ops.Resources = data.Resources.ValueStringPointer()
	}
	if !data.Tags.IsUnknown() && !data.Tags.IsNull() {
		tags, d := util.StringSetToApi(ctx, &data.Tags)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Tags = tags
	}
	if !data.Endpoints.IsUnknown() && !data.Endpoints.IsNull() {
		endpoints, d := util.SandboxEndpointsToApi(ctx, &data.Endpoints)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Endpoints = endpoints
	}
	_, d := r.waitForRunning(domain, sandboxId, false, waitForRunningTimeout)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, d = r.waitForConfigured(domain, sandboxId, waitForConfiguredTimeout)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	sandbox, _, err := r.client.SandboxService.Update(domain, sandboxId, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update sandbox", err))
		return
	}
	if waitForRunning {
		sb, d := r.waitForRunning(domain, sandbox.Id, true, waitForRunningTimeout)
		resp.Diagnostics.Append(d...)
		if sb == nil || resp.Diagnostics.HasError() {
			return
		}
		sandbox = sb
	}
	if waitForConfigured {
		sb, d := r.waitForConfigured(domain, sandbox.Id, waitForConfiguredTimeout)
		resp.Diagnostics.Append(d...)
		if sb == nil || resp.Diagnostics.HasError() {
			return
		}
		sandbox = sb
	}
	if waitForApps {
		sb, d := r.waitForApps(domain, sandbox.Id, waitForAppsTimeout)
		resp.Diagnostics.Append(d...)
		if sb == nil || resp.Diagnostics.HasError() {
			return
		}
		sandbox = sb
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, sandbox, waitForRunning, waitForRunningTimeout, waitForConfigured, waitForConfiguredTimeout, waitForApps, waitForAppsTimeout)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sandboxResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *sandboxResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, sandboxId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("sandbox", err))
		return
	}
	_, err = r.client.SandboxService.Delete(domain, sandboxId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete sandbox", err))
	}
}

func (r *sandboxResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

package resource

import (
	"context"
	"errors"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ resource.Resource                = &sandboxStatusResource{}
	_ resource.ResourceWithConfigure   = &sandboxStatusResource{}
	_ resource.ResourceWithImportState = &sandboxStatusResource{}
)

type sandboxStatusResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Domain               types.String `tfsdk:"domain"`
	SandboxId            types.String `tfsdk:"sandbox_id"`
	Status               types.String `tfsdk:"status"`
	WaitForStatus        types.Bool   `tfsdk:"wait_for_status"`
	SandboxStatus        types.String `tfsdk:"sandbox_status"`
	WaitForStatusTimeout types.Int32  `tfsdk:"wait_for_status_timeout"`
}

func (r *sandboxStatusResourceModel) decomposeId() (string, string, error) {
	if !r.ID.IsNull() && !r.ID.IsUnknown() {
		domain, sandboxId, err := util.DecomposeDoubleId(r.ID.ValueString())
		if err != nil {
			return "", "", err
		}
		return domain, sandboxId, nil
	}
	if !r.Domain.IsNull() && !r.Domain.IsUnknown() && !r.SandboxId.IsNull() && !r.SandboxId.IsUnknown() {
		return r.Domain.ValueString(), r.SandboxId.ValueString(), nil
	}
	return "", "", errors.New("can't decompose ID")
}

func (r *sandboxStatusResourceModel) loadAPI(domain string, sandbox *buddy.Sandbox, status string, waitForStatus bool, waitForStatusTimeout int32) {
	r.ID = types.StringValue(util.ComposeDoubleId(domain, sandbox.Id))
	r.Domain = types.StringValue(domain)
	r.SandboxId = types.StringValue(sandbox.Id)
	r.Status = types.StringValue(status)
	r.WaitForStatus = types.BoolValue(waitForStatus)
	r.SandboxStatus = types.StringValue(sandbox.Status)
	r.WaitForStatusTimeout = types.Int32Value(waitForStatusTimeout)
}

func NewSandboxStatusResource() resource.Resource {
	return &sandboxStatusResource{}
}

type sandboxStatusResource struct {
	client *buddy.Client
}

func (r *sandboxStatusResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sandbox_status"
}

func (r *sandboxStatusResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *sandboxStatusResource) update(ctx context.Context, diag *diag.Diagnostics, plan *tfsdk.Plan, state *tfsdk.State) {
	var data *sandboxStatusResourceModel
	diag.Append(plan.Get(ctx, &data)...)
	if diag.HasError() {
		return
	}
	domain, sandboxId, err := data.decomposeId()
	if err != nil {
		diag.Append(util.NewDiagnosticDecomposeError("sandbox", err))
		return
	}
	waitForStatusTimeout := data.WaitForStatusTimeout.ValueInt32()
	waitForStatus := data.WaitForStatus.ValueBool()
	status := data.Status.ValueString()
	sandbox, err := r.client.SandboxService.WaitForStatuses(domain, sandboxId, int(waitForStatusTimeout), []string{
		buddy.SandboxStatusRunning,
		buddy.SandboxStatusFailed,
		buddy.SandboxStatusStopped,
	})
	if err != nil {
		diag.Append(util.NewDiagnosticSandboxTimeout("timeout waiting for sandbox status"))
		return
	}
	if sandbox.Status == buddy.SandboxStatusFailed {
		diag.Append(util.NewDiagnosticSandboxTimeout("sandbox failed to start"))
		return
	}
	if status == buddy.SandboxStatusRunning && sandbox.Status == buddy.SandboxStatusStopped {
		sandbox, _, err = r.client.SandboxService.Start(domain, sandboxId)
		if err != nil {
			diag.Append(util.NewDiagnosticApiError("start sandbox", err))
			return
		}
		if waitForStatus {
			sandbox, err = r.client.SandboxService.WaitForStatuses(domain, sandboxId, int(waitForStatusTimeout), []string{
				buddy.SandboxStatusRunning,
			})
			if err != nil {
				diag.Append(util.NewDiagnosticSandboxTimeout("sandbox failed to start"))
				return
			}
		}
	} else if status == buddy.SandboxStatusStopped && sandbox.Status == buddy.SandboxStatusRunning {
		sandbox, _, err = r.client.SandboxService.Stop(domain, sandboxId)
		if err != nil {
			diag.Append(util.NewDiagnosticApiError("start sandbox", err))
			return
		}
		if waitForStatus {
			sandbox, err = r.client.SandboxService.WaitForStatuses(domain, sandboxId, int(waitForStatusTimeout), []string{
				buddy.SandboxStatusStopped,
			})
			if err != nil {
				diag.Append(util.NewDiagnosticSandboxTimeout("sandbox failed to stop"))
				return
			}
		}
	}
	data.loadAPI(domain, sandbox, status, waitForStatus, waitForStatusTimeout)
	diag.Append(state.Set(ctx, &data)...)
}

func (r *sandboxStatusResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.update(ctx, &resp.Diagnostics, &req.Plan, &resp.State)
}

func (r *sandboxStatusResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *sandboxStatusResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, sandboxId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("sandbox", err))
		return
	}
	sandbox, _, err := r.client.SandboxService.Get(domain, sandboxId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get sandbox", err))
		return
	}
	data.loadAPI(domain, sandbox, data.Status.ValueString(), data.WaitForStatus.ValueBool(), data.WaitForStatusTimeout.ValueInt32())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sandboxStatusResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.update(ctx, &resp.Diagnostics, &req.Plan, &resp.State)
}

func (r *sandboxStatusResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// do nothing
}

func (r *sandboxStatusResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *sandboxStatusResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage sandbox status\n\n" +
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
			"sandbox_id": schema.StringAttribute{
				MarkdownDescription: "The sandbox's ID",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The sandbox's status to achive",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						buddy.SandboxStatusStopped,
						buddy.SandboxStatusRunning,
					),
				},
			},
			"wait_for_status_timeout": schema.Int32Attribute{
				MarkdownDescription: "Seconds to wait until sandbox is in required status",
				Optional:            true,
				Default:             int32default.StaticInt32(120),
				Computed:            true,
			},
			"wait_for_status": schema.BoolAttribute{
				MarkdownDescription: "Wait until sandbox is in required status",
				Optional:            true,
				Default:             booldefault.StaticBool(true),
				Computed:            true,
			},
			"sandbox_status": schema.StringAttribute{
				MarkdownDescription: "The sandbox's real status",
				Computed:            true,
			},
		},
	}
}

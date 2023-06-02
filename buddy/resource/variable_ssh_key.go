package resource

import (
	"context"
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ resource.Resource                = &variableSshKeyResource{}
	_ resource.ResourceWithConfigure   = &variableSshKeyResource{}
	_ resource.ResourceWithImportState = &variableSshKeyResource{}
)

func NewVariableSshResource() resource.Resource {
	return &variableSshKeyResource{}
}

type variableSshKeyResource struct {
	client *buddy.Client
}

type variableSshKeyResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Domain         types.String `tfsdk:"domain"`
	Key            types.String `tfsdk:"key"`
	Value          types.String `tfsdk:"value"`
	FilePlace      types.String `tfsdk:"file_place"`
	FilePath       types.String `tfsdk:"file_path"`
	FileChmod      types.String `tfsdk:"file_chmod"`
	ProjectName    types.String `tfsdk:"project_name"`
	PipelineId     types.Int64  `tfsdk:"pipeline_id"`
	ActionId       types.Int64  `tfsdk:"action_id"`
	Settable       types.Bool   `tfsdk:"settable"`
	Description    types.String `tfsdk:"description"`
	VariableId     types.Int64  `tfsdk:"variable_id"`
	ValueProcessed types.String `tfsdk:"value_processed"`
	Encrypted      types.Bool   `tfsdk:"encrypted"`
	Checksum       types.String `tfsdk:"checksum"`
	KeyFingerprint types.String `tfsdk:"key_fingerprint"`
	PublicValue    types.String `tfsdk:"public_value"`
}

func (r *variableSshKeyResourceModel) loadAPI(domain string, variable *buddy.Variable) {
	r.ID = types.StringValue(util.ComposeDoubleId(domain, strconv.Itoa(variable.Id)))
	r.Domain = types.StringValue(domain)
	r.Key = types.StringValue(variable.Key)
	r.FilePlace = types.StringValue(variable.FilePlace)
	r.FilePath = types.StringValue(variable.FilePath)
	r.FileChmod = types.StringValue(variable.FileChmod)
	r.Settable = types.BoolValue(variable.Settable)
	r.Description = types.StringValue(variable.Description)
	r.VariableId = types.Int64Value(int64(variable.Id))
	r.ValueProcessed = types.StringValue(variable.Value)
	r.Encrypted = types.BoolValue(variable.Encrypted)
	r.Checksum = types.StringValue(variable.Checksum)
	r.KeyFingerprint = types.StringValue(variable.KeyFingerprint)
	r.PublicValue = types.StringValue(variable.PublicValue)
}

func (r *variableSshKeyResourceModel) decomposeId() (string, int, error) {
	domain, vid, err := util.DecomposeDoubleId(r.ID.ValueString())
	if err != nil {
		return "", 0, err
	}
	variableId, err := strconv.Atoi(vid)
	if err != nil {
		return "", 0, err
	}
	return domain, variableId, nil
}

func (r *variableSshKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_variable_ssh_key"
}

func (r *variableSshKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage a variable of SSH key type\n\n" +
			"Workspace administrator rights are required\n\n" +
			"Token scope required: `WORKSPACE`, `VARIABLE_ADD`, `VARIABLE_MANAGE`, `VARIABLE_INFO`",
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
			"key": schema.StringAttribute{
				MarkdownDescription: "The variable's name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "The variable's value",
				Required:            true,
				Sensitive:           true,
			},
			"file_place": schema.StringAttribute{
				MarkdownDescription: "Should the variable's be copied to an action's container in **file_path** (`CONTAINER`, `NONE`)",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						buddy.VariableSshKeyFilePlaceContainer,
						buddy.VariableSshKeyFilePlaceNone,
					),
				},
			},
			"file_path": schema.StringAttribute{
				MarkdownDescription: "The variable's path in the action's container",
				Required:            true,
			},
			"file_chmod": schema.StringAttribute{
				MarkdownDescription: "The variable's file permission in an action's container",
				Required:            true,
			},
			"project_name": schema.StringAttribute{
				MarkdownDescription: "The variable's project name",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"pipeline_id": schema.Int64Attribute{
				MarkdownDescription: "The variable's pipeline ID",
				Optional:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"action_id": schema.Int64Attribute{
				MarkdownDescription: "The variable's action ID",
				Optional:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"settable": schema.BoolAttribute{
				MarkdownDescription: "Is the variable's value changeable, always false for buddy_variable_ssh_key",
				Computed:            true,
			},
			"encrypted": schema.BoolAttribute{
				MarkdownDescription: "Is the variable's value encrypted, always true for buddy_variable_ssh_key",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The variable's description",
				Optional:            true,
				Computed:            true,
			},
			"variable_id": schema.Int64Attribute{
				MarkdownDescription: "The variable's ID",
				Computed:            true,
			},
			"value_processed": schema.StringAttribute{
				MarkdownDescription: "The variable's value, always encrypted for buddy_variable_ssh_key",
				Computed:            true,
				Sensitive:           true,
			},
			"checksum": schema.StringAttribute{
				MarkdownDescription: "The variable's checksum",
				Computed:            true,
			},
			"key_fingerprint": schema.StringAttribute{
				MarkdownDescription: "The variable's fingerprint",
				Computed:            true,
			},
			"public_value": schema.StringAttribute{
				MarkdownDescription: "The variable's public key",
				Computed:            true,
			},
		},
	}
}

func (r *variableSshKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *variableSshKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *variableSshKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	typ := buddy.VariableTypeSshKey
	encrypted := true
	settable := false
	ops := buddy.VariableOps{
		Key:       data.Key.ValueStringPointer(),
		Value:     data.Value.ValueStringPointer(),
		Type:      &typ,
		Encrypted: &encrypted,
		Settable:  &settable,
		FilePlace: data.FilePlace.ValueStringPointer(),
		FilePath:  data.FilePath.ValueStringPointer(),
		FileChmod: data.FileChmod.ValueStringPointer(),
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		ops.Description = data.Description.ValueStringPointer()
	}
	if !data.ProjectName.IsNull() && !data.ProjectName.IsUnknown() {
		ops.Project = &buddy.VariableProject{
			Name: data.ProjectName.ValueString(),
		}
	}
	if !data.PipelineId.IsNull() && !data.PipelineId.IsUnknown() {
		ops.Pipeline = &buddy.VariablePipeline{
			Id: int(data.PipelineId.ValueInt64()),
		}
	}
	if !data.ActionId.IsNull() && !data.ActionId.IsUnknown() {
		ops.Action = &buddy.VariableAction{
			Id: int(data.ActionId.ValueInt64()),
		}
	}
	variable, _, err := r.client.VariableService.Create(domain, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create variable ssh key", err))
		return
	}
	data.loadAPI(domain, variable)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *variableSshKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *variableSshKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, variableId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("variable ssh key", err))
		return
	}
	variable, httpResp, err := r.client.VariableService.Get(domain, variableId)
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get variable ssh key", err))
		return
	}
	if variable.Type != buddy.VariableTypeSshKey {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get variable ssh key", fmt.Errorf("variable not found")))
		return
	}
	data.loadAPI(domain, variable)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *variableSshKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *variableSshKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, variableId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("variable ssh key", err))
		return
	}
	typ := buddy.VariableTypeSshKey
	encrypted := true
	settable := false
	ops := buddy.VariableOps{
		Value:     data.Value.ValueStringPointer(),
		Type:      &typ,
		Encrypted: &encrypted,
		Settable:  &settable,
		FilePlace: data.FilePlace.ValueStringPointer(),
		FilePath:  data.FilePath.ValueStringPointer(),
		FileChmod: data.FileChmod.ValueStringPointer(),
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		ops.Description = data.Description.ValueStringPointer()
	}
	variable, _, err := r.client.VariableService.Update(domain, variableId, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update variable ssh key", err))
		return
	}
	data.loadAPI(domain, variable)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *variableSshKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *variableSshKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, variableId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("variable ssh key", err))
		return
	}
	_, err = r.client.VariableService.Delete(domain, variableId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete variable ssh key", err))
	}
}

func (r *variableSshKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

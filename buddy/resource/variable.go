package resource

import (
	"buddy-terraform/buddy/util"
	"context"
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
)

var (
	_ resource.Resource                = &variableResource{}
	_ resource.ResourceWithConfigure   = &variableResource{}
	_ resource.ResourceWithImportState = &variableResource{}
)

func NewVariableResource() resource.Resource {
	return &variableResource{}
}

type variableResource struct {
	client *buddy.Client
}

type variableResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Domain         types.String `tfsdk:"domain"`
	Key            types.String `tfsdk:"key"`
	Value          types.String `tfsdk:"value"`
	Encrypted      types.Bool   `tfsdk:"encrypted"`
	ProjectName    types.String `tfsdk:"project_name"`
	PipelineId     types.Int64  `tfsdk:"pipeline_id"`
	ActionId       types.Int64  `tfsdk:"action_id"`
	Settable       types.Bool   `tfsdk:"settable"`
	Description    types.String `tfsdk:"description"`
	ValueProcessed types.String `tfsdk:"value_processed"`
	VariableId     types.Int64  `tfsdk:"variable_id"`
}

func (r *variableResourceModel) decomposeId() (string, int, error) {
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

func (r *variableResourceModel) loadAPI(domain string, variable *buddy.Variable) {
	r.ID = types.StringValue(util.ComposeDoubleId(domain, strconv.Itoa(variable.Id)))
	r.Domain = types.StringValue(domain)
	r.Key = types.StringValue(variable.Key)
	r.Encrypted = types.BoolValue(variable.Encrypted)
	r.Settable = types.BoolValue(variable.Settable)
	r.Description = types.StringValue(variable.Description)
	r.ValueProcessed = types.StringValue(variable.Value)
	r.VariableId = types.Int64Value(int64(variable.Id))
}

func (r *variableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_variable"
}

func (r *variableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage a variable\n\n" +
			"Workspace administrator rights are required\n\n" +
			"Token scopes required: `WORKSPACE`, `VARIABLE_ADD`, `VARIABLE_MANAGE`, `VARIABLE_INFO`",
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
			"encrypted": schema.BoolAttribute{
				MarkdownDescription: "Is the variable's value encrypted",
				Optional:            true,
				Computed:            true,
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
				MarkdownDescription: "Is the variable's value changeable",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The variable's description",
				Optional:            true,
				Computed:            true,
			},
			"value_processed": schema.StringAttribute{
				MarkdownDescription: "The variable's processed value. Encrypted if **encrypted** == true",
				Computed:            true,
				Sensitive:           true,
			},
			"variable_id": schema.Int64Attribute{
				MarkdownDescription: "The variable's ID",
				Computed:            true,
			},
		},
	}
}

func (r *variableResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *variableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *variableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	typ := buddy.VariableTypeVar
	ops := buddy.VariableOps{
		Key:   data.Key.ValueStringPointer(),
		Value: data.Value.ValueStringPointer(),
		Type:  &typ,
	}
	if !data.Settable.IsNull() && !data.Settable.IsUnknown() {
		ops.Settable = data.Settable.ValueBoolPointer()
	}
	if !data.Encrypted.IsNull() && !data.Encrypted.IsUnknown() {
		ops.Encrypted = data.Encrypted.ValueBoolPointer()
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
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create variable", err))
		return
	}
	data.loadAPI(domain, variable)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *variableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *variableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, variableId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("variable", err))
		return
	}
	variable, httpResp, err := r.client.VariableService.Get(domain, variableId)
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get variable", err))
		return
	}
	if variable.Type != buddy.VariableTypeVar {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get variable", fmt.Errorf("variable not found")))
		return
	}
	data.loadAPI(domain, variable)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *variableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *variableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, variableId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("variable", err))
		return
	}
	ops := buddy.VariableOps{
		Value: data.Value.ValueStringPointer(),
	}
	if !data.Encrypted.IsNull() && !data.Encrypted.IsUnknown() {
		ops.Encrypted = data.Encrypted.ValueBoolPointer()
	}
	if !data.Settable.IsNull() && !data.Settable.IsUnknown() {
		ops.Settable = data.Settable.ValueBoolPointer()
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		ops.Description = data.Description.ValueStringPointer()
	}
	variable, _, err := r.client.VariableService.Update(domain, variableId, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update variable", err))
		return
	}
	data.loadAPI(domain, variable)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *variableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *variableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, variableId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("variable", err))
		return
	}
	_, err = r.client.VariableService.Delete(domain, variableId)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete variable", err))
	}
}

func (r *variableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

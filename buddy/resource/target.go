package resource

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ resource.Resource                = &targetResource{}
	_ resource.ResourceWithConfigure   = &targetResource{}
	_ resource.ResourceWithImportState = &targetResource{}
)

func NewTargetResource() resource.Resource { return &targetResource{} }

type targetResource struct{ client *buddy.Client }

type targetResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Domain        types.String `tfsdk:"domain"`
	HtmlUrl       types.String `tfsdk:"html_url"`
	TargetId      types.String `tfsdk:"target_id"`
	Name          types.String `tfsdk:"name"`
	Identifier    types.String `tfsdk:"identifier"`
	Type          types.String `tfsdk:"type"`
	Host          types.String `tfsdk:"host"`
	Repository    types.String `tfsdk:"repository"`
	Port          types.String `tfsdk:"port"`
	Path          types.String `tfsdk:"path"`
	Secure        types.Bool   `tfsdk:"secure"`
	Integration   types.String `tfsdk:"integration"`
	Disabled      types.Bool   `tfsdk:"disabled"`
	Tags          types.Set    `tfsdk:"tags"`
	ProjectName   types.String `tfsdk:"project_name"`
	PipelineId    types.Int64  `tfsdk:"pipeline_id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	Permissions   types.Set    `tfsdk:"permissions"`
}

func (m *targetResourceModel) decomposeId() (string, string, error) {
	domain, tid, err := util.DecomposeDoubleId(m.ID.ValueString())
	if err != nil {
		return "", "", err
	}
	return domain, tid, nil
}

func (m *targetResourceModel) loadAPI(ctx context.Context, domain string, target *buddy.Target) diag.Diagnostics {
	var diags diag.Diagnostics
	m.ID = types.StringValue(util.ComposeDoubleId(domain, target.Id))
	m.Domain = types.StringValue(domain)
	m.HtmlUrl = types.StringValue(target.HtmlUrl)
	m.TargetId = types.StringValue(target.Id)
	m.Name = types.StringValue(target.Name)
	m.Identifier = types.StringValue(target.Identifier)
	m.Type = types.StringValue(target.Type)
	m.Host = types.StringValue(target.Host)
	m.Repository = types.StringValue(target.Repository)
	m.Port = types.StringValue(target.Port)
	m.Path = types.StringValue(target.Path)
	m.Secure = types.BoolValue(target.Secure)
	m.Integration = types.StringValue(target.Integration)
	m.Disabled = types.BoolValue(target.Disabled)
	t, d := types.SetValueFrom(ctx, types.StringType, &target.Tags)
	diags.Append(d...)
	m.Tags = t
	if target.Project != nil {
		m.ProjectName = types.StringValue(target.Project.Name)
	} else {
		m.ProjectName = types.StringNull()
	}
	if target.Pipeline != nil {
		m.PipelineId = types.Int64Value(int64(target.Pipeline.Id))
	} else {
		m.PipelineId = types.Int64Null()
	}
	if target.Environment != nil {
		m.EnvironmentId = types.StringValue(target.Environment.Id)
	} else {
		m.EnvironmentId = types.StringNull()
	}
}

func (r *targetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_target"
}

func (r *targetResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *targetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage a deployment target",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"domain": schema.StringAttribute{
				Required:      true,
				Validators:    util.StringValidatorsDomain(),
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"html_url":  schema.StringAttribute{Computed: true},
			"target_id": schema.StringAttribute{Computed: true},
			"name":      schema.StringAttribute{Required: true},
			"identifier": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Validators:    util.StringValidatorIdentifier(),
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
			},
			"type":           schema.StringAttribute{Required: true},
			"host":           schema.StringAttribute{Optional: true, Computed: true},
			"repository":     schema.StringAttribute{Optional: true, Computed: true},
			"port":           schema.StringAttribute{Optional: true, Computed: true},
			"path":           schema.StringAttribute{Optional: true, Computed: true},
			"secure":         schema.BoolAttribute{Optional: true, Computed: true},
			"integration":    schema.StringAttribute{Optional: true, Computed: true},
			"disabled":       schema.BoolAttribute{Optional: true, Computed: true},
			"tags":           schema.SetAttribute{Optional: true, Computed: true, ElementType: types.StringType},
			"project_name":   schema.StringAttribute{Optional: true, Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"pipeline_id":    schema.Int64Attribute{Optional: true, Computed: true, PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()}},
			"environment_id": schema.StringAttribute{Optional: true, Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"permissions": schema.SetNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"others": schema.StringAttribute{Optional: true},
					},
					Blocks: map[string]schema.Block{
						"user":  schema.SetNestedBlock{NestedObject: schema.NestedBlockObject{Attributes: util.TargetPermissionsAccessModelAttributes()}},
						"group": schema.SetNestedBlock{NestedObject: schema.NestedBlockObject{Attributes: util.TargetPermissionsAccessModelAttributes()}},
					},
				},
			},
		},
	}
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
	if !data.Identifier.IsNull() && !data.Identifier.IsUnknown() {
		ops.Identifier = data.Identifier.ValueStringPointer()
	}
	if !data.Host.IsNull() && !data.Host.IsUnknown() {
		ops.Host = data.Host.ValueStringPointer()
	}
	if !data.Repository.IsNull() && !data.Repository.IsUnknown() {
		ops.Repository = data.Repository.ValueStringPointer()
	}
	if !data.Port.IsNull() && !data.Port.IsUnknown() {
		ops.Port = data.Port.ValueStringPointer()
	}
	if !data.Path.IsNull() && !data.Path.IsUnknown() {
		ops.Path = data.Path.ValueStringPointer()
	}
	if !data.Secure.IsNull() && !data.Secure.IsUnknown() {
		ops.Secure = data.Secure.ValueBoolPointer()
	}
	if !data.Integration.IsNull() && !data.Integration.IsUnknown() {
		ops.Integration = data.Integration.ValueStringPointer()
	}
	if !data.Disabled.IsNull() && !data.Disabled.IsUnknown() {
		ops.Disabled = data.Disabled.ValueBoolPointer()
	}
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		tags, d := util.StringSetToApi(ctx, &data.Tags)
		resp.Diagnostics.Append(d...)
		ops.Tags = tags
	}
	if !data.ProjectName.IsNull() && !data.ProjectName.IsUnknown() {
		ops.Project = &buddy.TargetProject{Name: data.ProjectName.ValueString()}
	}
	if !data.PipelineId.IsNull() && !data.PipelineId.IsUnknown() {
		ops.Pipeline = &buddy.TargetPipeline{Id: int(data.PipelineId.ValueInt64())}
	}
	if !data.EnvironmentId.IsNull() && !data.EnvironmentId.IsUnknown() {
		ops.Environment = &buddy.TargetEnvironment{Id: data.EnvironmentId.ValueString()}
	}
	if !data.Permissions.IsNull() && !data.Permissions.IsUnknown() {
		perm, d := util.TargetPermissionsModelToApi(ctx, &data.Permissions)
		resp.Diagnostics.Append(d...)
		ops.Permissions = perm
	}
	target, _, err := r.client.TargetService.Create(domain, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create target", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, target)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *targetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *targetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, tid, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("target", err))
		return
	}
	target, httpResp, err := r.client.TargetService.Get(domain, tid)
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get target", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, target)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *targetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *targetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, tid, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("target", err))
		return
	}
	ops := buddy.TargetOps{
		Name: data.Name.ValueStringPointer(),
		Type: data.Type.ValueStringPointer(),
	}
	if !data.Identifier.IsNull() && !data.Identifier.IsUnknown() {
		ops.Identifier = data.Identifier.ValueStringPointer()
	}
	if !data.Host.IsNull() && !data.Host.IsUnknown() {
		ops.Host = data.Host.ValueStringPointer()
	}
	if !data.Repository.IsNull() && !data.Repository.IsUnknown() {
		ops.Repository = data.Repository.ValueStringPointer()
	}
	if !data.Port.IsNull() && !data.Port.IsUnknown() {
		ops.Port = data.Port.ValueStringPointer()
	}
	if !data.Path.IsNull() && !data.Path.IsUnknown() {
		ops.Path = data.Path.ValueStringPointer()
	}
	if !data.Secure.IsNull() && !data.Secure.IsUnknown() {
		ops.Secure = data.Secure.ValueBoolPointer()
	}
	if !data.Integration.IsNull() && !data.Integration.IsUnknown() {
		ops.Integration = data.Integration.ValueStringPointer()
	}
	if !data.Disabled.IsNull() && !data.Disabled.IsUnknown() {
		ops.Disabled = data.Disabled.ValueBoolPointer()
	}
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		tags, d := util.StringSetToApi(ctx, &data.Tags)
		resp.Diagnostics.Append(d...)
		ops.Tags = tags
	}
	if !data.Permissions.IsNull() && !data.Permissions.IsUnknown() {
		perm, d := util.TargetPermissionsModelToApi(ctx, &data.Permissions)
		resp.Diagnostics.Append(d...)
		ops.Permissions = perm
	}
	target, _, err := r.client.TargetService.Update(domain, tid, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update target", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, target)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *targetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *targetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain, tid, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("target", err))
		return
	}
	_, err = r.client.TargetService.Delete(domain, tid)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete target", err))
	}
}

func (r *targetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

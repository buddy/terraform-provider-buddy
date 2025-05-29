package resource

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
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
	ID                  types.String `tfsdk:"id"`
	Domain              types.String `tfsdk:"domain"`
	ProjectName         types.String `tfsdk:"project_name"`
	TargetId            types.Int64  `tfsdk:"target_id"`
	Name                types.String `tfsdk:"name"`
	Type                types.String `tfsdk:"type"`
	Hostname            types.String `tfsdk:"hostname"`
	Port                types.Int64  `tfsdk:"port"`
	Username            types.String `tfsdk:"username"`
	Password            types.String `tfsdk:"password"`
	Passphrase          types.String `tfsdk:"passphrase"`
	KeyId               types.Int64  `tfsdk:"key_id"`
	FilePath            types.String `tfsdk:"file_path"`
	AuthMode            types.String `tfsdk:"auth_mode"`
	Tags                types.Set    `tfsdk:"tags"`
	Description         types.String `tfsdk:"description"`
	AllPipelinesAllowed types.Bool   `tfsdk:"all_pipelines_allowed"`
	AllowedPipelines    types.Set    `tfsdk:"allowed_pipelines"`
	HtmlUrl             types.String `tfsdk:"html_url"`
}

func (r *targetResourceModel) loadAPI(ctx context.Context, domain string, projectName string, target *buddy.Target) diag.Diagnostics {
	var diags diag.Diagnostics
	
	if projectName != "" {
		r.ID = types.StringValue(util.ComposeTripleId(domain, projectName, strconv.Itoa(target.Id)))
		r.ProjectName = types.StringValue(projectName)
	} else {
		r.ID = types.StringValue(util.ComposeDoubleId(domain, strconv.Itoa(target.Id)))
		r.ProjectName = types.StringNull()
	}
	
	r.Domain = types.StringValue(domain)
	r.TargetId = types.Int64Value(int64(target.Id))
	r.Name = types.StringValue(target.Name)
	r.Type = types.StringValue(target.Type)
	r.HtmlUrl = types.StringValue(target.HtmlUrl)
	
	if target.Hostname != "" {
		r.Hostname = types.StringValue(target.Hostname)
	} else {
		r.Hostname = types.StringNull()
	}
	
	if target.Port > 0 {
		r.Port = types.Int64Value(int64(target.Port))
	} else {
		r.Port = types.Int64Null()
	}
	
	if target.Username != "" {
		r.Username = types.StringValue(target.Username)
	} else {
		r.Username = types.StringNull()
	}
	
	if target.Password != "" {
		r.Password = types.StringValue(target.Password)
	} else {
		r.Password = types.StringNull()
	}
	
	if target.Passphrase != "" {
		r.Passphrase = types.StringValue(target.Passphrase)
	} else {
		r.Passphrase = types.StringNull()
	}
	
	if target.KeyId > 0 {
		r.KeyId = types.Int64Value(int64(target.KeyId))
	} else {
		r.KeyId = types.Int64Null()
	}
	
	if target.FilePath != "" {
		r.FilePath = types.StringValue(target.FilePath)
	} else {
		r.FilePath = types.StringNull()
	}
	
	if target.AuthMode != "" {
		r.AuthMode = types.StringValue(target.AuthMode)
	} else {
		r.AuthMode = types.StringNull()
	}
	
	if target.Description != "" {
		r.Description = types.StringValue(target.Description)
	} else {
		r.Description = types.StringNull()
	}
	
	r.AllPipelinesAllowed = types.BoolValue(target.AllPipelinesAllowed)
	
	tags, d := types.SetValueFrom(ctx, types.StringType, &target.Tags)
	diags.Append(d...)
	r.Tags = tags
	
	if len(target.AllowedPipelines) > 0 {
		allowedPipelines, d := types.SetValueFrom(ctx, types.Int64Type, &target.AllowedPipelines)
		diags.Append(d...)
		r.AllowedPipelines = allowedPipelines
	} else {
		r.AllowedPipelines = types.SetNull(types.Int64Type)
	}
	
	return diags
}

func (r *targetResourceModel) decomposeId() (string, string, int, error) {
	if r.ProjectName.IsNull() {
		domain, tid, err := util.DecomposeDoubleId(r.ID.ValueString())
		if err != nil {
			return "", "", 0, err
		}
		targetId, err := strconv.Atoi(tid)
		if err != nil {
			return "", "", 0, err
		}
		return domain, "", targetId, nil
	}
	
	domain, projectName, tid, err := util.DecomposeTripleId(r.ID.ValueString())
	if err != nil {
		return "", "", 0, err
	}
	targetId, err := strconv.Atoi(tid)
	if err != nil {
		return "", "", 0, err
	}
	return domain, projectName, targetId, nil
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
			"project_name": schema.StringAttribute{
				MarkdownDescription: "The project's name. Required if the target should be created in project scope",
				Optional:            true,
				Validators:          util.StringValidatorsSlug(),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_id": schema.Int64Attribute{
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
			"hostname": schema.StringAttribute{
				MarkdownDescription: "The target's hostname or IP address",
				Optional:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "The target's port",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The target's username",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The target's password",
				Optional:            true,
				Sensitive:           true,
			},
			"passphrase": schema.StringAttribute{
				MarkdownDescription: "The SSH key passphrase",
				Optional:            true,
				Sensitive:           true,
			},
			"key_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the SSH key to use for authentication",
				Optional:            true,
			},
			"file_path": schema.StringAttribute{
				MarkdownDescription: "The remote file path on the target",
				Optional:            true,
			},
			"auth_mode": schema.StringAttribute{
				MarkdownDescription: "The authentication mode",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("PASS"),
				Validators: []validator.String{
					stringvalidator.OneOf("PASS", "KEY", "BOTH"),
				},
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "The target's tags",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The target's description",
				Optional:            true,
			},
			"all_pipelines_allowed": schema.BoolAttribute{
				MarkdownDescription: "Whether all pipelines can use this target",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"allowed_pipelines": schema.SetAttribute{
				MarkdownDescription: "The list of pipeline IDs that are allowed to use this target",
				Optional:            true,
				ElementType:         types.Int64Type,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
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
	projectName := data.ProjectName.ValueString()
	
	ops := buddy.TargetOps{
		Name: data.Name.ValueStringPointer(),
		Type: data.Type.ValueStringPointer(),
	}
	
	if !data.Hostname.IsNull() {
		ops.Hostname = data.Hostname.ValueStringPointer()
	}
	
	if !data.Port.IsNull() {
		port := int(data.Port.ValueInt64())
		ops.Port = &port
	}
	
	if !data.Username.IsNull() {
		ops.Username = data.Username.ValueStringPointer()
	}
	
	if !data.Password.IsNull() {
		ops.Password = data.Password.ValueStringPointer()
	}
	
	if !data.Passphrase.IsNull() {
		ops.Passphrase = data.Passphrase.ValueStringPointer()
	}
	
	if !data.KeyId.IsNull() {
		keyId := int(data.KeyId.ValueInt64())
		ops.KeyId = &keyId
	}
	
	if !data.FilePath.IsNull() {
		ops.FilePath = data.FilePath.ValueStringPointer()
	}
	
	if !data.AuthMode.IsNull() {
		ops.AuthMode = data.AuthMode.ValueStringPointer()
	}
	
	if !data.Description.IsNull() {
		ops.Description = data.Description.ValueStringPointer()
	}
	
	if !data.AllPipelinesAllowed.IsNull() {
		allPipelinesAllowed := data.AllPipelinesAllowed.ValueBool()
		ops.AllPipelinesAllowed = &allPipelinesAllowed
	}
	
	if !data.Tags.IsNull() {
		var tags []string
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Tags = &tags
	}
	
	if !data.AllowedPipelines.IsNull() {
		var allowedPipelines []int64
		resp.Diagnostics.Append(data.AllowedPipelines.ElementsAs(ctx, &allowedPipelines, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		intPipelines := make([]int, len(allowedPipelines))
		for i, p := range allowedPipelines {
			intPipelines[i] = int(p)
		}
		ops.AllowedPipelines = &intPipelines
	}
	
	var target *buddy.Target
	var err error
	
	if projectName != "" {
		target, _, err = r.client.TargetService.CreateInProject(domain, projectName, &ops)
	} else {
		target, _, err = r.client.TargetService.CreateInWorkspace(domain, &ops)
	}
	
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("create target", err))
		return
	}
	
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, target)...)
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
	
	domain, projectName, targetId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("target", err))
		return
	}
	
	var target *buddy.Target
	
	if projectName != "" {
		target, _, err = r.client.TargetService.GetInProject(domain, projectName, targetId)
	} else {
		target, _, err = r.client.TargetService.GetInWorkspace(domain, targetId)
	}
	
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get target", err))
		return
	}
	
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, target)...)
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
	
	domain, projectName, targetId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("target", err))
		return
	}
	
	ops := buddy.TargetOps{}
	
	if !data.Name.IsNull() {
		ops.Name = data.Name.ValueStringPointer()
	}
	
	if !data.Hostname.IsNull() {
		ops.Hostname = data.Hostname.ValueStringPointer()
	}
	
	if !data.Port.IsNull() {
		port := int(data.Port.ValueInt64())
		ops.Port = &port
	}
	
	if !data.Username.IsNull() {
		ops.Username = data.Username.ValueStringPointer()
	}
	
	if !data.Password.IsNull() {
		ops.Password = data.Password.ValueStringPointer()
	}
	
	if !data.Passphrase.IsNull() {
		ops.Passphrase = data.Passphrase.ValueStringPointer()
	}
	
	if !data.KeyId.IsNull() {
		keyId := int(data.KeyId.ValueInt64())
		ops.KeyId = &keyId
	}
	
	if !data.FilePath.IsNull() {
		ops.FilePath = data.FilePath.ValueStringPointer()
	}
	
	if !data.AuthMode.IsNull() {
		ops.AuthMode = data.AuthMode.ValueStringPointer()
	}
	
	if !data.Description.IsNull() {
		ops.Description = data.Description.ValueStringPointer()
	}
	
	if !data.AllPipelinesAllowed.IsNull() {
		allPipelinesAllowed := data.AllPipelinesAllowed.ValueBool()
		ops.AllPipelinesAllowed = &allPipelinesAllowed
	}
	
	if !data.Tags.IsNull() {
		var tags []string
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		ops.Tags = &tags
	}
	
	if !data.AllowedPipelines.IsNull() {
		var allowedPipelines []int64
		resp.Diagnostics.Append(data.AllowedPipelines.ElementsAs(ctx, &allowedPipelines, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		intPipelines := make([]int, len(allowedPipelines))
		for i, p := range allowedPipelines {
			intPipelines[i] = int(p)
		}
		ops.AllowedPipelines = &intPipelines
	}
	
	var target *buddy.Target
	
	if projectName != "" {
		target, _, err = r.client.TargetService.UpdateInProject(domain, projectName, targetId, &ops)
	} else {
		target, _, err = r.client.TargetService.UpdateInWorkspace(domain, targetId, &ops)
	}
	
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update target", err))
		return
	}
	
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, target)...)
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
	
	domain, projectName, targetId, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("target", err))
		return
	}
	
	if projectName != "" {
		_, err = r.client.TargetService.DeleteInProject(domain, projectName, targetId)
	} else {
		_, err = r.client.TargetService.DeleteInWorkspace(domain, targetId)
	}
	
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete target", err))
		return
	}
}

func (r *targetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
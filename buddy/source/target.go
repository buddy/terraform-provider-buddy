package source

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ datasource.DataSource              = &targetSource{}
	_ datasource.DataSourceWithConfigure = &targetSource{}
)

func NewTargetSource() datasource.DataSource {
	return &targetSource{}
}

type targetSource struct {
	client *buddy.Client
}

type targetSourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Domain              types.String `tfsdk:"domain"`
	ProjectName         types.String `tfsdk:"project_name"`
	TargetId            types.Int64  `tfsdk:"target_id"`
	Name                types.String `tfsdk:"name"`
	Type                types.String `tfsdk:"type"`
	Hostname            types.String `tfsdk:"hostname"`
	Port                types.Int64  `tfsdk:"port"`
	Username            types.String `tfsdk:"username"`
	FilePath            types.String `tfsdk:"file_path"`
	AuthMode            types.String `tfsdk:"auth_mode"`
	Tags                types.Set    `tfsdk:"tags"`
	Description         types.String `tfsdk:"description"`
	AllPipelinesAllowed types.Bool   `tfsdk:"all_pipelines_allowed"`
	AllowedPipelines    types.Set    `tfsdk:"allowed_pipelines"`
	HtmlUrl             types.String `tfsdk:"html_url"`
}

func (s *targetSourceModel) loadAPI(ctx context.Context, domain string, projectName string, target *buddy.Target) diag.Diagnostics {
	var diags diag.Diagnostics
	
	if projectName != "" {
		s.ID = types.StringValue(util.ComposeTripleId(domain, projectName, strconv.Itoa(target.Id)))
		s.ProjectName = types.StringValue(projectName)
	} else {
		s.ID = types.StringValue(util.ComposeDoubleId(domain, strconv.Itoa(target.Id)))
		s.ProjectName = types.StringNull()
	}
	
	s.Domain = types.StringValue(domain)
	s.TargetId = types.Int64Value(int64(target.Id))
	s.Name = types.StringValue(target.Name)
	s.Type = types.StringValue(target.Type)
	s.HtmlUrl = types.StringValue(target.HtmlUrl)
	
	if target.Hostname != "" {
		s.Hostname = types.StringValue(target.Hostname)
	} else {
		s.Hostname = types.StringNull()
	}
	
	if target.Port > 0 {
		s.Port = types.Int64Value(int64(target.Port))
	} else {
		s.Port = types.Int64Null()
	}
	
	if target.Username != "" {
		s.Username = types.StringValue(target.Username)
	} else {
		s.Username = types.StringNull()
	}
	
	if target.FilePath != "" {
		s.FilePath = types.StringValue(target.FilePath)
	} else {
		s.FilePath = types.StringNull()
	}
	
	if target.AuthMode != "" {
		s.AuthMode = types.StringValue(target.AuthMode)
	} else {
		s.AuthMode = types.StringNull()
	}
	
	if target.Description != "" {
		s.Description = types.StringValue(target.Description)
	} else {
		s.Description = types.StringNull()
	}
	
	s.AllPipelinesAllowed = types.BoolValue(target.AllPipelinesAllowed)
	
	tags, d := types.SetValueFrom(ctx, types.StringType, &target.Tags)
	diags.Append(d...)
	s.Tags = tags
	
	if len(target.AllowedPipelines) > 0 {
		allowedPipelines, d := types.SetValueFrom(ctx, types.Int64Type, &target.AllowedPipelines)
		diags.Append(d...)
		s.AllowedPipelines = allowedPipelines
	} else {
		s.AllowedPipelines = types.SetNull(types.Int64Type)
	}
	
	return diags
}

func (s *targetSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_target"
}

func (s *targetSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get target by name or ID\n\n" +
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
			"project_name": schema.StringAttribute{
				MarkdownDescription: "The project's name. Required if the target is in project scope",
				Optional:            true,
				Validators:          util.StringValidatorsSlug(),
			},
			"target_id": schema.Int64Attribute{
				MarkdownDescription: "The target's ID",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The target's name",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The target's type",
				Computed:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "The target's hostname or IP address",
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "The target's port",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The target's username",
				Computed:            true,
			},
			"file_path": schema.StringAttribute{
				MarkdownDescription: "The remote file path on the target",
				Computed:            true,
			},
			"auth_mode": schema.StringAttribute{
				MarkdownDescription: "The authentication mode",
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "The target's tags",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The target's description",
				Computed:            true,
			},
			"all_pipelines_allowed": schema.BoolAttribute{
				MarkdownDescription: "Whether all pipelines can use this target",
				Computed:            true,
			},
			"allowed_pipelines": schema.SetAttribute{
				MarkdownDescription: "The list of pipeline IDs that are allowed to use this target",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"html_url": schema.StringAttribute{
				MarkdownDescription: "The target's URL",
				Computed:            true,
			},
		},
	}
}

func (s *targetSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *targetSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("target_id"),
			path.MatchRoot("name"),
		),
	}
}

func (s *targetSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *targetSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	domain := data.Domain.ValueString()
	projectName := data.ProjectName.ValueString()
	
	var targets *buddy.Targets
	var err error
	
	if projectName != "" {
		// Get targets from project
		if data.TargetId.IsNull() && data.Name.IsNull() {
			resp.Diagnostics.AddError("Missing required argument", "Either target_id or name must be specified")
			return
		}
		
		targets, _, err = s.client.TargetService.GetListInProject(domain, projectName)
	} else {
		// Get targets from workspace
		if data.TargetId.IsNull() && data.Name.IsNull() {
			resp.Diagnostics.AddError("Missing required argument", "Either target_id or name must be specified")
			return
		}
		
		targets, _, err = s.client.TargetService.GetListInWorkspace(domain)
	}
	
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get targets", err))
		return
	}
	
	var target *buddy.Target
	
	if !data.TargetId.IsNull() {
		targetId := int(data.TargetId.ValueInt64())
		for _, t := range targets.Targets {
			if t.Id == targetId {
				target = &t
				break
			}
		}
		if target == nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("target"))
			return
		}
	} else if !data.Name.IsNull() {
		name := data.Name.ValueString()
		for _, t := range targets.Targets {
			if t.Name == name {
				target = &t
				break
			}
		}
		if target == nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("target"))
			return
		}
	}
	
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, target)...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
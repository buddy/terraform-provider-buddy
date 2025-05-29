package source

import (
	"context"

	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ datasource.DataSource              = &targetsSource{}
	_ datasource.DataSourceWithConfigure = &targetsSource{}
)

// NewTargetsSource creates a new instance of the targets data source
func NewTargetsSource() datasource.DataSource {
	return &targetsSource{}
}

type targetsSource struct {
	client *buddy.Client
}

type targetsSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Domain      types.String `tfsdk:"domain"`
	ProjectName types.String `tfsdk:"project_name"`
	NameRegex   types.String `tfsdk:"name_regex"`
	Targets     types.Set    `tfsdk:"targets"`
}

func (s *targetsSourceModel) loadAPI(ctx context.Context, domain string, projectName string, targets *[]buddy.Target) diag.Diagnostics {
	var diags diag.Diagnostics
	
	if projectName != "" {
		s.ID = types.StringValue(util.ComposeDoubleId(domain, projectName))
		s.ProjectName = types.StringValue(projectName)
	} else {
		s.ID = types.StringValue(domain)
		s.ProjectName = types.StringNull()
	}
	
	s.Domain = types.StringValue(domain)
	
	targetsList := make([]targetModel, 0)
	for _, target := range *targets {
		t := targetModel{
			Name:                types.StringValue(target.Name),
			TargetId:            types.Int64Value(int64(target.Id)),
			Type:                types.StringValue(target.Type),
			HtmlUrl:             types.StringValue(target.HtmlUrl),
			AllPipelinesAllowed: types.BoolValue(target.AllPipelinesAllowed),
		}
		
		if target.Hostname != "" {
			t.Hostname = types.StringValue(target.Hostname)
		} else {
			t.Hostname = types.StringNull()
		}
		
		if target.Port > 0 {
			t.Port = types.Int64Value(int64(target.Port))
		} else {
			t.Port = types.Int64Null()
		}
		
		if target.Username != "" {
			t.Username = types.StringValue(target.Username)
		} else {
			t.Username = types.StringNull()
		}
		
		if target.FilePath != "" {
			t.FilePath = types.StringValue(target.FilePath)
		} else {
			t.FilePath = types.StringNull()
		}
		
		if target.AuthMode != "" {
			t.AuthMode = types.StringValue(target.AuthMode)
		} else {
			t.AuthMode = types.StringNull()
		}
		
		if target.Description != "" {
			t.Description = types.StringValue(target.Description)
		} else {
			t.Description = types.StringNull()
		}
		
		tags, d := types.SetValueFrom(ctx, types.StringType, &target.Tags)
		diags.Append(d...)
		t.Tags = tags
		
		if len(target.AllowedPipelines) > 0 {
			allowedPipelines, d := types.SetValueFrom(ctx, types.Int64Type, &target.AllowedPipelines)
			diags.Append(d...)
			t.AllowedPipelines = allowedPipelines
		} else {
			t.AllowedPipelines = types.SetNull(types.Int64Type)
		}
		
		targetsList = append(targetsList, t)
	}
	
	targetsSet, d := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: targetModelAttrTypes()}, &targetsList)
	diags.Append(d...)
	s.Targets = targetsSet
	
	return diags
}

func (s *targetsSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_targets"
}

func (s *targetsSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List targets\n\n" +
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
				MarkdownDescription: "The project's name. If specified, only targets from this project will be returned",
				Optional:            true,
				Validators:          util.StringValidatorsSlug(),
			},
			"name_regex": schema.StringAttribute{
				MarkdownDescription: "The regular expression to match target names",
				Optional:            true,
				Validators:          util.StringValidatorsRegexp(),
			},
			"targets": schema.SetNestedAttribute{
				MarkdownDescription: "List of targets",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: targetModelAttributes(),
				},
			},
		},
	}
}

func (s *targetsSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *targetsSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *targetsSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	domain := data.Domain.ValueString()
	projectName := data.ProjectName.ValueString()
	
	var targets *buddy.Targets
	var err error
	
	if projectName != "" {
		targets, _, err = s.client.TargetService.GetListInProject(domain, projectName)
	} else {
		targets, _, err = s.client.TargetService.GetListInWorkspace(domain)
	}
	
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get targets", err))
		return
	}
	
	var result []buddy.Target
	
	if !data.NameRegex.IsNull() {
		result, err = util.FilterTargetsByName(&targets.Targets, data.NameRegex.ValueString())
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("name_regex"),
				"Invalid regular expression",
				err.Error(),
			)
			return
		}
	} else {
		result = targets.Targets
	}
	
	resp.Diagnostics.Append(data.loadAPI(ctx, domain, projectName, &result)...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type targetModel struct {
	Name                types.String `tfsdk:"name"`
	TargetId            types.Int64  `tfsdk:"target_id"`
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

func targetModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			MarkdownDescription: "The target's name",
			Computed:            true,
		},
		"target_id": schema.Int64Attribute{
			MarkdownDescription: "The target's ID",
			Computed:            true,
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
	}
}

func targetModelAttrTypes() map[string]types.Type {
	return map[string]types.Type{
		"name":                  types.StringType,
		"target_id":             types.Int64Type,
		"type":                  types.StringType,
		"hostname":              types.StringType,
		"port":                  types.Int64Type,
		"username":              types.StringType,
		"file_path":             types.StringType,
		"auth_mode":             types.StringType,
		"tags":                  types.SetType{ElemType: types.StringType},
		"description":           types.StringType,
		"all_pipelines_allowed": types.BoolType,
		"allowed_pipelines":     types.SetType{ElemType: types.Int64Type},
		"html_url":              types.StringType,
	}
}
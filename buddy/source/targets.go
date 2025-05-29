package source

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ datasource.DataSource              = &targetsSource{}
	_ datasource.DataSourceWithConfigure = &targetsSource{}
)

func NewTargetsSource() datasource.DataSource {
	return &targetsSource{}
}

type targetsSource struct {
	client *buddy.Client
}

type targetsSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Domain    types.String `tfsdk:"domain"`
	NameRegex types.String `tfsdk:"name_regex"`
	Targets   types.List   `tfsdk:"targets"`
}

type targetModel struct {
	Name         types.String `tfsdk:"name"`
	TargetId     types.String `tfsdk:"target_id"`
	Type         types.String `tfsdk:"type"`
	Host         types.String `tfsdk:"host"`
	Port         types.String `tfsdk:"port"`
	Path         types.String `tfsdk:"path"`
	Secure       types.Bool   `tfsdk:"secure"`
	Scope        types.String `tfsdk:"scope"`
	Repository   types.String `tfsdk:"repository"`
	Integration  types.String `tfsdk:"integration"`
	Tags         types.Set    `tfsdk:"tags"`
	Disabled     types.Bool   `tfsdk:"disabled"`
	AuthMethod   types.String `tfsdk:"auth_method"`
	AuthUsername types.String `tfsdk:"auth_username"`
	HtmlUrl      types.String `tfsdk:"html_url"`
}

func (tm *targetModel) loadAPI(ctx context.Context, target *buddy.Target) diag.Diagnostics {
	var diags diag.Diagnostics

	tm.TargetId = types.StringValue(target.Id)
	tm.Name = types.StringValue(target.Name)
	tm.Type = types.StringValue(target.Type)
	tm.HtmlUrl = types.StringValue(target.HtmlUrl)
	tm.Scope = types.StringValue(target.Scope)
	tm.Disabled = types.BoolValue(target.Disabled)
	tm.Secure = types.BoolValue(target.Secure)

	if target.Host != "" {
		tm.Host = types.StringValue(target.Host)
	} else {
		tm.Host = types.StringNull()
	}

	if target.Port != "" {
		tm.Port = types.StringValue(target.Port)
	} else {
		tm.Port = types.StringNull()
	}

	if target.Path != "" {
		tm.Path = types.StringValue(target.Path)
	} else {
		tm.Path = types.StringNull()
	}

	if target.Repository != "" {
		tm.Repository = types.StringValue(target.Repository)
	} else {
		tm.Repository = types.StringNull()
	}

	if target.Integration != "" {
		tm.Integration = types.StringValue(target.Integration)
	} else {
		tm.Integration = types.StringNull()
	}

	tags, d := types.SetValueFrom(ctx, types.StringType, &target.Tags)
	diags.Append(d...)
	tm.Tags = tags

	// Handle auth fields if auth is present
	if target.Auth != nil {
		tm.AuthMethod = types.StringValue(target.Auth.Method)

		if target.Auth.Username != "" {
			tm.AuthUsername = types.StringValue(target.Auth.Username)
		} else {
			tm.AuthUsername = types.StringNull()
		}
	} else {
		tm.AuthMethod = types.StringNull()
		tm.AuthUsername = types.StringNull()
	}

	return diags
}

func (s *targetsSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_targets"
}

func (s *targetsSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *targetsSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List targets and optionally filter by name\n\n" +
			"Token scope required: `WORKSPACE`",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Terraform resource identifier for this item",
				Computed:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The workspace's URL handle",
				Required:            true,
			},
			"name_regex": schema.StringAttribute{
				MarkdownDescription: "The target's name regular expression to match",
				Optional:            true,
			},
			"targets": schema.ListNestedAttribute{
				MarkdownDescription: "List of targets",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: targetSchemaAttributes(),
				},
			},
		},
	}
}

func (s *targetsSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data targetsSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := data.Domain.ValueString()

	targets, _, err := s.client.TargetService.GetList(domain, &buddy.TargetGetListQuery{})
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get targets", err))
		return
	}

	filteredTargets := targets.Targets

	if !data.NameRegex.IsNull() && data.NameRegex.ValueString() != "" {
		filteredTargets = util.FilterTargetListByNameRegex(targets.Targets, data.NameRegex.ValueString())
	}

	targetModels := make([]targetModel, len(filteredTargets))
	for i, target := range filteredTargets {
		resp.Diagnostics.Append(targetModels[i].loadAPI(ctx, target)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	data.Targets, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: targetModelAttrTypes()}, &targetModels)

	data.ID = types.StringValue(domain)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func targetSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			MarkdownDescription: "The target's name",
			Computed:            true,
		},
		"target_id": schema.StringAttribute{
			MarkdownDescription: "The target's ID",
			Computed:            true,
		},
		"type": schema.StringAttribute{
			MarkdownDescription: "The target's type",
			Computed:            true,
		},
		"host": schema.StringAttribute{
			MarkdownDescription: "The target's hostname or IP address",
			Computed:            true,
		},
		"port": schema.StringAttribute{
			MarkdownDescription: "The target's port",
			Computed:            true,
		},
		"path": schema.StringAttribute{
			MarkdownDescription: "The remote path on the target",
			Computed:            true,
		},
		"secure": schema.BoolAttribute{
			MarkdownDescription: "Whether to use secure connection",
			Computed:            true,
		},
		"scope": schema.StringAttribute{
			MarkdownDescription: "The target's scope",
			Computed:            true,
		},
		"repository": schema.StringAttribute{
			MarkdownDescription: "The repository for registry targets",
			Computed:            true,
		},
		"integration": schema.StringAttribute{
			MarkdownDescription: "The integration ID used for cloud targets",
			Computed:            true,
		},
		"tags": schema.SetAttribute{
			MarkdownDescription: "The target's tags",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"disabled": schema.BoolAttribute{
			MarkdownDescription: "Whether the target is disabled",
			Computed:            true,
		},
		"auth_method": schema.StringAttribute{
			MarkdownDescription: "The authentication method",
			Computed:            true,
		},
		"auth_username": schema.StringAttribute{
			MarkdownDescription: "The authentication username",
			Computed:            true,
		},
		"html_url": schema.StringAttribute{
			MarkdownDescription: "The target's URL",
			Computed:            true,
		},
	}
}

func targetModelAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":          types.StringType,
		"target_id":     types.StringType,
		"type":          types.StringType,
		"host":          types.StringType,
		"port":          types.StringType,
		"path":          types.StringType,
		"secure":        types.BoolType,
		"scope":         types.StringType,
		"repository":    types.StringType,
		"integration":   types.StringType,
		"tags":          types.SetType{ElemType: types.StringType},
		"disabled":      types.BoolType,
		"auth_method":   types.StringType,
		"auth_username": types.StringType,
		"html_url":      types.StringType,
	}
}

package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	sourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type pipelineModel struct {
	Name                    types.String `tfsdk:"name"`
	PipelineId              types.Int64  `tfsdk:"pipeline_id"`
	HtmlUrl                 types.String `tfsdk:"html_url"`
	On                      types.String `tfsdk:"on"`
	Cpu                     types.String `tfsdk:"cpu"`
	Priority                types.String `tfsdk:"priority"`
	Disabled                types.Bool   `tfsdk:"disabled"`
	DisablingReason         types.String `tfsdk:"disabling_reason"`
	LastExecutionStatus     types.String `tfsdk:"last_execution_status"`
	LastExecutionRevision   types.String `tfsdk:"last_execution_revision"`
	Refs                    types.Set    `tfsdk:"refs"`
	Tags                    types.Set    `tfsdk:"tags"`
	Event                   types.Set    `tfsdk:"event"`
	GitConfigRef            types.String `tfsdk:"git_config_ref"`
	GitConfig               types.Object `tfsdk:"git_config"`
	ConcurrentPipelineRuns  types.Bool   `tfsdk:"concurrent_pipeline_runs"`
	DescriptionRequired     types.Bool   `tfsdk:"description_required"`
	GitChangesetBase        types.String `tfsdk:"git_changeset_base"`
	FilesystemChangesetBase types.String `tfsdk:"filesystem_changeset_base"`
	DefinitionSource        types.String `tfsdk:"definition_source"`
	RemoteProjectName       types.String `tfsdk:"remote_project_name"`
	RemoteBranch            types.String `tfsdk:"remote_branch"`
	RemotePath              types.String `tfsdk:"remote_path"`
	RemoteParameter         types.Set    `tfsdk:"remote_parameter"`
}

func pipelineModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"name":                      types.StringType,
		"pipeline_id":               types.Int64Type,
		"html_url":                  types.StringType,
		"on":                        types.StringType,
		"cpu":                       types.StringType,
		"priority":                  types.StringType,
		"disabled":                  types.BoolType,
		"disabling_reason":          types.StringType,
		"last_execution_status":     types.StringType,
		"last_execution_revision":   types.StringType,
		"refs":                      types.SetType{ElemType: types.StringType},
		"tags":                      types.SetType{ElemType: types.StringType},
		"event":                     types.SetType{ElemType: types.ObjectType{AttrTypes: eventModelAttrs()}},
		"concurrent_pipeline_runs":  types.BoolType,
		"description_required":      types.BoolType,
		"git_changeset_base":        types.StringType,
		"filesystem_changeset_base": types.StringType,
		"git_config_ref":            types.StringType,
		"git_config":                types.ObjectType{AttrTypes: GitConfigModelAttrs()},
		"definition_source":         types.StringType,
		"remote_project_name":       types.StringType,
		"remote_branch":             types.StringType,
		"remote_path":               types.StringType,
		"remote_parameter":          types.SetType{ElemType: types.ObjectType{AttrTypes: remoteParameterModelAttrs()}},
	}
}

func (p *pipelineModel) loadAPI(ctx context.Context, pipeline *buddy.Pipeline) diag.Diagnostics {
	var diags diag.Diagnostics
	p.Name = types.StringValue(pipeline.Name)
	p.PipelineId = types.Int64Value(int64(pipeline.Id))
	p.HtmlUrl = types.StringValue(pipeline.HtmlUrl)
	p.On = types.StringValue(pipeline.On)
	p.Cpu = types.StringValue(pipeline.Cpu)
	p.Priority = types.StringValue(pipeline.Priority)
	p.Disabled = types.BoolValue(pipeline.Disabled)
	p.DisablingReason = types.StringValue(pipeline.DisabledReason)
	p.LastExecutionStatus = types.StringValue(pipeline.LastExecutionStatus)
	p.LastExecutionRevision = types.StringValue(pipeline.LastExecutionRevision)
	p.ConcurrentPipelineRuns = types.BoolValue(pipeline.ConcurrentPipelineRuns)
	p.DescriptionRequired = types.BoolValue(pipeline.DescriptionRequired)
	p.GitChangesetBase = types.StringValue(pipeline.GitChangesetBase)
	p.FilesystemChangesetBase = types.StringValue(pipeline.FilesystemChangesetBase)
	r, d := types.SetValueFrom(ctx, types.StringType, &pipeline.Refs)
	diags.Append(d...)
	p.Refs = r
	t, d := types.SetValueFrom(ctx, types.StringType, &pipeline.Tags)
	diags.Append(d...)
	p.Tags = t
	e, d := EventsModelFromApi(ctx, &pipeline.Events)
	diags.Append(d...)
	p.Event = e
	p.GitConfigRef = types.StringValue(pipeline.GitConfigRef)
	gitConfig, d := GitConfigModelFromApi(ctx, pipeline.GitConfig)
	diags.Append(d...)
	p.GitConfig = gitConfig
	p.DefinitionSource = types.StringValue(GetPipelineDefinitionSource(pipeline))
	p.RemoteProjectName = types.StringValue(pipeline.RemoteProjectName)
	p.RemoteBranch = types.StringValue(pipeline.RemoteBranch)
	p.RemotePath = types.StringValue(pipeline.RemotePath)
	rp, d := RemoteParametersModelFromApi(ctx, &pipeline.RemoteParameters)
	diags.Append(d...)
	p.RemoteParameter = rp
	return diags
}

func SourcePipelineModelAttributes() map[string]sourceschema.Attribute {
	return map[string]sourceschema.Attribute{
		"name": sourceschema.StringAttribute{
			Computed: true,
		},
		"pipeline_id": sourceschema.Int64Attribute{
			Computed: true,
		},
		"html_url": sourceschema.StringAttribute{
			Computed: true,
		},
		"on": sourceschema.StringAttribute{
			Computed: true,
		},
		"cpu": sourceschema.StringAttribute{
			Computed: true,
		},
		"priority": sourceschema.StringAttribute{
			Computed: true,
		},
		"disabled": sourceschema.BoolAttribute{
			Computed: true,
		},
		"disabling_reason": sourceschema.StringAttribute{
			Computed: true,
		},
		"concurrent_pipeline_runs": sourceschema.BoolAttribute{
			Computed: true,
		},
		"description_required": sourceschema.BoolAttribute{
			Computed: true,
		},
		"git_changeset_base": sourceschema.StringAttribute{
			Computed: true,
		},
		"filesystem_changeset_base": sourceschema.StringAttribute{
			Computed: true,
		},
		"last_execution_status": sourceschema.StringAttribute{
			Computed: true,
		},
		"last_execution_revision": sourceschema.StringAttribute{
			Computed: true,
		},
		"refs": sourceschema.SetAttribute{
			Computed:    true,
			ElementType: types.StringType,
		},
		"event": sourceschema.SetNestedAttribute{
			Computed: true,
			NestedObject: sourceschema.NestedAttributeObject{
				Attributes: SourceEventModelAttributes(),
			},
		},
		"tags": sourceschema.SetAttribute{
			Computed:    true,
			ElementType: types.StringType,
		},
		"git_config_ref": sourceschema.StringAttribute{
			Computed: true,
		},
		"git_config": sourceschema.ObjectAttribute{
			Computed:       true,
			AttributeTypes: GitConfigModelAttrs(),
		},
		"definition_source": sourceschema.StringAttribute{
			Computed: true,
		},
		"remote_project_name": sourceschema.StringAttribute{
			Computed: true,
		},
		"remote_branch": sourceschema.StringAttribute{
			Computed: true,
		},
		"remote_path": sourceschema.StringAttribute{
			Computed: true,
		},
		"remote_parameter": sourceschema.SetNestedAttribute{
			Computed: true,
			NestedObject: sourceschema.NestedAttributeObject{
				Attributes: SourceRemoteParameterModelAttributes(),
			},
		},
	}
}

func PipelinesModelFromApi(ctx context.Context, pipelines *[]*buddy.Pipeline) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	p := make([]*pipelineModel, len(*pipelines))
	for i, v := range *pipelines {
		p[i] = &pipelineModel{}
		diags.Append(p[i].loadAPI(ctx, v)...)
	}
	r, d := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: pipelineModelAttrs()}, &p)
	diags.Append(d...)
	return r, diags
}

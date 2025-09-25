package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	sourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type eventModel struct {
	Type      types.String `tfsdk:"type"`
	Refs      types.Set    `tfsdk:"refs"`
	Branches  types.Set    `tfsdk:"branches"`
	Events    types.Set    `tfsdk:"events"`
	StartDate types.String `tfsdk:"start_date"`
	Delay     types.Int64  `tfsdk:"delay"`
	Cron      types.String `tfsdk:"cron"`
	Totp      types.Bool   `tfsdk:"totp"`
	Timezone  types.String `tfsdk:"timezone"`
}

func eventModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"type":       types.StringType,
		"refs":       types.SetType{ElemType: types.StringType},
		"events":     types.SetType{ElemType: types.StringType},
		"branches":   types.SetType{ElemType: types.StringType},
		"start_date": types.StringType,
		"delay":      types.Int64Type,
		"cron":       types.StringType,
		"totp":       types.BoolType,
		"timezone":   types.StringType,
	}
}

func (e *eventModel) loadAPI(ctx context.Context, event *buddy.PipelineEvent) diag.Diagnostics {
	var diags diag.Diagnostics
	e.Type = types.StringValue(event.Type)
	e.StartDate = types.StringValue(event.StartDate)
	e.Delay = types.Int64Value(int64(event.Delay))
	e.Cron = types.StringValue(event.Cron)
	e.Timezone = types.StringValue(event.Timezone)
	e.Totp = types.BoolValue(event.Totp)
	r, d1 := types.SetValueFrom(ctx, types.StringType, &event.Refs)
	e.Refs = r
	diags.Append(d1...)
	b, d2 := types.SetValueFrom(ctx, types.StringType, &event.Branches)
	diags.Append(d2...)
	e.Branches = b
	ev, d3 := types.SetValueFrom(ctx, types.StringType, &event.Events)
	diags.Append(d3...)
	e.Events = ev
	return diags
}

func SourceEventModelAttributes() map[string]sourceschema.Attribute {
	return map[string]sourceschema.Attribute{
		"type": sourceschema.StringAttribute{
			Computed: true,
		},
		"start_date": sourceschema.StringAttribute{
			Computed: true,
		},
		"delay": sourceschema.Int64Attribute{
			Computed: true,
		},
		"cron": sourceschema.StringAttribute{
			Computed: true,
		},
		"timezone": sourceschema.StringAttribute{
			Computed: true,
		},
		"refs": sourceschema.SetAttribute{
			Computed:    true,
			ElementType: types.StringType,
		},
		"totp": sourceschema.BoolAttribute{
			Computed: true,
		},
		"branches": sourceschema.SetAttribute{
			Computed:    true,
			ElementType: types.StringType,
		},
		"events": sourceschema.SetAttribute{
			Computed:    true,
			ElementType: types.StringType,
		},
	}
}

func ResourceEventModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					buddy.PipelineEventTypePush,
					buddy.PipelineEventTypeSchedule,
					buddy.PipelineEventTypeCreateRef,
					buddy.PipelineEventTypeDeleteRef,
					buddy.PipelineEventTypePullRequest,
					buddy.PipelineEventTypeWebhook,
				),
			},
		},
		"start_date": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.Expressions{
					path.MatchRelative().AtParent().AtName("cron"),
				}...),
				stringvalidator.AlsoRequires(path.Expressions{
					path.MatchRelative().AtParent().AtName("delay"),
				}...),
			},
		},
		"delay": schema.Int64Attribute{
			Optional: true,
			Validators: []validator.Int64{
				int64validator.ConflictsWith(path.Expressions{
					path.MatchRelative().AtParent().AtName("cron"),
				}...),
				int64validator.AlsoRequires(path.Expressions{
					path.MatchRelative().AtParent().AtName("start_date"),
				}...),
			},
		},
		"cron": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.Expressions{
					path.MatchRelative().AtParent().AtName("delay"),
					path.MatchRelative().AtParent().AtName("start_date"),
				}...),
			},
		},
		"timezone": schema.StringAttribute{
			Optional: true,
		},
		"refs": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Validators: []validator.Set{
				setvalidator.ConflictsWith(path.Expressions{
					path.MatchRelative().AtParent().AtName("branches"),
					path.MatchRelative().AtParent().AtName("events"),
				}...),
			},
		},
		"branches": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Validators: []validator.Set{
				setvalidator.ConflictsWith(path.Expressions{
					path.MatchRelative().AtParent().AtName("refs"),
				}...),
			},
		},
		"totp": schema.BoolAttribute{
			Optional: true,
		},
		"events": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Validators: []validator.Set{
				setvalidator.ConflictsWith(path.Expressions{
					path.MatchRelative().AtParent().AtName("refs"),
				}...),
				setvalidator.ValueStringsAre(stringvalidator.OneOf(
					buddy.PipelinePullRequestEventOpened,
					buddy.PipelinePullRequestEventEdited,
					buddy.PipelinePullRequestEventClosed,
					buddy.PipelinePullRequestEventReopened,
					buddy.PipelinePullRequestEventSynchronize,
					buddy.PipelinePullRequestEventConvertedToDraft,
					buddy.PipelinePullRequestEventLocked,
					buddy.PipelinePullRequestEventUnlocked,
					buddy.PipelinePullRequestEventEnqueued,
					buddy.PipelinePullRequestEventDequeued,
					buddy.PipelinePullRequestEventMilestoned,
					buddy.PipelinePullRequestEventDemilestoned,
					buddy.PipelinePullRequestEventReadyForReview,
					buddy.PipelinePullRequestEventReviewRequested,
					buddy.PipelinePullRequestEventReviewRequestRemoved,
					buddy.PipelinePullRequestEventAutoMergeEnabled,
					buddy.PipelinePullRequestEventAutoMergeDisabled,
				)),
			},
		},
	}
}

func EventsModelFromApi(ctx context.Context, events *[]*buddy.PipelineEvent) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	r := make([]*eventModel, len(*events))
	for i, v := range *events {
		r[i] = &eventModel{}
		diags.Append(r[i].loadAPI(ctx, v)...)
	}
	result, d := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: eventModelAttrs()}, &r)
	diags.Append(d...)
	return result, d
}

func EventsModelToApi(ctx context.Context, s *types.Set) (*[]*buddy.PipelineEvent, diag.Diagnostics) {
	var em []eventModel
	diags := s.ElementsAs(ctx, &em, false)
	pipelineEvents := make([]*buddy.PipelineEvent, len(em))
	for i, v := range em {
		pe := &buddy.PipelineEvent{}
		if !v.Type.IsNull() && !v.Type.IsUnknown() {
			pe.Type = v.Type.ValueString()
		}
		if !v.StartDate.IsNull() && !v.StartDate.IsUnknown() {
			pe.StartDate = v.StartDate.ValueString()
		}
		if !v.Delay.IsNull() && !v.Delay.IsUnknown() {
			pe.Delay = int(v.Delay.ValueInt64())
		}
		if !v.Cron.IsNull() && !v.Cron.IsUnknown() {
			pe.Cron = v.Cron.ValueString()
		}
		if !v.Timezone.IsNull() && !v.Timezone.IsUnknown() {
			pe.Timezone = v.Timezone.ValueString()
		}
		if !v.Refs.IsNull() && !v.Refs.IsUnknown() {
			refs, d := StringSetToApi(ctx, &v.Refs)
			diags.Append(d...)
			pe.Refs = *refs
		} else {
			pe.Refs = []string{}
		}
		if !v.Branches.IsNull() && !v.Branches.IsUnknown() {
			branches, d := StringSetToApi(ctx, &v.Branches)
			diags.Append(d...)
			pe.Branches = *branches
		} else {
			pe.Branches = []string{}
		}
		if !v.Events.IsNull() && !v.Events.IsUnknown() {
			events, d := StringSetToApi(ctx, &v.Events)
			diags.Append(d...)
			pe.Events = *events
		} else {
			pe.Events = []string{}
		}
		if !v.Totp.IsNull() && !v.Totp.IsUnknown() {
			pe.Totp = v.Totp.ValueBool()
		}
		pipelineEvents[i] = pe
	}
	return &pipelineEvents, diags
}

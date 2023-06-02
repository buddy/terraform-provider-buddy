package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	sourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type eventModel struct {
	Type types.String `tfsdk:"type"`
	Refs types.Set    `tfsdk:"refs"`
}

func eventModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"type": types.StringType,
		"refs": types.SetType{ElemType: types.StringType},
	}
}

func (e *eventModel) loadAPI(ctx context.Context, event *buddy.PipelineEvent) diag.Diagnostics {
	e.Type = types.StringValue(event.Type)
	r, d := types.SetValueFrom(ctx, types.StringType, &event.Refs)
	e.Refs = r
	return d
}

func SourceEventModelAttributes() map[string]sourceschema.Attribute {
	return map[string]sourceschema.Attribute{
		"type": sourceschema.StringAttribute{
			Computed: true,
		},
		"refs": sourceschema.SetAttribute{
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
					buddy.PipelineEventTypeCreateRef,
					buddy.PipelineEventTypeDeleteRef,
				),
			},
		},
		"refs": schema.SetAttribute{
			ElementType: types.StringType,
			Required:    true,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
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
		if !v.Refs.IsNull() && !v.Refs.IsUnknown() {
			refs, d := StringSetToApi(ctx, &v.Refs)
			diags.Append(d...)
			pe.Refs = *refs
		} else {
			pe.Refs = []string{}
		}
		pipelineEvents[i] = pe
	}
	return &pipelineEvents, diags
}

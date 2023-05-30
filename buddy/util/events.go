package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type eventModel struct {
	Type types.String `tfsdk:"type"`
	Refs types.Set    `tfsdk:"refs"`
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

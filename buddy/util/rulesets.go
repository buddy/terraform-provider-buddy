package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	sourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type rulesetBypassEntryModel struct {
	Id     types.Int64  `tfsdk:"id"`
	Access types.String `tfsdk:"access"`
}

func rulesetBypassEntryModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"id":     types.Int64Type,
		"access": types.StringType,
	}
}

func ResourceRulesetBypassEntryModelAttributes(of string) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			MarkdownDescription: "The ID of the " + of,
			Required:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
		"access": schema.StringAttribute{
			MarkdownDescription: "`ALLOWED` - may bypass the ruleset, `DENIED` - may not, even when an entry checked later allows it",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.OneOf(buddy.RulesetBypassAccessAllowed, buddy.RulesetBypassAccessDenied),
			},
		},
	}
}

func SourceRulesetBypassEntryModelAttributes() map[string]sourceschema.Attribute {
	return map[string]sourceschema.Attribute{
		"id": sourceschema.Int64Attribute{
			Computed: true,
		},
		"access": sourceschema.StringAttribute{
			Computed: true,
		},
	}
}

func RulesetBypassEntriesToApi(ctx context.Context, s *types.Set) ([]*buddy.RulesetBypassEntry, diag.Diagnostics) {
	if s.IsNull() || s.IsUnknown() {
		return nil, nil
	}
	var models []*rulesetBypassEntryModel
	diags := s.ElementsAs(ctx, &models, false)
	entries := make([]*buddy.RulesetBypassEntry, len(models))
	for i, m := range models {
		entries[i] = &buddy.RulesetBypassEntry{
			Id:     int(m.Id.ValueInt64()),
			Access: m.Access.ValueString(),
		}
	}
	return entries, diags
}

// RulesetBypassEntriesFromApi always returns a set, empty without entries: Terraform keeps set blocks left out of the
// configuration as an empty set, never null
func RulesetBypassEntriesFromApi(ctx context.Context, entries []*buddy.RulesetBypassEntry) (basetypes.SetValue, diag.Diagnostics) {
	t := types.ObjectType{AttrTypes: rulesetBypassEntryModelAttrs()}
	l := make([]*rulesetBypassEntryModel, len(entries))
	for i, e := range entries {
		l[i] = &rulesetBypassEntryModel{
			Id:     types.Int64Value(int64(e.Id)),
			Access: types.StringValue(e.Access),
		}
	}
	return types.SetValueFrom(ctx, t, &l)
}

type rulesetSourceModel struct {
	Type       types.String `tfsdk:"type"`
	Identifier types.String `tfsdk:"identifier"`
}

func rulesetSourceModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"type":       types.StringType,
		"identifier": types.StringType,
	}
}

func ResourceRulesetSourceModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"type": schema.StringAttribute{
			MarkdownDescription: "`PIPELINE`, `UNIT_TESTS_SUITE` or `VISUAL_TESTS_SUITE`",
			Computed:            true,
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier of the source",
			Computed:            true,
		},
	}
}

func SourceRulesetSourceModelAttributes() map[string]sourceschema.Attribute {
	return map[string]sourceschema.Attribute{
		"type": sourceschema.StringAttribute{
			Computed: true,
		},
		"identifier": sourceschema.StringAttribute{
			Computed: true,
		},
	}
}

func RulesetSourcesFromApi(ctx context.Context, sources []*buddy.RulesetSource) (basetypes.SetValue, diag.Diagnostics) {
	l := make([]*rulesetSourceModel, len(sources))
	for i, s := range sources {
		l[i] = &rulesetSourceModel{
			Type:       types.StringValue(s.Type),
			Identifier: types.StringValue(s.Identifier),
		}
	}
	return types.SetValueFrom(ctx, types.ObjectType{AttrTypes: rulesetSourceModelAttrs()}, &l)
}

// StringSetFromApiOrNull returns null for an empty list, the configuration cannot hold an empty set (SizeAtLeast(1))
func StringSetFromApiOrNull(ctx context.Context, l []string) (basetypes.SetValue, diag.Diagnostics) {
	if len(l) == 0 {
		return types.SetNull(types.StringType), nil
	}
	return types.SetValueFrom(ctx, types.StringType, &l)
}

// StringFromApiOrNull returns null for an empty string, the API leaves out a cleared value
func StringFromApiOrNull(s string) basetypes.StringValue {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

func StringSetToApiOrNil(ctx context.Context, s *types.Set) ([]string, diag.Diagnostics) {
	if s.IsNull() || s.IsUnknown() {
		return nil, nil
	}
	var l []string
	diags := s.ElementsAs(ctx, &l, false)
	return l, diags
}

type rulesetModel struct {
	HtmlUrl    types.String `tfsdk:"html_url"`
	RulesetId  types.String `tfsdk:"ruleset_id"`
	Identifier types.String `tfsdk:"identifier"`
	Name       types.String `tfsdk:"name"`
	Disabled   types.Bool   `tfsdk:"disabled"`
	Refs       types.Set    `tfsdk:"refs"`
}

func rulesetModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"html_url":   types.StringType,
		"ruleset_id": types.StringType,
		"identifier": types.StringType,
		"name":       types.StringType,
		"disabled":   types.BoolType,
		"refs":       types.SetType{ElemType: types.StringType},
	}
}

func SourceRulesetModelAttributes() map[string]sourceschema.Attribute {
	return map[string]sourceschema.Attribute{
		"html_url": sourceschema.StringAttribute{
			Computed: true,
		},
		"ruleset_id": sourceschema.StringAttribute{
			Computed: true,
		},
		"identifier": sourceschema.StringAttribute{
			Computed: true,
		},
		"name": sourceschema.StringAttribute{
			Computed: true,
		},
		"disabled": sourceschema.BoolAttribute{
			Computed: true,
		},
		"refs": sourceschema.SetAttribute{
			ElementType: types.StringType,
			Computed:    true,
		},
	}
}

func RulesetsModelFromApi(ctx context.Context, rulesets []*buddy.Ruleset) (basetypes.SetValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	l := make([]*rulesetModel, len(rulesets))
	for i, r := range rulesets {
		refs, d := types.SetValueFrom(ctx, types.StringType, &r.Refs)
		diags.Append(d...)
		l[i] = &rulesetModel{
			HtmlUrl:    types.StringValue(r.HtmlUrl),
			RulesetId:  types.StringValue(r.Id),
			Identifier: types.StringValue(r.Identifier),
			Name:       StringFromApiOrNull(r.Name),
			Disabled:   types.BoolValue(r.Disabled),
			Refs:       refs,
		}
	}
	s, d := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: rulesetModelAttrs()}, &l)
	diags.Append(d...)
	return s, diags
}

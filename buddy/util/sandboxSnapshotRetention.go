package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	sourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type sandboxSnapshotRetentionModel struct {
	KeepLast    types.Int32 `tfsdk:"keep_last"`
	KeepForDays types.Int32 `tfsdk:"keep_for_days"`
}

func SandboxSnapshotRetentionModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"keep_last":     types.Int32Type,
		"keep_for_days": types.Int32Type,
	}
}

func ResourceSandboxSnapshotRetentionModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"keep_last": schema.Int32Attribute{
			MarkdownDescription: "Keep only the N most recent snapshots. Conflicts with `keep_for_days`",
			Optional:            true,
			Validators: []validator.Int32{
				int32validator.AtLeast(1),
				int32validator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("keep_for_days")),
			},
		},
		"keep_for_days": schema.Int32Attribute{
			MarkdownDescription: "Keep only the snapshots taken within the last N days. Conflicts with `keep_last`",
			Optional:            true,
			Validators: []validator.Int32{
				int32validator.AtLeast(1),
			},
		},
	}
}

func SourceSandboxSnapshotRetentionModelAttributes() map[string]sourceschema.Attribute {
	return map[string]sourceschema.Attribute{
		"keep_last": sourceschema.Int32Attribute{
			MarkdownDescription: "Keep only the N most recent snapshots",
			Computed:            true,
		},
		"keep_for_days": sourceschema.Int32Attribute{
			MarkdownDescription: "Keep only the snapshots taken within the last N days",
			Computed:            true,
		},
	}
}

// SandboxSnapshotRetentionToApi returns an empty retention for a null object, which removes the retention
func SandboxSnapshotRetentionToApi(ctx context.Context, o *types.Object) (*buddy.SandboxSnapshotRetention, diag.Diagnostics) {
	retention := &buddy.SandboxSnapshotRetention{}
	if o.IsNull() || o.IsUnknown() {
		return retention, nil
	}
	var m sandboxSnapshotRetentionModel
	diags := o.As(ctx, &m, basetypes.ObjectAsOptions{})
	if !m.KeepLast.IsNull() && !m.KeepLast.IsUnknown() {
		keepLast := int(m.KeepLast.ValueInt32())
		retention.KeepLast = &keepLast
	}
	if !m.KeepForDays.IsNull() && !m.KeepForDays.IsUnknown() {
		keepForDays := int(m.KeepForDays.ValueInt32())
		retention.KeepForDays = &keepForDays
	}
	return retention, diags
}

// SandboxSnapshotRetentionFromApi returns a null object when the sandbox has no retention
func SandboxSnapshotRetentionFromApi(ctx context.Context, retention *buddy.SandboxSnapshotRetention) (basetypes.ObjectValue, diag.Diagnostics) {
	if retention == nil || (retention.KeepLast == nil && retention.KeepForDays == nil) {
		return types.ObjectNull(SandboxSnapshotRetentionModelAttrs()), nil
	}
	m := sandboxSnapshotRetentionModel{
		KeepLast:    types.Int32Null(),
		KeepForDays: types.Int32Null(),
	}
	if retention.KeepLast != nil {
		m.KeepLast = types.Int32Value(int32(*retention.KeepLast))
	}
	if retention.KeepForDays != nil {
		m.KeepForDays = types.Int32Value(int32(*retention.KeepForDays))
	}
	return types.ObjectValueFrom(ctx, SandboxSnapshotRetentionModelAttrs(), &m)
}

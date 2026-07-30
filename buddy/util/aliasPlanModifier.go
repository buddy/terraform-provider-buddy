package util

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.String = aliasPlanModifier{}

type aliasPlanModifier struct {
	alias path.Path
}

func (m aliasPlanModifier) Description(_ context.Context) string {
	return fmt.Sprintf("value follows %s if %s is set in the configuration", m.alias, m.alias)
}

func (m aliasPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m aliasPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// resource destroy
	if req.Config.Raw.IsNull() {
		return
	}
	// attribute set in the configuration, keep its value
	if !req.ConfigValue.IsNull() {
		return
	}
	var alias types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, m.alias, &alias)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// neither the attribute nor its alias set in the configuration, keep the value from the state
	if alias.IsNull() {
		return
	}
	resp.PlanValue = alias
}

// AliasPlanModifier keeps two attributes backed by the same API field in sync while planning.
// If the attribute is not set in the configuration but its alias is, the planned value follows
// the alias instead of the value kept in the state
func AliasPlanModifier(alias string) planmodifier.String {
	return aliasPlanModifier{alias: path.Root(alias)}
}

// StringValidatorsAlias makes setting an attribute together with its alias an error
func StringValidatorsAlias(alias string) []validator.String {
	return []validator.String{
		stringvalidator.ConflictsWith(path.MatchRoot(alias)),
	}
}

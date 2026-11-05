package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type allowedEnvironmentModel struct {
	Project     types.String `tfsdk:"project"`
	Environment types.String `tfsdk:"environment"`
	AccessLevel types.String `tfsdk:"access_level"`
}

func allowedEnvironmentsModelToApi(ctx context.Context, s *types.Set) ([]allowedEnvironmentModel, diag.Diagnostics) {
	var ss []allowedEnvironmentModel
	diags := s.ElementsAs(ctx, &ss, false)
	return ss, diags
}

func TargetAllowedEnvironmentsModelToApi(ctx context.Context, s *types.Set) (*[]*buddy.TargetAllowedEnvironment, diag.Diagnostics) {
	ss, diags := allowedEnvironmentsModelToApi(ctx, s)
	result := make([]*buddy.TargetAllowedEnvironment, len(ss))
	for i, e := range ss {
		result[i] = &buddy.TargetAllowedEnvironment{
			Project:     e.Project.ValueString(),
			Environment: e.Environment.ValueString(),
			AccessLevel: e.AccessLevel.ValueString(),
		}
	}
	return &result, diags
}

func VariableAllowedEnvironmentsModelToApi(ctx context.Context, s *types.Set) (*[]*buddy.VariableAllowedEnvironment, diag.Diagnostics) {
	ss, diags := allowedEnvironmentsModelToApi(ctx, s)
	result := make([]*buddy.VariableAllowedEnvironment, len(ss))
	for i, e := range ss {
		result[i] = &buddy.VariableAllowedEnvironment{
			Project:     e.Project.ValueString(),
			Environment: e.Environment.ValueString(),
			AccessLevel: e.AccessLevel.ValueString(),
		}
	}
	return &result, diags
}

// AllowedEnvironmentsBlock is the allowed_environments block shared by targets and variables
func AllowedEnvironmentsBlock(description string, accessLevelDescription string, accessLevels ...string) schema.SetNestedBlock {
	return schema.SetNestedBlock{
		MarkdownDescription: description,
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"project": schema.StringAttribute{
					MarkdownDescription: "The environment's project name. Leave empty for a workspace environment",
					Optional:            true,
				},
				"environment": schema.StringAttribute{
					MarkdownDescription: "The environment's identifier",
					Required:            true,
				},
				"access_level": schema.StringAttribute{
					MarkdownDescription: accessLevelDescription,
					Required:            true,
					Validators: []validator.String{
						stringvalidator.OneOf(accessLevels...),
					},
				},
			},
		},
		Validators: []validator.Set{
			setvalidator.SizeAtMost(30),
		},
	}
}

var allowedEnvironmentAttrTypes = map[string]attr.Type{
	"project":      types.StringType,
	"environment":  types.StringType,
	"access_level": types.StringType,
}

// AllowedEnvironment is a rule of allowed_environments, common to targets and variables
type AllowedEnvironment struct {
	Project     string
	Environment string
	AccessLevel string
}

func TargetAllowedEnvironments(rules []*buddy.TargetAllowedEnvironment) []AllowedEnvironment {
	out := make([]AllowedEnvironment, 0, len(rules))
	for _, r := range rules {
		out = append(out, AllowedEnvironment{r.Project, r.Environment, r.AccessLevel})
	}
	return out
}

func VariableAllowedEnvironments(rules []*buddy.VariableAllowedEnvironment) []AllowedEnvironment {
	out := make([]AllowedEnvironment, 0, len(rules))
	for _, r := range rules {
		out = append(out, AllowedEnvironment{r.Project, r.Environment, r.AccessLevel})
	}
	return out
}

// AllowedEnvironmentsApiToModel converts the rules returned by the API to the allowed_environments
// set, no rules give a null set - the same as a config without the blocks
func AllowedEnvironmentsApiToModel(rules []AllowedEnvironment) types.Set {
	elemType := types.ObjectType{AttrTypes: allowedEnvironmentAttrTypes}
	if len(rules) == 0 {
		return types.SetNull(elemType)
	}
	elems := make([]attr.Value, len(rules))
	for i, r := range rules {
		project := types.StringNull()
		if r.Project != "" {
			project = types.StringValue(r.Project)
		}
		elems[i] = types.ObjectValueMust(allowedEnvironmentAttrTypes, map[string]attr.Value{
			"project":      project,
			"environment":  types.StringValue(r.Environment),
			"access_level": types.StringValue(r.AccessLevel),
		})
	}
	return types.SetValueMust(elemType, elems)
}

// CheckAllowedEnvironments compares rules regardless of their order, the API does not keep it
func CheckAllowedEnvironments(field string, got []AllowedEnvironment, want []AllowedEnvironment) error {
	if err := CheckIntFieldEqual("len("+field+")", len(got), len(want)); err != nil {
		return err
	}
	levels := map[string]string{}
	for _, r := range got {
		levels[r.Project+"/"+r.Environment] = r.AccessLevel
	}
	for _, r := range want {
		key := r.Project + "/" + r.Environment
		if err := CheckFieldEqualAndSet(field+"["+key+"].AccessLevel", levels[key], r.AccessLevel); err != nil {
			return err
		}
	}
	return nil
}
